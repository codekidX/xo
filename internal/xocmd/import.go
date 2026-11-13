package xocmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"xo/internal/types"
	"xo/internal/utils"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) ImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "import a xo project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				utils.Fatal("To import from a URL or external source, use 'xo setup <url>'. 'xo import' is only for local directories.")
			}

			p, _ := os.Getwd()
			xoFilePath := filepath.Join(p, "xo.json")

			// check if xo file is there in this path
			if _, err := os.Stat(xoFilePath); os.IsNotExist(err) {
				utils.Fatal("Not a xo project. Run 'xo init' first.")
			}

			// TODO :: take xo.json to a const string, because we are using this too much
			xoFileBytes, _ := os.ReadFile(xoFilePath)
			var xoFile types.XOFile
			if err := json.Unmarshal(xoFileBytes, &xoFile); err != nil {
				utils.Fatal(err)
			}

			// check if there are some commands inside xofile
			if len(xoFile.Commands) == 0 {
				utils.Fatal("No commands inside this project.")
			}

			// check if the project path exists
			if _, ok := xo.store.GetProjectPaths()[p]; ok {
				utils.Fatal("Project path already exists in xo workspace.")
			}

			// check if project name is clashing with some other project in the workspace
			if _, ok := xo.store.GetCommands()[xoFile.Name]; ok {
				conflictProjectPath := xo.store.GetProjects()[xoFile.Name]
				err := fmt.Errorf("project name is clashing with another project at this location: %s", conflictProjectPath)
				utils.Fatal(err)
			}

			// all good here!! lets add this project to workspace
			if err := xo.store.AddProject(p, xoFile); err != nil {
				utils.Fatal(err)
			}

			fmt.Printf("added project: %s\n", xoFile.Name)
		},
	}
}
