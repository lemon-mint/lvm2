package lvm2

import (
	"errors"
	"os"

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

	vm.Registers[REGISTER_SYS35] = 0

	err = vm.Memory.GetMemoryFunc(p, n, func(_ uint64, b []byte) error {
		written, err := file.Write(b)
		if err != nil {
			return err
		}

		vm.Registers[REGISTER_SYS35] += uint64(written)
		return nil
	})
	if err != nil {
		return 1, err
	}

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

	vm.Registers[REGISTER_SYS35] = 0

	err = vm.Memory.GetMemoryFunc(p, n, func(_ uint64, b []byte) error {
		read, err := file.Read(b)
		if err != nil {
			return err
		}
		vm.Registers[REGISTER_SYS35] += uint64(read)

		if read != len(b) {
			return errBreak
		}

		return nil
	})

	if err != nil && err != errBreak {
		return 1, err
	}

	return 0, nil
}

func _syscall_open(vm *VM, _, _, _ uint64) (errno uint64, err error) {
	// func Open(path string, flags uint64, mode uint64) (fd uintptr, errno uint64)
	// SYS32[in]: path
	// SYS33[in]: flags
	// SYS34[in]: mode
	// SYS35[out]: fd

	path := vm.Registers[REGISTER_SYS32]
	flags := vm.Registers[REGISTER_SYS33]
	mode := vm.Registers[REGISTER_SYS34]

	var filename []byte
	err = vm.Memory.GetMemoryFunc(path, vm.Memory.MaxAddress-path, func(_ uint64, b []byte) error {
		for _, c := range b {
			if c == 0 {
				return errBreak
			}
			filename = append(filename, c)
		}
		return nil
	})
	if err != nil && err != errBreak {
		return 1, err
	}
	//log.Println(string(filename), int(flags), os.FileMode(mode), err)

	f, err := os.OpenFile(string(filename), int(flags), os.FileMode(mode))
	if err != nil {
		return 1, nil
	}

	fd := vm.FileCounter
	vm.FileCounter++
	vm.Files[fd] = f
	vm.Registers[REGISTER_SYS35] = fd

	return 0, nil
}

func _syscall_close(vm *VM, _, _, _ uint64) (errno uint64, err error) {
	// func Close(fd uintptr) (errno uint64)
	// SYS32[in]: fd

	fd := vm.Registers[REGISTER_SYS32]

	file, ok := vm.Files[fd]
	if !ok {
		return errs.EINVALIDFD.Errno(), nil
	}

	err = file.Close()
	if err != nil {
		return 1, err
	}

	delete(vm.Files, fd)
	return 0, nil
}

var errBreak = errors.New("break")

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
	syscall_Function_Table[SYS_OPEN] = _syscall_open
	syscall_Function_Table[SYS_CLOSE] = _syscall_close
	syscall_Function_Table[SYS_EXIT] = _syscall_exit
	return true
}()
