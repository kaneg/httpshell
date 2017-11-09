// +build windows

package main

import "github.com/gorilla/websocket"

func runShell(conn *websocket.Conn, command []string,row, column int) {
	panic("Doesn't support windows")
}
