package hubidcmd

import (
	"fmt"

	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New(cmdp *cmdflags.Persistent) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "hub-id",
		Short: "For getting the hub ID",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hubClient, err := cmdp.HubClient()
			if err != nil {
				return errors.Wrap(err, "getting hub client")
			}

			id, err := hubClient.HubID(ctx)
			if err != nil {
				return errors.Wrap(err, "fetching")
			}

			fmt.Fprintf(cmdp.Stdout, "%s\n", id)

			return nil
		},
	}

	return cmd
}
