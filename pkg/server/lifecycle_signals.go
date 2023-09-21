package server

import "sync"

type lifecycleSignal interface {
	// Signal signals the event, indicating that the event has occurred.
	// Signal is idempotent, once signaled the event stays signaled and
	// it immediately unblocks any goroutine waiting for this event.
	Signal()

	// Singaled returns a channel that is closed when the underling event
	// has been signaled. Successive calls to Signaled return the same value.
	Signaled() <-chan struct{}

	// Name returns the name of the signal, useful for logging.
	Name() string
}

// lifecycleSignals provides an abstraction of the events that
// transpire during the lifecycle of the apiserver. This abstraction makes it esay
// for use to write unit tests that can verify expected gracefull termination bahavior.
//
// GenericAPIServer can use these to either:
//   - Signal that a particular termination event has transpired
//   - wait for a designated termination event to transpire and do some action.
type lifecycleSignals struct {
	// ShutdownInitiated event is signaled when an server shutdown has been initiated.
	// It is signaled when the `stopCh` provided by the main gouroutine
	// receives a KILL signal and is closed as a consequence.
	ShutdownInitiated lifecycleSignal

	// AfterShutdownDelay event is signaled as soon as ShutdownDelayDuration
	// has elapsed since the ShutdownInitiated event.
	// ShutdownDelayDuration allows the server to delay shutdown for some time.
	AfterShutdownDelayDuration lifecycleSignal

	// PreShutdownHooksStopped event is signaled when all registered
	// preshutdown hook(s) have finished running.
	PreShutdownHooksStopped lifecycleSignal

	// NotAcceptingNewRequest event is signaled when the server is no
	// longer accepting any new request, from this point on any new
	// request will receive an error.
	NotAcceptingNewRequest lifecycleSignal
}

// ShuttingDown returns the lifecycle signal that is signaled when
// the server is not accepting any new requests.
// this is the lifecycle event that is exported to the request handler
// logic to indicate that the server is shutting down.
func (s lifecycleSignals) ShuttingDown() <-chan struct{} {
	return s.NotAcceptingNewRequest.Signaled()
}

// newLifecycleSignals returns an instance of lifecycleSignals interface to be used
// to coordinate lifecycle of the apiserver
func newLifecycleSignals() lifecycleSignals {
	return lifecycleSignals{
		ShutdownInitiated:          newNamedChannelWrapper("ShutdownInitiated"),
		AfterShutdownDelayDuration: newNamedChannelWrapper("AfterShutdownDelayDuration"),
		PreShutdownHooksStopped:    newNamedChannelWrapper("PreShutdownHooksStopped"),
		NotAcceptingNewRequest:     newNamedChannelWrapper("NotAcceptingNewRequest"),
	}
}

func newNamedChannelWrapper(name string) lifecycleSignal {
	ncw := &namedChannelWrapper{
		name: name,
		ch:   make(chan struct{}),
	}
	ncw.once = sync.OnceFunc(func() { close(ncw.ch) })
	return ncw
}

type namedChannelWrapper struct {
	name string
	once func()
	ch   chan struct{}
}

func (e *namedChannelWrapper) Signal() {
	e.once()
}

func (e *namedChannelWrapper) Signaled() <-chan struct{} {
	return e.ch
}

func (e *namedChannelWrapper) Name() string {
	return e.name
}
