package xocmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) AtCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "@",
		Short: "provides the project path to a command",
		Args:  cobra.MatchAll(cobra.ExactArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			binary := args[0]
			projectName := args[1]
			pathMap := xo.store.GetProjects()

			// check if project exists in workspace
			if projectPath, ok := pathMap[projectName]; ok {
				runCmd(binary+" .", projectPath)
				return
			}

			fmt.Printf("no project named: %s\n", projectName)
		},
	}
}
