package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var dbg_file io.Writer

func SetDebug(path string) error {

	if dbg_file != nil {
		return fmt.Errorf("already opened")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbg_file = f
	return nil
}

func Dbg(s string, args ...interface{}) {
	if dbg_file == nil {
		return
	}
	msg := fmt.Sprintf(s, args...)
	log.Output(2, msg)

}
