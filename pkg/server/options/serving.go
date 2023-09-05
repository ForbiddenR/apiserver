package options

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/ForbiddenR/apiserver/pkg/server"
)

type ServingOptions struct {
	BindAddress net.IP
	// BindPort is ignored when Listener is set, will serve https even with 0.
	BindPort int
	// BindNetwork is the type of network to bind to - defaults to "tcp". accepts "tcp",
	// "tcp4", and "tcp6"
	BindNetwork string
	// Required set to true mean that BindPort cannot be zero.
	Required bool
	// ExternalAddress is the addrss advertised. even if BindAddress is a loopback. By default this
	// is set to BindAddress if the later no loopback. or to the fist host interface addres.
	ExternalAddress net.IP
	// Listener is a server network listener.
	// either Listener or BindAddress/Bindport/BindNetwork is,
	// if Listener is set, use it and omit BindAddress/BindPort/BindNetwork.
	Listener net.Listener
}

func NewServingOptions() *ServingOptions {
	return &ServingOptions{
		BindAddress: net.IPv4(0, 0, 0, 0),
		BindPort:    8080,
	}
}

func (s *ServingOptions) ApplyTo(config **server.ServingInfo) error {
	if s == nil {
		return nil
	}
	if s.BindPort <= 0 && s.Listener == nil {
		return nil
	}

	if s.Listener == nil {
		var err error
		addr := net.JoinHostPort(s.BindAddress.String(), strconv.Itoa(s.BindPort))

		c := net.ListenConfig{}

		s.Listener, s.BindPort, err = CreateListener(s.BindNetwork, addr, c)
		if err != nil {
			return fmt.Errorf("failed to create listener: %v", err)
		}
	} else {
		if _, ok := s.Listener.Addr().(*net.TCPAddr); !ok {
			return fmt.Errorf("failed to parse ip and port from listener")
		}
		s.BindPort = s.Listener.Addr().(*net.TCPAddr).Port
		s.BindAddress = s.Listener.Addr().(*net.TCPAddr).IP
	}

	*config = &server.ServingInfo{
		Listener: s.Listener,
	}

	return nil
}

func CreateListener(network, addr string, config net.ListenConfig) (net.Listener, int, error) {
	if len(network) == 0 {
		network = "tcp"
	}

	ln, err := config.Listen(context.TODO(), network, addr)
	if err != nil {
		return nil, 0,fmt.Errorf("failed to listen on %v: %v", addr, err)
	}

	// get port
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		ln.Close()
		return nil, 0, fmt.Errorf("invalid listen address: %q", ln.Addr().String())
	}

	return ln, tcpAddr.Port, nil
}