package tui

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"xo/internal/store"
	"xo/internal/sys"
	"xo/internal/types"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	screenProjects screen = iota
	screenCommands
	screenAddForm
	screenAddProject
	screenReadme
	screenSetupForm
	screenSetupList
	screenAddSetupForm
)

type AppModel struct {
	store        *store.Store
	screen       screen
	projects     projectsModel
	commands     commandsModel
	addForm      formModel
	addProj      addProjectFormModel
	setupForm    setupFormModel
	setupList    setupListModel
	addSetupForm addSetupFormModel
	readme       readmeModel
	width        int
	height       int
	toastMsg     string
	toastID      int64
}

func NewAppModel(s *store.Store) AppModel {
	return AppModel{
		store:    s,
		screen:   screenProjects,
		projects: newProjectsModel(s),
	}
}

func (m AppModel) Init() tea.Cmd { return nil }

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.projects.resize(m.width, m.height)
		if m.commands.projectName != "" {
			m.commands.resize(m.width, m.height)
		}
		if m.readme.projectName != "" {
			m.readme.resize(m.width, m.height)
		}
		if m.setupList.list.Title != "" {
			m.setupList.resize(m.width, m.height)
		}

	case tea.KeyMsg:
		// Global quit
		if key.Matches(msg, Keys.Quit) && m.screen != screenAddForm && m.screen != screenAddProject {
			return m, tea.Quit
		}
		// Global back (from commands to projects)
		if key.Matches(msg, Keys.Back) && m.screen == screenCommands {
			// Don't go back if commands list is filtering
			if m.commands.list.FilterState() != 1 { // 1 = filtering
				m.screen = screenProjects
				return m, nil
			}
		}
		// Global back from setup list
		if key.Matches(msg, Keys.Back) && m.screen == screenSetupList {
			if m.setupList.list.FilterState() != 1 {
				m.screen = screenProjects
				return m, nil
			}
		}

	case ProjectSelectedMsg:
		cmdMap := m.store.GetCommands()[msg.Name]
		projPath := m.store.GetProjects()[msg.Name]
		m.commands = newCommandsModel(msg.Name, projPath, cmdMap)
		m.commands.resize(m.width, m.height)
		m.screen = screenCommands
		return m, nil

	case CDProjectMsg:
		shell := sys.DefaultShell()
		c := exec.Command(shell)
		c.Dir = msg.Path
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			// Back to projects screen
			return nil
		})

	case OpenAddFormMsg:
		m.addForm = newFormModel(m.commands.projectName)
		m.screen = screenAddForm
		return m, nil

	case OpenAddProjectFormMsg:
		m.addProj = newAddProjectFormModel()
		m.screen = screenAddProject
		return m, nil

	case FormCancelMsg:
		if m.screen == screenAddForm {
			m.screen = screenCommands
		} else if m.screen == screenAddProject || m.screen == screenSetupForm {
			m.screen = screenProjects
		} else if m.screen == screenReadme {
			m.screen = screenProjects
		} else if m.screen == screenAddSetupForm {
			m.screen = screenSetupList
		}
		return m, nil

	case OpenSetupListMsg:
		m.setupList = newSetupListModel(m.store)
		m.setupList.resize(m.width, m.height)
		m.screen = screenSetupList
		return m, nil

	case OpenAddSetupFormMsg:
		m.addSetupForm = newAddSetupFormModel()
		m.screen = screenAddSetupForm
		return m, nil

	case DeleteSetupCmdMsg:
		m.store.RemoveSetupCommand(msg.OSName, msg.CmdStr)
		m.setupList = newSetupListModel(m.store)
		m.setupList.resize(m.width, m.height)
		return m, ShowToast("Deleted setup command")

	case SetupFormSubmitMsg:
		if msg.OSName != "" && msg.CmdStr != "" {
			m.store.AddSetupCommand(msg.OSName, msg.CmdStr)
			m.setupList = newSetupListModel(m.store)
			m.setupList.resize(m.width, m.height)
		}
		m.screen = screenSetupList
		return m, nil

	case RunSetupCmdMsg:
		c := sys.ShellCommand(msg.CmdStr)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			// Just return to the setup list
			return nil
		})

	case SetupWorkspaceMsg:
		exe, _ := os.Executable()
		c := exec.Command(exe, "setup")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			m.projects = newProjectsModel(m.store)
			m.projects.resize(m.width, m.height)
			return nil
		})

	case ExportWorkspaceMsg:
		exe, _ := os.Executable()
		c := exec.Command(exe, "export")
		// we can just run this silently in the background
		_ = c.Run()
		return m, ShowToast("Exported workspace to xo-workspace.json")

	case OpenBootstrapFormMsg:
		m.setupForm = newSetupFormModel()
		m.screen = screenSetupForm
		return m, nil

	case BootstrapFormSubmitMsg:
		exe, _ := os.Executable()
		c := exec.Command(exe, "setup", msg.URLOrPath)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			m.projects = newProjectsModel(m.store)
			m.projects.resize(m.width, m.height)
			m.screen = screenProjects
			return nil
		})

	case ShowToastMsg:
		m.toastMsg = msg.Message
		m.toastID = msg.ID
		return m, nil

	case ClearToastMsg:
		if m.toastID == msg.ID {
			m.toastMsg = ""
		}
		return m, nil

	case OpenReadmeMsg:
		readmePath := filepath.Join(msg.Path, "README.md")
		content, err := os.ReadFile(readmePath)
		if err != nil {
			// check lowercase
			readmePath = filepath.Join(msg.Path, "readme.md")
			content, err = os.ReadFile(readmePath)
			if err != nil {
				// check capitalized
				readmePath = filepath.Join(msg.Path, "Readme.md")
				content, err = os.ReadFile(readmePath)
			}
		}

		if err != nil {
			return m, ShowToast("No README.md found in " + msg.Name)
		}

		m.readme = newReadmeModel(msg.Name, string(content), m.width, m.height)
		m.screen = screenReadme
		return m, nil

	case FormSubmitMsg:
		if msg.Command.Name != "" && msg.Command.CmdStr != "" {
			err := m.store.AddCommandToProject(m.commands.projectName, msg.Command)
			if err == nil {
				cmdMap := m.store.GetCommands()[m.commands.projectName]
				m.commands = newCommandsModel(m.commands.projectName, m.commands.projectPath, cmdMap)
				m.commands.resize(m.width, m.height)
			}
		}
		m.screen = screenCommands
		return m, nil

	case ProjectFormSubmitMsg:
		dir := msg.Path
		if strings.HasPrefix(dir, "~/") {
			home, _ := os.UserHomeDir()
			dir = filepath.Join(home, dir[2:])
		}
		dir = path.Clean(dir)

		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			m.addProj.errorMessage = "Directory does not exist."
			return m, nil
		}

		xoFilePath := filepath.Join(dir, "xo.json")
		if _, err := os.Stat(xoFilePath); os.IsNotExist(err) {
			if err := m.store.CreateXOFile(dir); err != nil {
				m.addProj.errorMessage = "Failed to create xo.json: " + err.Error()
				return m, nil
			}
		}

		b, _ := os.ReadFile(xoFilePath)
		var xoFile types.XOFile
		json.Unmarshal(b, &xoFile)

		if err := m.store.AddProject(dir, xoFile); err != nil {
			m.addProj.errorMessage = "Failed to import project: " + err.Error()
			return m, nil
		}

		m.projects = newProjectsModel(m.store)
		m.projects.resize(m.width, m.height)
		m.screen = screenProjects
		return m, nil

	case CommandSelectedMsg:
		workDir := msg.Path
		if msg.WorkDir != "" {
			workDir = filepath.Join(msg.Path, msg.WorkDir)
		}

		c := sys.ShellCommand(msg.CmdStr)
		c.Dir = workDir
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return tea.QuitMsg{}
		})
	}

	switch m.screen {
	case screenProjects:
		var cmd tea.Cmd
		m.projects, cmd = m.projects.Update(msg)
		return m, cmd
	case screenCommands:
		var cmd tea.Cmd
		m.commands, cmd = m.commands.Update(msg)
		return m, cmd
	case screenAddForm:
		var cmd tea.Cmd
		m.addForm, cmd = m.addForm.Update(msg)
		return m, cmd
	case screenAddProject:
		var cmd tea.Cmd
		m.addProj, cmd = m.addProj.Update(msg)
		return m, cmd
	case screenSetupForm:
		var cmd tea.Cmd
		m.setupForm, cmd = m.setupForm.Update(msg)
		return m, cmd
	case screenSetupList:
		var cmd tea.Cmd
		m.setupList, cmd = m.setupList.Update(msg)
		return m, cmd
	case screenAddSetupForm:
		var cmd tea.Cmd
		m.addSetupForm, cmd = m.addSetupForm.Update(msg)
		return m, cmd
	case screenReadme:
		var cmd tea.Cmd
		m.readme, cmd = m.readme.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m AppModel) View() string {
	view := ""
	switch m.screen {
	case screenProjects:
		view = m.projects.View()
	case screenCommands:
		view = m.commands.View()
	case screenAddForm:
		view = m.addForm.View()
	case screenAddProject:
		view = m.addProj.View()
	case screenSetupForm:
		view = m.setupForm.View()
	case screenSetupList:
		view = m.setupList.View()
	case screenAddSetupForm:
		view = m.addSetupForm.View()
	case screenReadme:
		view = m.readme.View()
	}

	toastLine := ""
	if m.toastMsg != "" {
		toastLine = toastStyle.Render(m.toastMsg)
	} else {
		toastLine = "\n" // placeholder to prevent jitter
	}

	return lipgloss.JoinVertical(lipgloss.Left, view, toastLine)
}

func Launch(s *store.Store) error {
	p := tea.NewProgram(NewAppModel(s), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
