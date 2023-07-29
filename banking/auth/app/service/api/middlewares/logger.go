package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/auth/telemetry"
)

func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(telemetry.WithLogger(r.Context(), logger))
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
