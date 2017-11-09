package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

const Row = 50
const Column = 120

var listenAddress = ":5000"
var command = []string{"bash"}
var debug = false

func init() {
	flag.StringVar(&listenAddress, "l", listenAddress, "Listen address")
	flag.BoolVar(&debug, "d", debug, "Debug")
	flag.Parse()
}

func createHandler(command [] string) func(http.ResponseWriter, *http.Request) {
	wsHandler := func(writer http.ResponseWriter, request *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
		}
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			panic(nil)
		}
		if debug {
			fmt.Println(conn.RemoteAddr())
			fmt.Println(request.RequestURI)
		}
		rowS := request.FormValue("row")
		columnS := request.FormValue("column")
		row := Row
		column := Column
		if rowS != "" {
			row, _ = strconv.Atoi(rowS)
			if row == 0 {
				row = Row
			}
		}
		if columnS != "" {
			column, _ = strconv.Atoi(columnS)
			if column == 0 {
				column = Column
			}
		}
		if debug {
			fmt.Println("Row:", row)
			fmt.Println("Column:", column)
			fmt.Println("Shell Start")
		}
		runShell(conn, command, row, column)
		if debug {
			fmt.Println("Shell End")
		}
	}
	return wsHandler
}

func main() {
	cmd := flag.Args()
	if flag.NArg() > 0 {
		command = cmd
	}

	if debug {
		fmt.Println("Command:", command)
	}
	http.HandleFunc("/", createHandler(command))
	fmt.Println("Server Started on ", listenAddress)
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		fmt.Println(err)
	}
}
