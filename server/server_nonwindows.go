// +build !windows

package main

import (
	"github.com/kr/pty"
	"os/exec"
	"github.com/gorilla/websocket"
	"github.com/lxc/lxd/shared"
	"os"
	"syscall"
	"unsafe"
)

func resizeTerminal(width int, height int, f *os.File) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(height),
		uint16(width),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return errno
	} else {
		return nil
	}
}

func runShell(conn *websocket.Conn, command []string, row, column int) {
	c := exec.Command(command[0], command[1:]...)
	f, err := pty.Start(c)

	if err != nil {
		panic(err)
	}

	resizeTerminal(column, row, f)

	go shared.WebsocketSendStream(conn, f, -1)
	go shared.WebsocketRecvStream(f, conn)
	c.Wait()
}
