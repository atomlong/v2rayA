# Project Context

## Purpose
v2rayA is a web-based GUI client for Xray (formerly V2Ray) on Linux. It provides a user-friendly interface to manage proxies, subscriptions, and routing rules, with a strong focus on simplifying transparent proxy configuration (TProxy/Redirect) and rule management. It aims to make advanced proxy features accessible while retaining flexibility.

## Tech Stack
- **Backend**: Go (Golang)
  - Web Framework: Gin, Beego (legacy parts)
  - Database: BoltDB (KV store)
- **Frontend**:
  - `gui/`: Vue.js 2 (Element UI)
  - `ngui/`: Nuxt 3 (Vue 3, Naive UI) - *Next generation UI*
- **Core**: XTLS/Xray-core (migrated from v2fly/v2ray-core)
- **Networking**:
  - iptables / nftables (for transparent proxy)
  - ipset
  - systemd (service management)

## Project Conventions

### Code Style
- **Go**: Standard `gofmt` style.
- **Directories**:
  - `service/`: Backend source code.
  - `gui/` / `ngui/`: Frontend source code.
  - `install/`: Installation scripts and systemd units.
  - `core/`: Core logic wrapping Xray APIs and system network calls.

### Architecture Patterns
- **Process Management**: v2rayA acts as a supervisor, spawning and managing the Xray-core process. It generates configuration files dynamically based on user settings and database state.
- **Separation of Concerns**: Backend handles privileged operations (network config, process management), Frontend handles user interaction.
- **Embedded Assets**: Frontend static files are typically embedded into the backend binary or served from a directory.

### Testing Strategy
- **Unit Tests**: Go standard `testing` package for core logic (e.g., config generation, string parsing).
- **Integration Tests**: Script-based testing (e.g., `test.sh`) to verify binary compilation, startup, and core process stability.

### Git Workflow
- Standard Pull Request workflow.
- Dependencies managed via `go.mod`.

## Domain Context
- **Proxy Protocols**: VMess, VLESS (including REALITY/Vision), Trojan, Shadowsocks, Socks5, HTTP.
- **Routing**: Advanced routing rules using GeoIP/GeoSite dat files.
- **Subscriptions**: parsing and updating node lists from subscription links.
- **Transparent Proxy**: TProxy (UDP/TCP) and Redirect (TCP only) modes.

## Important Constraints
- **Root Privileges**: Requires root access to manipulate network tables (iptables/nftables) and manage system services.
- **Core Compatibility**: Strictly depends on `xray-core` APIs (Protobuf).
- **Linux-focused**: Primarily designed for Linux environments (systemd, iptables).

## External Dependencies
- **Xray-core**: The underlying proxy engine.
- **Geo Data**: `geoip.dat` and `geosite.dat` (downloaded from upstream sources like Loyalsoldier).