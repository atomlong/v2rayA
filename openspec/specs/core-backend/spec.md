# core-backend Specification

## Purpose
TBD - created by archiving change migrate-to-xray-core. Update Purpose after archive.
## Requirements
### Requirement: Xray-core Backend
The system SHALL use `xtls/xray-core` as the primary core library for V2Ray/Xray protocol support.

#### Scenario: Core library dependency
- **WHEN** building the application
- **THEN** the `go.mod` file SHALL depend on `github.com/xtls/xray-core`
- **AND** the `go.mod` file SHALL NOT depend on `github.com/v2fly/v2ray-core`

#### Scenario: Configuration generation
- **WHEN** generating configuration files for the core process
- **THEN** the configuration format SHALL be compatible with Xray-core
- **AND** the configuration SHALL support Xray-specific features (e.g., REALITY) where applicable

#### Scenario: Process management
- **WHEN** managing the core process
- **THEN** the system SHALL correctly interact with the Xray-core binary or library interface

### Requirement: Sing-box Library Compatibility
The system SHALL be compatible with the version of `sagernet/sing` library required by `xtls/xray-core`.

#### Scenario: Tun module compatibility
- **WHEN** compiling the `core/tun` module
- **THEN** the code SHALL compile successfully with the `sagernet/sing` version pulled by `xray-core`

### Requirement: Temporary Configuration Management
The system SHALL provide utilities to create and manage temporary Xray-core configuration files for isolated process scenarios.

#### Scenario: Create temporary configuration
- **WHEN** a component needs a temporary Xray configuration file
- **THEN** the system SHALL write the configuration to a unique temporary file path
- **AND** the system SHALL provide a cleanup function to delete the temporary file

#### Scenario: Isolated process startup
- **WHEN** starting an isolated Xray-core process
- **THEN** the system SHALL specify a custom configuration file path using the `--config` argument
- **AND** the system SHALL return a process handle with cancelation context

### Requirement: Minimum Template Generation
The system SHALL generate minimal Xray templates suitable for latency testing without unnecessary overhead.

#### Scenario: Minimal template structure
- **WHEN** generating a template for latency testing
- **THEN** the template SHALL include only: the test node outbound, a test inbound (SOCKS5), basic DNS settings, and routing rules
- **AND** the template SHALL NOT include API endpoints, observatory, or other production features

