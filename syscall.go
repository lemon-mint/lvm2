package lvm2

const (
	SYS_WRITE = 16
	SYS_READ  = 17
)

type SYSCALLFunc func(vm *VM, R0, R1, R2 uint64)

var syscall_Function_Table map[uint64]SYSCALLFunc = make(map[uint64]SYSCALLFunc)

func _syscall_write(vm *VM, R0, R1, R2 uint64) {

}

func _syscall_read(vm *VM, R0, R1, R2 uint64) {

}

var _ = func() bool {
	syscall_Function_Table[SYS_WRITE] = _syscall_write
	syscall_Function_Table[SYS_READ] = _syscall_read
	return true
}()
