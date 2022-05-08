ROOT_PATH=cmd/webapp/main.go

default: pre-build
	go build $(ROOT_PATH)

pre-build:
	rm -rf main

start:
	go run $(ROOT_PATH)

dev:
	air

test:
	go test ./...

integration-test:
	go test -tags="gpintegration" cmd/integration/*.go

setup-git-hooks:
	cp git-hooks/pre-push.sh .git/hooks/pre-push