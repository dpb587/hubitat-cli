package hub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (c *Client) Reboot(ctx context.Context) error {
	c.log.V(2).Info("requesting reboot")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/hub/reboot", nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	res, err := c.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending request")
	} else if res.StatusCode != http.StatusOK {
		if cerr := res.Body.Close(); cerr != nil {
			c.log.Error(cerr, "closing errored response")
		}

		return fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	c.log.V(1).Info("requested reboot")

	return nil
}
