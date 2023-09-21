package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	defaultKeepAlivePeriod = 3 * time.Minute
)

// Serve runs the http server.
// The actual server loop (stoppable by closing stopCh) runs in a go routine, i.e. Serve does not block.
// It returns a stoppedCh that is closed when all non-hijacked active requests have been processed.
// It returns a listenerStoppedCh that is closed when the underlying http Server has stopped listening.
func (s *ServingInfo) Serve(handler *APIServerHandler, shutdownTimeout time.Duration, stopCh <-chan struct{}) (<-chan struct{}, <-chan struct{}, error) {
	if s.Listener == nil {
		return nil, nil, fmt.Errorf("listener must not be nil")
	}
	return RunServer(handler.GoRestfulApp, s.Listener, shutdownTimeout, stopCh)
}

// RunServer spawns a go-routine continously serving until the stopCh is
// closed.
// It returns a stoppedCh that is closed when all non-hijacked active requests
// have been processed.
// This function does not block.
func RunServer(
	server *fiber.App,
	ln net.Listener,
	shutDownTimeout time.Duration,
	stopCh <-chan struct{},
) (<-chan struct{}, <-chan struct{}, error) {
	if ln == nil {
		return nil, nil, fmt.Errorf("listener must not be nil")
	}

	// Shutdown server gracefully.
	serverShutdownCh, listenerStoppedCh := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(serverShutdownCh)
		<-stopCh
		ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
		server.ShutdownWithContext(ctx)
		cancel()
	}()

	go func() {
		defer close(listenerStoppedCh)

		listener := tcpKeepAliveListener{ln}

		err := server.Listener(listener)

		msg := fmt.Sprintf("Stopped listening on %s", ln.Addr().String())
		select {
		case <-stopCh:
			fmt.Println(msg)
		default:
			panic(fmt.Sprintf("%s due to error: %v", msg, err))
		}
	}()
	return serverShutdownCh, listenerStoppedCh, nil
}

type tcpKeepAliveListener struct {
	net.Listener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	c, err := ln.Listener.Accept()
	if err != nil {
		return nil, err
	}
	if tc, ok := c.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(defaultKeepAlivePeriod)
	}
	return c, nil
}
