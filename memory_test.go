package lvm2

import (
	"fmt"
	"testing"
)

func TestMemory_GetMemoryFunc(t *testing.T) {
	m := NewMemory()
	m.Blocks = append(m.Blocks, MemoryBlock{
		Start: 0,
		End:   9,
		Block: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	})
	m.Blocks = append(m.Blocks, MemoryBlock{
		Start: 10,
		End:   19,
		Block: []byte{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
	})

	var data []byte
	m.GetMemoryFunc(0, 15, func(addr uint64, b []byte) error {
		fmt.Printf("addr: %d, b: %v\n", addr, b)
		data = append(data, b...)
		return nil
	})

	fmt.Printf("data: %v\n", data)
	if string(data) != string([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}) {
		t.Error("data error")
	}
}
