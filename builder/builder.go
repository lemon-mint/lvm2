package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const CACHE_DIR = "./cache"

func BuildGolang(GOARCH, GOOS, PATH, VERSION, OUTPUT string) (log []byte, err error) {
	cmd := exec.Command(
		"go", "build",
		"-o", OUTPUT,
		"-ldflags", "-w -s -X main.Version="+VERSION,
		"-tags", "production",
		PATH,
	)

	ABSCACHE, err := filepath.Abs(CACHE_DIR)
	if err != nil {
		return nil, err
	}

	cmd.Env = append(os.Environ(), "GOARCH="+GOARCH, "GOOS="+GOOS, "GOCACHE="+ABSCACHE)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err = cmd.Run()
	return b.Bytes(), err
}

func main() {
	const PATH = "./cmd/lvm2"
	var VERSION = os.Getenv("VERSION")
	if VERSION == "" {
		VERSION = "v0.0.0"
	}

	Targets := []string{
		"darwin/amd64",
		"darwin/arm64",
		"freebsd/386",
		"freebsd/amd64",
		"freebsd/arm64",
		"freebsd/arm",
		"linux/386",
		"linux/amd64",
		"linux/arm64",
		"linux/arm",
		"linux/ppc64le",
		"linux/ppc64",
		"openbsd/386",
		"openbsd/amd64",
		"openbsd/arm64",
		"openbsd/arm",
		"windows/386",
		"windows/amd64",
		"windows/arm64",
	}

	// Clear dist
	os.RemoveAll("./dist")

	flowControl := make(chan bool, runtime.NumCPU())
	var wg sync.WaitGroup

	for _, Target := range Targets {
		flowControl <- true
		wg.Add(1)
		go func(Target string) {
			v := strings.Split(Target, "/")
			GOOS := v[0]
			GOARCH := v[1]
			OUTPUT := "./dist/lvm2_" + GOOS + "_" + GOARCH + ".exe"
			log.Println("Building", GOOS+"/"+GOARCH)
			log, err := BuildGolang(GOARCH, GOOS, PATH, VERSION, OUTPUT)
			if err != nil {
				os.Stderr.Write(log)
				os.Exit(1)
			}
			<-flowControl
			wg.Done()
		}(Target)
	}

	wg.Wait()
}
