# Release 0.34.2

## What's New

* The circuit id is now available in the SDK on the client and hosting side
   * Requires 0.34.2+ routers
   * Requests SDK support. Currently supported in the Go SDK 0.23.11+
* Bug fixes

## Component Updates and Bug Fixes
* github.com/hanzozt/edge-api: [v0.26.13 -> v0.26.14](https://github.com/hanzozt/edge-api/compare/v0.26.13...v0.26.14)
* github.com/hanzozt/sdk-golang: [v0.23.14 -> v0.23.15](https://github.com/hanzozt/sdk-golang/compare/v0.23.14...v0.23.15)
* github.com/hanzozt/secretstream: [v0.1.17 -> v0.1.18](https://github.com/hanzozt/secretstream/compare/v0.1.17...v0.1.18)
    * [Issue #24](https://github.com/hanzozt/secretstream/issues/24) - Potential side channel issue
    * [Issue #25](https://github.com/hanzozt/secretstream/issues/25) - Reads from crypto/rand not checked for errors

* github.com/hanzozt/zt: [v0.34.1 -> v0.34.2](https://github.com/hanzozt/zt/compare/v0.34.1...v0.34.2)
    * [Issue #1831](https://github.com/hanzozt/zt/issues/1831) - Circuit ID should be returned in the response to a Dial request
    * [Issue #1873](https://github.com/hanzozt/zt/issues/1873) - host.v1 health check time.Duration unconvertible

# Release 0.34.1

## What's New

* Updates version of go to 1.22.x
    * As usual when updating the go version, this is the only change in this release

# Release 0.34.0

## What's New

* Bug fixes and performance enhancements
* Version number is bumped as a large chunk of HA was merged up. The next version bump is likely to bring HA to alpha status.

## Component Updates and Bug Fixes

* github.com/hanzozt/channel/v2: [v2.0.119 -> v2.0.122](https://github.com/hanzozt/channel/compare/v2.0.119...v2.0.122)
* github.com/hanzozt/edge-api: [v0.26.12 -> v0.26.14](https://github.com/hanzozt/edge-api/compare/v0.26.12...v0.26.14)
* github.com/hanzozt/foundation/v2: [v2.0.37 -> v2.0.40](https://github.com/hanzozt/foundation/compare/v2.0.37...v2.0.40)
* github.com/hanzozt/identity: [v1.0.70 -> v1.0.73](https://github.com/hanzozt/identity/compare/v1.0.70...v1.0.73)
* github.com/hanzozt/metrics: [v1.2.45 -> v1.2.48](https://github.com/hanzozt/metrics/compare/v1.2.45...v1.2.48)
* github.com/hanzozt/runzmd: [v1.0.38 -> v1.0.41](https://github.com/hanzozt/runzmd/compare/v1.0.38...v1.0.41)
* github.com/hanzozt/sdk-golang: [v0.22.28 -> v0.23.14](https://github.com/hanzozt/sdk-golang/compare/v0.22.28...v0.23.14)
    * [Issue #524](https://github.com/hanzozt/sdk-golang/issues/524) - Add circuit id to edge.Conn, so sdk connections can be correlated with network traffic
    * [Issue #515](https://github.com/hanzozt/sdk-golang/issues/515) - Service hosting improvements
    * [Issue #501](https://github.com/hanzozt/sdk-golang/issues/501) - Improve hosting session management

* github.com/hanzozt/secretstream: [v0.1.16 -> v0.1.17](https://github.com/hanzozt/secretstream/compare/v0.1.16...v0.1.17)
* github.com/hanzozt/storage: [v0.2.30 -> v0.2.33](https://github.com/hanzozt/storage/compare/v0.2.30...v0.2.33)
* github.com/hanzozt/transport/v2: [v2.0.122 -> v2.0.125](https://github.com/hanzozt/transport/compare/v2.0.122...v2.0.125)
* github.com/hanzozt/zt: [v0.33.1 -> v0.34.0](https://github.com/hanzozt/zt/compare/v0.33.1...v0.34.0)
    * [Issue #1858](https://github.com/hanzozt/zt/issues/1858) - add option to create a generic env file instead of a BASH script
    * [Issue #1428](https://github.com/hanzozt/zt/issues/1428) - Investigate policy integrity performance
    * [Issue #1854](https://github.com/hanzozt/zt/issues/1854) - Controller can try to send unroute to router which has since disconnected, causing panic
    * [Issue #1576](https://github.com/hanzozt/zt/issues/1576) - Don't scan for posture checks if there are no posture checks 
    * [Issue #1849](https://github.com/hanzozt/zt/issues/1849) - Session Sync shouldn't be able to block the control channel
    * [Issue #1846](https://github.com/hanzozt/zt/issues/1846) - Looking up api session certs for api sessions is inefficient
