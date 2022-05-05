
default: pre-build
	go build cmd/webapp/main.go 

pre-build: 
	rm -rf main

dev:
	air

setup-git-hooks:
	cp git-hooks/pre-push.sh .git/hooks/pre-push 