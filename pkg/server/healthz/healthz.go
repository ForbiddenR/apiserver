package healthz

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type HealthzChecker interface {
	Name() string
	Check(req *fasthttp.Request) error
}

var LivenessHealthz HealthzChecker = liveness{}

type liveness struct{}

func (liveness) Name() string {
	return "liveness"
}

func (liveness) Check(req *fasthttp.Request) error {
	return nil
}

var ReadinessHealthz HealthzChecker = readiness{}

type readiness struct{}

func (readiness) Name() string {
	return "readiness"
}

func (readiness) Check(req *fasthttp.Request) error {
	return nil
}

// InstallHandler registers handlers for health checking on the
// "/healthz" to mux. *All handlers* for the path must be specified in
// exactly one call to InstallHandler. Calling InstallHandler more
// than once for the same path and mux will result in a panic.
func InstallHandler(mux mux, checks ...HealthzChecker) {
	InstallPathHandler(mux, "/healthz", checks...)
}

// InstallReadyzHandler registers handlers for health checking on the path
// "/readiness" to mux. *All handlers* for the path must be specified in
// exactly one call to InstallReadyzHandler. Calling InstallReadyzHandler more
// than once for the same path and mux will result in a panic.
func InstallReadyzHandler(mux mux, checks ...HealthzChecker) {
	InstallPathHandler(mux, "/readiness", checks...)
}

// InstallLivezHandler registers handlers for health checking on the path
// "/liveness" to mux. *All handlers* for the path must be specified in
// exactly one call to InstallLivezHandler. Calling InstallLivezHandler more
// than once for the same path and mux will result in a panic.
func InstallLivezHandler(mux mux, checks ...HealthzChecker) {
	InstallPathHandler(mux, "/liveness", checks...)
}

// InstallPathHandler registers handlers for health checking on
// a specific path to mux. *All handlers* for the path must be
// specified in exactly one call to InstallPathHandler. Calling
// InstallPathHandler more than once for the same path and mux will
// result in a panic.
func InstallPathHandler(mux mux, path string, checks ...HealthzChecker) {
	InstallPathHandlerWithHealthyFunc(mux, path, checks...)
}

func InstallPathHandlerWithHealthyFunc(mux mux, path string, checks ...HealthzChecker) {
	name := strings.Split(strings.TrimPrefix(path, "/"), "/")[0]
	mux.Add(fiber.MethodGet, path, handleRootHealth(name, checks...))
}

// mux is a interface describing the methods InstallHandler requires.
type mux interface {
	Add(verb string, pattern string, handlers ...fiber.Handler) fiber.Router
}

func handleRootHealth(name string, checks ...HealthzChecker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// var individualCheckOutput bytes.Buffer
		for _, check := range checks {
			if err := check.Check(ctx.Request()); err != nil {
				return ctx.Status(fasthttp.StatusNotFound).JSON(GetHealthzResponse(check.Name(), 1))
			}
		}
		return ctx.Status(fasthttp.StatusOK).JSON(GetHealthzResponse(name, 0))
	}
}

func GetHealthzResponse(name string, code int) *HealthzResponse {
	switch name {
	case "liveness":
		return NewHealthzResponse((&LivenessComponent{}).SetStatus(code))
	case "readiness":
		return NewHealthzResponse((&ReadinessComponent{}).SetStatus(code))
	}
	return nil
}

type HealthzResponse struct {
	Status     string     `json:"status"`
	Components Components `json:"components"`
}

func NewHealthzResponse(components Components) *HealthzResponse {
	return &HealthzResponse{
		Status:     components.GetStatus(),
		Components: components,
	}
}

type Components interface {
	SetStatus(code int) Components
	GetStatus() string
}

var _ Components = &LivenessComponent{}

type LivenessComponent struct {
	LivenessState Status `json:"livenessstate"`
}

func (l *LivenessComponent) SetStatus(code int) Components {
	if code == 0 {
		l.LivenessState.Status = "UP"
	} else {
		l.LivenessState.Status = "DOWN"
	}
	return l
}

func (l *LivenessComponent) GetStatus() string {
	return l.LivenessState.Status
}

var _ Components = &ReadinessComponent{}

type ReadinessComponent struct {
	ReadinessState Status `json:"readinessstate"`
}

func (r *ReadinessComponent) SetStatus(code int) Components {
	if code == 0 {
		r.ReadinessState.Status = "UP"
	} else {
		r.ReadinessState.Status = "OUT_OF_SERVICE"
	}
	return r
}

func (r *ReadinessComponent) GetStatus() string {
	return r.ReadinessState.Status
}

type Status struct {
	Status string `json:"status"`
}
