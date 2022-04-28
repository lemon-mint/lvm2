package main

import (
	"encoding/binary"
	"io"
	"os"
)

func main() {
	runtime := os.Args[1]
	prog := os.Args[2]
	outpath := os.Args[3]

	runf, err := os.Open(runtime)
	if err != nil {
		panic(err)
	}
	defer runf.Close()

	progf, err := os.Open(prog)
	if err != nil {
		panic(err)
	}
	defer progf.Close()

	outf, err := os.Create(outpath)
	if err != nil {
		panic(err)
	}
	defer outf.Close()

	_, err = io.Copy(outf, runf)
	if err != nil {
		panic(err)
	}
	progsize, err := io.Copy(outf, progf)
	if err != nil {
		panic(err)
	}

	// Write the magic header
	_, err = outf.WriteString("@%LVM2%\n")
	if err != nil {
		panic(err)
	}

	var sizebuf [8]byte
	binary.LittleEndian.PutUint64(sizebuf[:], uint64(progsize))

	_, err = outf.Write(sizebuf[:])
	if err != nil {
		panic(err)
	}
}
