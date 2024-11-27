package uyulala

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type Client struct {
	ClientID     string
	ClientSecret string
	API          string
	Redirect     string
}

type (
	Error struct {
		Code             int    `json:"code"`
		Message          string `json:"error_description"`
		Err              string `json:"error"`
		TechnicalMessage string `json:"technicalMsg"`
		Status           string `json:"status"`
	}
	Challenge struct {
		ChallengeID string    `json:"challenge_id"`
		Secret      string    `json:"secret"`
		StartTime   time.Time `json:"-"`
		api         string    `json:"-"`
	}
)

func (e Error) Error() string {
	return e.Message
}

func (c *Challenge) QR() string {
	dur := time.Since(c.StartTime)
	builder := jwt.NewBuilder()
	builder.
		Claim("challenge_id", c.ChallengeID).
		Claim("duration", int(dur.Seconds()))
	tok, _ := builder.Build()
	stoken, _ := jwt.Sign(tok, jwa.HS256, []byte(c.Secret))
	return fmt.Sprintf("%s/authenticator?token=%s", c.api, stoken)
}

func NewClient(api, redirect, clientID, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		API:          api,
		Redirect:     redirect,
	}
}

func (c *Client) newRequest(method, path string, body any) (*http.Request, error) {
	data, _ := json.Marshal(body)
	url := fmt.Sprintf("%s%s", c.API, path)
	var req *http.Request
	var err error
	if body != nil {
		body := strings.NewReader(string(data))
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.ClientID, c.ClientSecret)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, out any) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var e Error
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, e
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) CreateUser(name string) (*Challenge, error) {
	values := map[string]any{
		"suggestedName": name,
		"timeout":       380,
	}
	req, err := c.newRequest(http.MethodPost, "/api/v1/service/create/user", values)
	if err != nil {
		return &Challenge{}, err
	}
	ch := &Challenge{}
	_, err = c.do(req, ch)
	if err != nil {
		return &Challenge{}, err
	}
	ch.StartTime = time.Now()
	ch.api = c.API
	return ch, nil
}

func (c *Client) Collect(challenge string) (map[string]any, error) {
	values := map[string]any{
		"challengeId": challenge,
	}
	req, err := c.newRequest(http.MethodPost, "/api/v1/collect", values)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	_, err = c.do(req, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
