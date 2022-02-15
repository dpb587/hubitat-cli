package backupcmd

import (
	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/spf13/cobra"
)

func New(cmdp *cmdflags.Persistent) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "backup",
		Short: "For managing backups",
	}

	cmd.AddCommand(NewDownloadCommand(cmdp))

	return cmd
}
