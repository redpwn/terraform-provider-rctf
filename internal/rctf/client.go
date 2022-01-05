package rctf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redpwn/terraform-provider-rctf/internal/version"
	"net/http"
	"net/url"
)

type Client struct {
	BaseUrl   string
	AuthToken string
	http      *http.Client
}

func New(baseUrl, authToken string) (*Client, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	c := &Client{
		BaseUrl:   fmt.Sprintf("%s://%s/api/v1/", u.Scheme, u.Host),
		AuthToken: authToken,
		http:      http.DefaultClient,
	}
	return c, nil
}

type response struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

func (r response) String() string {
	return fmt.Sprintf("rctf response: %s: %s", r.Kind, r.Message)
}

func (c *Client) req(ctx context.Context, method, uri string, reqBody, resBody interface{}) error {
	body := new(bytes.Buffer)
	if reqBody != nil {
		if err := json.NewEncoder(body).Encode(reqBody); err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseUrl+uri, body)
	if err != nil {
		return err
	}
	req.Header.Set("user-agent", "terraform-provider-rctf/"+version.Version)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+c.AuthToken)
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request %s %s: %w", method, uri, err)
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(resBody)
}
