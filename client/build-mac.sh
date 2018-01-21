GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o httpshell-mac
tar czvf httpshell-mac.tar.gz httpshell-mac