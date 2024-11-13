package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/allim132/filesystem/internal/filesystem"
)

type CLI struct {
    fs *filesystem.FileSystem
}

func NewCLI() *CLI {
    return &CLI{}
}

func (c *CLI) Run() {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("FS> ")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)
        args := strings.Fields(input)

        if len(args) == 0 {
            continue
        }

        switch args[0] {
        case "createfs":
            // Handle createfs command
        case "formatfs":
            // Handle formatfs command
        // Add cases for other commands
        case "exit":
            return
        default:
            fmt.Println("Unknown command")
        }
    }
}

// Implement methods for each command (createfs, formatfs, etc.)