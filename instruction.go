package lvm2

type InstructionType byte

const WORD_SIZE = 8
const HALF_WORD_SIZE = WORD_SIZE / 2
const BYTE_SIZE = 1

const (
	InstructionType_NOP InstructionType = iota

	InstructionType_ADD // R0 = R1 + R2
	InstructionType_SUB // R0 = R1 - R2
	InstructionType_MUL // R0 = R1 * R2
	InstructionType_DIV // R0 = R1 / R2
	InstructionType_MOD // R0 = R1 % R2

	InstructionType_AND // R0 = R1 & R2
	InstructionType_OR  // R0 = R1 | R2
	InstructionType_XOR // R0 = R1 ^ R2
	InstructionType_NOT // R0 = ~R1

	InstructionType_SHL // R0 = R1 << R2
	InstructionType_SHR // R0 = R1 >> R2

	InstructionType_CMP // R0 = R1 - R2 (signed)
	InstructionType_JMP // PC = R0

	InstructionType_JG  // if R0 > 0; PC = R1
	InstructionType_JL  // if R0 < 0; PC = R1
	InstructionType_JE  // if R0 == 0; PC = R1
	InstructionType_JNE // if R0 != 0; PC = R1
	InstructionType_JGE // if R0 >= 0; PC = R1
	InstructionType_JLE // if R0 <= 0; PC = R1

	InstructionType_LOAD  // R0 = [MEM[R1 + R2]] (Load Register from Memory (WORD_SIZE))
	InstructionType_LOADH // R0 = [MEM[R1 + R2]] (Load Register from Memory (HALF_WORD_SIZE))
	InstructionType_LOADB // R0 = [MEM[R1 + R2]] (Load Register from Memory (BYTE_SIZE))

	InstructionType_STORE  // [MEM[R1 + R2]] = R0 (Store Register to Memory (WORD_SIZE))
	InstructionType_STOREH // [MEM[R1 + R2]] = R0 (Store Register to Memory (HALF_WORD_SIZE))
	InstructionType_STOREB // [MEM[R1 + R2]] = R0 (Store Register to Memory (BYTE_SIZE))

	InstructionType_MOV  // R0 = R1 (Move Register by WORD_SIZE)
	InstructionType_MOVH // R0 = R1 (Move Register by HALF_WORD_SIZE)
	InstructionType_MOVB // R0 = R1 (Move Register by BYTE_SIZE)

	InstructionType_PUSH // SP = SP - WORD_SIZE (stack.push(R0))
	InstructionType_POP  // SP = SP + WORD_SIZE (R0 = stack.pop())

	InstructionType_CALL // SP = SP - WORD_SIZE; [SP] = PC; PC = R0
	InstructionType_RET  // PC = [SP]; SP = SP + WORD_SIZE

	InstructionType_SYSCALL // R0 = syscall(R1, R2) (System Call) R0: errno, R1: syscall number, R2: register parameter
)

func (v InstructionType) String() string {
	switch v {
	case InstructionType_NOP:
		return "NOP"
	case InstructionType_ADD:
		return "ADD"
	case InstructionType_SUB:
		return "SUB"
	case InstructionType_MUL:
		return "MUL"
	case InstructionType_DIV:
		return "DIV"
	case InstructionType_AND:
		return "AND"
	case InstructionType_OR:
		return "OR"
	case InstructionType_XOR:
		return "XOR"
	case InstructionType_NOT:
		return "NOT"
	case InstructionType_SHL:
		return "SHL"
	case InstructionType_SHR:
		return "SHR"
	case InstructionType_CMP:
		return "CMP"
	case InstructionType_JMP:
		return "JMP"
	case InstructionType_JG:
		return "JG"
	case InstructionType_JL:
		return "JL"
	case InstructionType_JE:
		return "JE"
	case InstructionType_JNE:
		return "JNE"
	case InstructionType_JGE:
		return "JGE"
	case InstructionType_JLE:
		return "JLE"
	case InstructionType_LOAD:
		return "LOAD"
	case InstructionType_LOADH:
		return "LOADH"
	case InstructionType_LOADB:
		return "LOADB"
	case InstructionType_STORE:
		return "STORE"
	case InstructionType_STOREH:
		return "STOREH"
	case InstructionType_STOREB:
		return "STOREB"
	case InstructionType_MOV:
		return "MOV"
	case InstructionType_MOVH:
		return "MOVH"
	case InstructionType_MOVB:
		return "MOVB"
	case InstructionType_PUSH:
		return "PUSH"
	case InstructionType_POP:
		return "POP"
	case InstructionType_CALL:
		return "CALL"
	case InstructionType_RET:
		return "RET"
	case InstructionType_SYSCALL:
		return "SYSCALL"
	}
	return "UNKNOWN"
}

var Instructions = map[string]InstructionType{
	"NOP":     InstructionType_NOP,
	"ADD":     InstructionType_ADD,
	"SUB":     InstructionType_SUB,
	"MUL":     InstructionType_MUL,
	"DIV":     InstructionType_DIV,
	"AND":     InstructionType_AND,
	"OR":      InstructionType_OR,
	"XOR":     InstructionType_XOR,
	"NOT":     InstructionType_NOT,
	"SHL":     InstructionType_SHL,
	"SHR":     InstructionType_SHR,
	"CMP":     InstructionType_CMP,
	"JMP":     InstructionType_JMP,
	"JG":      InstructionType_JG,
	"JL":      InstructionType_JL,
	"JE":      InstructionType_JE,
	"JNE":     InstructionType_JNE,
	"JGE":     InstructionType_JGE,
	"JLE":     InstructionType_JLE,
	"LOAD":    InstructionType_LOAD,
	"LOADH":   InstructionType_LOADH,
	"LOADB":   InstructionType_LOADB,
	"STORE":   InstructionType_STORE,
	"STOREH":  InstructionType_STOREH,
	"STOREB":  InstructionType_STOREB,
	"MOV":     InstructionType_MOV,
	"MOVH":    InstructionType_MOVH,
	"MOVB":    InstructionType_MOVB,
	"PUSH":    InstructionType_PUSH,
	"POP":     InstructionType_POP,
	"CALL":    InstructionType_CALL,
	"RET":     InstructionType_RET,
	"SYSCALL": InstructionType_SYSCALL,
}
