package rebootcmd

import (
	"time"

	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/dpb587/hubitat-cli/hub/hubutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New(cmdp *cmdflags.Persistent) *cobra.Command {
	var fFollow bool

	var cmd = &cobra.Command{
		Use:   "reboot",
		Short: "For rebooting the hub",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hubClient, err := cmdp.HubClient()
			if err != nil {
				return errors.Wrap(err, "getting hub client")
			}

			err = hubClient.Reboot(ctx)
			if err != nil {
				return errors.Wrap(err, "rebooting")
			}

			cmdp.Logger.V(0).Info("requested reboot")

			if !fFollow {
				return nil
			}

			cmdp.Logger.V(0).Info("waiting for reboot")

			for {
				time.Sleep(2 * time.Second)

				if _, err := hubClient.HubID(ctx); err != nil {
					// assume rebooting
					break
				}
			}

			err = hubutil.WaitForStart(ctx, cmdp.Logger, hubClient)
			if err != nil {
				return errors.Wrap(err, "waiting for start")
			}

			cmdp.Logger.V(0).Info("reboot complete")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&fFollow, "follow", "f", false, "follow progress")

	return cmd
}
