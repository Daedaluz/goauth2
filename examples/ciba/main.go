package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/daedaluz/goauth2/ciba"
	"github.com/daedaluz/goauth2/examples/ciba/uyulala"
	"github.com/daedaluz/goauth2/oidc"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/mdp/qrterminal"
	"rsc.io/qr"
)

var (
	// Yes, these are intentionally leaked secrets.
	clientID     = "9a8292bb-6f89-4b16-9113-71e5514e96b2"
	clientSecret = "b93ec9fa-c516-4262-adf8-ea87784987da"
)

func register(name string) {
	client := uyulala.NewClient("https://idp.inits.se", clientID, clientSecret)
	challengeID, err := client.CreateUser(name)
	if err != nil {
		fmt.Println("Error creating user", err)
		return
	}
	fmt.Println("\033[2J")
	fmt.Println("Challenge ID: ", challengeID)
	fmt.Printf("Please register the key at https://idp.inits.se/authenticator?id=%s\n", challengeID)
}

func main() {
	issuer, err := oidc.NewIssuer(context.Background(),
		"https://idp.inits.se",
		oidc.NewBasicAuthClient(clientID, clientSecret, ""),
	)
	if err != nil {
		panic(err)
	}
	var opts []ciba.Option
	if len(os.Args) == 3 && os.Args[1] == "register" {
		register(os.Args[2])
		return
	} else if len(os.Args) == 2 {
		opts = append(opts, ciba.WithLoginHint(os.Args[1]), ciba.WithBindingMessage("CIBA Example wants user "+os.Args[1]+" to authenticate"))
	} else {
		opts = append(opts, ciba.WithBindingMessage("CIBA Example wants to identify you"))
	}
	x, err := ciba.StartAuthentication(context.Background(), issuer, opts...)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("\033[2J")
		qrBuilder := strings.Builder{}
		qrterminal.Generate(x.Request.QRData, qr.L, &qrBuilder)
		fmt.Println(qrBuilder.String())
		fmt.Println(x.Request.QRData)
		res, err := x.Poll(context.Background())
		if err != nil {
			e := &oidc.ErrorResponse{}
			isOIDCErr := errors.As(err, &e)
			if isOIDCErr {
				switch e.Err {
				case oidc.ErrAccessDenied:
					fmt.Println("Access Denied: ", e.ErrorDescription)
					return
				case oidc.ErrAuthorizationPending:
					fmt.Println("Authorization Pending: ", e.ErrorDescription)
				default:
					fmt.Println("Error: ", e.ErrorDescription)
					return
				}
			} else {
				fmt.Println("Poll error", "error", err, fmt.Sprintf("%T", err))
			}
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
		time.Sleep(time.Second / 2)
	}
}
