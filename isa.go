package lvm2

type InstructionType byte

const WORD_SIZE = 8

const (
	InstructionType_NOP InstructionType = iota

	InstructionType_ADD // R0 = R1 + R2
	InstructionType_SUB // R0 = R1 - R2
	InstructionType_MUL // R0 = R1 * R2
	InstructionType_DIV // R0 = R1 / R2

	InstructionType_AND // R0 = R1 & R2
	InstructionType_OR  // R0 = R1 | R2
	InstructionType_XOR // R0 = R1 ^ R2
	InstructionType_NOT // R0 = ~R1

	InstructionType_SHL // R0 = R1 << R2
	InstructionType_SHR // R0 = R1 >> R2

	InstructionType_CMP // R0 = R1 - R2
	InstructionType_JMP // PC = R0

	InstructionType_JG  // if R0 > 0; PC = R1
	InstructionType_JL  // if R0 < 0; PC = R1
	InstructionType_JE  // if R0 == 0; PC = R1
	InstructionType_JNE // if R0 != 0; PC = R1
	InstructionType_JGE // if R0 >= 0; PC = R1
	InstructionType_JLE // if R0 <= 0; PC = R1

	InstructionType_PUSH // SP = SP - WORD_SIZE
	InstructionType_POP  // SP = SP + WORD_SIZE

	InstructionType_LOAD  // R0 = [MEM[R1 + R2]]
	InstructionType_STORE // [MEM[R1 + R2]] = R0
	InstructionType_MOV   // R0 = R1

	InstructionType_CALL // SP = SP - WORD_SIZE; [SP] = PC; PC = R0
	InstructionType_RET  // PC = [SP]; SP = SP + WORD_SIZE

	InstructionType_SYSCALL // R0 = syscall(R1, R2)
)
