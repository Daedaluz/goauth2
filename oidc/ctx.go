package oidc

import (
	"context"
	"net/http"

	"github.com/lestrrat-go/jwx/jwk"
)

type ctxKey int

const (
	keyHTTPClient ctxKey = iota
	keyJwksFetchOptions
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func WithHTTPClient(c HTTPClient) context.Context {
	return context.WithValue(context.Background(), keyHTTPClient, c)
}

func httpclientFromContext(ctx context.Context) HTTPClient {
	if c, ok := ctx.Value(keyHTTPClient).(*http.Client); ok {
		return c
	}
	return http.DefaultClient
}

func WithJwkFetchOptions(opts ...jwk.FetchOption) context.Context {
	return context.WithValue(context.Background(), keyJwksFetchOptions, opts)
}

func jwkFetchOptionsFromContext(ctx context.Context) []jwk.FetchOption {
	if opts, ok := ctx.Value(keyJwksFetchOptions).([]jwk.FetchOption); ok {
		return opts
	}
	return nil
}
