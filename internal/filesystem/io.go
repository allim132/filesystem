package filesystem

import (
	"fmt"
	"os"
)

func writeBlock(file *os.File, blockNum int, data []byte) error {
    // Size Validation
    if len(data) != BlockSize {
        return fmt.Errorf("Data size and block size do not match: Data size must be exactly %d", BlockSize)
    }

    offset := int64(blockNum) * BlockSize
    _, err := file.WriteAt(data, offset)
    if err != nil {
        return fmt.Errorf("Failed to write block %d: %v", blockNum, err)
    }

    return nil
}

func readBlock(file *os.File, blockNum int) ([]byte, error) {
    // Create Buffer to hold block Data
    data := make([]byte, BlockSize)

    // Calculate the offset of the block
    offset := int64(blockNum) * BlockSize

    // Read block from offset buffer
    _, err := file.ReadAt(data, offset)
    if err != nil {
        return nil, fmt.Errorf("Failed to read block %d: %v", blockNum, err)
    }

    return data, nil
}