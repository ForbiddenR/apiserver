package healthz

import "net/http"

type HealthzChecker interface {
	Name() string
	Check(req *http.Request) error
}

var LivenessHealthz HealthzChecker = liveness{}

type liveness struct {}

func (liveness) Name() string {
	return "liveness"
}

func (liveness) Check(req *http.Request) error {
	return nil
}

var ReadinessHealthz HealthzChecker = readiness{}

type readiness struct {}

func (readiness) Name() string {
	return "readiness"
}

func (readiness) Check(req *http.Request) error {
	return nil
}