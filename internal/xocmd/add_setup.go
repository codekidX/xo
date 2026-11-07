package xocmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (xo *XOCmd) AddSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-setup [os] [command]",
		Short: "Add a global setup command (e.g., darwin 'brew install go')",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			xo.store.AddSetupCommand(args[0], args[1])
			fmt.Printf("Added setup command for %s: %s\n", args[0], args[1])
		},
	}
}
