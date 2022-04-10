package lvm2

type VM struct {
	Memory *Memory

	Registers [64 + 32]uint64
}

const Op0TypeMask = 0b11000000
const Op1TypeMask = 0b00110000
const Op2TypeMask = 0b00001100
const Op3TypeMask = 0b00000011
