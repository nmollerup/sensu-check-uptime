# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial Go port of [sensu-plugins-uptime-checks](https://github.com/sensu-plugins/sensu-plugins-uptime-checks)
- `check-uptime`: checks system uptime and warns if the system has recently rebooted; supports `--warn`/`-w` threshold and `--greater-than`/`-g` comparison direction; linux and windows support
- `metrics-uptime`: outputs uptime and idle time metrics in Graphite plaintext format from `/proc/uptime`; supports `--scheme`/`-s` metric naming prefix; linux only, matching the original Ruby plugin's scope
