package xgress_sdk

import (
	"net"

	"github.com/hanzozt/sdk-golang/zt"
)

// The Fabric API allows users to dial services and stay up-to-date with service/identity changes.
// The API is experimental and will either evolve or go away.
type Fabric interface {
	TunnelWithOptions(serviceName string, options *zt.DialOptions, conn net.Conn, halfClose bool) error
}
