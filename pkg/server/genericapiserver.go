package server

import (
	"fmt"
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

	// ShutdownDelayDuration allows to block shutdown for some time.
	// during this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration

	// lifecycleSignals provides access to teh various signals that happen during the life cycle of the apiserver.
	lifecycleSignals lifecycleSignals
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

	// s.installHealthz()
	s.installLivez()

	// as soon as shutdown is initiated, readiness should start failing
	readinessStopch := s.lifecycleSignals.ShutdownInitiated.Signaled()
	err := s.addReadyzShutdownCheck(readinessStopch)
	if err != nil {
		fmt.Printf("Failed to install readyz shutdown check %s", err)
	}
	s.installReadyz()

	return preparedGenericAPIServer{s}
}

func (s preparedGenericAPIServer) Run(stopCh <-chan struct{}) error {
	delayedStopCh := s.lifecycleSignals.AfterShutdownDelayDuration
	shutdownInitiatedCh := s.lifecycleSignals.ShutdownInitiated

	// Clean up resources on shutdown.
	defer s.Destory()

	go func() {
		defer delayedStopCh.Signal()

		<-stopCh
		// As soon as shutdown is initiated, /readyz should start returning faiure.
		// This gives the load balancer a window defined by ShutdownDelayDuration to detect that /readyz is red
		// and stop sending traffic to this server.
		shutdownInitiatedCh.Signal()

		time.Sleep(s.ShutdownDelayDuration)
	}()

	shutdownTimeout := s.ShutdownTimeout

	notAcceptingNewRequestCh := s.lifecycleSignals.NotAcceptingNewRequest
	stopHttpServerCh := make(chan struct{})
	go func() {
		defer close(stopHttpServerCh)

		timeToStopHttpServerCh := notAcceptingNewRequestCh.Signaled()

		<-timeToStopHttpServerCh
	}()

	stoppedCh, listenerStoppedCh, err := s.NonBlockingRun(stopHttpServerCh, shutdownTimeout)
	if err != nil {
		return err
	}

	httpServerStoppedListeningCh := s.lifecycleSignals.HTTPServerStoppedListening
	go func() {
		<-listenerStoppedCh
		httpServerStoppedListeningCh.Signal()
	}()

	preShutdownHooksHasStoppedCh := s.lifecycleSignals.PreShutdownHooksStopped
	go func() {
		defer notAcceptingNewRequestCh.Signal()

		// wait for the delayed stopch before closing the handler chain
		<-delayedStopCh.Signaled()

		// Additionally wait for preshutdown hooks to also be finished, as some of them need
		// to send API calls to clean up after themselves (e.g. lease reconcilers removing
		// itself from the lease servers).
		<-preShutdownHooksHasStoppedCh.Signaled()
	}()

	<-stopCh

	// run shutdown hooks directly.
	func() {
		defer preShutdownHooksHasStoppedCh.Signal()
	}()

	// wait for stoppedCh that is closed when the graceful termination (server.Shutdown) is finished.
	<-listenerStoppedCh
	<-stoppedCh
	return nil
}

func (s preparedGenericAPIServer) NonBlockingRun(stopCh <-chan struct{}, shutdownTimeout time.Duration) (<-chan struct{}, <-chan struct{}, error) {
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

	// Now that listener have bound successfully, it is the
	// reponsiblity of the caller to close the provided channel to
	// ensure cleanup.
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

// Destory cleans up all its resources on shutdown.
// It starts with destroying its own resources.
func (s *GenericAPIServer) Destory() {
}
