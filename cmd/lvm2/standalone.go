package main

import (
	"encoding/binary"
	"os"

	"github.com/lemon-mint/lvm2"
	"github.com/lemon-mint/lvm2/binf"
)

func main() {
	myexec := os.Args[0]
	f, err := os.Open(myexec)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	exesize := fi.Size()
	var hbuf [16]byte
	f.ReadAt(hbuf[:], exesize-16)
	var lvm2p binf.Program

	if string(hbuf[:8]) == "@%LVM2%\n" {
		// Magic header found, this is a lvm2 file
		var size uint64
		size = binary.LittleEndian.Uint64(hbuf[8:16])
		if uint64(exesize) <= size {
			panic("Invalid lvm2 file")
		}

		prog := make([]byte, size)
		f.ReadAt(prog, exesize-16-int64(size))

		lvm2p = binf.Program(prog)
		if !lvm2p.Vstruct_Validate() {
			panic("Invalid lvm2 file")
		}
	} else {
		// Use argv[1] as input file
		if len(os.Args) < 2 {
			panic("No input file specified")
		}

		prog, err := os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		lvm2p = binf.Program(prog)
		if !lvm2p.Vstruct_Validate() {
			panic("Invalid lvm2 file")
		}
	}

	vm := lvm2.VM{
		Memory: lvm2.NewMemory(),
		Files: map[uint64]lvm2.VMFile{
			0: os.Stdin,
			1: os.Stdout,
			2: os.Stderr,
		},
		FileCounter: 3,
	}
	vm.Registers[lvm2.REGISTER_SP] = vm.Memory.MaxAddress
	vm.Registers[lvm2.REGISTER_SB] = vm.Memory.MaxAddress

	if lvm2p.Encoding() == binf.EncodingType_RAW {
		vm.Memory.SetProgram(lvm2p.Code())
	} else {
		panic("Unsupported encoding")
	}

	vm.SetProgramCounter(lvm2p.Header().EntryPoint())

	ret, err := vm.Run()
	if err != nil {
		panic(err)
	}

	os.Exit(int(ret))
}
