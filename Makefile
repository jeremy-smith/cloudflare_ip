.PHONY: darwin linux

darwin: 
	go build -o cf_ip *.go

linux:
	env GOOS=linux GOARCH=amd64 go build -o cf_ip_linux *.go
