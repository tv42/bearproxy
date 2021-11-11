package flagx

import (
	"flag"
	"net"
)

// IP is a flag.Value accepts IP addresses.
type IP struct {
	net.IP
}

var _ flag.Value = (*IP)(nil)

func (a *IP) Set(value string) error {
	if err := a.UnmarshalText([]byte(value)); err != nil {
		return err
	}
	return nil
}
