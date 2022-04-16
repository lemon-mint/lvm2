package lvm2

import (
	"errors"

	"github.com/lemon-mint/lvm2/errs"
)

const (
	SYS_READ  = 0
	SYS_WRITE = 1
	SYS_OPEN  = 2
	SYS_CLOSE = 3
	SYS_EXIT  = 60
)

type SYSCALLFunc func(vm *VM, R0, R1, R2 uint64) (errno uint64, err error)

var syscall_Function_Table map[uint64]SYSCALLFunc = make(map[uint64]SYSCALLFunc)

func _syscall_write(vm *VM, _, _, _ uint64) (errno uint64, err error) {
	// func Write(fd uintptr, p uintptr, n uint64) (written uint64, errno uint64)
	// SYS32[in]: fd
	// SYS33[in]: p
	// SYS34[in]: n
	// SYS35[out]: written

	fd := vm.Registers[REGISTER_SYS32]
	p := vm.Registers[REGISTER_SYS33]
	n := vm.Registers[REGISTER_SYS34]

	file, ok := vm.Files[fd]
	if !ok {
		return errs.EINVALIDFD.Errno(), nil
	}

	var buffer []byte = make([]byte, n)
	_, err = vm.Memory.ReadAt(p, buffer)
	if err != nil {
		return 1, err
	}

	written, err := file.Write(buffer)
	if err != nil {
		return errs.EFILEWRITE.Errno(), nil
	}

	vm.Registers[REGISTER_SYS35] = uint64(written)
	return 0, nil
}

func _syscall_read(vm *VM, _, _, _ uint64) (errno uint64, err error) {
	// func Read(fd uintptr, p uintptr, n uint64) (read uint64, errno uint64)
	// SYS32[in]: fd
	// SYS33[in]: p
	// SYS34[in]: n
	// SYS35[out]: read

	fd := vm.Registers[REGISTER_SYS32]
	p := vm.Registers[REGISTER_SYS33]
	n := vm.Registers[REGISTER_SYS34]

	file, ok := vm.Files[fd]
	if !ok {
		return errs.EINVALIDFD.Errno(), nil
	}

	var buffer []byte = make([]byte, n)
	read, err := file.Read(buffer)
	if err != nil {
		return errs.EFILEREAD.Errno(), nil
	}

	_, err = vm.Memory.WriteAt(p, buffer)
	if err != nil {
		return 1, err
	}

	vm.Registers[REGISTER_SYS35] = uint64(read)
	return 0, nil
}

func _syscall_exit(vm *VM, _, _, _ uint64) (errno uint64, err error) {
	// func Exit(code uint64)
	// SYS32[in]: code

	code := vm.Registers[REGISTER_SYS32]
	return code, ErrExited
}

var ErrExited = errors.New("exited")

var _ = func() bool {
	syscall_Function_Table[SYS_WRITE] = _syscall_write
	syscall_Function_Table[SYS_READ] = _syscall_read
	syscall_Function_Table[SYS_EXIT] = _syscall_exit
	return true
}()
