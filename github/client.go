package github

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	HeaderXRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderXRateLimitReset     = "X-RateLimit-Reset"
)

// Token -
type Token struct {
	Tag   string
	Token string
}

// Client -
type Client struct {
	Remaining int
	Reset     int64
	Tag       string
	*github.Client
}

// NewClient -
func NewClient(token *Token) (*Client, error) {
	if token == nil {
		return nil, errors.New("token cannot be nil")
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.Token},
	)

	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		Client: github.NewClient(tc),
		Tag:    token.Tag,
	}, nil
}

// HandleResponse -
func (c *Client) HandleResponse(resp *github.Response) error {
	remaining, err := strconv.Atoi(resp.Header.Get(HeaderXRateLimitRemaining))
	if err != nil {
		return err
	}
	c.Remaining = remaining
	reset, err := strconv.ParseInt(resp.Header.Get(HeaderXRateLimitReset), 10, 64)
	if err != nil {
		return err
	}
	c.Reset = reset
	return nil
}
