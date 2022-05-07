PATH=cmd/webapp/main.go

default: pre-build
	go build $(PATH)

pre-build: 
	rm -rf main

start:
	go run $(PATH)

dev:
	air

setup-git-hooks:
	cp git-hooks/pre-push.sh .git/hooks/pre-push 