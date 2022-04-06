package lvm2

type InstructionType byte

const WORD_SIZE = 8

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

	InstructionType_CMP
	InstructionType_JMP

	InstructionType_JG
	InstructionType_JL
	InstructionType_JE
	InstructionType_JNE
	InstructionType_JGE
	InstructionType_JLE

	InstructionType_PUSH
	InstructionType_POP

	InstructionType_LOAD
	InstructionType_STORE
	InstructionType_MOV

	InstructionType_CALL
	InstructionType_RET

	InstructionType_SYSCALL
)
