package advancedcmd

import (
	"github.com/dpb587/hubitat-cli/cmd/advancedcmd/certificatecmd"
	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/spf13/cobra"
)

func New(cmdp *cmdflags.Persistent) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "advanced",
		Short: "For advanced features",
	}

	cmd.AddCommand(certificatecmd.New(cmdp))

	return cmd
}
