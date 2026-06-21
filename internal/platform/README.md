# Platform Abstraction

Platform-specific logic for network interface access is implemented directly in `internal/capture` using Go build tags, rather than in a separate package.

`internal/capture/pcap_source.go` (build tag `!windows`) handles Linux and macOS using libpcap.
`internal/capture/pcap_source_windows.go` (build tag `windows`) provides a stub for Windows, pending Npcap integration.

This keeps the platform-specific code colocated with the single consumer of that abstraction, avoiding an extra layer of indirection for a single use case.

This directory is reserved for future platform-specific functionality that extends beyond packet capture, such as system service integration or OS-level permission handling.
