package xocmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"xo/internal/types"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) ImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "import a xo project",
		Run: func(cmd *cobra.Command, args []string) {
			p, _ := os.Getwd()
			xoFilePath := path.Join(p, "xo.json")

			// check if xo file is there in this path
			if _, err := os.Stat(xoFilePath); os.IsNotExist(err) {
				panic("not a xo project. do 'xo i'")
			}

			// TODO :: take xo.json to a const string, because we are using this too much
			xoFileBytes, _ := os.ReadFile(xoFilePath)
			var xoFile types.XOFile
			if err := json.Unmarshal(xoFileBytes, &xoFile); err != nil {
				panic(err.Error())
			}

			// check if there are some commands inside xofile
			if len(xoFile.Commands) == 0 {
				panic("no commands inside this project")
			}

			// check if the project path exists
			if _, ok := xo.store.GetProjectPaths()[p]; ok {
				panic("project path already exists in xo workspace")
			}

			// check if project name is clashing with some other project in the workspace
			if _, ok := xo.store.GetCommands()[xoFile.Name]; ok {
				conflictProjectPath := xo.store.GetProjects()[xoFile.Name]
				err := fmt.Errorf("project name is clashing with another project at this location: %s", conflictProjectPath)
				panic(err)
			}

			// all good here!! lets add this project to workspace
			if err := xo.store.AddProject(p, xoFile); err != nil {
				panic(err)
			}

			fmt.Printf("added project: %s\n", xoFile.Name)
		},
	}
}
