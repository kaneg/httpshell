GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o httpshell.exe
zip httpshell-windows.zip httpshell.exe