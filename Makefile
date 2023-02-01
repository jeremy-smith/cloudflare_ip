.PHONY: darwin darwin_arm linux raspberry

darwin: 
	GOOS=darwin GOARCH=amd64 go build -o cf_ip_darwin_amd64 main.go

darwin_arm:
	GOOS=darwin GOARCH=arm64 go build -o cf_ip_darwin_arm64 main.go

linux:
	env GOOS=linux GOARCH=amd64 go build -o cf_ip_linux_amd64 main.go

raspberry:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o cf_ip_raspberry main.go