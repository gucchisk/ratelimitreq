/*
Copyright © 2024 gucchisk <gucchi_sk@yahoo.co.jp>
*/
package client

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/time/rate"
)

type Client struct {
	limiters      map[string]*rate.Limiter
	globalLimiter *rate.Limiter
	httpClient    *http.Client
}

func NewClientByRateLimit(r rate.Limit) *Client {
	return &Client{
		globalLimiter: rate.NewLimiter(r, 1),
		httpClient:    new(http.Client),
	}
}

func NewClient(globalLimit rate.Limit, rates map[string]float64) *Client {
	limiters := make(map[string]*rate.Limiter)
	for d, r := range rates {
		limiters[d] = rate.NewLimiter(rate.Limit(r), 1)
	}
	return &Client{
		limiters:      limiters,
		globalLimiter: rate.NewLimiter(globalLimit, 1),
		httpClient:    new(http.Client),
	}
}

func (c *Client) getLimiter(domain string) *rate.Limiter {
	limiter, ok := c.limiters[domain]
	if ok {
		return limiter
	}
	return c.globalLimiter
}

func (c *Client) getLimiterByURLString(urlString string) *rate.Limiter {
	u, err := url.Parse(urlString)
	if err != nil {
		return c.globalLimiter
	}
	return c.getLimiter(u.Host)
}

func (c *Client) CloseIdleConnections() {
	c.httpClient.CloseIdleConnections()
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	err := c.getLimiter(req.URL.Host).Wait(context.Background())
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func (c *Client) Get(urlString string) (*http.Response, error) {
	err := c.getLimiterByURLString(urlString).Wait(context.Background())
	if err != nil {
		return nil, err
	}
	return c.httpClient.Get(urlString)
}

func (c *Client) Head(urlString string) (*http.Response, error) {
	err := c.getLimiterByURLString(urlString).Wait(context.Background())
	if err != nil {
		return nil, err
	}
	return c.httpClient.Head(urlString)
}

func (c *Client) Post(urlString, contentType string, body io.Reader) (*http.Response, error) {
	err := c.getLimiterByURLString(urlString).Wait(context.Background())
	if err != nil {
		return nil, err
	}
	return c.httpClient.Post(urlString, contentType, body)
}

func (c *Client) PostForm(urlString string, data url.Values) (*http.Response, error) {
	err := c.getLimiterByURLString(urlString).Wait(context.Background())
	if err != nil {
		return nil, err
	}
	return c.httpClient.PostForm(urlString, data)
}
