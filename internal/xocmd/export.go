package xocmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"xo/internal/types"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) ExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export [filename.json]",
		Short: "Export imported projects to a shareable workspace file",
		Run: func(cmd *cobra.Command, args []string) {
			filename := "xo-workspace.json"
			if len(args) > 0 {
				filename = args[0]
			}

			projects := xo.store.GetProjects()
			gitURLs := xo.store.GetGitURLs()
			setupCmds := xo.store.GetSetupCommands()

			workspace := types.SharedWorkspace{
				Name:         "exported-workspace",
				Setup:        setupCmds,
				Repositories: []types.SharedRepository{},
			}

			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				return
			}

			for name, p := range projects {
				gitURL := gitURLs[name]
				
				// Try to make the path relative to the current directory
				relPath, err := filepath.Rel(cwd, p)
				if err != nil {
					// Fallback to absolute if Rel fails
					relPath = p
				}

				workspace.Repositories = append(workspace.Repositories, types.SharedRepository{
					Name: name,
					Git:  gitURL,
					Path: relPath,
				})
			}

			b, err := json.MarshalIndent(workspace, "", "  ")
			if err != nil {
				fmt.Println("Error encoding workspace:", err)
				return
			}

			if err := os.WriteFile(filename, b, 0644); err != nil {
				fmt.Println("Error writing file:", err)
				return
			}

			fmt.Printf("Successfully exported workspace to %s\n", filename)
		},
	}
}
