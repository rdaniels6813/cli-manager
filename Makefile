binary_name=cli-manager
golangci_lint_version=1.17.1

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(binary_name)-darwin-amd64 cmd/cli/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(binary_name)-windows-amd64.exe cmd/cli/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/$(binary_name)-linux-amd64 cmd/cli/main.go

test:
	go test -cover -v ./...

lint:
	# Install golangci-lint if it's not already installed
	command -v golangci-lint || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b `go env GOPATH`/bin v$(golangci_lint_version)
	golangci-lint run

release:
	semantic-release

