package xocmd

import (
	"fmt"
	"xo/internal/display"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) InfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "!",
		Short: "show info about the xo project",
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			pathMap := xo.store.GetProjects()
			if _, ok := pathMap[projectName]; !ok {
				fmt.Println("project:", projectName, "is not present please do 'xo i && xo import'!")
				return
			}

			commands := xo.store.GetCommands()
			display.ProjectInfo(projectName, commands[projectName])
		},
	}
}
