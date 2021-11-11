package flagx

import (
	"flag"
	"net"
)

// TCPAddr is a flag.Value that accepts a TCP address, for example "127.0.0.1:8000".
type TCPAddr struct {
	Addr *net.TCPAddr
}

var _ flag.Value = (*TCPAddr)(nil)

func (a TCPAddr) String() string {
	return a.Addr.String()
}

func (a *TCPAddr) Set(value string) error {
	if value == "" {
		a.Addr = nil
		return nil
	}

	addr, err := net.ResolveTCPAddr("tcp", value)
	if err != nil {
		return err
	}
	a.Addr = addr
	return nil
}
