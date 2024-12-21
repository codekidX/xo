package main

import (
	"xo/internal/xocmd"
)

func main() {
	if err := xocmd.Run(); err != nil {
		panic(err)
	}
}
