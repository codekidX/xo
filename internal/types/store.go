package types

type Command struct {
	Name    string   `json:"name"`
	CmdStr  string   `json:"cmd"`
	EnvVars []string `json:"env"`
	Help    string   `json:"help"`
}

type CommandMap map[string]Command

// WorkspaceFile resides in `~/.xorc`. This file is responsible running user
// commands in the specific project directories
type WorkspaceFile struct {
	// Projects has key=project_name, value=project_path
	Projects map[string]string `json:"projects"`

	// ProjectPaths has key=project_path, value=project_name
	// this is an inverse map of projects, this is to quickly check if
	// a path already exists in thr workspace during `import` command
	ProjectPaths map[string]string `json:"paths"`

	// Commands has key=project_name, value=map[command_name]Command
	Commands map[string]CommandMap `json:"cmds"`
}

// XOFile is xo.json file which can be created in the project root folder
// for easy import of a xo project commands
type XOFile struct {
	Name     string    `json:"name"`
	Commands []Command `json:"commands"`
}
