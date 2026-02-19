# Release 0.29.0

## What's New

### Deprecated Binary Removal
This release removes the following deprecated binaries from the release archives.

* `ziti-controller` - replaced by `ziti controller`
* `ziti-router`     - replaced by `ziti router`
* `ziti-tunnel`     - replaced by `ziti tunnel`

The release archives now only contain the `ziti` executable. This executable is now at the root of the archive instead of nested under a `ziti` directory.

### Ziti CLI Demo Consolidation

The ziti CLI functions under `ziti learn`, namely `ziti learn demo` and `ziti learn tutorial` have been consolidated under `ziti demo`.

### Continued Quickstart Changes

The quickstart continues to evolve. A breaking change has occurred as numerous environment variables used to customize the quickstart
have changed again. A summary of changes is below

* All `ZITI_EDGE_ROUTER_` variables have been changed to just `ZITI_ROUTER_`.
  * `ZITI_EDGE_ROUTER_NAME` -> `ZITI_ROUTER_NAME`
  * `ZITI_EDGE_ROUTER_PORT` -> `ZITI_ROUTER_PORT`
  * `ZITI_EDGE_ROUTER_ADVERTISED_HOST` -> `ZITI_ROUTER_ADVERTISED_HOST`
  * `ZITI_EDGE_ROUTER_IP_OVERRIDE` -> `ZITI_ROUTER_IP_OVERRIDE`
  * `ZITI_EDGE_ROUTER_ENROLLMENT_DURATION` -> `ZITI_ROUTER_ENROLLMENT_DURATION`
  * `ZITI_EDGE_ROUTER_ADVERTISED_HOST` -> `ZITI_ROUTER_ADVERTISED_HOST`
  * `ZITI_EDGE_ROUTER_LISTENER_BIND_PORT` -> `ZITI_ROUTER_LISTENER_BIND_PORT`
* Additional variables have been added to support "alternative addresses" and "alternative PKI", for example
  to support using Let's Encrypt certificates easily in the quickstarts.
* New variables were introduced to allow automatic generation of the `alt_server_certs` section. Both variables
  must be supplied for the variables to impact the configurations.
  * `ZITI_PKI_ALT_SERVER_CERT` - "Alternative server certificate. Must be specified with ZITI_PKI_ALT_SERVER_KEY"
  * `ZITI_PKI_ALT_SERVER_KEY` - "Key to use with the alternative server certificate. Must be specified with ZITI_PKI_ALT_SERVER_CERT"
* New variables were introduced to allow one to override and customize the CSR section of routers which is used during enrollment.
  * `ZITI_ROUTER_CSR_C` - "The country (C) to use for router CSRs"
  * `ZITI_ROUTER_CSR_ST` - "The state/province (ST) to use for router CSRs"
  * `ZITI_ROUTER_CSR_L` - "The locality (L) to use for router CSRs"
  * `ZITI_ROUTER_CSR_O` - "The organization (O) to use for router CSRs"
  * `ZITI_ROUTER_CSR_OU` - "The organization unit to use for router CSRs"
  *	`ZITI_ROUTER_CSR_SANS_DNS` - "The DNS name used in the CSR request"
* New variable `ZITI_CTRL_EDGE_BIND_ADDRESS` allows controlling the IP the edge API uses

## Component Updates and Bug Fixes

* github.com/hanzozt/channel/v2: [v2.0.81 -> v2.0.84](https://github.com/hanzozt/channel/compare/v2.0.81...v2.0.84)
* github.com/hanzozt/edge: [v0.24.348 -> v0.24.364](https://github.com/hanzozt/edge/compare/v0.24.348...v0.24.364)
  * [Issue #1543](https://github.com/hanzozt/edge/issues/1543) - controller ca normalization can go into infinite loop on startup with bad certs

* github.com/hanzozt/edge-api: [v0.25.29 -> v0.25.30](https://github.com/hanzozt/edge-api/compare/v0.25.29...v0.25.30)
* github.com/hanzozt/fabric: [v0.23.39 -> v0.23.45](https://github.com/hanzozt/fabric/compare/v0.23.39...v0.23.45)
* github.com/hanzozt/foundation/v2: [v2.0.26 -> v2.0.28](https://github.com/hanzozt/foundation/compare/v2.0.26...v2.0.28)
* github.com/hanzozt/identity: [v1.0.57 -> v1.0.59](https://github.com/hanzozt/identity/compare/v1.0.57...v1.0.59)
* github.com/hanzozt/runzmd: [v1.0.26 -> v1.0.28](https://github.com/hanzozt/runzmd/compare/v1.0.26...v1.0.28)
* github.com/hanzozt/sdk-golang: [v0.20.67 -> v0.20.78](https://github.com/hanzozt/sdk-golang/compare/v0.20.67...v0.20.78)
* github.com/hanzozt/storage: [v0.2.8 -> v0.2.11](https://github.com/hanzozt/storage/compare/v0.2.8...v0.2.11)
* github.com/hanzozt/transport/v2: [v2.0.91 -> v2.0.93](https://github.com/hanzozt/transport/compare/v2.0.91...v2.0.93)
* github.com/hanzozt/metrics: [v1.2.27 -> v1.2.30](https://github.com/hanzozt/metrics/compare/v1.2.27...v1.2.30)
* github.com/hanzozt/secretstream: [v0.1.9 -> v0.1.10](https://github.com/hanzozt/secretstream/compare/v0.1.9...v0.1.10)
* github.com/hanzozt/ziti: [v0.28.4 -> v0.29.0](https://github.com/hanzozt/ziti/compare/v0.28.4...v0.29.0)
  * [Issue #1180](https://github.com/hanzozt/ziti/issues/1180) - Add ability to debug failed smoketests
  * [Issue #1169](https://github.com/hanzozt/ziti/issues/1169) - Consolidate demo and tutorial under demo
  * [Issue #1168](https://github.com/hanzozt/ziti/issues/1168) - Remove ziti-controller, ziti-router and ziti-tunnel executables from build
  * [Issue #1158](https://github.com/hanzozt/ziti/issues/1158) - Add iperf tests to ziti smoketest
