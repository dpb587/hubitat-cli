package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type CloudUpdate struct {
	ReleaseNotesContent string

	// UPDATE_AVAILABLE
	Status  string
	Upgrade bool
	Version string
}

func (c *Client) CloudCheckForUpdate(ctx context.Context) (CloudUpdate, error) {
	c.log.V(2).Info("requesting status")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/hub/cloud/checkForUpdate", nil)
	if err != nil {
		return CloudUpdate{}, errors.Wrap(err, "creating request")
	}

	req.Header.Set("Accept", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return CloudUpdate{}, errors.Wrap(err, "sending request")
	} else if res.StatusCode != http.StatusOK {
		if cerr := res.Body.Close(); cerr != nil {
			c.log.Error(cerr, "closing errored response")
		}

		return CloudUpdate{}, fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	var meta CloudUpdate

	err = json.NewDecoder(res.Body).Decode(&meta)
	if err != nil {
		return CloudUpdate{}, errors.Wrap(err, "decoding response")
	}

	c.log.V(1).Info("requested status")

	return meta, nil
}

type CloudUpdateStatus struct {
	// IDLE, DOWNLOAD_IN_PROGRESS, DOWNLOAD_VERIFY, EXTRACT_IN_PROGRESS
	Status  string `json:"status"`
	Percent int    `json:"percent"`
}

func (c *Client) CloudCheckUpdateStatus(ctx context.Context) (CloudUpdateStatus, error) {
	c.log.V(2).Info("requesting update")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/hub/cloud/checkUpdateStatus", nil)
	if err != nil {
		return CloudUpdateStatus{}, errors.Wrap(err, "creating request")
	}

	req.Header.Set("Accept", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return CloudUpdateStatus{}, errors.Wrap(err, "sending request")
	} else if res.StatusCode != http.StatusOK {
		if cerr := res.Body.Close(); cerr != nil {
			c.log.Error(cerr, "closing errored response")
		}

		return CloudUpdateStatus{}, fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	var meta CloudUpdateStatus

	err = json.NewDecoder(res.Body).Decode(&meta)
	if err != nil {
		return CloudUpdateStatus{}, errors.Wrap(err, "decoding response")
	}

	c.log.V(1).Info("requested update")

	return meta, nil
}

func (c *Client) CloudUpdatePlatform(ctx context.Context) error {
	c.log.V(2).Info("requesting update")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/hub/cloud/updatePlatform", nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	req.Header.Set("Accept", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending request")
	} else if res.StatusCode != http.StatusOK {
		if cerr := res.Body.Close(); cerr != nil {
			c.log.Error(cerr, "closing errored response")
		}

		return fmt.Errorf("unexpected response status code: %d", res.StatusCode)
	}

	var meta struct {
		Success string `json:"success"`
	}

	err = json.NewDecoder(res.Body).Decode(&meta)
	if err != nil {
		return errors.Wrap(err, "decoding response")
	}

	if meta.Success != "true" {
		return errors.New("update request failed")
	}

	c.log.V(1).Info("requested update")

	return nil
}
