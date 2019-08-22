binary_name=cli-manager

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(binary_name)-darwin-amd64 cmd/cli/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(binary_name)-windows-amd64.exe cmd/cli/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/$(binary_name)-linux-amd64 cmd/cli/main.go

test:
	go test -cover -v cmd/cli/main.go

release:
	semantic-release

