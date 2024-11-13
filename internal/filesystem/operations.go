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
        return fmt.Errorf("not enough blocks for %d filenames and %d DABPT entries", numFilenames, numDABPTEntries)
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
            return fmt.Errorf("failed to write FNT entry: %v", err)
        }
    }

    // Write DABPT
    for _, entry := range fs.DABPT {
        err := binary.Write(file, binary.LittleEndian, entry);
        if err != nil {
            return fmt.Errorf("failed to write DABPT entry: %v", err)
        }
    }

    // Write DataBlocks
    for _, block := range fs.DataBlocks {
        err := binary.Write(file, binary.LittleEndian, block);
        if err != nil {
            return fmt.Errorf("failed to write DataBlocks entry: %v", err)
        }
    }

    // Write FreeBlocks
    for _, isFree := range fs.FreeBlocks {
        err := binary.Write(file, binary.LittleEndian, isFree);
        if err != nil {
            return fmt.Errorf("failed to write FreeBlocks entry: %v", err);
        }
    }

    return nil;
}

// Use an existing disk image
func OpenFS(name string) (*FileSystem, error) {
	// Open file
    file, err := os.Open(name)
    if err != nil{
        return nil, fmt.Errorf("failed to open file: %v", err);
    }
    defer file.Close()

    fs := &FileSystem{}

    // Read total number of blocks
    var totalNumberOfBlocks int32
    err = binary.Read(file, binary.LittleEndian, &totalNumberOfBlocks)
    if err != nil{
        return nil, fmt.Errorf("failed to read total blocks: %v", err)
    }
    fs.TotalBlocks = int(totalNumberOfBlocks)

    // Read FNT
    fs.FNT = make([]FNTEntry, fs.TotalBlocks/4)
    for i := range fs.FNT {
        err = binary.Read(file, binary.LittleEndian, &fs.FNT[i])
        if err != nil {
            return nil, fmt.Errorf("failed to read FNT: %v", err)
        }
    } 

    // Read DABPT
    fs.DABPT = make([]DABPTEntry, fs.TotalBlocks/4)
    for i := range fs.DABPT {
        err = binary.Read(file, binary.LittleEndian, &fs.DABPT[i])
        if err != nil {
            return nil, fmt.Errorf("failed to read DABPT entry: %v", err)
        }
    }

    // Read DataBlocks
    fs.DataBlocks = make([][]byte, fs.TotalBlocks /2)
    for i := range fs.DataBlocks {
        fs.DataBlocks[i] = make([]byte, BlockSize)
        _, err = file.Read(fs.DataBlocks[i]) 
        if err != nil {
            return nil, fmt.Errorf("failed to read DataBlock entry: %v", err)
        }
    }

    // Read FreeBlocks
    fs.FreeBlocks = make([]bool, fs.TotalBlocks)
    for i := range fs.FreeBlocks {
        err = binary.Read(file, binary.LittleEndian, &fs.FreeBlocks[i])
        if err != nil {
            return nil, fmt.Errorf("failed to read FreeBlock entry: %v", err)
        }
    }

    return fs, nil
}

// Implement other operations (List, Remove, Rename, Put, Get, User)