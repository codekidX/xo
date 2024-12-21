// / display is a package for showing all outputs in colored format!
// / because we are not living in the 40s!
package display

import (
	"os"
	"xo/internal/types"

	"github.com/jedib0t/go-pretty/v6/table"
)

func ProjectInfo(name string, commands types.CommandMap) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{name})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Command", "Help"})
	t.AppendSeparator()
	for commandName, cmdInfo := range commands {
		t.AppendRow(table.Row{commandName, cmdInfo.Help})
	}

	t.Render()
}
