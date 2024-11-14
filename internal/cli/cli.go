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
			c.formatfs() // Call formatfs method to format file system
		case "savefs":
			c.savefs(args)
		case "openfs":
			c.openfs(args)
		case "list":
			c.listFiles() // Call listFiles method to list files
		case "remove":
			c.remove(args)
		case "rename":
			c.rename(args)
		case "put":
			c.put(args)
		case "get":
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
	fmt.Println("openfs (diskname) - Open existing file system")
	fmt.Println("list - List files")
	fmt.Println("remove (name) - Removes given file")
	fmt.Println("rename (currentname) (newname) - Renames a given file")
	fmt.Println("put (externalfile) - Stores a file into the disk")
	fmt.Println("get (internalfile) - Gets a file from the file system to host's OS file system")
}

func createfs(reader *bufio.Reader) *filesystem.FileSystem {
	fmt.Print("Enter a name for the disk (e.g., disk01): ")
	diskName, _ := reader.ReadString('\n')
	diskName = strings.TrimSpace(diskName)

	fmt.Print("Enter your username: ")
	currentUser, _ := reader.ReadString('\n')
	currentUser = strings.TrimSpace(currentUser)

	fmt.Printf("Creating File System...\nEnter number of blocks: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	number, err := strconv.ParseInt(input, 10, 32)

	if err != nil {
			fmt.Println("Error: Input must be an integer!")
			return nil
	}

	// Create the filesystem and set DiskName and CurrentUser
	fs := filesystem.CreateFS(int(number), currentUser)
	fs.DiskName = diskName // Optionally set DiskName here

	fmt.Printf("File system with %d blocks successfully created!\n", number)
	return fs
}

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

func (c *CLI) formatfs() {
    // Check if the filesystem is loaded
    if c.fs == nil {
        fmt.Println("No filesystem loaded. Please create or open a filesystem first.")
        return
    }

    // Get total number of blocks from the filesystem
    totalBlocks := c.fs.TotalBlocks

    // Prompt user for number of entries (for both FNT and DABPT)
    fmt.Printf("Enter number of entries (for filenames and DABPT). Max number of entries is %d: ", totalBlocks)
    
    reader := bufio.NewReader(os.Stdin)
    inputEntries, _ := reader.ReadString('\n')
    inputEntries = strings.TrimSpace(inputEntries)
    numEntries, err := strconv.Atoi(inputEntries)
    if err != nil || numEntries <= 0 || numEntries > totalBlocks {
        fmt.Println("Invalid input for number of entries. Please enter a positive integer within the limit.")
        return
    }

    // Call FormatFS function to format the filesystem
    err = filesystem.FormatFS(c.fs, numEntries, numEntries) // Same number for both FNT and DABPT
    if err != nil {
        fmt.Printf("Failed to format filesystem: %v\n", err)
        return
    }

    fmt.Println("Filesystem formatted successfully.")
}

func (c *CLI) put(args []string) {
    // Check if the filesystem is loaded
    if c.fs == nil {
        fmt.Println("No filesystem loaded. Please create or open a filesystem first.")
        return
    }

    // Check if a filename argument is provided
    if len(args) < 2 {
        fmt.Println("Usage: put <filename>")
        return
    }

    // Get the external file name from the command arguments
    externalFileName := args[1]

    // Call PutFS function to store the external file in the filesystem
    err := filesystem.PutFS(c.fs, externalFileName)
    if err != nil {
        fmt.Printf("Failed to put file into filesystem: %v\n", err)
        return
    }

    fmt.Println("File successfully stored in the filesystem.")
}

func (c *CLI) remove(args []string) {
	// Check if the filesystem is loaded
	if c.fs == nil {
		fmt.Println("No filesystem loaded. Please create or open a filesystem first.")
		return
	}
	
	// Check if a filename argument is provided
	if len(args) < 2 {
		fmt.Println("Usage: remove <filename>")
		return
	}

	// Get the internal file name from the command arguments
	internalFileName := args[1]
	
	// Call RemoveFS function to remove the internal file from the filesystem
	err := filesystem.RemoveFS(c.fs, internalFileName)
	if err != nil {
		fmt.Printf("Failed to remove file from filesystem: %v\n", err)
		return
	}
	
	fmt.Println("File successfully removed from the filesystem.")
}

func (c *CLI) savefs(args []string) {
	// Check if the filesystem is loaded
	if c.fs == nil {
		fmt.Println("No filesystem loaded. Please create or open a filesystem first.")
		return
	}
	
	// Call SaveFS function to save the filesystem
	err := filesystem.SaveFS(c.fs, c.fs.DiskName)
	if err != nil {
		fmt.Printf("Failed to save filesystem: %v\n", err)
		return
	}
	
	fmt.Println("File system successfully saved.")
}

func (c *CLI) openfs(args []string) {
	// Check if the filesystem is loaded
	if c.fs != nil {
		fmt.Println("File system already loaded. Please close the current file system first.")
		return
	}
	
	// Check if a filename argument is provided
	if len(args) < 2 {
		fmt.Println("Usage: openfs <filename>")
		return
	}
	
	// Get the file name from the command arguments
	fileName := args[1]
	
	// Call OpenFS function to open the file system
	fmt.Printf("Trying to open file system: %s\n", fileNameo)
	fs, err := filesystem.OpenFS(fileName)
	if err != nil {
		fmt.Printf("Failed to open file system: %v\n", err)
		return
	}

	c.fs = fs
	fmt.Println("File system successfully opened.")
}

func (c *CLI) rename(args []string) {
	// Check if the filesystem is loaded
	if c.fs == nil {
		fmt.Println("No filesystem loaded. Please create or open a filesystem first.")
		return
	}
	
	// Check if a filename argument is provided
	if len(args) < 3 {
		fmt.Println("Usage: rename <currentfilename> <newfilename>")
		return
	}
	
	// Get the current file name and new file name from the command arguments
	currentFileName := args[1]
	newFileName := args[2]
	
	// Call RenameFS function to rename the internal file in the filesystem
	err := filesystem.RenameFS(c.fs, currentFileName, newFileName)
	if err != nil {
		fmt.Printf("Failed to rename file in filesystem: %v\n", err)
		return
	}
	
	fmt.Println("File successfully renamed in the filesystem.")
}