package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/v2rayA/v2rayA/common/httpClient"
	"github.com/v2rayA/v2rayA/common/resolv"
	"github.com/v2rayA/v2rayA/core/coreObj"
	"github.com/v2rayA/v2rayA/core/serverObj"
	"github.com/v2rayA/v2rayA/core/v2ray"
	"github.com/v2rayA/v2rayA/db/configure"
	"github.com/v2rayA/v2rayA/pkg/util/log"
)

const HttpTestURL = "https://gstatic.com/generate_204"

func Ping(which []*configure.Which, timeout time.Duration) (_ []*configure.Which, err error) {
	var whiches = configure.NewWhiches(which)
	//对要Ping的which去重
	which = whiches.GetNonDuplicated()
	//暂时关闭透明代理
	v2ray.ProcessManager.CheckAndStopTransparentProxy(nil)
	defer func() {
		if e := v2ray.ProcessManager.CheckAndSetupTransparentProxy(true, nil); e != nil {
			err = fmt.Errorf("Ping: %v: %v", e, err)
		}
	}()
	//多线程异步ping
	wg := new(sync.WaitGroup)
	for i, v := range which {
		if v.TYPE == configure.SubscriptionType { //subscription不能ping
			continue
		}
		wg.Add(1)
		go func(i int) {
			// Broadcast running status
			v2ray.BroadcastLatencyProgress("ping", which[i], v2ray.LatencyStatusRunning, "")
			_ = which[i].Ping(timeout)
			// Broadcast done status with latency
			latency := which[i].Latency
			status := v2ray.LatencyStatusDone
			if latency == "" || latency == "TIMEOUT" || latency == "SYSTEM ERROR" {
				status = v2ray.LatencyStatusError
			}
			v2ray.BroadcastLatencyProgress("ping", which[i], status, latency)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := len(which) - 1; i >= 0; i-- {
		if which[i].TYPE == configure.SubscriptionType { //不返回subscriptionType
			which = append(which[:i], which[i+1:]...)
		}
	}
	return which, nil
}

func addHosts(tmpl *v2ray.Template, vms []serverObj.ServerObj) {
	if tmpl.DNS == nil {
		tmpl.DNS = new(coreObj.DNS)
	}
	if tmpl.DNS.Hosts == nil {
		tmpl.DNS.Hosts = make(coreObj.Hosts)
	}
	const concurrency = 5
	var mu sync.Mutex
	var limit = make(chan struct{}, concurrency)
	var wg = sync.WaitGroup{}
	for _, v := range vms {
		if net.ParseIP(v.GetHostname()) == nil {
			wg.Add(1)
			go func(addr string) {
				limit <- struct{}{}
				defer func() {
					wg.Done()
					<-limit
				}()
				ips, err := resolv.LookupHost(addr)
				if err != nil {
					return
				}
				if len(ips) > 0 {
					ips = v2ray.FilterIPs(ips)
					if len(ips) == 0 {
						return
					}
					mu.Lock()
					tmpl.DNS.Hosts[addr] = ips
					mu.Unlock()
				}
			}(v.GetHostname())
		}
	}
	wg.Wait()
}

func TestHttpLatency(which []*configure.Which, timeout time.Duration, maxParallel int, showLog bool, customTestUrl string) ([]*configure.Which, error) {
	var whiches = configure.NewWhiches(which)
	for i := len(which) - 1; i >= 0; i-- {
		if which[i].TYPE == configure.SubscriptionType { //去掉subscriptionType
			which = append(which[:i], which[i+1:]...)
		}
	}
	which = whiches.Get()
	vms := make([]serverObj.ServerObj, len(which))
	//init vmessInfos
	for i := range which {
		which[i].Latency = ""
		sr, err := which[i].LocateServerRaw()
		if err != nil {
			which[i].Latency = err.Error()
			continue
		}
		vms[i] = sr.ServerObj
	}

	//limit the concurrency
	wg := new(sync.WaitGroup)
	cc := make(chan interface{}, maxParallel)
	for i := range which {
		if which[i].Latency != "" {
			if showLog {
				log.Warn("Error[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
			}
			continue
		}
		wg.Add(1)
		go func(i int) {
			cc <- nil
			defer func() { <-cc; wg.Done() }()

			// Broadcast running status
			v2ray.BroadcastLatencyProgress("http", which[i], v2ray.LatencyStatusRunning, "")

			// Create minimal template for this node only
			tmpl := v2ray.NewEmptyTemplate(&configure.Setting{
				RulePortMode:  configure.WhitelistMode,
				TcpFastOpen:   configure.Default,
				MuxOn:         configure.No,
				Transparent:   configure.TransparentClose,
				SpecialMode:   configure.SpecialModeNone,
				AntiPollution: configure.AntipollutionClosed,
			})
			tmpl.SetAPI(nil)

			// Find available port for this test
			l, err := net.Listen("tcp", "0.0.0.0:0")
			if err != nil {
				which[i].Latency = fmt.Sprintf("SYSTEM ERROR: %v", err)
				if showLog {
					log.Warn("Test failed[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
				}
				return
			}
			port := l.Addr().(*net.TCPAddr).Port
			l.Close()

			v2rayInboundPort := strconv.Itoa(port)
			pluginPort := 0
			if vms[i].NeedPluginPort() {
				l2, err := net.Listen("tcp", "127.0.0.1:0")
				if err != nil {
					which[i].Latency = fmt.Sprintf("SYSTEM ERROR: %v", err)
					if showLog {
						log.Warn("Test failed[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
					}
					return
				}
				pluginPort = l2.Addr().(*net.TCPAddr).Port
				l2.Close()
			}

			// Insert the outbound for this node
			err = tmpl.InsertMappingOutbound(vms[i], v2rayInboundPort, false, pluginPort, "socks")
			if err != nil {
				if strings.Contains(err.Error(), "unsupported") {
					which[i].Latency = "UNSUPPORTED PROTOCOL"
				} else {
					which[i].Latency = fmt.Sprintf("CONFIG ERROR: %v", err)
				}
				if showLog {
					log.Warn("Test failed[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
				}
				return
			}

			// Set minimal routing and DNS
			tmpl.Routing.DomainStrategy = "AsIs"
			// Add hosts for this node only
			if net.ParseIP(vms[i].GetHostname()) == nil {
				if tmpl.DNS == nil {
					tmpl.DNS = new(coreObj.DNS)
				}
				if tmpl.DNS.Hosts == nil {
					tmpl.DNS.Hosts = make(coreObj.Hosts)
				}
				ips, err := resolv.LookupHost(vms[i].GetHostname())
				if err == nil && len(ips) > 0 {
					ips = v2ray.FilterIPs(ips)
					if len(ips) > 0 {
						tmpl.DNS.Hosts[vms[i].GetHostname()] = ips
					}
				}
			}
			// Force mark to 0x80 to avoid transparent proxy loopback
			for j := range tmpl.Outbounds {
				if tmpl.Outbounds[j].StreamSettings == nil {
					tmpl.Outbounds[j].StreamSettings = new(coreObj.StreamSettings)
				}
				if tmpl.Outbounds[j].StreamSettings.Sockopt == nil {
					tmpl.Outbounds[j].StreamSettings.Sockopt = new(coreObj.Sockopt)
				}
				mark := 0x80
				tmpl.Outbounds[j].StreamSettings.Sockopt.Mark = &mark
			}

			// Write temporary config
			configPath, cleanupConfig, err := v2ray.WriteTempConfig(tmpl.ToConfigBytes())
			if err != nil {
				which[i].Latency = fmt.Sprintf("SYSTEM ERROR: %v", err)
				if showLog {
					log.Warn("Test failed[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
				}
				return
			}
			defer cleanupConfig()

			// Start plugins if any
			if len(tmpl.Plugins) > 0 {
				if err := tmpl.ServePlugins(); err != nil {
					which[i].Latency = fmt.Sprintf("PLUGIN ERROR: %v", err)
					if showLog {
						log.Warn("Test failed[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
					}
					tmpl.Close()
					return
				}
				defer tmpl.Close()
			}

			// Start isolated process
			ctx, cancel := context.WithTimeout(context.Background(), timeout+30*time.Second)
			defer cancel()

			proc, cancelProc, err := v2ray.StartIsolatedProcess(ctx, configPath)
			if err != nil {
				which[i].Latency = fmt.Sprintf("PROCESS ERROR: %v", err)
				if showLog {
					log.Warn("Test failed[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
				}
				return
			}
			defer cancelProc()

			// Wait for process to be ready with actual port check
			if !waitForPortReady("127.0.0.1", v2rayInboundPort, 5*time.Second) {
				which[i].Latency = "PROCESS NOT READY"
				if showLog {
					log.Warn("Test failed[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
				}
				return
			}

			// Check if process is still running
			if proc == nil {
				which[i].Latency = "PROCESS CRASHED"
				if showLog {
					log.Warn("Test failed[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
				}
				return
			}

			// Perform latency test with retry for NOT STABLE errors
			httpLatencyWithRetry(which[i], v2rayInboundPort, timeout, customTestUrl, 2, showLog)
			if showLog {
				log.Info("Test done[%v]%v: %v", i+1, which[i].Latency, which[i].Link)
			}

			// Broadcast done status with latency
			latency := which[i].Latency
			status := v2ray.LatencyStatusDone
			if latency == "" || !strings.HasSuffix(latency, "ms") {
				status = v2ray.LatencyStatusError
			}
			v2ray.BroadcastLatencyProgress("http", which[i], status, latency)

			// Process will be cleaned up by defer cancelProc()
		}(i)
	}
	wg.Wait()

	if err := configure.NewWhiches(which).SaveLatencies(); err != nil {
		return nil, fmt.Errorf("failed to save the latency test result: %v", err)
	}
	return which, nil
}
func waitForPortReady(host string, portStr string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, portStr), 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func httpLatencyWithRetry(which *configure.Which, port string, timeout time.Duration, customTestUrl string, maxRetries int, showLog bool) {
	for retry := 0; retry <= maxRetries; retry++ {
		httpLatency(which, port, timeout, customTestUrl)
		if which.Latency != "NOT STABLE" {
			// Success or other error, no need to retry
			break
		}
		if retry < maxRetries {
			if showLog {
				log.Info("Retry[%v] for NOT STABLE: %v", retry+1, which.Link)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func httpLatency(which *configure.Which, port string, timeout time.Duration, customTestUrl string) {
	c, err := httpClient.GetHttpClientWithProxy("socks5://127.0.0.1:" + port)
	if err != nil {
		which.Latency = "SYSTEM ERROR"
		return
	}
	defer c.CloseIdleConnections()
	c.Timeout = timeout
	t := time.Now()
	// NOT follow redirects
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	testUrl := HttpTestURL
	if len(customTestUrl) != 0 {
		testUrl = customTestUrl
	}
	req, _ := http.NewRequest("GET", testUrl, nil)
	//req, _ := http.NewRequest("GET", "http://www.gstatic.com/generate_204", nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "curl/7.70.0")
	resp, err := c.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			es := strings.ToLower(err.Error())
			switch {
			case strings.Contains(es, "eof"):
				which.Latency = "NOT STABLE"
			case strings.Contains(es, "does not look like a tls handshake"):
				which.Latency = "INVALID"
			case strings.Contains(es, "timeout"):
				which.Latency = "TIMEOUT"
			default:
				which.Latency = err.Error()
			}
		} else {
			which.Latency = "BAD RESPONSE"
		}
		return
	}
	_ = resp.Body.Close()
	which.Latency = fmt.Sprintf("%.0fms", time.Since(t).Seconds()*1000)
}

func IsSupported(which configure.Which) (bool, error) {
	var (
		tmpl *v2ray.Template
		err  error
	)

	tmpl = v2ray.NewEmptyTemplate(&configure.Setting{
		RulePortMode:  configure.WhitelistMode,
		TcpFastOpen:   configure.Default,
		MuxOn:         configure.No,
		Transparent:   configure.TransparentClose,
		SpecialMode:   configure.SpecialModeNone,
		AntiPollution: configure.AntipollutionClosed,
	})
	tmpl.SetAPI(nil)
	serverRaw, _ := which.LocateServerRaw()
	err = tmpl.InsertMappingOutbound(serverRaw.ServerObj, "0", false, 0, "socks")
	if err != nil {
		if strings.Contains(err.Error(), "unsupported") {
			return false, err
		}
	}
	return true, nil
}
