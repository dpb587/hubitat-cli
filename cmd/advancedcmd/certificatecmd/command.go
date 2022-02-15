package certificatecmd

import (
	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/spf13/cobra"
)

func New(cmdp *cmdflags.Persistent) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "certificate",
		Short: "For managing hub certificate",
	}

	cmd.AddCommand(NewUpdateCommand(cmdp))

	return cmd
}
