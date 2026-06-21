# WardenGrid

![License](https://img.shields.io/badge/License-PolyForm--Noncommercial-blue)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)
![Language](https://img.shields.io/badge/Language-Go-00ADD8)
![Status](https://img.shields.io/badge/Status-Early%20Development-orange)
![CI](https://github.com/DeveloperBatuhanALGUL/WardenGrid/actions/workflows/ci.yml/badge.svg)

**Cross-platform ICS/SCADA security monitoring and anomaly detection system.**

WardenGrid is an open-source monitoring tool designed to detect anomalous behavior in Industrial Control System (ICS) and SCADA network traffic, with a primary focus on critical infrastructure protection, including power grid environments.

## Status

Early development. Core architecture and detection modules are in progress.

## Rationale

Critical infrastructure such as power generation and distribution systems increasingly relies on legacy industrial protocols, including Modbus and DNP3, that were not designed with modern cybersecurity threats in mind. WardenGrid provides a lightweight, auditable, cross-platform tool to monitor this traffic and flag suspicious activity, without dependency on proprietary or closed-source software.

## Design Principles

| Principle | Description |
| :--- | :--- |
| Cross-Platform | Native builds for Windows, macOS, and Linux from a single codebase. |
| Defensive Only | Detection and monitoring only. No offensive capability of any kind. |
| Protocol-Aware | Parses industrial protocols, starting with Modbus/TCP, with DNP3 planned. |
| Layered Architecture | Capture, protocol parsing, detection, and reporting are independently testable modules. |
| No Vendor Lock-In | Built entirely on open standards with no proprietary dependencies. |

## Architecture

Full system design is documented in [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Repository Structure

```text
WardenGrid/
├── cmd/wardengrid/        Entry point and CLI
├── internal/capture/      Network capture layer
├── internal/protocol/     Industrial protocol parsers
├── internal/detector/     Anomaly detection engine
├── internal/platform/     OS-specific abstractions
├── internal/simulator/    Synthetic ICS traffic generator for testing
├── internal/report/       Logging and alert reporting
├── web/dashboard/         Monitoring dashboard
├── docs/                  Architecture and design documentation
└── test/fixtures/         Test data and fixtures
```

## Installation

Build instructions will be published once the first functional module is complete.

## License

Licensed under the [PolyForm Noncommercial License 1.0.0](LICENSE). Free for personal, academic, research, government, and nonprofit use. Commercial use is not permitted under this license.

## Author

Batuhan ALGUL
[GitHub](https://github.com/DeveloperBatuhanALGUL)
