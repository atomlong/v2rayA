## ADDED Requirements
### Requirement: Latency Testing Isolation
The system SHALL perform HTTP latency testing using isolated Xray-core processes for each server node to prevent failure propagation and avoid interrupting the main proxy service.

#### Scenario: Node-specific isolated process
- **WHEN** testing HTTP latency for a server node
- **THEN** the system SHALL create a minimal Xray configuration for that node only
- **AND** the system SHALL start a separate Xray-core process with a temporary configuration file
- **AND** the system SHALL perform the latency test through that isolated process
- **AND** the system SHALL stop the isolated process and clean up the temporary configuration file after testing

#### Scenario: Main service continuity
- **WHEN** latency testing is in progress
- **THEN** the main proxy service SHALL continue running without interruption
- **AND** the system SHALL NOT call `ProcessManager.Stop()` during latency testing

#### Scenario: Single node failure isolation
- **WHEN** a node's configuration is invalid or causes Xray-core to crash
- **THEN** the system SHALL mark only that node's latency as an error
- **AND** the system SHALL continue testing remaining nodes
- **AND** the crash SHALL NOT affect latency testing of other nodes

#### Scenario: Concurrent testing with parallel limit
- **WHEN** testing multiple nodes with `maxParallel` > 1
- **THEN** the system SHALL limit concurrent isolated processes to `maxParallel`
- **AND** each concurrent test SHALL use independent temporary configuration files
- **AND** process shutdown and cleanup SHALL be performed per node

## ADDED Requirements
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