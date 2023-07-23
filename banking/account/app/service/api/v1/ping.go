package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/account/app/service/api/responses"
	"github.com/SousaGLucas/apps/banking/account/gateways/auth"
	"github.com/SousaGLucas/apps/banking/account/gateways/ledger"
	"github.com/SousaGLucas/apps/banking/account/telemetry"
)

func Ping(authCli *auth.Client, ledgerCli *ledger.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger, err := telemetry.Logger(ctx)
		if err != nil {
			responses.InternalServerError(w, r)
			return
		}

		authPing, err := authCli.Ping(ctx)
		if err != nil {
			logger.Error("ping to auth service", zap.Error(err))

			responses.InternalServerError(w, r)
			return
		}

		ledgerPing, err := ledgerCli.Ping(ctx)
		if err != nil {
			logger.Error("ping to ledger service", zap.Error(err))

			responses.InternalServerError(w, r)
			return
		}

		accountPing := "Banking Account Microservice V1"

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, fmt.Sprintf("%s\n%s\n%s\n", authPing, ledgerPing, accountPing))
	}
}
