package rifx

import (
	"encoding/binary"
	"fmt"
	"io"
)

// FromReader parses the RIFX format from a reader
func FromReader(r io.Reader) (*List, error) {
	b4 := make([]byte, 4)
	io.ReadFull(r, b4)
	if fmt.Sprintf("%s", b4) != "RIFX" {
		return nil, fmt.Errorf("Unknown AEP file format")
	}
	io.ReadFull(r, b4)
	fileSize := binary.BigEndian.Uint32(b4)
	list, _, err := readList(r, fileSize)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func readList(r io.Reader, limit uint32) (*List, uint32, error) {
	listBlock := &List{}
	readBytes := uint32(0)
	blocks := make([]*Block, 0)

	// Read the list identifier
	idBytes := make([]byte, 4)
	n, err := io.ReadFull(r, idBytes)
	readBytes += uint32(n)
	if err != nil {
		return nil, readBytes, err
	}

	// Read all blocks in the list up to byte limit
	numBlocks := 0
	for readBytes < limit {
		block, n, err := readBlock(r)
		readBytes += n
		if err != nil {
			return nil, readBytes, err
		}
		blocks = append(blocks, block)
		numBlocks++
	}

	listBlock.Identifier = fmt.Sprintf("%s", idBytes)
	listBlock.NumBlocks = numBlocks
	listBlock.Blocks = blocks
	return listBlock, readBytes, nil
}

func readBlock(r io.Reader) (*Block, uint32, error) {
	block := &Block{}
	bytesRead := uint32(0)
	b4 := make([]byte, 4)

	// Read the type of block
	n, err := io.ReadFull(r, b4)
	bytesRead += uint32(n)
	if err != nil {
		return nil, bytesRead, err
	}
	block.Type = fmt.Sprintf("%s", b4)

	// Read the number of bytes contained in the block
	n, err = io.ReadFull(r, b4)
	bytesRead += uint32(n)
	if err != nil {
		return nil, bytesRead, err
	}
	block.Size = binary.BigEndian.Uint32(b4)

	// Read the bock data
	switch block.Type {
	case "LIST":
		group, n, err := readList(r, block.Size)
		bytesRead += n
		if err != nil {
			return nil, bytesRead, err
		}
		block.Data = group
	default:
		blockData := make([]byte, block.Size)
		n, err = io.ReadFull(r, blockData)
		bytesRead += uint32(n)
		if err != nil {
			return nil, bytesRead, err
		}
		block.Data = blockData
	}

	// Read padding if data is odd length
	if (block.Size % 2) != 0 {
		io.ReadFull(r, []byte{0})
		bytesRead++
	}

	return block, bytesRead, nil
}
