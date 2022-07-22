package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
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

func (c *Client) HubID(ctx context.Context) (string, error) {
	c.log.V(2).Info("requesting id")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/hubId", nil)
	if err != nil {
		return "", errors.Wrap(err, "creating request")
	}

	res, err := c.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "sending request")
	} else if res.StatusCode != http.StatusOK {
		if cerr := res.Body.Close(); cerr != nil {
			c.log.Error(cerr, "closing errored response")
		}

		return "", fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "reading body")
	}

	id := strings.TrimSpace(string(buf))
	if v := len(id); v != 36 {
		return "", fmt.Errorf("unexpected response body (length %d)", v)
	}

	c.log.V(1).Info("requested id")

	return id, nil
}

type HubStatus struct {
	Status               string `json:"status"`
	ServerInitPercentage string `json:"serverInitPercentage"`
	ServerInitDetails    string `json:"serverInitDetails"`
}

func (c *Client) HubStatus(ctx context.Context) (HubStatus, error) {
	c.log.V(2).Info("requesting status")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/hubStatus", nil)
	if err != nil {
		return HubStatus{}, errors.Wrap(err, "creating request")
	}

	res, err := c.Do(req)
	if err != nil {
		return HubStatus{}, errors.Wrap(err, "sending request")
	} else if res.StatusCode != http.StatusOK {
		if cerr := res.Body.Close(); cerr != nil {
			c.log.Error(cerr, "closing errored response")
		}

		return HubStatus{}, fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	var meta HubStatus

	err = json.NewDecoder(res.Body).Decode(&meta)
	if err != nil {
		return HubStatus{}, errors.Wrap(err, "decoding response")
	}

	c.log.V(1).Info("requested status")

	return meta, nil
}
