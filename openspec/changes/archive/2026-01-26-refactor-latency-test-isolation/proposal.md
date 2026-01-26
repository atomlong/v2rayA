# Change: Refactor Latency Testing to Use Isolated Processes

## Why
当前 `TestHttpLatency` 函数会调用 `ProcessManager.Start(tmpl)` 启动单一的 Xray 实例来测试所有节点。这导致两个严重问题：
1. **主服务中断**: 测试期间会停止用户正在使用的主代理服务。
2. **单节点崩溃影响全部**: 任一节点配置无效导致 Xray 崩溃，所有节点测试都会失败。

目标：借鉴 v2rayN/v2rayNG 的做法，为每个节点（或小批量）启动独立进程进行延迟测试，实现故障隔离。

## What Changes
- 新增 **临时配置文件写入** 与 **独立进程启动** 辅助函数。
- 重构 `TestHttpLatency`，不再调用 `ProcessManager`，改用独立进程：
  - 每个节点生成最小模板。
  - 启动临时 Xray 进程。
  - 测速完成后关闭进程、清理临时文件。
- 单节点配置无效 / 进程崩溃仅标记该节点 `Latency` 错误，不中断其他节点测试。
- 主服务在测试期间**保持运行**（不调用 `ProcessManager.Stop`）。

## Impact
- Affected specs: `core-backend`
- Affected code:
  - `service/server/service/latency.go` (主改造)
  - `service/core/v2ray/process.go` 或新文件（辅助函数）
- 用户体验提升：测速不再中断代理；失败节点不拖累成功节点。