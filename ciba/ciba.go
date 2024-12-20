package ciba

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/daedaluz/goauth2/oidc"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type AuthSession struct {
	issuer         *oidc.Issuer
	hint           string
	loginHintToken string
	idTokenHint    string
	loginHint      string
	acrValues      []string
	scope          []string
	values         url.Values

	bindingMessage  string
	requestedExpiry time.Duration
	pollInterval    time.Duration

	StartTime time.Time

	Request *Request
}

type Request struct {
	AuthReqID string `json:"auth_req_id,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
	Interval  int    `json:"interval,omitempty"`
	QRData    string `json:"qr_data,omitempty"`
	QRSecret  string `json:"qr_secret,omitempty"`
}

func (r *Request) QRToken(startTime time.Time) string {
	dur := time.Since(startTime)
	builder := jwt.NewBuilder()
	builder.
		Claim("challenge_id", r.QRData).
		Claim("duration", int(dur.Seconds()))
	tok, _ := builder.Build()
	stoken, _ := jwt.Sign(tok, jwa.HS256, []byte(r.QRSecret))
	return string(stoken)
}

func Authenticate(ctx context.Context, issuer *oidc.Issuer, opts ...Option) (*oidc.Result, error) {
	sess, err := StartAuthentication(ctx, issuer, opts...)
	if err != nil {
		return nil, err
	}
	return sess.Complete(ctx)
}

func StartAuthentication(ctx context.Context, issuer *oidc.Issuer, opts ...Option) (*AuthSession, error) {
	cibaURL := issuer.Meta.CIBAURL
	if cibaURL == "" {
		return nil, fmt.Errorf("issuer %s does not support CIBA", issuer.Meta.Issuer)
	}

	sess := &AuthSession{issuer: issuer, Request: &Request{}}
	// Apply the options
	for _, opt := range opts {
		opt.Apply(sess)
	}
	values := url.Values{}
	// Prepare the request variables
	if sess.loginHintToken != "" {
		values.Set("login_hint_token", sess.loginHintToken)
	} else if sess.idTokenHint != "" {
		values.Set("id_token_hint", sess.idTokenHint)
	} else if sess.loginHint != "" {
		values.Set("login_hint", sess.loginHint)
	}
	if len(sess.scope) > 0 {
		values.Set("scope", strings.Join(sess.scope, " "))
	} else {
		values.Set("scope", "openid")
	}
	if sess.bindingMessage != "" {
		values.Set("binding_message", sess.bindingMessage)
	}
	if sess.requestedExpiry > 0 {
		nSecs := int(math.Floor(sess.requestedExpiry.Seconds()))
		values.Set("requested_expiry", fmt.Sprintf("%d", nSecs))
	}
	if len(sess.acrValues) > 0 {
		values.Set("acr_values", strings.Join(sess.acrValues, " "))
	}

	// Start the authentication
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cibaURL, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := sess.issuer.Client.Do(req, values)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		err = oidc.ParseError(resp.Body)
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(sess.Request); err != nil {
		return nil, err
	}
	sess.StartTime = time.Now()
	return sess, nil
}

func (a *AuthSession) QrCode() string {
	return fmt.Sprintf("%s?token=%s", a.issuer.Meta.CIBAQRURL, a.Request.QRToken(a.StartTime))
}

func (a *AuthSession) Poll(ctx context.Context) (*oidc.Result, error) {
	pollURL := a.issuer.Meta.TokenURL
	if pollURL == "" {
		return nil, fmt.Errorf("issuer %s does not support CIBA", a.issuer.Meta.Issuer)
	}
	values := url.Values{}
	values.Set("auth_req_id", a.Request.AuthReqID)
	values.Set("grant_type", "urn:openid:params:grant-type:ciba")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pollURL, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := a.issuer.Client.Do(req, values)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		err = oidc.ParseError(resp.Body)
		return nil, err
	}
	res := &oidc.Result{}
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (a *AuthSession) Complete(ctx context.Context) (*oidc.Result, error) {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			res, err := a.Poll(ctx)
			if err != nil {
				return nil, err
			}
			if res != nil {
				return res, nil
			}
		}
	}
}
