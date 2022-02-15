package hub

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func (c *Client) UpdateAdvancedCertificates(ctx context.Context, certificate, privateKey []byte) error {
	c.log.V(2).Info("requesting certificate update")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/hub/advanced/certificate/save",
		strings.NewReader(url.Values{
			"certificate":    []string{string(certificate)},
			"privateKey":     []string{string(privateKey)},
			"_action_update": []string{"Save Certificate and Key"},
		}.Encode()),
	)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending request")
	}

	err = res.Body.Close()
	if err != nil {
		return errors.Wrap(err, "closing response body")
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	c.log.V(2).Info("requested certificate update")

	return nil
}
