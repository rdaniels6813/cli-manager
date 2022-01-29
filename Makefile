binary_name=cli-manager
golangci_lint_version=1.31.0

macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(binary_name)-darwin-amd64 cmd/cli/main.go

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(binary_name)-darwin-amd64 cmd/cli/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(binary_name)-windows-amd64.exe cmd/cli/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/$(binary_name)-linux-amd64 cmd/cli/main.go

test:
	go test -coverprofile=coverage.txt -covermode=atomic ./...

lint:
	# Install golangci-lint if it's not already installed
	test -f .tools/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .tools v$(golangci_lint_version)
	.tools/golangci-lint run

release:
	semantic-release

ci:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/rdaniels6813/cli-manager/internal/version.version=${VERSION_TAG}" -o bin/$(binary_name)-darwin-amd64 cmd/cli/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/rdaniels6813/cli-manager/internal/version.version=${VERSION_TAG}" -o bin/$(binary_name)-windows-amd64.exe cmd/cli/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/rdaniels6813/cli-manager/internal/version.version=${VERSION_TAG}" -o bin/$(binary_name)-linux-amd64 cmd/cli/main.go