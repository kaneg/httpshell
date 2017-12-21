package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"syscall"

	"crypto/tls"

	"github.com/gorilla/websocket"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/termios"
	"github.com/mattn/go-colorable"
	"../certgen"
	"os/user"
	"path"
)

func loadPublicKey(tlsConfig *tls.Config) {
	usr, _ := user.Current()
	home := usr.HomeDir
	crtPath := path.Join(home, ".httpshell/crt.pem")
	keyPath := path.Join(home, ".httpshell/key.pem")
	certificate, e := tls.LoadX509KeyPair(crtPath, keyPath)
	if e == nil {
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0] = certificate
	}
}

func rawWebSocket(url string) (*websocket.Conn, error) {
	tlsConfig := tls.Config{InsecureSkipVerify: true}
	loadPublicKey(&tlsConfig)
	httpTransport := http.Transport{TLSClientConfig: &tlsConfig}
	dialer := websocket.Dialer{
		TLSClientConfig: httpTransport.TLSClientConfig,
		Proxy:           httpTransport.Proxy,
	}

	headers := http.Header{}

	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return nil, err
	}

	return conn, err

}

func getPatchStdout() io.Writer {
	if runtime.GOOS == "windows" {
		return colorable.NewColorableStdout()
	} else {
		return os.Stdout
	}
}

var server = ""
var debug = false
var genCert = false

func init() {
	flag.BoolVar(&debug, "d", debug, "Debug")
	flag.BoolVar(&genCert, "g", genCert, "genCert")
	flag.Parse()
}

func main() {
	if genCert {
		certgen.CreateNewKeyPair("")
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println("No server address specified.")
		flag.PrintDefaults()
		os.Exit(1)
	} else {
		server = flag.Arg(0)
	}

	var err error

	cfd := int(syscall.Stdin)
	column, row, err := termios.GetSize(cfd)
	if err == nil {
		if debug {
			fmt.Printf("Current window size is row=%d, column=%d\n", row, column)
		}
	}

	address, _ := url.Parse(server)
	if address.Scheme == "https" {
		address.Scheme = "wss"
	} else {
		address.Scheme = "ws"
	}
	address.RawQuery = fmt.Sprintf("row=%d&column=%d", row, column)
	serverAddr := address.String()

	if debug {
		fmt.Println("Connecting to ", serverAddr)
	}

	conn, err := rawWebSocket(serverAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if debug {
		fmt.Println("Connected.")
	}

	var oldTtyState *termios.State
	oldTtyState, err = termios.MakeRaw(cfd)
	if err != nil {
		panic(err)
	}
	defer termios.Restore(cfd, oldTtyState)
	shared.WebsocketSendStream(conn, os.Stdin, -1)
	<-shared.WebsocketRecvStream(getPatchStdout(), conn)
}
