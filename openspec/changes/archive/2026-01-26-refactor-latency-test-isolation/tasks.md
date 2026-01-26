## 1. 准备工作
- [x] 1.1 确认当前 `TestHttpLatency` 行为及调用链
- [x] 1.2 确认 `InsertMappingOutbound` 返回的配置最小结构

## 2. 新增辅助函数
- [x] 2.1 在 `service/core/v2ray/` 新增 `isolated.go`（或在 `process.go` 追加）
  - `WriteTempConfig(content []byte) (path string, cleanup func(), err error)`
  - `StartIsolatedProcess(ctx context.Context, configPath string) (*os.Process, error)`
- [x] 2.2 单元测试验证辅助函数正常工作

## 3. 重构 `TestHttpLatency`
- [x] 3.1 移除对 `ProcessManager.Start` / `Stop` / `UpdateV2RayConfig` 的调用
- [x] 3.2 循环每节点：
  - 生成最小模板 (`NewEmptyTemplate` + `InsertMappingOutbound` + `SetOutboundSockopt` + 路由/hosts)
  - 调用 `WriteTempConfig` / `StartIsolatedProcess`
  - 执行 `httpLatency`
  - 关闭进程、清理临时文件
- [x] 3.3 并发控制复用现有 `maxParallel` / channel 机制
- [x] 3.4 错误处理：进程启动失败或崩溃仅标记节点 `Latency` 错误，不中断

## 4. 修复透明代理和插件问题
- [x] 4.1 强制设置 mark=0x80 到所有 outbound（避免透明代理回环）
- [x] 4.2 调用 `tmpl.ServePlugins()` 启动插件链
- [x] 4.3 添加 `waitForPortReady()` 端口就绪检测
- [x] 4.4 添加 `httpLatencyWithRetry()` 对 NOT STABLE 错误重试

## 5. 回归与验证
- [x] 5.1 编译通过 (`go build`)
- [x] 5.2 运行 `test.sh` / 手动测试订阅，确认：
  - 主服务不中断
  - 单节点崩溃不影响其他节点
  - 延迟结果正确保存
  - NOT STABLE 错误减少
  - 插件节点测试待后续验证（已记录到 memory-bank）
- [x] 5.3 更新 memory-bank `progress.md` 与 `activeContext.md`
