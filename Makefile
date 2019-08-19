binary_name=cli-manager

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(binary_name)-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o bin/$(binary_name)-windows-amd64.exe
	GOOS=linux GOARCH=amd64 go build -o bin/$(binary_name)-linux-amd64

test:
	go test -cover -v

release:
	semantic-release

