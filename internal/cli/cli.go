package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
            fmt.Printf("Creating File System...")
            fmt.Printf("Enter number of blocks:")
            input, _ := reader.ReadString('\n')
            input = strings.TrimSpace(input)
            number, err := strconv.ParseInt(input, 10, 32)
            
            if (err != nil){
                fmt.Println("Error: Input must be a integer!")
                continue
            }

            filesystem.CreateFS(int(number))

            fmt.Printf("File with blocks %d succesfully created!", number)
            
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