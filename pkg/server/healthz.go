package server

import "github.com/ForbiddenR/apiserver/pkg/server/healthz"

// installHealthz creates the halthz endpoint for this server.
func (s *GenericAPIServer) installHealthz() {
	s.healthzLock.Lock()
	defer s.healthzLock.Unlock()
	s.healthzChecksInstalled = true
	healthz.InstallHandler(s.Handler.GoRestfulApp.Group("/actuator/health"), s.healthzChecks...)
}