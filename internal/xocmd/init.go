package xocmd

import (
	"os"

	"github.com/spf13/cobra"
)

func (xo *XOCmd) InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "initialize xo in this folder",
		Run: func(cmd *cobra.Command, args []string) {
			p, _ := os.Getwd()
			if err := xo.store.CreateXOFile(p); err != nil {
				panic(err)
			}
		},
	}
}
