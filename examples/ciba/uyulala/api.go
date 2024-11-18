package uyulala

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	ClientID     string
	ClientSecret string
	Api          string
}

type (
	challengeID struct {
		ChallengeID string `json:"challenge_id"`
	}
)

func NewClient(api, clientID, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Api:          api,
	}
}

func (c *Client) newRequest(method, path string, body any) (*http.Request, error) {
	data, _ := json.Marshal(body)
	url := fmt.Sprintf("%s%s", c.Api, path)
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
	//teeReader := io.TeeReader(resp.Body, os.Stdout)
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) CreateUser(name string) (string, error) {
	values := map[string]any{
		"suggestedName": name,
		"timeout":       380,
		"redirect":      "https://idp.inits.se/demo",
	}
	req, err := c.newRequest(http.MethodPost, "/api/v1/service/create/user", values)
	if err != nil {
		return "", err
	}
	var ch challengeID
	_, err = c.do(req, &ch)
	if err != nil {
		return "", err
	}
	return ch.ChallengeID, nil
}
