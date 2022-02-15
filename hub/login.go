package hub

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func (c *Client) Login(ctx context.Context, username, password string) error {
	c.log.V(2).Info("logging in")

	if c.httpClient.Jar == nil {
		return errors.New("cannot login when cookie jar is missing")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL.ResolveReference(&url.URL{Path: "login"}).String(),
		strings.NewReader(url.Values{
			"username": []string{username},
			"password": []string{password},
			"submit":   []string{"Login"},
		}.Encode()),
	)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending request")
	}

	err = res.Body.Close()
	if err != nil {
		return errors.Wrap(err, "closing response body")
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	c.log.V(1).Info("logged in")

	return nil
}
