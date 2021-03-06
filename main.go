package main

import (
	"fmt"
	"os"

	"github.com/dpb587/hubitat-cli/cmd/advancedcmd"
	"github.com/dpb587/hubitat-cli/cmd/backupcmd"
	"github.com/dpb587/hubitat-cli/cmd/cloudcmd"
	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/dpb587/hubitat-cli/cmd/hubidcmd"
	"github.com/dpb587/hubitat-cli/cmd/rebootcmd"
	"github.com/spf13/cobra"
)

func main() {
	var cmd = &cobra.Command{
		Use:           "hubitat-cli",
		Short:         "For interacting with Hubitat",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       fmt.Sprintf("%s (commit %s, built %s)", cmdflags.VersionName, cmdflags.VersionCommit, cmdflags.VersionBuilt),
	}

	// simplify output; --help still exists, and unlikely to need completion
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: true,
	}

	cmdp := cmdflags.NewPersistent(cmd)

	cmd.AddCommand(advancedcmd.New(cmdp))
	cmd.AddCommand(backupcmd.New(cmdp))
	cmd.AddCommand(cloudcmd.New(cmdp))
	cmd.AddCommand(hubidcmd.New(cmdp))
	// cmd.AddCommand(curlcmd.New(cmdp))
	cmd.AddCommand(rebootcmd.New(cmdp))

	if err := cmd.Execute(); err != nil {
		var code = 1

		if ee, ok := err.(cmdflags.ErrorCode); ok {
			err = ee.Err
			code = ee.Code
		}

		if err != nil {
			fmt.Printf("%s: error: %s\n", cmd.Use, err)
		}

		os.Exit(code)
	}
}
