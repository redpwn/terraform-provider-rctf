package rctf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/redpwn/terraform-provider-rctf/internal/version"
)

type Client struct {
	BaseUrl string
	Header  http.Header
	http    *http.Client
}

func New(baseUrl, authToken string) (*Client, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	h := http.Header{}
	h.Set("user-agent", "terraform-provider-rctf/"+version.Version)
	h.Set("content-type", "application/json")
	h.Set("authorization", "Bearer "+authToken)
	c := &Client{
		BaseUrl: fmt.Sprintf("%s://%s/api/v1/", u.Scheme, u.Host),
		Header:  h,
		http:    http.DefaultClient,
	}
	return c, nil
}

type response struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

type Error response

func (e Error) Error() string {
	return fmt.Sprintf("rctf error: %s: %s", e.Kind, e.Message)
}

func (r response) err() error {
	return Error(r)
}

func (c *Client) req(ctx context.Context, method, uri string, reqBody, resBody interface{}) error {
	body := &bytes.Buffer{}
	if reqBody != nil {
		if err := json.NewEncoder(body).Encode(reqBody); err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseUrl+uri, body)
	if err != nil {
		return err
	}
	req.Header = c.Header
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request %s %s: %w", method, uri, err)
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(resBody)
}
