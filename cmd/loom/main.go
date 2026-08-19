package main

import (
	loom "github.com/orieken/loom"
	"github.com/orieken/loom/cmd/loom/cmd"
)

func main() {
	_ = loom.FrameworkFS
	cmd.Execute()
}
