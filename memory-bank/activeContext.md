# Active Context

## 当前工作焦点

### 延迟测试隔离重构（2026-01-26）

**问题背景**：
- v2rayA 延迟测试使用主 Xray 服务进程
- 一个节点配置错误会导致 Xray 崩溃，影响所有节点测试
- 透明代理开启时测试会干扰正常代理服务

**已完成修复**：
1. **强制 mark 标记**：所有 outbound 设置 mark=0x80，避免透明代理回环
2. **启动插件链**：调用 `tmpl.ServePlugins()` 支持含插件的节点（SSR、obfs 等）
3. **端口就绪检测**：`waitForPortReady()` 确保端口真正监听后再测试
4. **错误重试机制**：`httpLatencyWithRetry()` 对 NOT STABLE 错误自动重试 2 次

**修改文件**：
- `service/server/service/latency.go`：核心测试逻辑重构

**预期效果**：
- 单节点测试失败不影响其他节点
- 支持插件节点测试
- 减少因进程未就绪导致的误报
- 透明代理环境测试更稳定

## 最近变更

### 2026-01-26
- 重构 HTTP 延迟测试为隔离进程模式
- 添加 `service/core/v2ray/isolated.go` 实现独立进程管理
- 修复透明代理 mark 标记缺失问题
- 添加插件支持和端口就绪检测

## 重要决策

### 延迟测试架构
- **决策**：每个节点使用独立 Xray 进程测试
- **理由**：
  - 隔离失败：一个节点崩溃不影响其他节点
  - 不干扰主服务：测试不影响正常代理流量
  - 参考实现：v2rayN/v2rayNG 采用类似方案
- **权衡**：
  - 启动开销增加，但测试并发度可控（maxParallel 参数）
  - 内存占用增加，但每个进程生命周期短（约 10-30 秒）

### 透明代理处理
- **决策**：测试时强制设置 mark=0x80
- **理由**：
  - 避免透明代理规则回环重定向
  - 保持 iptables 规则不变（不写临时规则）
  - 与旧版本 `config_old.json` 行为一致

## 模式和偏好

### 代码组织
- 进程管理放在 `service/core/v2ray/` 目录
- 服务层逻辑在 `service/server/service/`
- 配置模板在 `v2ray.Template` 结构

### 错误处理
- 网络错误分类明确（NOT STABLE、TIMEOUT、INVALID）
- 支持重试机制提高成功率
- 错误信息国际化友好

## 下一步

- [ ] 验证修复效果（用户测试）
- [ ] 更新 OpenSpec proposal 状态
- [ ] 归档 change proposal

## 学习和项目洞察

### Xray 进程管理
- 每个进程需要独立的临时配置文件
- 使用 `cmd.Start()` 而非 `cmd.Run()` 实现异步管理
- context.Context 用于超时控制

### Go 并发模式
- 使用 sync.WaitGroup 控制并发测试
- channel 用于限制并发数（maxParallel）
- defer 确保资源清理（进程、配置文件、插件）

### 插件系统
- `tmpl.Plugins` 存储插件实例
- `tmpl.ServePlugins()` 启动插件进程
- `tmpl.Close()` 清理插件资源