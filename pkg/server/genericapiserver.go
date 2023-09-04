package server

import (
	"sync"

	"github.com/ForbiddenR/apiserver/pkg/server/healthz"
)

type APIGroupInfo struct {
	Version string
}

type GenericAPIServer struct {
	// Handler holds the handlers being used by this API server.
	Handler *APIServerHandler

	// healthz checks
	healthzLock            sync.Mutex
	healthzChecks          []healthz.HealthzChecker
	healthzChecksInstalled bool
	// livez checks
	livezLock            sync.Mutex
	livezChecks          []healthz.HealthzChecker
	livezChecksInstalled bool
	// readyz checks
	readyzLock            sync.Mutex
	readyzChecks          []healthz.HealthzChecker
	readyzChecksInstalled bool
}

func (s *GenericAPIServer) InstallAPIGroups(apiGroupInfos ...*APIGroupInfo) error {
	for range apiGroupInfos {
		s.Handler.GoRestfulApp.Group("")
	}
	return nil
}

func (s *GenericAPIServer) InstallAPIGroup(apiGroupInfo *APIGroupInfo) error {
	return s.InstallAPIGroups(apiGroupInfo)
}

type preparedGenericAPIServer struct {
	*GenericAPIServer
}

func (s *GenericAPIServer) PrepareRun() preparedGenericAPIServer {
	
	return preparedGenericAPIServer{s}
}

func NewDefaultAPIGroupInfo(group string) APIGroupInfo {
	return APIGroupInfo{
		Version: group,
	}
}
