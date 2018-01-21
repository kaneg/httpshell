GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o httpshell-linux
tar czvf httpshell-linux.tar.gz httpshell-linux