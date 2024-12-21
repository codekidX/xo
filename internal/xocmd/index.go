package xocmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"xo/internal/store"

	"github.com/spf13/cobra"
)

type XOCmd struct {
	store *store.Store
}

func Run() error {
	root := &cobra.Command{
		Use:   "xo",
		Short: "xo is aliasing on steroids",
	}

	xoapp := &XOCmd{store: store.New()}

	root.AddCommand(xoapp.InitCmd())
	root.AddCommand(xoapp.RemoveCmd())
	root.AddCommand(xoapp.InfoCmd())
	root.AddCommand(xoapp.AtCmd())
	root.AddCommand(xoapp.ImportCmd())
	root.AddCommand(xoapp.UpdateCmd())

	isProjectCommand := checkAndRunProjectCommand(root, xoapp.store)

	if isProjectCommand {
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
		if len(os.Args) < 3 {
			fmt.Println("not enough argument to run xo command.")
			return true
		}

		projectName := os.Args[1]
		commandName := os.Args[2]
		commandMap := s.GetCommands()

		// check if project exists in workspace
		if _, ok := commandMap[projectName]; !ok {
			fmt.Printf("no such project: %s\n", projectName)
			return true
		}

		// check if command exists in project
		if _, ok := commandMap[projectName][commandName]; !ok {
			fmt.Printf("no such command in this project: %s\n", commandName)
			return true
		}

		// all good here! run the command!
		projectPath := s.GetProjects()[projectName]
		runCmd(commandMap[projectName][commandName].CmdStr, projectPath)

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

	// if it contains multiple commands -- split it!
	splittedCmds := strings.Split(cmd, "&&")
	for _, commandStr := range splittedCmds {
		cmdArr := strings.Split(strings.TrimSpace(commandStr), " ")
		bin := cmdArr[0]

		cmdHandle := exec.Command(bin, cmdArr[1:]...)
		stdout, _ := cmdHandle.Output()

		// TODO :: we need to know when error happens and paint it RED!

		fmt.Println(string(stdout))
	}
}
