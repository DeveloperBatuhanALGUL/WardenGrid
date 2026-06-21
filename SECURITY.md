# Security Policy

## Scope

WardenGrid is a defensive monitoring tool for Industrial Control Systems and SCADA environments. This policy covers vulnerabilities in the WardenGrid codebase itself, including the protocol parsers, detection engine, capture layer, and CLI.

## Supported Versions

WardenGrid is in early development. Security fixes are applied to the latest commit on the `main` branch. There are no maintained release branches at this stage.

## Reporting a Vulnerability

Do not open a public GitHub issue for security vulnerabilities.

Report vulnerabilities privately through GitHub Security Advisories on this repository, or by contacting the author directly at batuhanalgul@proton.me.

Include the following in your report.

A clear description of the vulnerability.
Steps to reproduce the issue, including any relevant network traffic samples or configuration.
The potential impact, including whether it affects detection accuracy, parsing correctness, or system stability.

## Disclosure Process

Reports are acknowledged as soon as possible after receipt. A fix or mitigation plan will be communicated once the issue is confirmed. Public disclosure happens only after a fix is available, in coordination with the reporter.

## Out of Scope

This project does not control the security of the underlying operating system, network infrastructure, or third-party industrial equipment being monitored. Vulnerabilities in upstream Go dependencies should be reported to the respective upstream projects.

## Author

Batuhan ALGUL
