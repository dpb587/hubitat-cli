package backupcmd

import (
	"fmt"
	"io"
	"os"

	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewDownloadCommand(cmdp *cmdflags.Persistent) *cobra.Command {
	var fOutput string

	var cmd = &cobra.Command{
		Use:   "download [BACKUP-FILE-NAME]",
		Short: "For downloading backup files",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hubClient, err := cmdp.HubClient()
			if err != nil {
				return errors.Wrap(err, "getting hub client")
			}

			var remoteFile = "latest"

			if len(args) > 0 {
				remoteFile = args[0]
			}

			meta, r, err := hubClient.DownloadBackupFile(ctx, remoteFile)
			if err != nil {
				return errors.Wrap(err, "fetching download")
			}

			defer r.Close()

			var localFile = fOutput

			if len(localFile) == 0 {
				localFile = meta.Name
			}

			cmdp.Logger.V(2).Info("writing backup file", "name", localFile)

			fh, err := os.OpenFile(localFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				return errors.Wrap(err, "opening local file")
			}

			defer fh.Close()

			_, err = io.Copy(fh, r)
			if err != nil {
				return errors.Wrap(err, "downloading")
			}

			cmdp.Logger.V(1).Info("wrote backup file", "name", localFile)

			cmdp.Logger.V(0).Info(fmt.Sprintf("downloaded backup file to %s", localFile))

			return nil
		},
	}

	cmd.Flags().StringVarP(&fOutput, "output", "o", "", "output file path")

	return cmd
}
