package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"xo/internal/types"
	"xo/internal/utils"
)

// Workspace: holds the aliases and commands of all projects on your computer
// Project: points to a specific path on your computer
// Commands: are tagged to a project and can be run according to the XOFile spec
// XOFile: holds the data of only a single project

type Store struct {
	workspaceFilePath string
}

func (s *Store) AddProject(projectPath string, xoFile types.XOFile) error {
	workspaceFile := s.getConfig()
	workspaceFile.ProjectPaths[projectPath] = xoFile.Name
	workspaceFile.Projects[xoFile.Name] = projectPath

	// set commands as map, so that we can access it quickly
	commands := make(map[string]types.Command, len(xoFile.Commands))
	for _, cmd := range xoFile.Commands {
		commands[cmd.Name] = cmd
	}
	workspaceFile.Commands[xoFile.Name] = commands

	if workspaceFile.GitURLs == nil {
		workspaceFile.GitURLs = make(map[string]string)
	}
	if xoFile.GitURL != "" {
		workspaceFile.GitURLs[xoFile.Name] = xoFile.GitURL
	}

	// workspace file as json bytes
	b, err := json.MarshalIndent(workspaceFile, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.workspaceFilePath, b, 0644)
}

func (s *Store) RemoveProject(p string) {
	config := s.getConfig()
	projectName := config.ProjectPaths[p]

	delete(config.Projects, projectName)
	delete(config.Commands, projectName)
	delete(config.ProjectPaths, p)
	if config.GitURLs != nil {
		delete(config.GitURLs, projectName)
	}

	s.saveWorkspaceFile(config)
}

func (s *Store) saveWorkspaceFile(c types.WorkspaceFile) {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		utils.Fatal(err)
	}

	err = os.WriteFile(s.workspaceFilePath, b, 0644)
	if err != nil {
		utils.Fatal(err)
	}
}

func (s *Store) saveXOFile(c types.XOFile, xoFilePath string) {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		utils.Fatal(err)
	}
	err = os.WriteFile(xoFilePath, b, 0644)
	if err != nil {
		utils.Fatal(err) // do we panic here really?
	}
}

func (s *Store) getConfig() types.WorkspaceFile {
	if _, err := os.Stat(s.workspaceFilePath); os.IsNotExist(err) {
		newConfig := types.WorkspaceFile{
			Projects:     make(map[string]string),
			ProjectPaths: make(map[string]string),
			Commands:     make(map[string]types.CommandMap),
			GitURLs:      make(map[string]string),
			Setup:        make(map[string][]string),
		}
		s.saveWorkspaceFile(newConfig)
		return newConfig
	}

	b, err := os.ReadFile(s.workspaceFilePath)
	if err != nil {
		utils.Fatal(err)
	}

	var config types.WorkspaceFile
	err = json.Unmarshal(b, &config)
	if err != nil {
		utils.Fatal(err)
	}
	return config
}

func (s *Store) GetProjects() map[string]string {
	config := s.getConfig()
	return config.Projects
}

func (s *Store) GetProjectPaths() map[string]string {
	config := s.getConfig()
	return config.ProjectPaths
}

func (s *Store) GetCommands() map[string]types.CommandMap {
	config := s.getConfig()
	return config.Commands
}

func (s *Store) GetGitURLs() map[string]string {
	config := s.getConfig()
	if config.GitURLs == nil {
		return make(map[string]string)
	}
	return config.GitURLs
}

func (s *Store) GetSetupCommands() map[string][]string {
	config := s.getConfig()
	if config.Setup == nil {
		return make(map[string][]string)
	}
	return config.Setup
}

func (s *Store) AddSetupCommand(osName, command string) error {
	config := s.getConfig()
	if config.Setup == nil {
		config.Setup = make(map[string][]string)
	}
	config.Setup[osName] = append(config.Setup[osName], command)
	s.saveWorkspaceFile(config)
	return nil
}

func (s *Store) RemoveSetupCommand(osName, command string) error {
	config := s.getConfig()
	if config.Setup == nil {
		return nil
	}
	cmds := config.Setup[osName]
	for i, c := range cmds {
		if c == command {
			config.Setup[osName] = append(cmds[:i], cmds[i+1:]...)
			break
		}
	}
	if len(config.Setup[osName]) == 0 {
		delete(config.Setup, osName)
	}
	s.saveWorkspaceFile(config)
	return nil
}

func (s *Store) CreateXOFile(p string) error {
	xoFilePath := filepath.Join(p, "xo.json")
	if f, _ := os.Stat(xoFilePath); f != nil {
		return errors.New("xo is already initialized here")
	}

	s.saveXOFile(types.XOFile{
		Name:     filepath.Base(p),
		Commands: []types.Command{},
	}, xoFilePath)

	return nil
}

func New() *Store {
	home, _ := os.UserHomeDir()
	xoFilePath := filepath.Join(home, ".xorc")
	return &Store{
		workspaceFilePath: xoFilePath,
	}
}

func (s *Store) AddCommandToProject(projectName string, cmd types.Command) error {
	// 1. Get project path
	projPath, ok := s.GetProjects()[projectName]
	if !ok {
		return fmt.Errorf("project not found: %s", projectName)
	}

	// 2. Read existing xo.json
	xoFilePath := filepath.Join(projPath, "xo.json")
	b, err := os.ReadFile(xoFilePath)
	if err != nil {
		return err
	}

	var xoFile types.XOFile
	if err := json.Unmarshal(b, &xoFile); err != nil {
		return err
	}

	// 3. Append command and save xo.json
	xoFile.Commands = append(xoFile.Commands, cmd)
	s.saveXOFile(xoFile, xoFilePath)

	// 4. Re-sync workspace (~/.xorc)
	return s.AddProject(projPath, xoFile)
}

func (s *Store) Sync() {
	config := s.getConfig()
	changed := false

	for name, p := range config.Projects {
		xoFilePath := filepath.Join(p, "xo.json")
		b, err := os.ReadFile(xoFilePath)
		if err != nil {
			// Skip missing or unreadable files instead of deleting to be safe
			continue
		}

		var xoFile types.XOFile
		if err := json.Unmarshal(b, &xoFile); err == nil {
			commands := make(map[string]types.Command, len(xoFile.Commands))
			for _, cmd := range xoFile.Commands {
				commands[cmd.Name] = cmd
			}

			config.Commands[name] = commands
			if xoFile.GitURL != "" {
				if config.GitURLs == nil {
					config.GitURLs = make(map[string]string)
				}
				config.GitURLs[name] = xoFile.GitURL
			}
			changed = true
		}
	}

	if changed {
		s.saveWorkspaceFile(config)
	}
}
