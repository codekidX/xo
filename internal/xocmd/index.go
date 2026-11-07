package xocmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"xo/internal/store"
	"xo/internal/sys"

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
	}

	root.AddCommand(xoapp.InitCmd())
	root.AddCommand(xoapp.RemoveCmd())
	root.AddCommand(xoapp.InfoCmd())
	root.AddCommand(xoapp.AtCmd())
	root.AddCommand(xoapp.ImportCmd())

	isProjectCommand := checkAndRunProjectCommand(root, xoapp.store)

	if isProjectCommand {
		// skip cobra execution
		return nil
	}

	return root.Execute()
}

func checkAndRunProjectCommand(root *cobra.Command, s *store.Store) bool {
	commandMap := s.GetCommands()

	// check if project exists
	if len(os.Args) > 1 {
		projectName := os.Args[1]

		if _, ok := commandMap[projectName]; !ok {
			return false
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

		if _, ok := commandMap[projectName][commandName]; !ok {
			fmt.Println("no command named:", commandName, "in project:", projectName)
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
	if inDirectory != "" {
		os.Chdir(inDirectory)
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
