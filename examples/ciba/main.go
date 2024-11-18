package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/daedaluz/goauth2/ciba"
	"github.com/daedaluz/goauth2/oidc"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/mdp/qrterminal"
	"rsc.io/qr"
)

func main() {
	issuer, err := oidc.NewIssuer(context.Background(),
		"https://idp.inits.se",
		oidc.NewBasicAuthClient("9a8292bb-6f89-4b16-9113-71e5514e96b2", "b93ec9fa-c516-4262-adf8-ea87784987da", "http://localhost:8080/oauth2/callback"),
	)
	if err != nil {
		panic(err)
	}
	var opts []ciba.Option
	if len(os.Args) > 1 {
		opts = append(opts, ciba.WithLoginHint(os.Args[1]))
	}
	x, err := ciba.StartAuthentication(context.Background(), issuer, opts...)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("\033[2J")
		qrBuilder := strings.Builder{}
		qrterminal.GenerateHalfBlock(x.Request.QRData, qr.L, &qrBuilder)
		fmt.Println(qrBuilder.String())
		fmt.Println(x.Request.QRData)
		res, err := x.Poll(context.Background())
		if err != nil {
			fmt.Println("Poll error", "error", err)
		} else {
			fmt.Println("\033[2J")
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "    ")
			enc.Encode(res)
			fmt.Println()
			token, err := jwt.ParseString(res.IDToken)
			if err != nil {
				fmt.Println("Couldn't parse IDToken", err)
				return
			}
			enc.Encode(token)
			break
		}
		time.Sleep(1 * time.Second)
	}
}
