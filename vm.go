package lvm2

import (
	"encoding/binary"
	"fmt"
)

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
	v.Registers[64] = pc
}

func (v *VM) parseOpcode() (instructionType InstructionType, op0Type OpType, op1Type OpType, op2Type OpType, op0Value uint64, op1Value uint64, op2Value uint64, err error) {
	var buffer [instructionBytecodeSize]byte
	_, err = v.Memory.ReadAt(v.Registers[64], buffer[:])
	if err != nil {
		return
	}
	v.Registers[64] += instructionBytecodeSize

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
			v.Registers[64] = op0Value
		case InstructionType_JG:
			// JG
			if int64(op0Value) > 0 {
				v.Registers[64] = op1Value
			}
		case InstructionType_JL:
			// JL
			if int64(op0Value) < 0 {
				v.Registers[64] = op1Value
			}
		case InstructionType_JE:
			// JE
			if int64(op0Value) == 0 {
				v.Registers[64] = op1Value
			}
		case InstructionType_JNE:
			// JNE
			if int64(op0Value) != 0 {
				v.Registers[64] = op1Value
			}
		case InstructionType_JGE:
			// JGE
			if int64(op0Value) >= 0 {
				v.Registers[64] = op1Value
			}
		case InstructionType_JLE:
			// JLE
			if int64(op0Value) <= 0 {
				v.Registers[64] = op1Value
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
			v.Registers[65] -= 8
			var buffer [8]byte
			binary.LittleEndian.PutUint64(buffer[:], op0Value)
			_, err = v.Memory.WriteAt(v.Registers[65], buffer[:])
			if err != nil {
				return 0, err
			}
		case InstructionType_POP:
			// POP
			var buffer [8]byte
			_, err = v.Memory.ReadAt(v.Registers[65], buffer[:])
			v.Registers[65] += 8
			if err != nil {
				return 0, err
			}
			v.Registers[op0Value] = binary.LittleEndian.Uint64(buffer[:])

		case InstructionType_CALL:
			// CALL
			v.Registers[65] -= 8
			var buffer [8]byte
			binary.LittleEndian.PutUint64(buffer[:], v.Registers[64])
			_, err = v.Memory.WriteAt(v.Registers[65], buffer[:])
			if err != nil {
				return 0, err
			}
			v.Registers[64] = op0Value
		case InstructionType_RET:
			// RET
			var buffer [8]byte
			_, err = v.Memory.ReadAt(v.Registers[65], buffer[:])
			v.Registers[65] += 8
			if err != nil {
				return 0, err
			}
			v.Registers[64] = binary.LittleEndian.Uint64(buffer[:])

		default:
			return 0, ErrInvalidInstruction
		}
	}

	return 0, nil
}
