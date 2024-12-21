package xocmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"xo/internal/types"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) UpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update the commands of a xo project",
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			commandMap := xo.store.GetCommands()

			// check if project exists in workspace
			if _, ok := commandMap[projectName]; !ok {
				fmt.Println("no project names: ", projectName, "in xo workspace")
				return
			}

			projectPath := xo.store.GetProjects()[projectName]
			xoFilePath := path.Join(projectPath, "xo.json")

			// check if xo file is there in this path
			if _, err := os.Stat(xoFilePath); os.IsNotExist(err) {
				panic("not a xo project. do 'xo i'")
			}

			xoFileBytes, _ := os.ReadFile(xoFilePath)
			var xoFile types.XOFile
			if err := json.Unmarshal(xoFileBytes, &xoFile); err != nil {
				panic(err.Error())
			}

			if err := xo.store.AddProject(projectPath, xoFile); err != nil {
				panic(err)
			}
		},
	}
}
