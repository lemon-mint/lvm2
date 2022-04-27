package main

import (
	"os"

	"github.com/lemon-mint/lvm2"
	"github.com/lemon-mint/lvm2/asm"
)

func main() {
	e := asm.NewEncoder()
	buf := e.Encode(asm.DATA(make([]byte, 20)))
	ep := e.Encode(asm.INST(lvm2.InstructionType_MOV, asm.OPCONST(lvm2.REGISTER_SYS32), asm.OPCONST(0)))                         // MOV %SYS32, 0
	e.Encode(asm.INST(lvm2.InstructionType_MOV, asm.OPCONST(lvm2.REGISTER_SYS33), asm.OPCONST(buf)))                             // MOV %SYS33, #buf
	e.Encode(asm.INST(lvm2.InstructionType_MOV, asm.OPCONST(lvm2.REGISTER_SYS34), asm.OPCONST(20)))                              // MOV %SYS34, 20
	e.Encode(asm.INST(lvm2.InstructionType_SYSCALL, asm.OPCONST(lvm2.REGISTER_R0), asm.OPCONST(lvm2.SYS_READ), asm.OPCONST(0)))  // SYSCALL %R0, 0, 0
	e.Encode(asm.INST(lvm2.InstructionType_MOV, asm.OPCONST(lvm2.REGISTER_SYS32), asm.OPCONST(1)))                               // MOV %SYS32, 1
	e.Encode(asm.INST(lvm2.InstructionType_MOV, asm.OPCONST(lvm2.REGISTER_SYS33), asm.OPCONST(buf)))                             // MOV %SYS33, #buf
	e.Encode(asm.INST(lvm2.InstructionType_MOV, asm.OPCONST(lvm2.REGISTER_SYS34), asm.OPREG(lvm2.REGISTER_SYS35)))               // MOV %SYS34, %SYS35
	e.Encode(asm.INST(lvm2.InstructionType_SYSCALL, asm.OPCONST(lvm2.REGISTER_R0), asm.OPCONST(lvm2.SYS_WRITE), asm.OPCONST(1))) // SYSCALL %R0, 1, 0
	e.Encode(asm.INST(lvm2.InstructionType_MOV, asm.OPCONST(lvm2.REGISTER_SYS32), asm.OPREG(lvm2.REGISTER_R0)))                  // MOV %SYS32, %R0
	e.Encode(asm.INST(lvm2.InstructionType_SYSCALL, asm.OPCONST(lvm2.REGISTER_R0), asm.OPCONST(60), asm.OPCONST(0)))             // SYSCALL %R0, 60, 0

	//fmt.Println("Entry point:", ep)
	//fmt.Println("Opcodes:")
	//fmt.Println(e.Bytes())

	vm := lvm2.VM{
		Memory: lvm2.NewMemory(),
		Files: map[uint64]lvm2.VMFile{
			0: os.Stdin,
			1: os.Stdout,
			2: os.Stderr,
		},
		FileCounter: 3,
	}
	vm.Memory.SetProgram(e.Bytes())

	// Set Entry Point
	vm.SetProgramCounter(ep)
	ret, err := vm.Run()
	if err != nil {
		panic(err)
	}
	os.Exit(int(ret))
}
