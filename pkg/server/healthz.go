package server

import (
	"fmt"

	"github.com/ForbiddenR/apiserver/pkg/server/healthz"
	"github.com/valyala/fasthttp"
)

// AddReadyzChecks allows you to add a HealthzCheck to readyz.
func (s *GenericAPIServer) AddReadyzChecks(checks ...healthz.HealthzChecker) error {
	s.readyzLock.Lock()
	defer s.readyzLock.Unlock()
	if s.readyzChecksInstalled {
		return fmt.Errorf("unable to add because the readyz endpoint has already been created")
	}
	s.readyzChecks = append(s.readyzChecks, checks...)
	return nil
}

// installHealthz creates the halthz endpoint for this server.
func (s *GenericAPIServer) installHealthz() {
	s.healthzLock.Lock()
	defer s.healthzLock.Unlock()
	s.healthzChecksInstalled = true
	healthz.InstallHandler(s.Handler.NonGoRestfulMux, s.healthzChecks...)
}

func (s *GenericAPIServer) installReadyz() {
	s.readyzLock.Lock()
	defer s.readyzLock.Unlock()
	s.readyzChecksInstalled = true
	healthz.InstallReadyzHandler(s.Handler.NonGoRestfulMux, s.readyzChecks...)
}

func (s *GenericAPIServer) addReadyzShutdownCheck(stopCh <-chan struct{}) error {
	return s.AddReadyzChecks(shutdownCheck{stopCh})
}

// installLivez creates the livez endpoint for this server.
func (s *GenericAPIServer) installLivez() {
	s.livezLock.Lock()
	defer s.livezLock.Unlock()
	healthz.InstallLivezHandler(s.Handler.NonGoRestfulMux, s.livezChecks...)
}

type shutdownCheck struct {
	StopCh <-chan struct{}
}

func (shutdownCheck) Name() string {
	return "shutdown"
}

func (c shutdownCheck) Check(req *fasthttp.Request) error {
	select {
	case <-c.StopCh:
		return fmt.Errorf("process is shutting down")
	default:
	}
	return nil
}