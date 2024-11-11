package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
)

type Result struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
	Scope        []string
	TokenType    string
	Other        json.RawMessage
}

type IssuerMeta struct {
	Issuer    string `json:"issuer"`
	JWKSetURL string `json:"jwks_url"`

	TokenURL string `json:"token_endpoint"`

	CIBAURL           string   `json:"backchannel_authentication_endpoint"`
	CIBADeliveryModes []string `json:"backchannel_authentication_delivery_modes_supported"`
	meta              json.RawMessage
}

type Issuer struct {
	issuer string
	Meta   *IssuerMeta
	jwks   jwk.Set
	Client Client
}

func fetchMetadata(ctx context.Context, issuer string) (*IssuerMeta, error) {
	wellKnown := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wellKnown, http.NoBody)
	if err != nil {
		return nil, err
	}
	client := httpclientFromContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching metadata failed with status %d", resp.StatusCode)
	}
	meta := &IssuerMeta{}
	if err := json.NewDecoder(resp.Body).Decode(meta); err != nil {
		return nil, err
	}
	if meta.Issuer != issuer {
		return nil, fmt.Errorf("issuer mismatch: expected %s, got %s", issuer, meta.Issuer)
	}
	return meta, nil
}

// NewIssuer creates a new Issuer object from the given issuer URL.
// The issuer URL must be a valid OIDC issuer URL.
// The function fetches the metadata from the issuer URL and returns the Issuer object.
func NewIssuer(ctx context.Context, issuer string, c Client) (*Issuer, error) {
	meta, err := fetchMetadata(ctx, issuer)
	if err != nil {
		return nil, err
	}

	opts := []jwk.FetchOption{
		jwk.WithHTTPClient(httpclientFromContext(ctx)),
	}
	opts = append(opts, jwkFetchOptionsFromContext(ctx)...)
	jwks, err := jwk.Fetch(ctx, meta.JWKSetURL, opts...)
	if err != nil {
		return nil, err
	}

	return &Issuer{
		issuer: issuer,
		Meta:   meta,
		jwks:   jwks,
		Client: c,
	}, nil
}

func (i *Issuer) Jwks() jwk.Set {
	return i.jwks
}

func (i *Issuer) Issuer() string {
	return i.issuer
}

func (i *Issuer) Do(req *http.Request, values url.Values) (*http.Response, error) {
	return i.Client.Do(req, values)
}
