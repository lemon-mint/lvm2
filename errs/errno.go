package errs

import "strconv"

type Errno uint64

const (
	EINVALIDFD Errno = iota + 1
	EFILEWRITE
	EFILEREAD
)

func (e Errno) Error() string {
	return "errno: " + strconv.FormatUint(uint64(e), 10)
}

func (e Errno) String() string {
	return "errno: " + strconv.FormatUint(uint64(e), 10)
}

func (e Errno) Errno() uint64 {
	return uint64(e)
}
