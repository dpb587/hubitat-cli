package cloudcmd

import (
	"fmt"

	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCheckForUpdateCommand(cmdp *cmdflags.Persistent) *cobra.Command {
	const exitStatusUpdateAvailable = 65

	var fExitCode bool

	var cmd = &cobra.Command{
		Use:   "check-for-update",
		Short: "For checking the latest platform version",
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

			fmt.Fprintf(cmdp.Stdout, "status\t%s\n", result.Status)

			if len(result.Version) > 0 {
				fmt.Fprintf(cmdp.Stdout, "version\t%s\n", result.Version)
				fmt.Fprintf(cmdp.Stdout, "upgrade\t%v\n", result.Upgrade)
			}

			if result.Status == "UPDATE_AVAILABLE" {
				cmdp.Logger.V(0).Info("update available", "update_version", result.Version)

				if fExitCode {
					return cmdflags.ErrorCode{Code: exitStatusUpdateAvailable}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&fExitCode, "exit-code", "", false, fmt.Sprintf("use detailed exit codes (%d: update available)", exitStatusUpdateAvailable))

	return cmd
}
