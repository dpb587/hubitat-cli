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
	} else if res.StatusCode != http.StatusOK {
		if cerr := res.Body.Close(); cerr != nil {
			c.log.Error(cerr, "closing errored response")
		}

		return fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	return nil
}
