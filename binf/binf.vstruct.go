package binf

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

type EncodingType uint8

const (
	EncodingType_RAW  EncodingType = 0
	EncodingType_GZIP EncodingType = 1
)

func (e EncodingType) String() string {
	switch e {
	case EncodingType_RAW:
		return "RAW"
	case EncodingType_GZIP:
		return "GZIP"
	}
	return ""
}

func (e EncodingType) Match(
	onRAW func(),
	onGZIP func(),
) {
	switch e {
	case EncodingType_RAW:
		onRAW()
	case EncodingType_GZIP:
		onGZIP()
	}
}

type Header []byte

func (s Header) Version() uint8 {
	_ = s[0]
	var __v uint8 = uint8(s[0])
	return uint8(__v)
}

func (s Header) EntryPoint() uint64 {
	_ = s[8]
	var __v uint64 = uint64(s[1]) |
		uint64(s[2])<<8 |
		uint64(s[3])<<16 |
		uint64(s[4])<<24 |
		uint64(s[5])<<32 |
		uint64(s[6])<<40 |
		uint64(s[7])<<48 |
		uint64(s[8])<<56
	return uint64(__v)
}

func (s Header) Vstruct_Validate() bool {
	return len(s) >= 9
}

func (s Header) String() string {
	if !s.Vstruct_Validate() {
		return "Header (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("Header {")
	__b.WriteString("Version: ")
	__b.WriteString(strconv.FormatUint(uint64(s.Version()), 10))
	__b.WriteString(", ")
	__b.WriteString("EntryPoint: ")
	__b.WriteString(strconv.FormatUint(uint64(s.EntryPoint()), 10))
	__b.WriteString("}")
	return __b.String()
}

type Program []byte

func (s Program) Encoding() EncodingType {
	return EncodingType(s[0])
}

func (s Program) Header() Header {
	return Header(s[1:10])
}

func (s Program) Code() []byte {
	_ = s[17]
	var __off0 uint64 = 18
	var __off1 uint64 = uint64(s[10]) |
		uint64(s[11])<<8 |
		uint64(s[12])<<16 |
		uint64(s[13])<<24 |
		uint64(s[14])<<32 |
		uint64(s[15])<<40 |
		uint64(s[16])<<48 |
		uint64(s[17])<<56
	return []byte(s[__off0:__off1])
}

func (s Program) Vstruct_Validate() bool {
	if len(s) < 18 {
		return false
	}

	_ = s[17]

	var __off0 uint64 = 18
	var __off1 uint64 = uint64(s[10]) |
		uint64(s[11])<<8 |
		uint64(s[12])<<16 |
		uint64(s[13])<<24 |
		uint64(s[14])<<32 |
		uint64(s[15])<<40 |
		uint64(s[16])<<48 |
		uint64(s[17])<<56
	var __off2 uint64 = uint64(len(s))
	return __off0 <= __off1 && __off1 <= __off2
}

func (s Program) String() string {
	if !s.Vstruct_Validate() {
		return "Program (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("Program {")
	__b.WriteString("Encoding: ")
	__b.WriteString(s.Encoding().String())
	__b.WriteString(", ")
	__b.WriteString("Header: ")
	__b.WriteString(s.Header().String())
	__b.WriteString(", ")
	__b.WriteString("Code: ")
	__b.WriteString(fmt.Sprint(s.Code()))
	__b.WriteString("}")
	return __b.String()
}

func Serialize_Header(dst Header, Version uint8, EntryPoint uint64) Header {
	_ = dst[8]
	dst[0] = byte(Version)
	dst[1] = byte(EntryPoint)
	dst[2] = byte(EntryPoint >> 8)
	dst[3] = byte(EntryPoint >> 16)
	dst[4] = byte(EntryPoint >> 24)
	dst[5] = byte(EntryPoint >> 32)
	dst[6] = byte(EntryPoint >> 40)
	dst[7] = byte(EntryPoint >> 48)
	dst[8] = byte(EntryPoint >> 56)

	return dst
}

func New_Header(Version uint8, EntryPoint uint64) Header {
	var __vstruct__size = 9
	var __vstruct__buf = make(Header, __vstruct__size)
	__vstruct__buf = Serialize_Header(__vstruct__buf, Version, EntryPoint)
	return __vstruct__buf
}

func Serialize_Program(dst Program, Encoding EncodingType, Header Header, Code []byte) Program {
	_ = dst[17]
	dst[0] = byte(Encoding)
	copy(dst[1:10], Header)

	var __index = uint64(18)
	__tmp_2 := uint64(len(Code)) + __index
	dst[10] = byte(__tmp_2)
	dst[11] = byte(__tmp_2 >> 8)
	dst[12] = byte(__tmp_2 >> 16)
	dst[13] = byte(__tmp_2 >> 24)
	dst[14] = byte(__tmp_2 >> 32)
	dst[15] = byte(__tmp_2 >> 40)
	dst[16] = byte(__tmp_2 >> 48)
	dst[17] = byte(__tmp_2 >> 56)
	copy(dst[__index:__tmp_2], Code)
	return dst
}

func New_Program(Encoding EncodingType, Header Header, Code []byte) Program {
	var __vstruct__size = 18 + len(Code)
	var __vstruct__buf = make(Program, __vstruct__size)
	__vstruct__buf = Serialize_Program(__vstruct__buf, Encoding, Header, Code)
	return __vstruct__buf
}
