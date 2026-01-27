# Active Context

## 当前工作焦点

### 编译修复与依赖升级（2026-01-27）

**问题背景**：
- CI 构建失败：`sing-tun v0.1.20` 与 `sing v0.5.1` 不兼容（`Addr` 字段未导出）
- 存在版本碎片化：`sing-tun` 过旧，`sing` 较新，`gvisor` 版本不匹配
- 存在代码质量问题：`tools_test.go` 循环依赖，`dns.go` 不可达代码

**已完成修复**：
1. **依赖升级**：
   - `sing-tun` -> `v0.7.5` (stable)
   - `sing` -> `v0.7.14` (stable)
   - `gvisor` -> `v0.0.0-20241123...` (auto upgraded)
2. **代码适配**：
   - `service/core/tun/singTun.go`：适配 `sing-tun` v0.7.x API
     - `Options`：`TableIndex` -> `IPRoute2TableIndex`
     - `StackOptions`：结构体字段变更，使用 `TunOptions`
     - `Handler`：实现 `NewConnectionEx`、`NewPacketConnectionEx`、`PrepareConnection`
     - `InterfaceFinder`：使用 `control.NewDefaultInterfaceFinder()`
3. **代码清理**：
   - 修复 `service/common/tools_test.go` 循环依赖（改为 `package common_test`）
   - 修复 `service/core/tun/dns.go` 不可达代码警告

**预期效果**：
- 本地编译通过
- CI 构建恢复正常
- 依赖库处于较新的稳定状态

## 最近变更

### 2026-01-27
- 升级 `sing-tun` 和 `sing` 到最新稳定版
- 重构 `singTun.go` 以适配新版 API
- 修复 `common` 包测试循环依赖
- 修复 `dns` 包代码警告

### 2026-01-26
- 重构 HTTP 延迟测试为隔离进程模式
- 添加 `service/core/v2ray/isolated.go` 实现独立进程管理
- 修复透明代理 mark 标记缺失问题
- 添加插件支持和端口就绪检测

## 重要决策

### 依赖版本管理
- **决策**：升级到 `sing-tun v0.7.5` 和 `sing v0.7.14`
- **理由**：
  - 解决旧版本间的 API 兼容性问题
  - 避免停留在不受支持的旧版本
  - `sing-tun v0.2.0` API 变更较大，不如直接升级到稳定版
- **权衡**：
  - 需要修改 `singTun.go` 适配新 API
  - 可能引入新的行为变更（需测试验证）

### 延迟测试架构
- **决策**：每个节点使用独立 Xray 进程测试
- **理由**：
  - 隔离失败：一个节点崩溃不影响其他节点
  - 不干扰主服务：测试不影响正常代理流量
- **权衡**：
  - 启动开销增加，但测试并发度可控
  - 内存占用增加，但每个进程生命周期短

## 模式和偏好

### 代码组织
- 进程管理放在 `service/core/v2ray/` 目录
- 服务层逻辑在 `service/server/service/`
- 配置模板在 `v2ray.Template` 结构

### 错误处理
- 网络错误分类明确（NOT STABLE、TIMEOUT、INVALID）
- 支持重试机制提高成功率

## 下一步

- [ ] 提交代码并关注 CI 状态
- [ ] 验证透明代理功能（因 API 变更）
- [ ] 验证 DNS 功能
- [ ] 归档 change proposal

## 学习和项目洞察

### sing-box 生态
- `sing` 和 `sing-tun` 版本紧密耦合，需保持同步升级
- `StackOptions` 和 `Options` 在 v0.2.0+ 有重大变更
- `Handler` 接口在 v0.5.0+ 引入了 `Ex` 后缀方法以支持源/目的地址

### 插件系统
- `tmpl.Plugins` 存储插件实例
- `tmpl.ServePlugins()` 启动插件进程
- `tmpl.Close()` 清理插件资源