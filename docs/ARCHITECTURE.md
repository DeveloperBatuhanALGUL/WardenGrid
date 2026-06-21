# Architecture

## Overview

WardenGrid follows a layered architecture. Each layer has a single responsibility and communicates with adjacent layers through defined interfaces. This allows every layer to be tested independently and extended without affecting the others.

## Layer Diagram

```text
+----------------------------------------------------------+
|  PRESENTATION LAYER                                       |
|  CLI (cmd/wardengrid) and Web Dashboard (web/dashboard)    |
+----------------------------------------------------------+
                        |
+----------------------------------------------------------+
|  DETECTION ENGINE (internal/detector)                      |
|  Rule-based matching and statistical anomaly scoring        |
+----------------------------------------------------------+
                        |
+----------------------------------------------------------+
|  PROTOCOL PARSER LAYER (internal/protocol)                  |
|  Modbus/TCP parser, DNP3 parser planned                     |
+----------------------------------------------------------+
                        |
+----------------------------------------------------------+
|  CAPTURE AND SIMULATION LAYER                                |
|  internal/capture for live traffic                           |
|  internal/simulator for synthetic traffic                    |
+----------------------------------------------------------+
                        |
+----------------------------------------------------------+
|  PLATFORM ABSTRACTION LAYER (internal/platform)              |
|  Windows, macOS, and Linux network interface handling        |
+----------------------------------------------------------+
```

## Layer Responsibilities

### Platform Abstraction Layer

Isolates operating system differences in network interface access. Windows requires Npcap, macOS and Linux use native BPF and libpcap bindings. This layer exposes a single interface so upper layers never depend on OS-specific code.

### Capture and Simulation Layer

Two interchangeable data sources implementing the same interface.

Capture reads live network traffic from a configured interface.
Simulator generates synthetic Modbus/TCP traffic for development and testing without requiring access to real industrial equipment.

### Protocol Parser Layer

Parses raw packet bytes into structured protocol messages. Each protocol is implemented as an independent package under internal/protocol, allowing new protocols to be added without modifying existing parsers.

### Detection Engine

Consumes structured protocol messages and evaluates them against two mechanisms.

Rule-based detection covers known-bad patterns, such as unexpected function codes or writes to protected registers.
Statistical detection covers baseline deviation, such as abnormal request frequency or timing patterns.

### Presentation Layer

The CLI provides operator commands and real-time console output. The web dashboard provides a visual view of detected events and system status.

## Data Flow

```text
Network Interface or Simulator
        |
        v
Platform Abstraction Layer
        |
        v
Capture Layer
        |
        v
Protocol Parser
        |
        v
Detection Engine
        |
        v
Report Layer
        |
        v
CLI Output, Web Dashboard, Log Files
```

## Design Constraints

No layer depends on a layer above it.
The Detection Engine does not depend on whether input came from live capture or simulation.
All protocol parsers implement the same internal interface.
No external network calls are made by the detection or protocol layers.

Author: Batuhan ALGUL
