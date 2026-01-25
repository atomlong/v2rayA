# Change: 迁移核心依赖至 Xray-core

## Why
v2rayA 目前依赖 `v2fly/v2ray-core` 作为底层核心库。然而：
1.  **活跃度**：`v2fly/v2ray-core` 更新频率较低，而 `xtls/xray-core` 社区非常活跃。
2.  **协议支持**：Xray-core 原生支持 REALITY、XTLS Vision 等新一代协议，而 v2ray-core 不支持或支持有限。
3.  **长期维护**：Xray-core 被视为 v2ray-core 的超集和未来方向。
4.  **依赖冲突**：尝试在项目中同时引入两个核心会导致严重的依赖冲突（如 `sing` 库版本不兼容）和运行时错误（Proto 命名空间冲突）。

为了支持更先进的协议并解决长期维护问题，我们需要将项目的核心依赖完全迁移到 `xtls/xray-core`。

## What Changes
-   **依赖替换**：在 `go.mod` 中移除 `github.com/v2fly/v2ray-core/v5`，替换为 `github.com/xtls/xray-core`。
-   **代码重构**：将所有代码中引用 `v2ray-core` 的包路径修改为 `xray-core` 对应的路径。
-   **适配升级**：升级 `v2rayA` 的相关模块（如 `core/tun`）以适配 `xray-core` 所依赖的新版第三方库（如 `sagernet/sing`）。

## Impact
-   **Affected specs**: `core-backend` (New)
-   **Affected code**: 全局影响。涉及 `core/`, `service/`, `db/` 等多个模块，凡是引用了 v2ray-core 类型定义的文件都需要修改。
-   **Breaking Changes**: 对内部 API 是破坏性变更。对用户而言，只要用户安装了兼容的 xray 二进制文件，功能应保持一致或增强。