package hubutil

import (
	"context"
	"fmt"
	"time"

	"github.com/dpb587/hubitat-cli/hub"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

func WaitForUpdate(ctx context.Context, log logr.Logger, hubClient *hub.Client) error {
	log.V(0).Info("waiting for update")

	var latestStatus hub.CloudUpdateStatus

	for {
		time.Sleep(2 * time.Second)

		status, err := hubClient.CloudCheckUpdateStatus(ctx)
		if err != nil {
			if _, err := hubClient.HubID(ctx); err != nil {
				// assume rebooting
				break
			}

			return errors.Wrap(err, "checking update status")
		}

		if status.Status != latestStatus.Status || status.Percent != latestStatus.Percent {
			statusString := status.Status

			if status.Percent > 0 {
				statusString = fmt.Sprintf("%s (%d%%)", statusString, status.Percent)
			}

			log.V(0).Info("waiting for update", "status", statusString)

			if status.Status == "IDLE" {
				// missed expected error earlier
				break
			}

			latestStatus = status
		}
	}

	return nil
}
