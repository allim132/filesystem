package main

import (
	"github.com/allim132/filesystem/internal/cli"
)

func main() {
    cli := cli.NewCLI()
    cli.Run()
}