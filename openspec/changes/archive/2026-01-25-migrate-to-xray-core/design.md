## Context
v2rayA 作为一个 v2ray/xray 的管理前端，其后端核心逻辑深度依赖于 `v2ray-core` 的 Go 库。随着 `xray-core` 的发展，它在协议支持和活跃度上已超越 `v2ray-core`。此外，尝试引入 `xray-core` 特定功能（如 REALITY）时遇到了难以解决的依赖和命名空间冲突。

## Goals / Non-Goals
**Goals**:
- 将所有 Go 依赖从 `v2fly/v2ray-core` 迁移到 `xtls/xray-core`。
- 解决因依赖升级带来的编译错误（主要是 `sagernet/sing` 库升级）。
- 确保迁移后 v2rayA 能够正常编译并运行。
- 保持对现有配置文件的兼容性（在 Xray-core 兼容范围内）。

**Non-Goals**:
- 不在此次迁移中重构非核心业务逻辑。
- 不在此次迁移中引入新的 UI 功能。

## Decisions

### Decision 1: 全局替换策略
**选择**：使用脚本或批量替换工具，将所有 `github.com/v2fly/v2ray-core/v5` 替换为 `github.com/xtls/xray-core`。

**理由**：
- 两个核心的包结构高度相似，大部分代码可以平滑迁移。
- 手动修改容易出错且效率低下。

### Decision 2: 解决 sing 库冲突
**选择**：升级 `v2rayA` 代码以适配 `xray-core` 所需的 `sagernet/sing v0.4.1+`。

**理由**：
- `xray-core` 依赖新版 `sing`，降级 `xray-core` 不可行（会丢失新功能）。
- `v2rayA` 的 `core/tun` 模块使用了旧版 `sing` API，必须升级适配。

### Decision 3: 处理 API 差异
**选择**：对于不兼容的 API（如 proto 定义差异、接口变更），进行针对性修改。

**理由**：虽然大部分兼容，但仍存在少量差异，需要人工介入修复。

## Risks / Trade-offs
| 风险 | 缓解措施 |
|-----|---------|
| 兼容性破坏 | 重点测试核心功能（连接、路由、DNS）；依靠 Xray-core 的兼容性承诺。 |
| 工作量大 | 分模块进行，先确保 `core` 模块编译通过，再处理 `service` 模块。 |
| 引入新 Bug | 尽量保持原有逻辑结构不变，只做必要的适配。 |

## Migration Plan
1.  **清理环境**：移除旧的 `go.mod` 依赖。
2.  **批量替换**：执行全局字符串替换。
3.  **依赖更新**：运行 `go mod tidy` 拉取 `xray-core`。
4.  **修复编译错误**：
    -   修复 `core/tun` 模块（sing 库适配）。
    -   修复 `core/serverObj` 和 `core/v2ray` 中因类型定义变化导致的错误。
    -   修复 `service` 模块中的调用错误。
5.  **验证**：编译并运行，进行基本功能测试。