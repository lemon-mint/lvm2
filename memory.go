package lvm2

import "errors"

type MemoryBlock struct {
	Start uint64
	End   uint64

	Block []byte
}

type Memory struct {
	Blocks []MemoryBlock

	MemoryHead uint64
	MaxAddress uint64

	Stack      MemoryBlock
	Cache      *MemoryBlock
	CacheIndex int
}

func NewMemory() *Memory {
	m := &Memory{}

	m.MaxAddress = 0xFFFFFFFFFFFFFFFF

	// Stack Size = 16MB
	m.Stack.Block = make([]byte, 1024*1024*16)
	m.Stack.End = m.MaxAddress
	m.Stack.Start = m.Stack.End - uint64(len(m.Stack.Block))

	m.Blocks = make([]MemoryBlock, 0, 32)

	return m
}

const PAGE_SIZE = 1 << 12 // 4KB

var (
	ErrInvalidSize    = errors.New("Invalid Size")
	ErrNoMemory       = errors.New("No Memory")
	ErrInvalidAddress = errors.New("Invalid Address")
)

func (m *Memory) Allocate(size uint64) uint64 {
	if size == 0 {
		return m.MemoryHead
	}

	var block MemoryBlock
	block.Block = make([]byte, size)
	block.Start = m.MemoryHead
	block.End = block.Start + size
	m.MemoryHead += size
	m.Blocks = append(m.Blocks, block)

	return block.Start
}

func (m *Memory) Free(start uint64) error {
	_, index, err := m.LoadBlockIndex(start)
	if err != nil {
		return err
	}
	if index == -1 {
		return ErrInvalidAddress
	}

	m.Blocks = append(m.Blocks[:index], m.Blocks[index+1:]...)
	return nil
}

var ErrSegmentationFault = errors.New("Segmentation Fault")

func (m *Memory) LoadBlock(address uint64) (MemoryBlock, error) {
	block, _, err := m.LoadBlockIndex(address)
	return block, err
}

func (m *Memory) LoadBlockIndex(address uint64) (MemoryBlock, int, error) {
	// Check if we have a cache
	if m.Cache != nil && m.Cache.Start <= address && m.Cache.End >= address {
		return *m.Cache, m.CacheIndex, nil
	}

	// Check if address is in stack
	if address >= m.Stack.Start && address < m.Stack.End {
		return m.Stack, -1, nil
	}

	// Check if address is in memory (binary search)
	low := 0
	high := len(m.Blocks) - 1
	// for i := range m.Blocks {
	// 	if m.Blocks[i].Start <= address && m.Blocks[i].End >= address {
	// 		m.Cache = &m.Blocks[i]
	// 		m.CacheIndex = i
	// 		return m.Blocks[i], i, nil
	// 	}
	// }
	for low <= high {
		mid := (low + high) / 2
		if m.Blocks[mid].Start <= address && m.Blocks[mid].End >= address {
			m.Cache = &m.Blocks[mid]
			m.CacheIndex = mid
			return m.Blocks[mid], mid, nil
		} else if m.Blocks[mid].Start > address {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return MemoryBlock{}, -1, ErrSegmentationFault
}

func (m *Memory) ReadAt(address uint64, p []byte) (int, error) {
	var read int
	for {
		block, err := m.LoadBlock(address)
		if err != nil {
			return 0, err
		}

		offset := address - block.Start
		n := copy(p, block.Block[offset:])
		read += n
		p = p[n:]
		address += uint64(n)

		if len(p) == 0 || n == 0 {
			return read, nil
		}
	}
}

func (m *Memory) WriteAt(address uint64, p []byte) (int, error) {
	var written int
	for {
		block, err := m.LoadBlock(address)
		if err != nil {
			return 0, err
		}

		offset := address - block.Start
		n := copy(block.Block[offset:], p)
		written += n
		p = p[n:]
		address += uint64(n)

		if len(p) == 0 || n == 0 {
			return written, nil
		}
	}
}

func (m *Memory) GetMemoryFunc(address uint64, size uint64, iterf func(addr uint64, b []byte) error) error {
	var r uint64 = size
	for {
		block, err := m.LoadBlock(address)
		if err != nil {
			return err
		}

		offset := address - block.Start
		b := block.Block[offset:]
		if r > uint64(len(b)) {
			err = iterf(address, b)
		} else {
			b = b[:r]
			err = iterf(address, b)
		}

		if err != nil {
			return err
		}

		r -= uint64(len(b))
		address += uint64(len(b))

		if r == 0 {
			return nil
		}
	}
}

func (m *Memory) SetProgram(p []byte) {
	m.MemoryHead = uint64(len(p))
	m.Blocks = append(m.Blocks, MemoryBlock{
		Start: 0,
		End:   uint64(len(p)),
		Block: p,
	})
}

func (m *Memory) Reset() {
	for i := range m.Blocks {
		m.Blocks[i] = MemoryBlock{}
	}
	m.Blocks = m.Blocks[:0]
	m.Cache = nil
	m.CacheIndex = 0
	m.MemoryHead = 0
	for i := range m.Stack.Block {
		m.Stack.Block[i] = 0
	}
}
