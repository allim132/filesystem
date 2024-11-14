package filesystem

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
	FileSize               int32
	LastModified           uint32 // Unix timestamp in seconds
	BlockPointerTableIndex int32
	Username               [MaxUsername]byte
}

type BlockPointerTable struct {
	Pointers [8]int32 // 7 data block pointers + 1 chaining pointer
}

type FileSystem struct {
	FNT         []FNTEntry
	DABPT       []DABPTEntry
	DataBlocks  [][]byte
	TotalBlocks int
	FreeBlocks  []bool
	CurrentUser [MaxUsername]byte
	DiskName    string
}