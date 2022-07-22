package cloudcmd

import (
	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/spf13/cobra"
)

func New(cmdp *cmdflags.Persistent) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cloud",
		Short: "For managing cloud capabilities",
	}

	cmd.AddCommand(NewCheckForUpdateCommand(cmdp))
	cmd.AddCommand(NewUpdatePlatformCommand(cmdp))

	return cmd
}
