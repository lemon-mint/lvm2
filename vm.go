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

func (v *VM) SetProgram(p []byte) {
	v.Memory.Reset()
	v.Memory.SetProgram(p)
}

func (v *VM) SetProgramCounter(pc uint64) {
	v.Registers[64] = pc
}

func (v *VM) parseOpcode() (instructionType uint8, op0Type uint8, op1Type uint8, op2Type uint8, op0Value uint64, op1Value uint64, op2Value uint64, err error) {
	var buffer [1 + 1 + 8*3]byte
	_, err = v.Memory.ReadAt(v.Registers[64], buffer[:])
	if err != nil {
		return
	}

	vs := InstructionOpcode(buffer[:])
	instructionType = vs.InstructionType()
	typeinfo := vs.OperandType()
	op0Type = (typeinfo & 0b11000000) >> 6
	op1Type = (typeinfo & 0b00110000) >> 4
	op2Type = (typeinfo & 0b00001100) >> 2
	op0Value = vs.Operand0()
	op1Value = vs.Operand1()
	op2Value = vs.Operand2()

	return
}

func (v *VM) Run() (uint64, error) {
	for {
		break
	}

	return 0, nil
}
