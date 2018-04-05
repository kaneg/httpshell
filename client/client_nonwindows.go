// +build !windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/kaneg/httpshell"
)

func notifySignalSIGWINCH(c chan<- os.Signal) {
	signal.Notify(c, syscall.SIGWINCH)
}

func resetSignalSIGWINCH() {
	signal.Reset(syscall.SIGWINCH)
}

func syscallTIOCGWINSZ() ([]byte, error) {
	ws := httpshell.Winsize{}

	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(0), uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))

	b, err := json.Marshal(ws)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal error: %v", err)
	}
	return b, err
}
