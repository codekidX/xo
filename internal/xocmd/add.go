package xocmd

import (
	"fmt"
	"xo/internal/types"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) AddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <project> <name> <command> [target_dir] [help]",
		Short: "add a new command to a project",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			
			newCmd := types.Command{
				Name:   args[1],
				CmdStr: args[2],
			}
			if len(args) > 3 {
				newCmd.WorkDir = args[3]
			}
			if len(args) > 4 {
				newCmd.Help = args[4]
			}

			if err := xo.store.AddCommandToProject(projectName, newCmd); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Printf("Added command '%s' to project '%s'\n", newCmd.Name, projectName)
		},
	}
}
