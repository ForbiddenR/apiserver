package server

import "github.com/ForbiddenR/apiserver/pkg/server/healthz"

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

// installLivez creates the livez endpoint for this server.
func (s *GenericAPIServer) installLivez() {
	s.livezLock.Lock()
	defer s.livezLock.Unlock()
	healthz.InstallLivezHandler(s.Handler.NonGoRestfulMux, s.livezChecks...)
}