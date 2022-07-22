package hubutil

import (
	"context"
	"fmt"
	"time"

	"github.com/dpb587/hubitat-cli/hub"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

func WaitForStart(ctx context.Context, log logr.Logger, hubClient *hub.Client) error {
	log.V(0).Info("waiting for start")

	var latestStatus hub.HubStatus
	var healthy bool

	for {
		time.Sleep(2 * time.Second)

		status, err := hubClient.HubStatus(ctx)
		if err != nil {
			if healthy {
				return errors.Wrap(err, "checking status")
			}

			// assume still initing
			continue
		}

		healthy = true

		if status.Status != latestStatus.Status || status.ServerInitPercentage != latestStatus.ServerInitPercentage || status.ServerInitDetails != latestStatus.ServerInitDetails {
			statusString := status.ServerInitPercentage

			if len(status.ServerInitDetails) > 0 {
				statusString = fmt.Sprintf("%s - %s", statusString, status.ServerInitDetails)
			}

			log.V(0).Info("waiting for start", "status", statusString)

			if status.Status == "running" {
				break
			}

			latestStatus = status
		}
	}

	return nil
}
