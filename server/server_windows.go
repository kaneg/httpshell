// +build windows

package main

import (
	"os"
	"github.com/gorilla/websocket"
)

func runShell(conn *websocket.Conn, command []string,row, column uint16) {
	panic("Doesn't support windows")
}
func resizeTerminal(width uint16, height uint16, f *os.File) error {
	return nil
}
