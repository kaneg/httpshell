GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o httpshelld
tar czvf httpshelld-linux.tar.gz httpshelld