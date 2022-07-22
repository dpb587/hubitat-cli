package cloudcmd

import (
	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/dpb587/hubitat-cli/hub/hubutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewUpdatePlatformCommand(cmdp *cmdflags.Persistent) *cobra.Command {
	var fFollow bool

	var cmd = &cobra.Command{
		Use:   "update-platform",
		Short: "For installing the latest platform version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hubClient, err := cmdp.HubClient()
			if err != nil {
				return errors.Wrap(err, "getting hub client")
			}

			result, err := hubClient.CloudCheckForUpdate(ctx)
			if err != nil {
				return errors.Wrap(err, "checking")
			}

			if result.Status != "UPDATE_AVAILABLE" {
				cmdp.Logger.V(0).Info("no update available", "update_status", result.Status)

				return nil
			}

			cmdp.Logger.V(0).Info("update available", "update_version", result.Version)

			err = hubClient.CloudUpdatePlatform(ctx)
			if err != nil {
				return errors.Wrap(err, "updating")
			}

			cmdp.Logger.V(0).Info("requested update")

			if !fFollow {
				return nil
			}

			cmdp.Logger.V(0).Info("applying update")

			err = hubutil.WaitForUpdate(ctx, cmdp.Logger, hubClient)
			if err != nil {
				return errors.Wrap(err, "waiting for update")
			}

			err = hubutil.WaitForStart(ctx, cmdp.Logger, hubClient)
			if err != nil {
				return errors.Wrap(err, "waiting for start")
			}

			cmdp.Logger.V(0).Info("update complete")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&fFollow, "follow", "f", false, "follow progress")

	return cmd
}
