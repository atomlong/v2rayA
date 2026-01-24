# v2rayA 当前活动上下文

## 最近完成的工作

### Vless 用户 ID 双重 URL 编码问题修复 (2026-01-24)

**问题描述：**
vless 链接中用户 ID 包含双重 URL 编码字符 `%2540`（`%40` 的编码），导致 xray 崩溃，错误：`encoding/hex: invalid byte: U+0025 '%'`

**问题链接示例：**
`vless://%2540X_Her0%2540X_Her0%2540X_Her0%2540X_Her0@star0.kharabetam.de:2053?...`

**根本原因：**
- `%2540` 是 `%40` 的 URL 编码，`%40` = `@`
- Go 的 `url.Parse()` 只做一次解码：`%2540` → `%40`
- 需要再解码一次：`%40` → `@`

**解决方案：**
在 `ParseVlessURL` 函数中对用户 ID 进行额外的 URL 解码处理：

**修改的文件：**
- `service/core/serverObj/v2ray.go` - `ParseVlessURL` 函数

**修改内容：**
```go
// 处理用户 ID，可能存在双重 URL 编码的情况
userID := u.User.String()
if strings.Contains(userID, "%") {
    if decoded, err := url.PathUnescape(userID); err == nil {
        userID = decoded
    }
}
```

---

### Vmess 协议 `raw` 传输类型支持修复 (2026-01-24)

**问题描述：**
测试 HTTP 延时时遇到错误 `unexpected transport type: raw`。vmess 链接中包含 `"net":"raw"` 的传输类型未被正确处理。

**根本原因：**
- `ParseVlessURL` 函数已经有 `raw`→`tcp` 的转换
- `ParseVmessURL` 函数只处理了 `none`→`tcp`，缺少对 `raw` 的处理

**解决方案：**
在 `ParseVmessURL` 函数中添加 `raw`→`tcp` 的别名转换：

**修改的文件：**
- `service/core/serverObj/v2ray.go` - `ParseVmessURL` 函数

**修改内容：**
```go
// 修改前
if info.Net == "" || info.Net == "none" {
    info.Net = "tcp"
}

// 修改后
if info.Net == "" || info.Net == "none" || info.Net == "raw" {
    info.Net = "tcp"
}
```

---

### URL 查询参数空格问题全面修复 (2026-01-22)

**问题描述：**
订阅链接中 `+` 在 `application/x-www-form-urlencoded` 编码中被解码为空格，导致多个参数解析失败，使 xray 核心崩溃。

**遇到的具体错误：**
1. `type=ws+` → `"unexpected transport type: ws "`
2. `type=tcp+` → `"unknown transport protocol: tcp "`
3. `fp=chrome+` → `"unknown \"fingerprint\": chrome "`
4. `net=none` → `"unexpected transport type: none"`
5. `type=raw` → `"unexpected transport type: raw"`
6. `type=splithttp` → `"unexpected transport type: splithttp"`
7. `sid=2404+` → `"invalid \"shortId\": 2404"` (REALITY 配置)

**解决方案：**
对所有从 `url.Query().Get()` 获取的参数统一应用 `strings.TrimSpace()` 清理首尾空格，同时添加传输协议别名转换。

**修改的文件：**

1. **service/core/serverObj/v2ray.go - ParseVlessURL 函数**
   - 所有字段应用 TrimSpace：`aid`, `type`, `headerType`, `host`, `sni`, `path`, `security`, `fp`, `pbk`, `sid`, `spx`, `flow`, `alpn`, `allowInsecure`, `key`
   - 后续赋值也应用 TrimSpace：`serviceName`, `host`, `seed`, `quicSecurity`
   - 添加传输协议别名转换：`raw`→`tcp`, `splithttp`→`xhttp`

2. **service/core/serverObj/v2ray.go - ParseVmessURL 函数**
   - 添加 `none`→`tcp` 转换

3. **service/core/serverObj/trojan.go - ParseTrojanURL 函数**
   - 所有字段应用 TrimSpace：`allowInsecure`, `peer`, `sni`, `alpn`, `type`, `path`, `serviceName`, `encryption`, `host`

**技术要点：**
- URL 中 `+` 在 `application/x-www-form-urlencoded` 编码中被解码为空格是标准行为
- `strings.TrimSpace()` 只去除首尾空白字符，不影响有效内容
- 这些参数值本身不应包含首尾空格，因此修复是安全的

## 当前状态

- 无活跃开发任务
- 代码编译验证通过

## 技术笔记

### URL 参数解析最佳实践
对于所有从 `url.Query().Get()` 获取的参数，应统一应用 `strings.TrimSpace()` 进行防御性处理，避免因订阅源编码问题导致的解析失败。

### 传输协议别名映射
- `raw` → `tcp`
- `none` → `tcp`
- `splithttp` → `xhttp`
- `websocket` → `ws`

## 上下文刷新时间
2026-01-24 12:17 CST
