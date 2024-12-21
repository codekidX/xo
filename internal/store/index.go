package store

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"xo/internal/types"
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

	s.saveWorkspaceFile(config)
}

func (s *Store) saveWorkspaceFile(c types.WorkspaceFile) {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(s.workspaceFilePath, b, 0644)
	if err != nil {
		panic(err)
	}
}

func (s *Store) saveXOFile(c types.XOFile, xoFilePath string) {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(xoFilePath, b, 0644)
	if err != nil {
		panic(err) // do we panic here really?
	}
}

func (s *Store) getConfig() types.WorkspaceFile {
	if _, err := os.Stat(s.workspaceFilePath); os.IsNotExist(err) {
		newConfig := types.WorkspaceFile{
			Projects:     make(map[string]string),
			ProjectPaths: make(map[string]string),
			Commands:     make(map[string]types.CommandMap),
		}
		s.saveWorkspaceFile(newConfig)
		return newConfig
	}

	b, err := os.ReadFile(s.workspaceFilePath)
	if err != nil {
		panic(err)
	}

	var config types.WorkspaceFile
	err = json.Unmarshal(b, &config)
	if err != nil {
		panic(err)
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

func (s *Store) CreateXOFile(p string) error {
	xoFilePath := path.Join(p, "xo.json")
	if f, _ := os.Stat(xoFilePath); f != nil {
		return errors.New("xo is already initialized here")
	}

	s.saveXOFile(types.XOFile{
		Name:     path.Base(p),
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
