package xocmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) RemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm",
		Short: "remove a xo project",
		Run: func(cmd *cobra.Command, args []string) {
			p, _ := os.Getwd()
			pathMap := xo.store.GetProjectPaths()
			if _, ok := pathMap[p]; !ok {
				fmt.Println("project:", p, "is not present please do 'xo i && xo import'!")
			} else {
				xo.store.RemoveProject(p)
				fmt.Println("removed project:", pathMap[p], "from 'xo'")
			}
		},
	}
}
