package lvm2

type VM struct {
	Memory *Memory

	// # Registers
	//
	// ## 64-bit General Purpose Registers
	//
	// R0-R31
	//
	// ## SYSCALL registers
	//
	// SYS32-SYS63
	//
	// ## VM registers
	//
	// Program Counter (PC) (Register ID: 64)
	// Stack Pointer (SP)   (Register ID: 65)
	// Stack Base (SB)      (Register ID: 66)
	Registers [32 + 32 + 3]uint64
}

const Op0TypeMask = 0b11000000
const Op1TypeMask = 0b00110000
const Op2TypeMask = 0b00001100

/*
Opcode Format:

| 0   | 1   | 2   | 3   | 4   | 5   | 6   | 7   |
|-----|-----|-----|-----|-----|-----|-----|-----|
| Instruction Type                              |
|-----|-----|-----|-----|-----|-----|-----|-----|
| Op0 Type  | Op1 Type  | Op2 Type  | Reserved  |
|-----|-----|-----|-----|-----|-----|-----|-----|
|                                               |
|                   Op0 Value                   |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|-----|-----|-----|-----|-----|-----|-----|-----|
|                                               |
|                   Op1 Value                   |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|-----|-----|-----|-----|-----|-----|-----|-----|
|                                               |
|                   Op2 Value                   |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|-----|-----|-----|-----|-----|-----|-----|-----|
*/
