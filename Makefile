ROOT_PATH=cmd/webapp/main.go
SETTINGS_PRIVATE=settings_private.yaml

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
	swagger generate spec -o swagger-ui/swagger.yaml --scan-models

serve-swagger:
	swagger serve -F=swagger swagger-ui/swagger.yaml

download-extract-ui:
	curl -L -o swagger-ui.tar.gz https://github.com/swagger-api/swagger-ui/archive/refs/tags/v4.1.3.tar.gz
	mkdir -p swagger-ui-bundle
	tar -xzf swagger-ui.tar.gz -C swagger-ui-bundle --strip-components 1
	mkdir -p swagger-ui
	mv swagger-ui-bundle/dist/* swagger-ui

download-swagger-ui:
	@make download-extract-ui
	sed 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/swagger.yaml/' swagger-ui/index.html > swagger-ui/index_temp.html
	mv swagger-ui/index_temp.html swagger-ui/index.html
	@make cleanup-download-swagger-ui

cleanup-download-swagger-ui:
	rm swagger-ui.tar.gz
	rm -rf swagger-ui-bundle


scan-serve-swagger: check-swagger scan-swagger serve-swagger

create-settings-private:
  ifeq ($(wildcard $(SETTINGS_PRIVATE)),) # only create if does not exist
		@touch $(SETTINGS_PRIVATE)
		@echo "WEATHER_BIT_API_KEY:" > $(SETTINGS_PRIVATE)
		$(info Created file: $(SETTINGS_PRIVATE))
  endif