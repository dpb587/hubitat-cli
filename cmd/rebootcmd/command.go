package rebootcmd

import (
	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New(cmdp *cmdflags.Persistent) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "reboot",
		Short: "For restarting the hub",
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

			return nil
		},
	}

	return cmd
}
