package main

import (
	loom "github.com/orieken/loom"
	"github.com/orieken/loom/cmd/loom/cmd"
)

func main() {
	cmd.Execute(loom.FrameworkFS, loom.MCPFS)
}
