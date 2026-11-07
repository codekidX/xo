package xocmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"xo/internal/store"
	"xo/internal/sys"
	"xo/internal/tui"

	"github.com/spf13/cobra"
)

type XOCmd struct {
	store *store.Store
}

func Run() error {
	xoapp := &XOCmd{store: store.New()}

	// Auto-sync the workspace with on-disk project files
	xoapp.store.Sync()

	root := &cobra.Command{
		Use:   "xo",
		Short: "xo is aliasing on steroids",
		Run: func(cmd *cobra.Command, args []string) {
			if err := tui.Launch(xoapp.store); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
			}
		},
	}

	root.AddCommand(xoapp.InitCmd())
	root.AddCommand(xoapp.RemoveCmd())
	root.AddCommand(xoapp.InfoCmd())
	root.AddCommand(xoapp.AtCmd())
	root.AddCommand(xoapp.ImportCmd())
	root.AddCommand(xoapp.AddCmd())
	root.AddCommand(xoapp.SetupCmd())
	root.AddCommand(xoapp.ExportCmd())
	root.AddCommand(xoapp.AddSetupCmd())

	isProjectCommand := checkAndRunProjectCommand(root, xoapp.store)

	if isProjectCommand {
		// skip cobra execution
		return nil
	}

	root.CompletionOptions.DisableDefaultCmd = true
	return root.Execute()
}

// checkAndRunProjectCommand checks if the command provided is one of the project
// name. This makes using this app a lot easier because now you can do this:
//
// xo boomerang run
//
// this will run the boomerang project.
// It returns `true` if it encountered command which is one of xo project.
func checkAndRunProjectCommand(root *cobra.Command, s *store.Store) bool {
	projectNameOrCommand, _, err := root.Find(os.Args[1:])

	if err != nil && projectNameOrCommand.Use == root.Use {
		// a command must have 1=project_name, 2=command_name, we can ignore 0 because it will be the
		// bin_name=xo because cobra is not handling these args for us!
		if len(os.Args) < 2 {
			fmt.Println("not enough argument to run xo command.")
			return true
		}

		projectName := os.Args[1]

		// check if project exists in workspace
		commandMap := s.GetCommands()
		if _, ok := commandMap[projectName]; !ok {
			fmt.Printf("no such project: %s\n", projectName)
			return true
		}

		projectPath := s.GetProjects()[projectName]

		if len(os.Args) == 2 {
			// this means that only project name is mentioned
			shell := sys.DefaultShell()
			c := exec.Command(shell)
			c.Dir = projectPath
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "Shell exited with error:", err)
			}
			return true
		}

		commandName := os.Args[2]

		// check if command exists in project
		if _, ok := commandMap[projectName][commandName]; !ok {
			fmt.Printf("no such command in this project: %s\n", commandName)
			return true
		}

		// all good here! run the command!
		targetCommand := commandMap[projectName][commandName]
		if targetCommand.WorkDir != "" {
			runCmd(targetCommand.CmdStr, filepath.Join(projectPath, targetCommand.WorkDir))
		} else {
			runCmd(targetCommand.CmdStr, projectPath)
		}

		// return without running root command.
		return true
	}

	return false
}

func runCmd(cmd string, inDirectory string) {
	if err := os.Chdir(inDirectory); err != nil {
		fmt.Println("error: ", err.Error())
		return
	}

	cmdHandle := sys.ShellCommand(cmd)
	cmdHandle.Stdin = os.Stdin
	cmdHandle.Stdout = os.Stdout
	cmdHandle.Stderr = os.Stderr

	if err := cmdHandle.Run(); err != nil {
		fmt.Printf("Command failed: %v\n", err)
	}
}
