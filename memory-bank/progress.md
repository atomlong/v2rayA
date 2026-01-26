# Progress

## What Works

### 核心功能
- ✅ 节点订阅与导入：支持多种协议（VMess、VLESS、Shadowsocks、Trojan、V2Ray/Xray）
- ✅ 多服务器管理：同时管理多个服务器和订阅
- ✅ 自动选择：从订阅自动选择延迟最低的服务器
- ✅ 规则模式：支持 GFWList、自定义规则、白名单模式
- ✅ 透明代理：支持透明代理（TPROXY）、TUN 模式
- ✅ 插件支持：支持 SSR、obfs、v2ray-plugin 等插件

### 延迟测试（已重构）
- ✅ HTTP 延迟测试：测试节点连通性
- ✅ 隔离进程模式：每个节点使用独立 Xray 进程测试
- ✅ 透明代理兼容：强制 mark=0x80 避免回环
- ✅ 插件节点支持：启动插件链进行测试
- ✅ 端口就绪检测：等待端口真正监听后再测试
- ✅ 错误重试机制：对 NOT STABLE 错误自动重试
- ✅ 单节点失败隔离：一个节点崩溃不影响其他节点

### Web 界面
- ✅ Vue.js 2.x 前端（gui/）和 Nuxt 3 前端（ngui/）
- ✅ 节点管理：添加、编辑、删除、导入导出
- ✅ 订阅管理：更新、查看、自动选择
- ✅ 实时日志：查看 v2ray-core 日志
- ✅ 设置管理：端口、规则、透明代理等配置

## What's Left to Build

### 测试与验证
- [x] 编译验证通过
- [x] 单元测试覆盖
  - 测试文件位置：`service/test/isolated_test.go`
  - 测试结果：7 个测试，5 个通过，2 个跳过
  - 修复：`service/conf/environmentConfig.go` 忽略测试 flag 冲突
- [x] 基础功能验证：主服务不中断、单节点崩溃不影响其他节点
- [ ] 插件节点测试：需要环境中有 v2ray-plugin、obfs-local 等插件可执行文件
  - 代码已支持：调用 `tmpl.ServePlugins()` 启动插件链
  - 待验证：实际运行时测试插件节点
- [ ] 性能测试：大量节点（100+）的并发测试

### 文档
- [ ] 更新 README 说明延迟测试新架构
- [ ] API 文档完善

## Current Status

### 2026-01-26
**延迟测试隔离重构完成**
- 问题：旧版本使用主 Xray 服务进程测试，一个节点崩溃影响所有节点
- 解决：每个节点使用独立进程测试，强制 mark=0x80，支持插件
- 结果：从 3 个节点恢复到 16 个节点检测

**修改文件**：
- `service/core/v2ray/isolated.go`：独立进程管理
- `service/test/isolated_test.go`：单元测试（从 `service/core/v2ray/` 移动）
- `service/server/service/latency.go`：核心测试逻辑重构

**已知问题**：
- 无

## Known Issues

### 测试相关
1. ~~**单元测试运行失败**：`pre.go` 初始化代码解析测试 flag 时出错~~
   - ~~影响：无法直接运行 `go test ./core/v2ray`~~
   - ~~临时方案：通过编译验证代码正确性~~
   - ~~长期方案：拆分测试初始化或使用 build tag~~
   - **已解决**：修复 `service/conf/environmentConfig.go`，使用前缀匹配忽略所有测试 flag

### 延迟测试
1. **需要 v2ray/Xray 二进制**：测试依赖系统安装的 v2ray-core
   - 影响：没有 v2ray 时测试会 skip
   - 解决方案：预期行为，测试代码已处理

## Evolution of Project Decisions

### 延迟测试架构（2026-01-26）
**决策**：从主进程测试改为隔离进程测试
**理由**：
- 隔离失败：一个节点崩溃不影响其他节点
- 不干扰主服务：测试不影响正常代理流量
- 参考实现：v2rayN/v2rayNG 采用类似方案
**权衡**：
- 启动开销增加，但测试并发度可控
- 内存占用增加，但每个进程生命周期短（约 10-30 秒）

### 透明代理处理（2026-01-26）
**决策**：测试时强制设置 mark=0x80
**理由**：
- 避免透明代理规则回环重定向
- 保持 iptables 规则不变
- 与旧版本 `config_old.json` 行为一致