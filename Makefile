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
	go test -tags="integration" cmd/integration/*.go

setup-git-hooks:
	cp git-hooks/pre-push.sh .git/hooks/pre-push

check-swagger:
	which swagger || echo "Please install go swagger"

scan-swagger:
	swagger generate spec -o swagger.yaml --scan-models

serve-swagger:
	swagger serve -F=swagger swagger.yaml

scan-serve-swagger: check-swagger scan-swagger serve-swagger

create-stettings-private:
  ifeq ($(wildcard $(SETTINGS_PRIVATE)),) # only create if does not exist
		@touch $(SETTINGS_PRIVATE)
		@echo "WEATHER_BIT_API_KEY:" > $(SETTINGS_PRIVATE)
		$(info Created file: $(SETTINGS_PRIVATE))
  endif