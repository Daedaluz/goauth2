package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/daedaluz/goauth2/ciba"
	"github.com/daedaluz/goauth2/oidc"
)

func main() {
	issuer, err := oidc.NewIssuer(context.Background(),
		"https://idp.inits.se",
		oidc.NewBasicAuthClient("9a8292bb-6f89-4b16-9113-71e5514e96b2", "b93ec9fa-c516-4262-adf8-ea87784987da", "http://localhost:8080/oauth2/callback"),
	)
	if err != nil {
		panic(err)
	}
	x, err := ciba.StartAuthentication(context.Background(), issuer)
	if err != nil {
		panic(err)
	}
	slog.Info("Polling for result", "authenticate here", x.Request.QRData)
	for {
		res, err := x.Poll(context.Background())
		if err != nil {
			slog.Error("Poll error", "error", err)
		} else {
			slog.Info("Poll result", "result", res)
			break
		}
		time.Sleep(1 * time.Second)
	}
}
