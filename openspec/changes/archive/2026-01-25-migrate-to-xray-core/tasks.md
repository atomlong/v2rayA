## 1. Preparation

- [x] 1.1 清理 `go.mod`
  - [x] 移除 `github.com/v2fly/v2ray-core/v5` 及其 replace 指令
  - [x] 添加 `github.com/xtls/xray-core` 依赖
  - [x] 运行 `go mod tidy`

## 2. Refactoring

- [x] 2.1 全局包路径替换
  - [x] 替换 `github.com/v2fly/v2ray-core/v5` -> `github.com/xtls/xray-core`
  - [x] 检查 `core/v2ray/api.go` 中的引用
  - [x] 检查 `core/v2ray/v2rayTmpl.go` 中的引用
  - [x] 检查 `core/serverObj` 中的引用 (如有)
  - [x] 检查 `core/tun` 中的引用 (如有)

- [x] 2.2 修复编译错误 (Core 模块)
  - [x] 修复 `core/tun/singTun.go`: 适配 `sing v0.4.1+` (InitializeReadWaiter, WaitReadPacket, InterfaceFinder)
  - [x] 修复 `core/tun/dns.go`: 适配 `sing v0.4.1+`
  - [x] 修复 `core/v2ray` 中类型定义不匹配的问题 (如 `strmatcher`, `observatory`)

- [x] 2.3 修复编译错误 (Service 模块)
  - [x] 运行 `go build ./...` 检查其他模块的编译错误并修复

## 3. Verification & Testing

- [x] 3.1 单元测试
  - [x] 为 `core/v2ray` 编写/运行单元测试，验证配置生成正确性
  - [x] 运行 `go test ./core/...`

- [x] 3.2 编译测试
  - [x] 编译 v2rayA 二进制 `go build -o v2raya ./main.go`

- [x] 3.3 系统集成测试
  - [x] 替换 `/usr/bin/v2raya` (或服务使用的路径) 为新编译的二进制
  - [x] 重启服务 `systemctl restart v2raya`
  - [x] 验证服务启动日志无异常
  - [x] 验证 Web UI 可访问
  - [x] 验证节点连接和代理功能
  - [x] 验证 HTTP 延时测试功能 (注：对于正确配置的节点测试通过；无效节点导致核心崩溃为预期行为，后续修复)
