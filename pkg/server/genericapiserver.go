package server

import (
	"sync"
	"time"

	"github.com/ForbiddenR/apiserver/pkg/server/healthz"
)

type APIGroupInfo struct {
	Version string
}

// GenericAPIServer contains state for a cluster api server.
type GenericAPIServer struct {

	// minRequestTimeout is how short the request timeout can be.  This is used to build the RESTHandler
	minRequestTimeout time.Duration

	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration

	// ServingInfo holds configuration of the server.
	ServingInfo *ServingInfo

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

	s.installHealthz()

	return preparedGenericAPIServer{s}
}

func (s preparedGenericAPIServer) Run(stopCh <-chan struct{}) error {
	stopHttpServerCh := make(chan struct{})

	shutdownTimeout := s.ShutdownTimeout

	stoppedCh, listenerStoppedCh, err := s.NonBlockingRun(stopHttpServerCh, shutdownTimeout)
	if err != nil {
		return err
	}

	go func() {
		<-listenerStoppedCh
	}()

	<-stopCh
	<-stoppedCh
	return nil
}

func (s preparedGenericAPIServer) NonBlockingRun(stopCh <-chan struct{}, shutdownTimeout time.Duration) (<-chan struct{}, <-chan struct{}, error){
	// Use an internal stop channel to allow cleanup of the listeners on error.
	internalStopCh := make(chan struct{})
	var stoppedCh <-chan struct{}
	var listenrerStoppedCh <-chan struct{}
	if s.ServingInfo != nil && s.Handler != nil {
		var err error
		stoppedCh, listenrerStoppedCh, err = s.ServingInfo.Serve(s.Handler, shutdownTimeout, internalStopCh)
		if err != nil {
			close(internalStopCh)
			return nil, nil, err
		}
	}

	go func() {
		<-stopCh
		close(internalStopCh)
	}()

	return stoppedCh, listenrerStoppedCh, nil
}

func NewDefaultAPIGroupInfo(group string) APIGroupInfo {
	return APIGroupInfo{
		Version: group,
	}
}
