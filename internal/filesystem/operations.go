package filesystem

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"time"
)

func CreateFS(numBlocks int, currentUser string) *FileSystem {
    // Convert currentUser to bytes
    var username [MaxUsername]byte
    copy(username[:], currentUser)

    if len(currentUser) > MaxUsername {
        username = [MaxUsername]byte{}
        copy(username[:], currentUser[:MaxUsername])
    }


    fs := &FileSystem{
        TotalBlocks: numBlocks,
        FNT:         make([]FNTEntry, 0),  // Will be initialized in FormatFS
        DABPT:       make([]DABPTEntry, 0), // Will be initialized in FormatFS
        DataBlocks:  make([][]byte, numBlocks),
        FreeBlocks:  make([]bool, numBlocks),
        DiskName:    "",  // Will be set when saving or opening a disk image
        CurrentUser: username, // Set the CurrentUser here
    }

    // Initialize all blocks as free initially
    for i := range fs.FreeBlocks {
        fs.FreeBlocks[i] = true
    }

    // Initialize all data blocks
    for i := range fs.DataBlocks {
        fs.DataBlocks[i] = make([]byte, BlockSize)
    }

    return fs
}

func FormatFS(fs *FileSystem, numFilenames, numDABPTEntries int) error {
    // Validate input parameters
    totalMetaBlocks := (numFilenames + 3) / 4 + (numDABPTEntries + 3) / 4 // 4 entries per block
    if totalMetaBlocks > fs.TotalBlocks {
        return fmt.Errorf("not enough blocks for %d filenames and %d DABPT entries", numFilenames, numDABPTEntries)
    }

    // Set up FNT
    fs.FNT = make([]FNTEntry, numFilenames)
    for i := range fs.FNT {
        fs.FNT[i] = FNTEntry{
            Filename:     [MaxFilename]byte{},
            InodePointer: -1, // Invalid pointer
        }
    }

    // Set up the Directory and Attribute/Block Pointer Table (DABPT)
    fs.DABPT = make([]DABPTEntry, numDABPTEntries)
    for i := range fs.DABPT {
        fs.DABPT[i] = DABPTEntry{
            FileSize:               0,
            LastModified:           uint32(time.Now().Unix()),
            BlockPointerTableIndex: -1, // Invalid pointer
            Username:               [MaxUsername]byte{},
        }
    }

    // Initialize FreeBlocks
    for i := range fs.FreeBlocks {
        if i < totalMetaBlocks {
            fs.FreeBlocks[i] = false // Metadata blocks are not free
        } else {
            fs.FreeBlocks[i] = true // Data blocks are initially free
        }
    }

    // Clear existing data in DataBlocks
    for i := range fs.DataBlocks {
        for j := range fs.DataBlocks[i] {
            fs.DataBlocks[i][j] = 0
        }
    }

    return nil // nil = no error
}

// Save the "Disk" in a file "name"
func SaveFS(fs *FileSystem, name string) error {
    // Open file
    fmt.Printf("Attempting to create file with name: %s\n", name)
    file, err := os.Create(name)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer file.Close()

    // Write total number of blocks
    err = binary.Write(file, binary.LittleEndian, int32(fs.TotalBlocks))
    if err != nil {
        return fmt.Errorf("failed to write total blocks: %v", err)
    }

    // Write FNT
    err = binary.Write(file, binary.LittleEndian, int32(len(fs.FNT)))
    if err != nil {
        return fmt.Errorf("failed to write FNT length: %v", err)
    }
    for _, entry := range fs.FNT {
        err := binary.Write(file, binary.LittleEndian, entry)
        if err != nil {
            return fmt.Errorf("failed to write FNT entry: %v", err)
        }
    }

    // Write DABPT
    err = binary.Write(file, binary.LittleEndian, int32(len(fs.DABPT)))
    if err != nil {
        return fmt.Errorf("failed to write DABPT length: %v", err)
    }
    for _, entry := range fs.DABPT {
        err := binary.Write(file, binary.LittleEndian, entry)
        if err != nil {
            return fmt.Errorf("failed to write DABPT entry: %v", err)
        }
    }

    // Write DataBlocks
    for _, block := range fs.DataBlocks {
        _, err := file.Write(block)
        if err != nil {
            return fmt.Errorf("failed to write DataBlock: %v", err)
        }
    }

    // Write FreeBlocks
    for _, isFree := range fs.FreeBlocks {
        err := binary.Write(file, binary.LittleEndian, isFree)
        if err != nil {
            return fmt.Errorf("failed to write FreeBlocks entry: %v", err)
        }
    }

    // Set current user to NULL
    _, err = file.Write(fs.CurrentUser[:])
    if err != nil {
        return fmt.Errorf("failed to write CurrentUser: %v", err)
    }

    // Write current user
    _, err = file.WriteString(fs.DiskName)
    if err != nil {
        return fmt.Errorf("failed to write DiskName: %v", err)
    }

    fs.DiskName = name
    return nil
}

// Use an existing disk image
func OpenFS(name string) (*FileSystem, error) {
    // Open file
    file, err := os.Open(name)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    fs := &FileSystem{}

    // Read total number of blocks
    var totalNumberOfBlocks int32
    err = binary.Read(file, binary.LittleEndian, &totalNumberOfBlocks)
    if err != nil {
        return nil, fmt.Errorf("failed to read total blocks: %v", err)
    }
    fs.TotalBlocks = int(totalNumberOfBlocks)

    // Read FNT length
    var fntLength int32
    err = binary.Read(file, binary.LittleEndian, &fntLength)
    if err != nil {
        return nil, fmt.Errorf("failed to read FNT length: %v", err)
    }

    // Read FNT
    fs.FNT = make([]FNTEntry, fntLength)
    for i := range fs.FNT {
        err = binary.Read(file, binary.LittleEndian, &fs.FNT[i])
        if err != nil {
            return nil, fmt.Errorf("failed to read FNT entry: %v", err)
        }
    }

    // Read DABPT length
    var dabptLength int32
    err = binary.Read(file, binary.LittleEndian, &dabptLength)
    if err != nil {
        return nil, fmt.Errorf("failed to read DABPT length: %v", err)
    }

    // Read DABPT
    fs.DABPT = make([]DABPTEntry, dabptLength)
    for i := range fs.DABPT {
        err = binary.Read(file, binary.LittleEndian, &fs.DABPT[i])
        if err != nil {
            return nil, fmt.Errorf("failed to read DABPT entry: %v", err)
        }
    }

    // Read DataBlocks
    fs.DataBlocks = make([][]byte, fs.TotalBlocks)
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

    // Read CurrentUser
    fs.CurrentUser = [MaxUsername]byte{}
    _, err = file.Read(fs.CurrentUser[:])
    if err != nil {
        return nil, fmt.Errorf("failed to read CurrentUser: %v", err)
    }

    // Read DiskName
    diskName, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read DiskName: %v", err)
    }
    fs.DiskName = string(diskName)

    return fs, nil
}

// Implement other operations (List, Remove, Rename, Put, Get, User)
func ListFS(fs *FileSystem) ([]string, error) {
    var fileList []string
    
    for _, entry := range fs.FNT {
        if entry.Filename != [MaxFilename]byte{} {
            filename := string(bytes.Trim(entry.Filename[:], "\x00"))
            
            // Get corresponding DABPT entry
            if int(entry.InodePointer) < len(fs.DABPT) {
                dabptEntry := fs.DABPT[entry.InodePointer]
                
                // Format file information
                fileInfo := fmt.Sprintf("File: %s, Size: %d bytes, Last Modified: %s, Owner: %s",
                    filename,
                    dabptEntry.FileSize,
                    time.Unix(int64(dabptEntry.LastModified), 0).Format(time.RFC3339),
                    string(bytes.Trim(dabptEntry.Username[:], "\x00")))
                
                fileList = append(fileList, fileInfo)
            } else {
                return nil, fmt.Errorf("invalid DABPT index for file %s", filename)
            }
        }
    }

    return fileList, nil
}
// getFreeBlockCount, addToFNT, allocateBlockPointerTable, allocateDataBlock, writeBlock, updateBlockPointerTable, updateDABPT, SaveFs
func PutFS(fs *FileSystem, externalFileName string) error {
    // Check if external file exists
    if _, err := os.Stat(externalFileName); os.IsNotExist(err) {
        return fmt.Errorf("external file does not exist: %v", err)
    }

    // Open and read the external file
    externalFile, err := os.Open(externalFileName)
    if err != nil {
        return fmt.Errorf("failed to open external file: %v", err)
    }
    defer externalFile.Close()


    fileInfo, err := externalFile.Stat()
    if err != nil {
        return fmt.Errorf("failed to get external file stats: %v", err)
    }

    // Validate available space in FS
    fileSize := fileInfo.Size()
    if fileSize == 0 {
        return fmt.Errorf("cannot add empty file")
    }
    requiredBlocks := int(math.Ceil(float64(fileSize) / float64(BlockSize)))
    if fs.getFreeBlockCount() < requiredBlocks {
        return fmt.Errorf("not enough space in the file system")
    }
    // Add file entry to FNT
    fntIndex, err := fs.addToFNT(filepath.Base(externalFileName))
    if err != nil {
        return fmt.Errorf("failed to add file to FNT: %v", err)
    }

    // Create DABPT entry with file metadata
    dabptEntry := DABPTEntry{
        FileSize:               int32(fileSize),
        LastModified:           uint32(fileInfo.ModTime().Unix()),
        BlockPointerTableIndex: 0, // Will be set later
        Username:               fs.CurrentUser,
    }

    // Allocate and update Block Pointer Table
    bptIndex, err := fs.allocateBlockPointerTable(requiredBlocks)
    if err != nil {
        return fmt.Errorf("failed to allocate Block Pointer Table: %v", err)
    }
    dabptEntry.BlockPointerTableIndex = int32(bptIndex)

    // Write file content to data blocks
    buffer := make([]byte, BlockSize)
    for i := 0; i < requiredBlocks; i++ {
        n, err := externalFile.Read(buffer)
        if err != nil && err != io.EOF {
            return fmt.Errorf("failed to read external file: %v", err)
        }
        
        blockIndex, err := fs.allocateDataBlock()
        if err != nil {
            return fmt.Errorf("failed to allocate data block: %v", err)
        }
        
        err = fs.writeBlock(blockIndex, buffer[:n])
        if err != nil {
            return fmt.Errorf("failed to write data block: %v", err)
        }
        
        err = fs.updateBlockPointerTable(bptIndex, i, blockIndex)
        if err != nil {
            return fmt.Errorf("failed to update Block Pointer Table: %v", err)
        }

        // Update FreeBlocks
        fs.FreeBlocks[blockIndex] = false
    }

    // Update DABPT
    err = fs.updateDABPT(fntIndex, dabptEntry)
    if err != nil {
        return fmt.Errorf("failed to update DABPT: %v", err)
    }

    // Save updated filesystem state
    err = SaveFS(fs, externalFileName)
    if err != nil {
        return fmt.Errorf("failed to save updated filesystem state: %v", err)
    }

    return nil
}

// getFreeBlockCount returns the number of free blocks in the filesystem
func (fs *FileSystem) getFreeBlockCount() int {
    count := 0
    for _, isFree := range fs.FreeBlocks {
        if isFree {
            count++
        }
    }
    return count
}

// addToFNT adds a new file entry to the FileNameTable
func (fs *FileSystem) addToFNT(filename string) (int, error) {
    for i, entry := range fs.FNT {
        if entry.Filename == [MaxFilename]byte{} {
            copy(fs.FNT[i].Filename[:], filename)
            fs.FNT[i].InodePointer = int32(i) // Assuming DABPT index matches FNT index
            return i, nil
        }
    }
    return -1, fmt.Errorf("FNT is full")
}

// allocateBlockPointerTable allocates a new Block Pointer Table
func (fs *FileSystem) allocateBlockPointerTable(requiredBlocks int) (int, error) {
    neededEntries := (requiredBlocks + 6) / 7 // 7 data pointers per entry, rounded up
    for i := 0; i < len(fs.DataBlocks); i++ {
        if fs.FreeBlocks[i] {
            fs.FreeBlocks[i] = false
            bpt := BlockPointerTable{
                Pointers: [8]int32{-1, -1, -1, -1, -1, -1, -1, -1}, // Initialize with -1 (invalid pointer)
            }
            binary.LittleEndian.PutUint32(fs.DataBlocks[i][:4], uint32(neededEntries))
            for j := 0; j < 8; j++ {
                binary.LittleEndian.PutUint32(fs.DataBlocks[i][4+j*4:], uint32(bpt.Pointers[j]))
            }
            return i, nil
        }
    }
    return -1, fmt.Errorf("no free blocks for Block Pointer Table")
}

// allocateDataBlock finds and allocates a free data block
func (fs *FileSystem) allocateDataBlock() (int, error) {
    for i := 0; i < len(fs.DataBlocks); i++ {
        if fs.FreeBlocks[i] {
            fs.FreeBlocks[i] = false
            return i, nil
        }
    }
    return -1, fmt.Errorf("no free data blocks")
}

// writeBlock writes data to a specific block
func (fs *FileSystem) writeBlock(blockIndex int, data []byte) error {
    if blockIndex < 0 || blockIndex >= len(fs.DataBlocks) {
        return fmt.Errorf("invalid block index")
    }
    copy(fs.DataBlocks[blockIndex], data)
    return nil
}

// updateBlockPointerTable updates a Block Pointer Table entry
func (fs *FileSystem) updateBlockPointerTable(bptIndex, entryIndex, blockIndex int) error {
    if bptIndex < 0 || bptIndex >= len(fs.DataBlocks) {
        return fmt.Errorf("invalid BPT index")
    }

    // Calculate entry offset
    entryOffset := 4 + (entryIndex % 7) * 4
    if entryOffset+4 > len(fs.DataBlocks[bptIndex]) {
        return fmt.Errorf("entry offset out of bounds")
    }

    // Write block index to the specified offset
    binary.LittleEndian.PutUint32(fs.DataBlocks[bptIndex][entryOffset:], uint32(blockIndex))

    // Check if we need to set a chaining pointer
    if entryIndex%7 == 6 { // If this is the last pointer in the group
        nextBPTCount := binary.LittleEndian.Uint32(fs.DataBlocks[bptIndex][:4]) // Assuming first 4 bytes store count
        if entryIndex/7+1 < int(nextBPTCount) { // Check if there are more entries needed
            nextBPTIndex, err := fs.allocateBlockPointerTable((entryIndex + 1) / 7)
            if err != nil {
                return err
            }
            // Set chaining pointer at the end of current block pointer table
            binary.LittleEndian.PutUint32(fs.DataBlocks[bptIndex][32:], uint32(nextBPTIndex))
        }
    }

    return nil
}

// updateDABPT updates a DABPT entry
func (fs *FileSystem) updateDABPT(fntIndex int, entry DABPTEntry) error {
    if fntIndex < 0 || fntIndex >= len(fs.DABPT) {
        return fmt.Errorf("invalid FNT index")
    }
    fs.DABPT[fntIndex] = entry
    return nil
}