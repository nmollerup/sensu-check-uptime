[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nmollerup/sensu-check-uptime)
[![Go Test](https://github.com/nmollerup/sensu-check-uptime/actions/workflows/test.yml/badge.svg)](https://github.com/nmollerup/sensu-check-uptime/actions/workflows/test.yml)
[![Go Lint](https://github.com/nmollerup/sensu-check-uptime/actions/workflows/lint.yml/badge.svg)](https://github.com/nmollerup/sensu-check-uptime/actions/workflows/lint.yml)

## sensu-check-uptime

A Go port of [sensu-plugins-uptime-checks](https://github.com/sensu-plugins/sensu-plugins-uptime-checks), distributed as standalone binaries (no Ruby runtime required).

- [Overview](#overview)
- [Commands](#commands)
  - [check-uptime](#check-uptime)
  - [metrics-uptime](#metrics-uptime)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)

### Overview

This plugin provides native instrumentation for collecting uptime and idle time metrics, and for alerting when a monitored system has recently rebooted.

### Commands

#### check-uptime

Checks the system's uptime and warns if the system has recently rebooted. Supported on Linux and Windows.

```
Usage:
  check-uptime [flags]

Flags:
  -g, --greater-than         Compare uptime > threshold. Default behavior is uptime < threshold
  -w, --warn int              Warn threshold in seconds (default 180)
```

#### metrics-uptime

Outputs uptime and idle time metrics read from `/proc/uptime`, in the Graphite plaintext protocol (`path value timestamp`). Linux only, matching the scope of the original Ruby plugin.

```
Usage:
  metrics-uptime [flags]

Flags:
  -s, --scheme string   Metric naming scheme, text to prepend to metric (default "<hostname>.uptime")
```

Example output:

```
myhost.uptime.uptime 123456 1700000000
myhost.uptime.idletime 654321 1700000000
```

### Configuration

#### Asset registration

Assets are the best way to make use of this plugin. If you're using sensuctl 5.13 or later, you can use the following command to add the asset:

```
sensuctl asset add nmollerup/sensu-check-uptime
```

If you're using an earlier version of sensuctl, you can download the asset definition from this project's [Bonsai asset index page](https://bonsai.sensu.io/assets/nmollerup/sensu-check-uptime).

#### Check definition

```yaml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-uptime
spec:
  command: "check-uptime -w 180"
  handlers: []
  high_flap_threshold: 0
  interval: 10
  low_flap_threshold: 0
  publish: true
  runtime_assets:
  - nmollerup/sensu-check-uptime
  subscriptions:
  - linux
```

```yaml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: metrics-uptime
spec:
  command: "metrics-uptime"
  handlers:
  - influxdb
  high_flap_threshold: 0
  interval: 60
  low_flap_threshold: 0
  publish: true
  runtime_assets:
  - nmollerup/sensu-check-uptime
  subscriptions:
  - linux
  output_metric_format: graphite_plaintext
  output_metric_handlers:
  - influxdb
```

### Installation from source

Download the latest release binaries from the [releases page](https://github.com/nmollerup/sensu-check-uptime/releases), or build from source:

```
go build -o bin/check-uptime ./cmd/check-uptime
go build -o bin/metrics-uptime ./cmd/metrics-uptime
```
