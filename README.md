# WardenGrid

![License](https://img.shields.io/badge/License-PolyForm--Noncommercial-blue)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)
![Language](https://img.shields.io/badge/Language-Go-00ADD8)
![Status](https://img.shields.io/badge/Status-Early%20Development-orange)
![CI](https://github.com/DeveloperBatuhanALGUL/WardenGrid/actions/workflows/ci.yml/badge.svg)

**Cross-platform ICS/SCADA security monitoring and anomaly detection system.**

WardenGrid is an open-source monitoring tool designed to detect anomalous behavior in Industrial Control System (ICS) and SCADA network traffic, with a primary focus on critical infrastructure protection, including power grid environments.

## Status

Early development. Core protocol parsing, detection, capture, and reporting layers are functional. Web dashboard and additional protocol support are in progress.

## Rationale

Critical infrastructure such as power generation and distribution systems increasingly relies on legacy industrial protocols, including Modbus and DNP3, that were not designed with modern cybersecurity threats in mind. WardenGrid provides a lightweight, auditable, cross-platform tool to monitor this traffic and flag suspicious activity, without dependency on proprietary or closed-source software.

## Design Principles

| Principle | Description |
| :--- | :--- |
| Cross-Platform | Native builds for Windows, macOS, and Linux from a single codebase. |
| Defensive Only | Detection and monitoring only. No offensive capability of any kind. |
| Protocol-Aware | Parses industrial protocols. Modbus/TCP and DNP3 are implemented. |
| Layered Architecture | Capture, protocol parsing, detection, and reporting are independently testable modules. |
| No Vendor Lock-In | Built entirely on open standards with no proprietary dependencies. |

## Architecture

Full system design is documented in [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Repository Structure

```text
WardenGrid/
├── cmd/wardengrid/             Entry point and CLI
├── internal/capture/           Live packet capture layer (pcap-based)
├── internal/protocol/modbus/   Modbus/TCP parser
├── internal/protocol/dnp3/     DNP3 parser
├── internal/detector/          Rule-based and statistical anomaly detection
├── internal/simulator/         Synthetic ICS traffic generator for testing
├── internal/report/            Alert reporting in text and JSON formats
├── web/dashboard/               Monitoring dashboard (planned)
├── docs/                        Architecture and design documentation
└── .github/workflows/           Continuous integration pipeline
```

## Installation

### Prerequisites

WardenGrid requires Go 1.26 or later. Live packet capture requires libpcap (Linux and macOS) or Npcap (Windows, not yet implemented).

**Linux:**

```bash
sudo apt-get install libpcap-dev
```

**macOS:**

```bash
brew install libpcap pkg-config
```

If libpcap is not found at build time on macOS, set the following environment variables before building:

```bash
export PKG_CONFIG_PATH="/opt/homebrew/opt/libpcap/lib/pkgconfig:$PKG_CONFIG_PATH"
export LDFLAGS="-L/opt/homebrew/opt/libpcap/lib"
export CPPFLAGS="-I/opt/homebrew/opt/libpcap/include"
```

**Windows:**

Live capture is not yet implemented on Windows. The CLI, protocol parsers, detector, and simulator all build and run normally.

### Build

```bash
git clone https://github.com/DeveloperBatuhanALGUL/WardenGrid.git
cd WardenGrid
go build -o bin/wardengrid ./cmd/wardengrid
```

### Run

WardenGrid currently runs against simulated traffic. Available scenarios:

```bash
./bin/wardengrid -scenario normal -count 5
./bin/wardengrid -scenario unknown-function -count 5
./bin/wardengrid -scenario protected-write -count 4
./bin/wardengrid -scenario write-frequency -count 10
```

Use the `-json` flag to emit alerts as JSON lines instead of text.

### Test

```bash
go test ./... -v
```

## License

Licensed under the [PolyForm Noncommercial License 1.0.0](LICENSE). Free for personal, academic, research, government, and nonprofit use. Commercial use is not permitted under this license.

## Author

Batuhan ALGUL
[GitHub](https://github.com/DeveloperBatuhanALGUL)
