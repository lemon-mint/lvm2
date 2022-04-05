package lvm2

type InstructionType uint16

const (
	InstructionType_NOP InstructionType = iota

	InstructionType_ADD
	InstructionType_SUB
	InstructionType_MUL
	InstructionType_DIV

	InstructionType_AND
	InstructionType_OR
	InstructionType_XOR
	InstructionType_NOT

	InstructionType_SHL
	InstructionType_SHR
)
