package hub

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func (c *Client) Login(username, password string) error {
	c.log.V(2).Info("logging into hub")

	if c.httpClient.Jar == nil {
		return errors.New("cannot login when cookie jar is missing")
	}

	res, err := c.httpClient.PostForm(
		c.baseURL.ResolveReference(&url.URL{Path: "login"}).String(),
		url.Values{
			"username": []string{username},
			"password": []string{password},
			"submit":   []string{"Login"},
		},
	)
	if err != nil {
		return errors.Wrap(err, "sending request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	c.log.V(1).Info("logged into hub")

	return nil
}
