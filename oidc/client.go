package oidc

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	Do(req *http.Request, values url.Values) (*http.Response, error)
	GetRedirectURL() string
}

type PostClient struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (p *PostClient) Do(req *http.Request, values url.Values) (*http.Response, error) {
	values.Set("client_id", p.ClientID)
	values.Set("client_secret", p.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(strings.NewReader(values.Encode()))
	return http.DefaultClient.Do(req)
}

func (p *PostClient) GetRedirectURL() string {
	return p.RedirectURL
}

func NewPostClient(clientID, clientSecret, redirectURL string) Client {
	return &PostClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

type BasicAuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func (b *BasicAuthClient) Do(req *http.Request, values url.Values) (*http.Response, error) {
	req.SetBasicAuth(b.ClientID, b.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(strings.NewReader(values.Encode()))
	return http.DefaultClient.Do(req)
}

func (b *BasicAuthClient) GetRedirectURL() string {
	return b.RedirectURL
}

func NewBasicAuthClient(clientID, clientSecret, redirectURL string) Client {
	return &BasicAuthClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

type MTLSClient struct {
	ClientID    string
	RedirectURL string
	HTTPClient  *http.Client
}

func (m *MTLSClient) Do(req *http.Request, values url.Values) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(strings.NewReader(values.Encode()))
	return http.DefaultClient.Do(req)
}

func (m *MTLSClient) GetRedirectURL() string {
	return m.RedirectURL
}

func NewMTLSClient(clientID, redirectURL string, config *http.Transport) Client {
	res := &MTLSClient{
		ClientID:    clientID,
		RedirectURL: redirectURL,
		HTTPClient: &http.Client{
			Transport: config,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar:     nil,
			Timeout: 0,
		},
	}
	return res
}
