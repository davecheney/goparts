package main

import (
	"flag"
	"fmt"
	"os"
	"exec"
	"path"
)

func Exitf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

type Visitor chan<- os.Error 

func (v Visitor) VisitDir(path string, f *os.FileInfo) bool {
	return true
}

func (v Visitor) VisitFile(name string, file *os.FileInfo) {
	fmt.Printf("path: %s, file: %#v\n", name, file)
	cmd, err := exec.Run(name, []string { name }, os.Environ(), "", exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if err != nil {
		Exitf("Unable to execute %s: %s\n", file.Name, err)
	}
	msg, err := cmd.Wait(0)
	if err != nil {
		Exitf("Unable to wait on %s: %s\n", cmd, err)
	}
	defer cmd.Close()
	if msg.ExitStatus() != 0 {
		v <- fmt.Errorf("%s exited with status %d", file.Name, msg.ExitStatus())
	}	
}

func errChan() chan<- os.Error {
	c := make(chan os.Error)
	go func() {
		for err := range c {
			fmt.Println(err)
		}
	}()
	return c
}

func main() {
	for _, dir := range flag.Args() {
		fmt.Printf("dir: %s\n", dir)
		errChan := errChan()
		path.Walk(dir, Visitor(errChan), errChan);
	}
}