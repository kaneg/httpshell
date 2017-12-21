package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
	"crypto/tls"
	"../certgen"
	"io/ioutil"
	"crypto/x509"
	"os/user"
	"path"
)

const Row = 30
const Column = 120

var listenAddress = ":5000"
var command = []string{"bash"}
var debug = false
var authByKey = false

func init() {
	flag.StringVar(&listenAddress, "l", listenAddress, "Listen address")
	flag.BoolVar(&debug, "d", debug, "Debug")
	flag.BoolVar(&authByKey, "k", authByKey, "authByKey")
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

func loadTrustedClientCerts(config *tls.Config) {
	usr, _ := user.Current()
	home := usr.HomeDir
	bytes, err := ioutil.ReadFile(path.Join(home, "/.httpshell/authorized.pem"))
	if err != nil {
		panic(err)
	}
	config.ClientCAs = x509.NewCertPool()
	config.ClientCAs.AppendCertsFromPEM(bytes)
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
	config := tls.Config{}
	certificate, err := certgen.CreateNewCert("")
	if err != nil {
		panic(err)
	}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = *certificate
	if authByKey {
		config.ClientAuth = tls.RequireAndVerifyClientCert
		loadTrustedClientCerts(&config)
	}
	server := &http.Server{Addr: listenAddress, TLSConfig: &config}
	server.ListenAndServeTLS("", "")

}
