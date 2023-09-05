package healthz

import (
	"bytes"
	"fmt"
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
		var individualCheckOutput bytes.Buffer
		for _, check := range checks {
			if err := check.Check(ctx.Request()); err != nil {
				fmt.Fprintf(&individualCheckOutput, "[-]%s failed: reason withheld\n", check.Name())
			} else {
				fmt.Fprintf(&individualCheckOutput, "[+]%s ok\n", check.Name())
			}
		}
		ctx.Response().Header.Set("Content-Type", "text/plain; charset=utf-8")
		ctx.Response().Header.Set("X-Content-Type-Options", "nosniff")
		return ctx.SendString(individualCheckOutput.String())
	}
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	var individualCheckOutput bytes.Buffer
	// 	for _, check := range checks {
	// 		if err := check.Check(r); err != nil {
	// 			fmt.Fprintf(&individualCheckOutput, "[-]%s failed: reason withheld\n", check.Name())
	// 		} else {
	// 			fmt.Fprintf(&individualCheckOutput, "[+]%s ok\n", check.Name())
	// 		}
	// 	}

	// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// 	w.Header().Set("X-Content-Type-Options", "nosniff")
	// 	individualCheckOutput.WriteTo(w)
	// 	fmt.Fprintf(w, "%s check passed\n", name)
	// }
}
