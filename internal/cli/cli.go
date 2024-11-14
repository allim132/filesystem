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
        fmt.Print("\nType \"commands\" for list of commands\n")
        fmt.Print("FS> ")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)
        input = strings.ToLower(input)
        args := strings.Fields(input)

        if len(args) == 0 {
            continue
        }

        switch args[0] {
        case "commands":
            listoperations();
        case "createfs":
            // Handle createfs command
            createfs(*reader)
            
        case "formatfs":
            // Handle formatfs command
        // Add cases for other commands
        case "quit":
            return
        case "exit":
            return
        default:
            fmt.Println("Unknown command")
        }
    }
}

// Implement methods for each command (createfs, formatfs, etc.)
func listoperations() {
    fmt.Println("\nOperations:")
    fmt.Println("createfs - Create file system")
    fmt.Println("formatfs - Format file system")
    fmt.Println("savefs - Save file system")
    fmt.Println("openfs - Format file system")
    fmt.Println("list - List files")
    fmt.Println("remove (name) - Removes given file")
    fmt.Println("rename (currentname) (newname) - Renames a given file")
    fmt.Println("put (externalfile) - Stores a file into the disk")
    fmt.Println("get (internalfile) - Gets a file from the file system to host's OS file system")
    
}
func createfs(reader bufio.Reader) {
    fmt.Printf("Creating File System...")
    fmt.Printf("Enter number of blocks: ")
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(input)
    number, err := strconv.ParseInt(input, 10, 32)
    
    if (err != nil){
        fmt.Println("Error: Input must be a integer!")
        return
    }

    filesystem.CreateFS(int(number))

    fmt.Printf("File with blocks %d successfully created!\n", number)
}