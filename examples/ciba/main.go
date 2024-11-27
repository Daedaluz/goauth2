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
	clientID     = "f0cbf686-5e81-4546-9b96-5a3acc045219"
	clientSecret = "d27144a6-2dac-458c-8357-5e37a635548d"
)

func register(name string) {
	client := uyulala.NewClient("https://idp.inits.se", "http://localhost:9094", clientID, clientSecret)
	challenge, err := client.CreateUser(name)
	if err != nil {
		fmt.Println("Error creating user", err)
		return
	}
	for {
		resp, err := client.Collect(challenge.ChallengeID)
		if err == nil {
			fmt.Println("User has completed registration")
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "    ")
			enc.Encode(resp)
			return
		}
		var uErr *uyulala.Error
		if errors.As(err, &uErr) {
			switch uErr.Status {
			case "pending", "viewed":
				fmt.Println("Waiting for user to complete registration")
			default:
				fmt.Println("Error", uErr)
			}
		}
		qrBuilder := strings.Builder{}
		qrterminal.Generate(challenge.QR(), qr.L, &qrBuilder)
		fmt.Println("\033[2J")
		fmt.Println("Follow the link in this QR code")
		fmt.Println(qrBuilder.String())
		fmt.Println("Or visit this URL in your browser")
		fmt.Println(challenge.QR())
		time.Sleep(1)
	}
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
		qrterminal.Generate(x.QrCode(), qr.L, &qrBuilder)
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
					fmt.Println("URL: ", x.QrCode())
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
