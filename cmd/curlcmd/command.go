package curlcmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New(cmdp *cmdflags.Persistent) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "curl PATH",
		Short: "For manual requests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hubClient, err := cmdp.HubClient()
			if err != nil {
				return errors.Wrap(err, "getting hub client")
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, args[0], nil)
			if err != nil {
				return errors.Wrap(err, "creating request")
			}

			res, err := hubClient.Do(req)
			if err != nil {
				return errors.Wrap(err, "sending request")
			} else if res.StatusCode != http.StatusOK {
				if cerr := res.Body.Close(); cerr != nil {
					cmdp.Logger.Error(cerr, "closing errored response")
				}

				return fmt.Errorf("unexpected response status code: %d", res.StatusCode)
			}

			defer res.Body.Close()

			_, err = io.Copy(cmdp.Stdout, res.Body)
			if err != nil {
				return errors.Wrap(err, "downloading")
			}

			return nil
		},
	}

	return cmd
}
