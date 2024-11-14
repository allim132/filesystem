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
			listoperations()
		case "createfs":
			c.fs = createfs(reader) // Update to store the returned FileSystem
			if c.fs != nil {
				fmt.Println("File system created successfully.")
			}
		case "formatfs":
			// Handle formatfs command
			fmt.Println("Format FS command not implemented yet.")
		case "list":
			c.listFiles() // Call listFiles method to list files
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
	fmt.Println("openfs - Open existing file system")
	fmt.Println("list - List files")
	fmt.Println("remove (name) - Removes given file")
	fmt.Println("rename (currentname) (newname) - Renames a given file")
	fmt.Println("put (externalfile) - Stores a file into the disk")
	fmt.Println("get (internalfile) - Gets a file from the file system to host's OS file system")
}

func createfs(reader *bufio.Reader) *filesystem.FileSystem {
	fmt.Print("Creating File System...\nEnter number of blocks: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	number, err := strconv.ParseInt(input, 10, 32)

	if err != nil {
		fmt.Println("Error: Input must be an integer!")
		return nil
	}

	fs := filesystem.CreateFS(int(number)) // Call CreateFS and assign to fs

	if fs == nil {
		fmt.Println("Error creating filesystem.")
		return nil
	}

	fmt.Printf("File system with %d blocks successfully created!\n", number)
	return fs // Return the created filesystem
}

// Implement listFiles method to list files in the filesystem
func (c *CLI) listFiles() {
	if c.fs == nil {
		fmt.Println("No filesystem created. Please create one first.")
		return
	}

	fileList, err := filesystem.ListFS(c.fs) // Assuming ListFS is a method in your filesystem package
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		return
	}

    // Empty check

    if len(fileList) != 0 {
        for _, file := range fileList {
            fmt.Println(file)
        } 
    } else {
        fmt.Println("File system is empty!")
    }
	
}