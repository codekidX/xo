package xocmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"xo/internal/sys"
	"xo/internal/types"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) SetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup [workspace.json | https://...]",
		Short: "Bootstrap or repair your workspace",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				xo.repairGlobalWorkspace()
			} else {
				xo.bootstrapSharedWorkspace(args[0])
			}
		},
	}
}

func (xo *XOCmd) repairGlobalWorkspace() {
	projects := xo.store.GetProjects()
	gitURLs := xo.store.GetGitURLs()
	repaired := 0

	for name, p := range projects {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			url, ok := gitURLs[name]
			if !ok || url == "" {
				fmt.Printf("Cannot repair '%s': No Git URL stored.\n", name)
				continue
			}
			fmt.Printf("Repairing '%s'...\n", name)
			cloneCmd := exec.Command("git", "clone", url, p)
			cloneCmd.Stdout = os.Stdout
			cloneCmd.Stderr = os.Stderr
			if err := cloneCmd.Run(); err != nil {
				fmt.Printf("Failed to clone '%s': %v\n", name, err)
			} else {
				repaired++
				fmt.Printf("Successfully restored '%s'.\n", name)
			}
		}
	}

	if repaired == 0 {
		fmt.Println("Workspace is intact. Nothing to repair.")
	} else {
		fmt.Printf("Repaired %d projects.\n", repaired)
	}
}

func (xo *XOCmd) bootstrapSharedWorkspace(arg string) {
	// TRUST PROMPT
	fmt.Printf("\nDo you trust this import path (%s)?\n", arg)
	fmt.Print("This workspace may execute system commands. There is no automated bypass for this step. (y/N): ")
	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Aborted workspace setup.")
		return
	}

	var b []byte
	var err error

	if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
		fmt.Printf("\nFetching workspace from %s...\n", arg)
		resp, reqErr := http.Get(arg)
		if reqErr != nil {
			fmt.Println("Failed to fetch workspace:", reqErr)
			return
		}
		defer resp.Body.Close()
		b, err = io.ReadAll(resp.Body)
	} else {
		b, err = os.ReadFile(arg)
	}

	if err != nil {
		fmt.Println("Failed to read workspace:", err)
		return
	}

	var workspace types.SharedWorkspace
	if err := json.Unmarshal(b, &workspace); err != nil {
		fmt.Println("Failed to parse workspace JSON:", err)
		return
	}

	fmt.Printf("\nBootstrapping workspace '%s'...\n", workspace.Name)

	// RUN SETUP
	if setupCmds, ok := workspace.Setup[runtime.GOOS]; ok && len(setupCmds) > 0 {
		fmt.Printf("Running platform setup for %s...\n", runtime.GOOS)
		for _, cmdStr := range setupCmds {
			fmt.Printf(" > %s\n", cmdStr)
			cmd := sys.ShellCommand(cmdStr)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("Setup command failed: %v\n", err)
				return // Stop bootstrapping if a setup step fails
			}
		}
	}

	cwd, _ := os.Getwd()

	// CLONE & IMPORT
	for _, repo := range workspace.Repositories {
		targetPath := filepath.Join(cwd, repo.Path)

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			fmt.Printf("\nCloning '%s' into %s...\n", repo.Name, repo.Path)
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				fmt.Println("Failed to create parent directory:", err)
				continue
			}
			cloneCmd := exec.Command("git", "clone", repo.Git, targetPath)
			cloneCmd.Stdout = os.Stdout
			cloneCmd.Stderr = os.Stderr
			if err := cloneCmd.Run(); err != nil {
				fmt.Printf("Failed to clone '%s': %v\n", repo.Name, err)
				continue
			}
		} else {
			fmt.Printf("\nDirectory %s already exists. Skipping clone.\n", repo.Path)
		}

		xoFilePath := filepath.Join(targetPath, "xo.json")
		if _, err := os.Stat(xoFilePath); os.IsNotExist(err) {
			if err := xo.store.CreateXOFile(targetPath); err != nil {
				fmt.Printf("Failed to init xo.json in %s: %v\n", repo.Path, err)
				continue
			}
		}

		xb, _ := os.ReadFile(xoFilePath)
		var xoFile types.XOFile
		json.Unmarshal(xb, &xoFile)

		// If the repo didn't have a git url saved in its xo.json, inject it from the shared workspace
		if xoFile.GitURL == "" && repo.Git != "" {
			xoFile.GitURL = repo.Git
		}

		if err := xo.store.AddProject(targetPath, xoFile); err != nil {
			fmt.Printf("Failed to import '%s': %v\n", repo.Name, err)
		} else {
			fmt.Printf("Imported '%s'.\n", repo.Name)
		}
	}

	fmt.Println("\nWorkspace successfully bootstrapped!")
}
