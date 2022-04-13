package lvm2

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

type _ = strings.Builder
type _ = unsafe.Pointer

var _ = math.Float32frombits
var _ = math.Float64frombits
var _ = strconv.FormatInt
var _ = strconv.FormatUint
var _ = strconv.FormatFloat
var _ = fmt.Sprint

type InstructionOpcode []byte

func (s InstructionOpcode) InstructionType() uint8 {
	_ = s[0]
	var __v uint8 = uint8(s[0])
	return uint8(__v)
}

func (s InstructionOpcode) OperandType() uint8 {
	_ = s[1]
	var __v uint8 = uint8(s[1])
	return uint8(__v)
}

func (s InstructionOpcode) Operand0() uint64 {
	_ = s[9]
	var __v uint64 = uint64(s[2]) |
		uint64(s[3])<<8 |
		uint64(s[4])<<16 |
		uint64(s[5])<<24 |
		uint64(s[6])<<32 |
		uint64(s[7])<<40 |
		uint64(s[8])<<48 |
		uint64(s[9])<<56
	return uint64(__v)
}

func (s InstructionOpcode) Operand1() uint64 {
	_ = s[17]
	var __v uint64 = uint64(s[10]) |
		uint64(s[11])<<8 |
		uint64(s[12])<<16 |
		uint64(s[13])<<24 |
		uint64(s[14])<<32 |
		uint64(s[15])<<40 |
		uint64(s[16])<<48 |
		uint64(s[17])<<56
	return uint64(__v)
}

func (s InstructionOpcode) Operand2() uint64 {
	_ = s[25]
	var __v uint64 = uint64(s[18]) |
		uint64(s[19])<<8 |
		uint64(s[20])<<16 |
		uint64(s[21])<<24 |
		uint64(s[22])<<32 |
		uint64(s[23])<<40 |
		uint64(s[24])<<48 |
		uint64(s[25])<<56
	return uint64(__v)
}

func (s InstructionOpcode) Vstruct_Validate() bool {
	return len(s) >= 26
}

func (s InstructionOpcode) String() string {
	if !s.Vstruct_Validate() {
		return "InstructionOpcode (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("InstructionOpcode {")
	__b.WriteString("InstructionType: ")
	__b.WriteString(strconv.FormatUint(uint64(s.InstructionType()), 10))
	__b.WriteString(", ")
	__b.WriteString("OperandType: ")
	__b.WriteString(strconv.FormatUint(uint64(s.OperandType()), 10))
	__b.WriteString(", ")
	__b.WriteString("Operand0: ")
	__b.WriteString(strconv.FormatUint(uint64(s.Operand0()), 10))
	__b.WriteString(", ")
	__b.WriteString("Operand1: ")
	__b.WriteString(strconv.FormatUint(uint64(s.Operand1()), 10))
	__b.WriteString(", ")
	__b.WriteString("Operand2: ")
	__b.WriteString(strconv.FormatUint(uint64(s.Operand2()), 10))
	__b.WriteString("}")
	return __b.String()
}

func Serialize_InstructionOpcode(dst InstructionOpcode, InstructionType uint8, OperandType uint8, Operand0 uint64, Operand1 uint64, Operand2 uint64) InstructionOpcode {
	_ = dst[25]
	dst[0] = byte(InstructionType)
	dst[1] = byte(OperandType)
	dst[2] = byte(Operand0)
	dst[3] = byte(Operand0 >> 8)
	dst[4] = byte(Operand0 >> 16)
	dst[5] = byte(Operand0 >> 24)
	dst[6] = byte(Operand0 >> 32)
	dst[7] = byte(Operand0 >> 40)
	dst[8] = byte(Operand0 >> 48)
	dst[9] = byte(Operand0 >> 56)
	dst[10] = byte(Operand1)
	dst[11] = byte(Operand1 >> 8)
	dst[12] = byte(Operand1 >> 16)
	dst[13] = byte(Operand1 >> 24)
	dst[14] = byte(Operand1 >> 32)
	dst[15] = byte(Operand1 >> 40)
	dst[16] = byte(Operand1 >> 48)
	dst[17] = byte(Operand1 >> 56)
	dst[18] = byte(Operand2)
	dst[19] = byte(Operand2 >> 8)
	dst[20] = byte(Operand2 >> 16)
	dst[21] = byte(Operand2 >> 24)
	dst[22] = byte(Operand2 >> 32)
	dst[23] = byte(Operand2 >> 40)
	dst[24] = byte(Operand2 >> 48)
	dst[25] = byte(Operand2 >> 56)

	return dst
}

func New_InstructionOpcode(InstructionType uint8, OperandType uint8, Operand0 uint64, Operand1 uint64, Operand2 uint64) InstructionOpcode {
	var __vstruct__size = 26
	var __vstruct__buf = make(InstructionOpcode, __vstruct__size)
	__vstruct__buf = Serialize_InstructionOpcode(__vstruct__buf, InstructionType, OperandType, Operand0, Operand1, Operand2)
	return __vstruct__buf
}
