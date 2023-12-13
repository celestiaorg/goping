# goping

**goping** is a command-line tool designed for Kubernetes environments to address the challenge of testing network connectivity to services when traditional `ping` is ineffective due to the nature of Kubernetes Services.

## Table of Contents

- [Background](#background)
- [Installation](#installation)
- [Usage](#usage)
  - [Serve Command](#serve-command)
  - [Ping Command](#ping-command)
- [Technical Explanation](#technical-explanation)

## Background

In Kubernetes, using `ping` to test connectivity to services is ineffective due to the port requirement for Services to function. This tool addresses this limitation by providing an alternative method to test network connectivity to Kubernetes Services.

### Why Regular Ping Fails in Kubernetes Services?

Kubernetes Services are stable networking endpoints in front of Pods, exposing a DNS name, virtual IP, and port. However, traditional `ping` does not utilize ports and thus cannot activate a Kubernetes Service.

## Installation

To build the `goping` tool, ensure you have Go installed. Then run:

```bash
make build
```

This will generate the executable binary goping in the bin directory.

For building a Docker image:

```bash
make docker
```

This command will build a Docker image named `goping-server`.

## Usage

### Serve Command

The serve command starts the goping server.

```bash
goping serve --serve-addr=:8080 --log-level=info --production-mode
```

- `--serve-addr` specifies the address to serve on. Default is :8080.
- `--log-level` sets the log level (e.g., debug, info, warn, error, dpanic, panic, fatal). Default is info.
- `--production-mode` toggles production mode to disable debug logs.

### Ping Command

The `ping` command tests the connectivity to the `goping` server.

```bash
goping ping <server-address> --packet-timeout=10ms --packets-count=4 --log-level=info --production-mode
```

- `<server-address>` is the address of the server to ping.
- `--packet-timeout` sets the timeout for each packet. Default is 10ms.
- `--packets-count` specifies the number of packets to send. Default is 4.
- `--log-level` sets the log level. Default is info.
- `--production-mode` toggles production mode.

## Technical Explanation

In Kubernetes environments, the conventional `ping` command proves ineffective for testing connectivity to Services due to limitations in ICMP (Internet Control Message Protocol), which `ping` relies on and doesn't support port specifications. Kubernetes Services demand connections to arrive on specific ports to activate, a requirement `ping` fails to fulfill as it lacks port-specific functionality. To tackle this challenge, `goping` introduces a tailored solution. It simulates `ping`-like functionality, purpose-built for Kubernetes. `goping` utilizes UDP (User Datagram Protocol) to send packets targeted at Service endpoints on specified ports. This custom approach enables accurate connectivity testing to Kubernetes Services, circumventing the port-based constraints inherent in traditional `ping` commands.
