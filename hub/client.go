package hub

import (
	"net/http"
	"net/url"

	"github.com/go-logr/logr"
)

type Client struct {
	log        logr.Logger
	httpClient *http.Client
	baseURL    *url.URL
}

func NewClient(log logr.Logger, httpClient *http.Client, baseURL *url.URL) *Client {
	return &Client{
		log:        log,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *Client) BaseURL() *url.URL {
	return c.baseURL
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if !req.URL.IsAbs() {
		req.URL = c.baseURL.ResolveReference(req.URL)
	}

	return c.httpClient.Do(req)
}
