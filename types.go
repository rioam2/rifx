package rifx

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Block represents a single block of binary data
type Block struct {
	Type string
	Size uint32
	Data interface{}
}

// ToStruct de-serializes block data into the provided struct pointer
func (b *Block) ToStruct(ptr interface{}) error {
	return binary.Read(bytes.NewReader(b.Data.([]byte)), binary.BigEndian, ptr)
}

// ToString returns block data as a string
func (b *Block) ToString() string {
	return fmt.Sprintf("%s", b.Data.([]byte))
}

// ToUint8 returns block data as uint8
func (b *Block) ToUint8() uint8 {
	return uint8(b.Data.([]byte)[0])
}

// ToUint16 returns block data as uint16
func (b *Block) ToUint16() uint16 {
	return binary.BigEndian.Uint16(b.Data.([]byte))
}

// ToUint32 returns block data as uint32
func (b *Block) ToUint32() uint32 {
	return binary.BigEndian.Uint32(b.Data.([]byte))
}

// ToUint64 returns block data as uint64
func (b *Block) ToUint64() uint64 {
	return binary.BigEndian.Uint64(b.Data.([]byte))
}

// List represents a collection of binary blocks
type List struct {
	Identifier string
	NumBlocks  int
	Blocks     []*Block
}

// ForEach iterates over the list's blocks and invokes a callback for each
func (l *List) ForEach(cb func(*Block)) {
	for _, block := range l.Blocks {
		cb(block)
	}
}

// Map performs a basic mapping operator over a list's blocks
func (l *List) Map(cb func(*Block) interface{}) []interface{} {
	var ret []interface{}
	l.ForEach(func(b *Block) {
		ret = append(ret, cb(b))
	})
	return ret
}

// Filter performs a basic filtering operator over a list's blocks
func (l *List) Filter(cb func(*Block) bool) *List {
	ret := &List{Identifier: l.Identifier}
	l.ForEach(func(b *Block) {
		if cb(b) {
			ret.Blocks = append(ret.Blocks, b)
		}
	})
	ret.NumBlocks = len(ret.Blocks)
	return ret
}

// SublistFilter returns a slice of child lists that have the provided identifier
func (l *List) SublistFilter(identifier string) []*List {
	filtered := l.Filter(func(b *Block) bool {
		return b.Type == "LIST" && b.Data.(*List).Identifier == identifier
	}).Map(func(b *Block) interface{} {
		return b.Data
	})
	ret := make([]*List, len(filtered))
	for idx, block := range filtered {
		ret[idx] = block.(*List)
	}
	return ret
}

// Find performs a basic find operation over a list's blocks
func (l *List) Find(cb func(*Block) bool) (*Block, error) {
	for _, block := range l.Blocks {
		if cb(block) {
			return block, nil
		}
	}
	return nil, fmt.Errorf("ENOTFOUND")
}

// SublistMerge filters sublists with the specified identifier and concatenates their blocks in a new list
func (l *List) SublistMerge(identifier string) *List {
	newList := &List{
		Identifier: identifier,
		Blocks:     make([]*Block, 0),
		NumBlocks:  0,
	}
	for _, sublist := range l.SublistFilter(identifier) {
		newList.Blocks = append(newList.Blocks, sublist.Blocks...)
		newList.NumBlocks += sublist.NumBlocks
	}
	return newList
}
