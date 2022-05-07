
default: pre-build
	go build cmd/webapp/main.go 

pre-build: 
	rm -rf main

start:
	go run cmd/webapp/main.go 

dev:
	air

setup-git-hooks:
	cp git-hooks/pre-push.sh .git/hooks/pre-push 