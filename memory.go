package lvm2

type MemoryBlock struct {
	Start uint64
	End   uint64

	Block []byte
}

type Memory struct {
	Blocks []MemoryBlock

	MemoryHead uint64
	MaxAddress uint64

	Stack MemoryBlock
	Cache *MemoryBlock
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

type StringError string

func (s StringError) Error() string {
	return string(s)
}

const (
	ErrInvalidSize = StringError("Invalid Size")
	ErrNoMemory    = StringError("No Memory")
)

func (m *Memory) Allocate(size uint64) (uint64, error) {
	if size == 0 {
		return 0, nil
	}

	if size%PAGE_SIZE != 0 {
		return 0, ErrInvalidSize
	}

	blockCount := size / PAGE_SIZE

	// Flush Cache
	m.Cache = nil

	var start uint64 = m.MemoryHead

	for i := 0; i < int(blockCount); i++ {
		var block MemoryBlock
		block.Start = m.MemoryHead
		block.End = m.MemoryHead + PAGE_SIZE
		m.MemoryHead += PAGE_SIZE
		block.Block = make([]byte, PAGE_SIZE)

		m.Blocks = append(m.Blocks, block)
	}

	return start, nil
}

func (m *Memory) Free(start uint64, size uint64) error {
	if size == 0 {
		return nil
	}

	if size%PAGE_SIZE != 0 {
		return ErrInvalidSize
	}

	blockCount := size / PAGE_SIZE

	// Flush Cache
	m.Cache = nil

	for i := 0; i < int(blockCount); i++ {
		m.Blocks = append(m.Blocks[:start], m.Blocks[start+1:]...)
	}

	return nil
}

func (m *Memory) LoadBlocks(address uint64, size uint64) ([]MemoryBlock, error) {
	// Check if we have a cache
	if m.Cache != nil && m.Cache.Start <= address && m.Cache.End >= address+size {
		return []MemoryBlock{*m.Cache}, nil
	}

	// TODO: Check if we have a address in the stack
	// TODO: Check if we have a address in the program
	// TODO: Check if we have a address in the heap
	return nil, nil
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
	m.MemoryHead = 0
	for i := range m.Stack.Block {
		m.Stack.Block[i] = 0
	}
}
