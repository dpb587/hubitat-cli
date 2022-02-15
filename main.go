package main

import (
	"fmt"
	"os"

	"github.com/dpb587/hubitat-cli/cmd/advancedcmd"
	"github.com/dpb587/hubitat-cli/cmd/backupcmd"
	"github.com/dpb587/hubitat-cli/cmd/cmdflags"
	"github.com/dpb587/hubitat-cli/cmd/rebootcmd"
	"github.com/spf13/cobra"
)

func main() {
	var cmd = &cobra.Command{
		Use:           "hubitat-cli",
		Short:         "For interacting with Hubitat",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmdp := cmdflags.NewPersistent(cmd)

	cmd.AddCommand(advancedcmd.New(cmdp))
	cmd.AddCommand(backupcmd.New(cmdp))
	// cmd.AddCommand(curlcmd.New(cmdp))
	cmd.AddCommand(rebootcmd.New(cmdp))

	if err := cmd.Execute(); err != nil {
		fmt.Printf("%s: error: %s\n", cmd.Use, err)
		os.Exit(1)
	}
}
