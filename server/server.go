package main

import (
	"os"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
	"crypto/tls"
	"github.com/kaneg/httpshell/certgen"
	"io/ioutil"
	"crypto/x509"
	"os/user"
	"path"
	"text/template"
	"bytes"
	"github.com/kaneg/httpshell"
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

func renderCommand(rawCommand string, context *map[string]string) string {
	if debug {
		fmt.Println("Parse:", rawCommand)
	}
	t := template.Must(template.New("rawCommand").Parse(rawCommand))
	var doc bytes.Buffer

	t.Execute(&doc, context)
	if debug {
		fmt.Println("Parse result:", rawCommand)
	}
	return doc.String()
}

func createHandler(rawCommands [] string) func(http.ResponseWriter, *http.Request) {
	wsHandler := func(writer http.ResponseWriter, request *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
		}
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			panic(nil)
		}
		request.ParseForm()
		if debug {
			fmt.Println(conn.RemoteAddr())
			fmt.Println(request.RequestURI)
			fmt.Println(request.Form)
		}
		context := make(map[string]string)
		for k, v := range request.Form {
			context[k] = v[0]
		}
		var command = make([]string, len(rawCommands))
		for i := 0; i < len(rawCommands); i++ {
			rawCommand := rawCommands[i]
			command[i] = renderCommand(rawCommand, &context)
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
			fmt.Println("Command:", command)
			fmt.Println("Shell Start")
		}

		runShell(conn, command, uint16(row), uint16(column))

		if debug {
			fmt.Println("Shell End")
		}
	}
	return wsHandler
}

func handlePing(conn *websocket.Conn, f *os.File) {
	ph := conn.PingHandler()
	ph2 := func(appData string) error {
		err := ph(appData)
		wsize := httpshell.Winsize{}
		err2 := json.Unmarshal([]byte(appData), &wsize)
		if err2 == nil {
			resizeTerminal(wsize.Columns, wsize.Rows, f)
		}

		return err
	}
	conn.SetPingHandler(ph2)
}

func loadTrustedClientCerts(config *tls.Config) {
	usr, _ := user.Current()
	home := usr.HomeDir
	buffer, err := ioutil.ReadFile(path.Join(home, "/.httpshell/authorized.pem"))
	if err != nil {
		panic(err)
	}
	config.ClientCAs = x509.NewCertPool()
	config.ClientCAs.AppendCertsFromPEM(buffer)
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
