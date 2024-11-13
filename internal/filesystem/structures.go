package filesystem

import "time"

const (
    BlockSize            = 256
    MaxFilename          = 56
    MaxUsername          = 40
    EntriesPerDABPTBlock = 4
)

type FNTEntry struct {
    Filename     [MaxFilename]byte
    InodePointer int32
}

type DABPTEntry struct {
    FileSize      int32
    LastModified  time.Time
    BlockPointers [8]int32
    Username      [MaxUsername]byte
}

type FileSystem struct {
    FNT         []FNTEntry
    DABPT       []DABPTEntry
    DataBlocks  [][]byte
    TotalBlocks int
    FreeBlocks  []bool
}
