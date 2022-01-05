package rctf

import (
	"context"
	"errors"
	"net/url"
)

type ChallengeFile struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ChallengePoints struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Challenge struct {
	Id               string          `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	Category         string          `json:"category"`
	Author           string          `json:"author"`
	Files            []ChallengeFile `json:"files"`
	Points           ChallengePoints `json:"points"`
	Flag             string          `json:"flag"`
	TiebreakEligible bool            `json:"tiebreakEligible"`
	SortWeight       int             `json:"sortWeight"`
}

type challengeResponse struct {
	response
	Data Challenge `json:"data"`
}

type challengeRequest struct {
	Data Challenge `json:"data"`
}

func challengeUrl(id string) string {
	return "admin/challs/" + url.QueryEscape(id)
}

func (c *Client) Challenge(ctx context.Context, id string) (Challenge, error) {
	res := challengeResponse{}
	if err := c.req(ctx, "GET", challengeUrl(id), nil, &res); err != nil {
		return Challenge{}, err
	}
	if res.Kind != "goodChallenges" {
		return Challenge{}, errors.New(res.String())
	}
	return res.Data, nil
}

func (c *Client) PutChallenge(ctx context.Context, challenge Challenge) error {
	res := response{}
	req := challengeRequest{
		Data: challenge,
	}
	if err := c.req(ctx, "PUT", challengeUrl(challenge.Id), &req, &res); err != nil {
		return err
	}
	if res.Kind != "goodChallengeUpdate" {
		return errors.New(res.String())
	}
	return nil
}

func (c *Client) DeleteChallenge(ctx context.Context, id string) error {
	res := response{}
	if err := c.req(ctx, "DELETE", challengeUrl(id), nil, &res); err != nil {
		return err
	}
	if res.Kind != "goodChallengeDelete" {
		return errors.New(res.String())
	}
	return nil
}
