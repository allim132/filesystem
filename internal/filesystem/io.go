package filesystem

import (
    "os"
    "encoding/binary"
)

func writeBlock(file *os.File, blockNum int, data []byte) error {
    // Implementation
}

func readBlock(file *os.File, blockNum int) ([]byte, error) {
    // Implementation
}