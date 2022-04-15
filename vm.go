package lvm2

import (
	"encoding/binary"
	"fmt"
)

type VMFile interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
}

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

	// File Descriptor Table
	Files []VMFile
}

const (
	REGISTER_R0    = 0
	REGISTER_R1    = 1
	REGISTER_R2    = 2
	REGISTER_R3    = 3
	REGISTER_R4    = 4
	REGISTER_R5    = 5
	REGISTER_R6    = 6
	REGISTER_R7    = 7
	REGISTER_R8    = 8
	REGISTER_R9    = 9
	REGISTER_R10   = 10
	REGISTER_R11   = 11
	REGISTER_R12   = 12
	REGISTER_R13   = 13
	REGISTER_R14   = 14
	REGISTER_R15   = 15
	REGISTER_R16   = 16
	REGISTER_R17   = 17
	REGISTER_R18   = 18
	REGISTER_R19   = 19
	REGISTER_R20   = 20
	REGISTER_R21   = 21
	REGISTER_R22   = 22
	REGISTER_R23   = 23
	REGISTER_R24   = 24
	REGISTER_R25   = 25
	REGISTER_R26   = 26
	REGISTER_R27   = 27
	REGISTER_R28   = 28
	REGISTER_R29   = 29
	REGISTER_R30   = 30
	REGISTER_R31   = 31
	REGISTER_SYS32 = 32
	REGISTER_SYS33 = 33
	REGISTER_SYS34 = 34
	REGISTER_SYS35 = 35
	REGISTER_SYS36 = 36
	REGISTER_SYS37 = 37
	REGISTER_SYS38 = 38
	REGISTER_SYS39 = 39
	REGISTER_SYS40 = 40
	REGISTER_SYS41 = 41
	REGISTER_SYS42 = 42
	REGISTER_SYS43 = 43
	REGISTER_SYS44 = 44
	REGISTER_SYS45 = 45
	REGISTER_SYS46 = 46
	REGISTER_SYS47 = 47
	REGISTER_SYS48 = 48
	REGISTER_SYS49 = 49
	REGISTER_SYS50 = 50
	REGISTER_SYS51 = 51
	REGISTER_SYS52 = 52
	REGISTER_SYS53 = 53
	REGISTER_SYS54 = 54
	REGISTER_SYS55 = 55
	REGISTER_SYS56 = 56
	REGISTER_SYS57 = 57
	REGISTER_SYS58 = 58
	REGISTER_SYS59 = 59
	REGISTER_SYS60 = 60
	REGISTER_SYS61 = 61
	REGISTER_SYS62 = 62
	REGISTER_SYS63 = 63
	REGISTER_PC    = 64
	REGISTER_SP    = 65
	REGISTER_SB    = 66
)

/*
Bytecode Format:

|-----|-----|-----|-----|-----|-----|-----|-----|
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
|-----|-----|-----|-----|-----|-----|-----|-----|
|                                               |
|                   Op2 Value                   |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|                                               |
|-----|-----|-----|-----|-----|-----|-----|-----|

Opx Type:

0b00: None
0b01: Register
0b10: Constant
0b11: Reserved

*/

const instructionBytecodeSize = 1 + 1 + 8*3

type OpType byte

const (
	OpTypeNone     OpType = 0b00
	OpTypeRegister OpType = 0b01
	OpTypeConstant OpType = 0b10
	OpTypeReserved OpType = 0b11
)

func (v *VM) SetProgram(p []byte) {
	v.Memory.Reset()
	v.Memory.SetProgram(p)
}

func (v *VM) SetProgramCounter(pc uint64) {
	v.Registers[REGISTER_PC] = pc
}

func (v *VM) parseOpcode() (instructionType InstructionType, op0Type OpType, op1Type OpType, op2Type OpType, op0Value uint64, op1Value uint64, op2Value uint64, err error) {
	var buffer [instructionBytecodeSize]byte
	_, err = v.Memory.ReadAt(v.Registers[REGISTER_PC], buffer[:])
	if err != nil {
		return
	}
	v.Registers[REGISTER_PC] += instructionBytecodeSize

	vs := InstructionOpcode(buffer[:])
	instructionType = InstructionType(vs.InstructionType())
	typeinfo := vs.OperandType()
	op0Type = OpType((typeinfo & 0b11000000) >> 6)
	op1Type = OpType((typeinfo & 0b00110000) >> 4)
	op2Type = OpType((typeinfo & 0b00001100) >> 2)
	op0Value = vs.Operand0()
	op1Value = vs.Operand1()
	op2Value = vs.Operand2()

	return
}

var ErrInvalidInstruction = fmt.Errorf("invalid instruction")

func (v *VM) Run() (uint64, error) {
	for {
		instructionType, op0Type, op1Type, op2Type, op0Value, op1Value, op2Value, err := v.parseOpcode()
		if err != nil {
			return 0, err
		}

		if op0Type == OpTypeRegister {
			if op0Value >= uint64(len(v.Registers)) {
				return 0, fmt.Errorf("invalid register: %d", op0Value)
			}
			op0Value = v.Registers[op0Value]
		}
		if op1Type == OpTypeRegister {
			if op1Value >= uint64(len(v.Registers)) {
				return 0, fmt.Errorf("invalid register: %d", op1Value)
			}
			op1Value = v.Registers[op1Value]
		}
		if op2Type == OpTypeRegister {
			if op2Value >= uint64(len(v.Registers)) {
				return 0, fmt.Errorf("invalid register: %d", op2Value)
			}
			op2Value = v.Registers[op2Value]
		}

		switch instructionType {
		case InstructionType_NOP:
			// NOP
		case InstructionType_ADD:
			// ADD
			v.Registers[op0Value] = op1Value + op2Value
		case InstructionType_SUB:
			// SUB
			v.Registers[op0Value] = op1Value - op2Value
		case InstructionType_MUL:
			// MUL
			v.Registers[op0Value] = op1Value * op2Value
		case InstructionType_DIV:
			// DIV
			v.Registers[op0Value] = op1Value / op2Value
		case InstructionType_MOD:
			// MOD
			v.Registers[op0Value] = op1Value % op2Value

		case InstructionType_AND:
			// AND
			v.Registers[op0Value] = op1Value & op2Value
		case InstructionType_OR:
			// OR
			v.Registers[op0Value] = op1Value | op2Value
		case InstructionType_XOR:
			// XOR
			v.Registers[op0Value] = op1Value ^ op2Value
		case InstructionType_NOT:
			// NOT
			v.Registers[op0Value] = ^op1Value

		case InstructionType_SHL:
			// SHL
			v.Registers[op0Value] = op1Value << op2Value
		case InstructionType_SHR:
			// SHR
			v.Registers[op0Value] = op1Value >> op2Value

		case InstructionType_CMP:
			// CMP
			diff := op1Value - op2Value
			if diff < 0 {
				v.Registers[op0Value] = ^uint64(0)
			} else if diff > 0 {
				v.Registers[op0Value] = 1
			} else {
				v.Registers[op0Value] = 0
			}

		case InstructionType_JMP:
			// JMP
			v.Registers[REGISTER_PC] = op0Value
		case InstructionType_JG:
			// JG
			if int64(op0Value) > 0 {
				v.Registers[REGISTER_PC] = op1Value
			}
		case InstructionType_JL:
			// JL
			if int64(op0Value) < 0 {
				v.Registers[REGISTER_PC] = op1Value
			}
		case InstructionType_JE:
			// JE
			if int64(op0Value) == 0 {
				v.Registers[REGISTER_PC] = op1Value
			}
		case InstructionType_JNE:
			// JNE
			if int64(op0Value) != 0 {
				v.Registers[REGISTER_PC] = op1Value
			}
		case InstructionType_JGE:
			// JGE
			if int64(op0Value) >= 0 {
				v.Registers[REGISTER_PC] = op1Value
			}
		case InstructionType_JLE:
			// JLE
			if int64(op0Value) <= 0 {
				v.Registers[REGISTER_PC] = op1Value
			}

		case InstructionType_LOAD:
			// LOAD
			var buffer [8]byte
			_, err = v.Memory.ReadAt(op1Value+op2Value, buffer[:])
			if err != nil {
				return 0, err
			}
			v.Registers[op0Value] = binary.LittleEndian.Uint64(buffer[:])
		case InstructionType_LOADH:
			// LOADH
			var buffer [4]byte
			_, err = v.Memory.ReadAt(op1Value+op2Value, buffer[:])
			if err != nil {
				return 0, err
			}
			v.Registers[op0Value] = uint64(binary.LittleEndian.Uint32(buffer[:]))
		case InstructionType_LOADB:
			// LOADB
			var buffer [1]byte
			_, err = v.Memory.ReadAt(op1Value+op2Value, buffer[:])
			if err != nil {
				return 0, err
			}
			v.Registers[op0Value] = uint64(buffer[0])

		case InstructionType_STORE:
			// STORE
			var buffer [8]byte
			binary.LittleEndian.PutUint64(buffer[:], op0Value)
			_, err = v.Memory.WriteAt(op1Value+op2Value, buffer[:])
			if err != nil {
				return 0, err
			}
		case InstructionType_STOREH:
			// STOREH
			var buffer [4]byte
			binary.LittleEndian.PutUint32(buffer[:], uint32(op0Value))
			_, err = v.Memory.WriteAt(op1Value+op2Value, buffer[:])
			if err != nil {
				return 0, err
			}
		case InstructionType_STOREB:
			// STOREB
			var buffer [1]byte
			buffer[0] = byte(op0Value)
			_, err = v.Memory.WriteAt(op1Value+op2Value, buffer[:])
			if err != nil {
				return 0, err
			}

		case InstructionType_MOV:
			// MOV
			v.Registers[op0Value] = op1Value
		case InstructionType_MOVH:
			// MOVH
			v.Registers[op0Value] = uint64(uint32(op1Value))
		case InstructionType_MOVB:
			// MOVB
			v.Registers[op0Value] = uint64(uint8(op1Value))

		case InstructionType_PUSH:
			// PUSH
			v.Registers[REGISTER_SP] -= 8
			var buffer [8]byte
			binary.LittleEndian.PutUint64(buffer[:], op0Value)
			_, err = v.Memory.WriteAt(v.Registers[REGISTER_SP], buffer[:])
			if err != nil {
				return 0, err
			}
		case InstructionType_POP:
			// POP
			var buffer [8]byte
			_, err = v.Memory.ReadAt(v.Registers[REGISTER_SP], buffer[:])
			v.Registers[REGISTER_SP] += 8
			if err != nil {
				return 0, err
			}
			v.Registers[op0Value] = binary.LittleEndian.Uint64(buffer[:])

		case InstructionType_CALL:
			// CALL
			v.Registers[REGISTER_SP] -= 8
			var buffer [8]byte
			binary.LittleEndian.PutUint64(buffer[:], v.Registers[REGISTER_PC])
			_, err = v.Memory.WriteAt(v.Registers[REGISTER_SP], buffer[:])
			if err != nil {
				return 0, err
			}
			v.Registers[REGISTER_PC] = op0Value
		case InstructionType_RET:
			// RET
			var buffer [8]byte
			_, err = v.Memory.ReadAt(v.Registers[REGISTER_SP], buffer[:])
			v.Registers[REGISTER_SP] += 8
			if err != nil {
				return 0, err
			}
			v.Registers[REGISTER_PC] = binary.LittleEndian.Uint64(buffer[:])

		case InstructionType_SYSCALL:
			// SYSCALL
		default:
			return 0, ErrInvalidInstruction
		}
	}
}
