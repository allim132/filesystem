package filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

func CreateFS(numBlocks int) *FileSystem {
	fs := &FileSystem{
		TotalBlocks: numBlocks,
		FNT:         make([]FNTEntry, numBlocks/4),
		DABPT:       make([]DABPTEntry, numBlocks/4),
		DataBlocks:  make([][]byte, numBlocks/2),
		FreeBlocks:  make([]bool, numBlocks),
	}

	for i := range fs.DataBlocks {
		fs.DataBlocks[i] = make([]byte, BlockSize)
	}

	// Initialize free space management
	for i := range fs.FreeBlocks {
		// Check if block is free
		if i < len(fs.FNT)+len(fs.DABPT) {
			fs.FreeBlocks[i] = false
		} else {
			fs.FreeBlocks[i] = true
		}
	}
	return fs
}
func FormatFS(fs *FileSystem, numFilenames, numDABPTEntries int) error {
	// Validate input paramaters
    totalMetaBlocks := (numFilenames + 3) / 4 + (numDABPTEntries + 3) / 4 // 4 entries per block
    if (totalMetaBlocks > fs.TotalBlocks) {
        return fmt.Errorf("Not enough blocks for %d filenames and %d DABPT entries", numFilenames, numDABPTEntries)
    }
    
    // Set up FNT
	fs.FNT = make([]FNTEntry, numFilenames)
	

	// Set up the Directory and Attribute/Block Pointer Table (DABPT)
    fs.DABPT = make([]DABPTEntry, numDABPTEntries)
	for i := range fs.DABPT {
		fs.DABPT[i].FileSize = 0
		fs.DABPT[i].LastModified = time.Now()
		fs.DABPT[i].BlockPointers = [8]int32{-1, -1, -1, -1, -1, -1, -1, -1} // Invalid pointers
		fs.DABPT[i].Username = [MaxUsername]byte{}
	}

    // Initialize DataBlocks
    // Clear existing data and reset free space management
    dataBlocksStart := totalMetaBlocks
    fs.DataBlocks = make([][]byte, fs.TotalBlocks-dataBlocksStart)
    for i := range fs.DataBlocks {
        fs.DataBlocks[i] = make([]byte, BlockSize)
    }

    // Initialize FreeBlocks
    // Reset free space bitmap
    for i := range fs.FreeBlocks {
        if i < dataBlocksStart {
            fs.FreeBlocks[i] = true    
        } else {
            fs.FreeBlocks[i] = false
        }
    }
    
	return nil // nil = no error
}

// Save the "Disk" in a file "name"
func SaveFS(fs *FileSystem, name string) error {
	

    // Open file
    file, err := os.Create(name)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write total number of blocks
    err = binary.Write(file, binary.LittleEndian, fs.TotalBlocks)
    if err != nil {
        return err
    }

    // Write FNT
    for entry := range fs.FNT {
        err := binary.Write(file, binary.LittleEndian, entry);
        if err != nil {
            return fmt.Errorf("Failed to write FNT entry: %v", err)
        }
    }

    // Write DABPT
    for _, entry := range fs.DABPT {
        err := binary.Write(file, binary.LittleEndian, entry);
        if err != nil {
            return fmt.Errorf("Failed to write DABPT entry: %v", err)
        }
    }

    // Write DataBlocks
    for _, block := range fs.DataBlocks {
        err := binary.Write(file, binary.LittleEndian, block);
        if err != nil {
            return fmt.Errorf("Failed to write DataBlocks entry: %v", err)
        }
    }

    // Write FreeBlocks
    for _, isFree := range fs.FreeBlocks {
        err := binary.Write(file, binary.LittleEndian, isFree);
        if err != nil {
            return fmt.Errorf("Failed to write FreeBlocks entry: %v", err);
        }
    }

    return nil;
}

func OpenFS(name string) (*FileSystem, error) {
	// Implementation
}

// Implement other operations (List, Remove, Rename, Put, Get, User)