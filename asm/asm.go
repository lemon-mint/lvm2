package asm

import (
	"strconv"
	"strings"

	"github.com/lemon-mint/lvm2"
)

type CodeType byte

const (
	CODE_INST CodeType = iota
	CODE_DATA
	CODE_LABEL
)

type Code struct {
	Type CodeType

	Instruction lvm2.InstructionType
	Operands    []Operand

	Label string
	Data  []byte
}

type Operand struct {
	Type        OperandType
	Value       uint64
	Value_Label string
}

type OperandType byte

const (
	OperandType_None OperandType = iota
	OperandType_RegisterValue
	OperandType_ConstantValue
	OperandType_Label
)

func INST(t lvm2.InstructionType, ops ...Operand) Code {
	return Code{
		Type:        CODE_INST,
		Instruction: t,
		Operands:    ops,
	}
}

func LABEL(l string) Code {
	return Code{
		Type:  CODE_LABEL,
		Label: l,
	}
}

func DATA(b []byte) Code {
	return Code{
		Type: CODE_DATA,
		Data: b,
	}
}

func OPCONST(v uint64) Operand {
	return Operand{
		Type:  OperandType_ConstantValue,
		Value: v,
	}
}

func OPREG(r uint64) Operand {
	return Operand{
		Type:  OperandType_RegisterValue,
		Value: uint64(r),
	}
}

func OPLABEL(l string) Operand {
	return Operand{
		Type:        OperandType_Label,
		Value_Label: l,
	}
}

func (v Operand) String() string {
	var sb strings.Builder
	switch v.Type {
	case OperandType_None:
		sb.WriteString("NONE")
	case OperandType_RegisterValue:
		sb.WriteString("REGISTER(0x")
		sb.WriteString(strconv.FormatUint(v.Value, 16))
		sb.WriteString(")")
	case OperandType_ConstantValue:
		sb.WriteString("CONST(0x")
		sb.WriteString(strconv.FormatUint(v.Value, 16))
		sb.WriteString(")")
	case OperandType_Label:
		sb.WriteString("LABEL(\"")
		sb.WriteString(strconv.Quote(v.Value_Label))
		sb.WriteString("\")")
	}
	return sb.String()
}

func (v Code) String() string {
	var sb strings.Builder
	switch v.Type {
	case CODE_INST:
		sb.WriteString(v.Instruction.String())
		for _, op := range v.Operands {
			sb.WriteByte(' ')
			sb.WriteString(op.String())
		}
	case CODE_DATA:
		sb.WriteString("DATA [")
		for i, b := range v.Data {
			if i != 0 {
				sb.WriteString(", ")
			}

			sb.WriteString("0x")
			sb.WriteString(strconv.FormatUint(uint64(b), 16))
		}
		sb.WriteString("]")
	case CODE_LABEL:
		sb.WriteString("LABEL \"")
		sb.WriteString(strconv.Quote(v.Label))
		sb.WriteString("\"")
	}
	return sb.String()
}

type Encoder struct {
	Dst []byte
	PC  uint64

	Labels map[string]uint64
}

func NewEncoder() *Encoder {
	return &Encoder{
		Labels: make(map[string]uint64),
	}
}

func (e *Encoder) Encode(c Code) uint64 {
	switch c.Type {
	case CODE_INST:
		return e.encodeInstruction(c)
	case CODE_DATA:
		return e.encodeData(c)
	case CODE_LABEL:
		return e.encodeLabel(c)
	}
	return e.PC
}

func (e *Encoder) encodeInstruction(c Code) (out uint64) {
	out = e.PC

	var ops [3]uint64
	var opt byte
	for i, op := range c.Operands {
		switch op.Type {
		case OperandType_RegisterValue:
			opt |= byte(lvm2.OpTypeRegister) << (6 - i*2)
			ops[i] = op.Value
		case OperandType_ConstantValue:
			opt |= byte(lvm2.OpTypeConstant) << (6 - i*2)
			ops[i] = op.Value
		case OperandType_Label:
			opt |= byte(lvm2.OpTypeConstant) << (6 - i*2)
			ops[i] = e.Labels[op.Value_Label]
		}
	}
	opcode := lvm2.New_InstructionOpcode(uint8(c.Instruction), opt, ops[0], ops[1], ops[2])
	e.Dst = append(e.Dst, opcode...)
	e.PC += uint64(len(opcode))
	return
}

func (e *Encoder) encodeData(c Code) (out uint64) {
	out = e.PC
	e.Dst = append(e.Dst, c.Data...)
	e.PC += uint64(len(c.Data))
	return
}

func (e *Encoder) encodeLabel(c Code) (out uint64) {
	e.Labels[c.Label] = e.PC
	e.Dst = append(e.Dst, lvm2.New_InstructionOpcode(
		uint8(lvm2.InstructionType_NOP),
		0,
		0,
		0,
		0,
	)...)
	return e.PC
}

func (e *Encoder) Bytes() []byte {
	return e.Dst
}
