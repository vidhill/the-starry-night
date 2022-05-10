ROOT_PATH=cmd/webapp/main.go
SETTINGS_PRIVATE=settings_private.yaml
SWAGGER_UI_FOLDER=swagger-ui

default: pre-build download-swagger-ui scan-swagger
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
	which swagger || echo "Please install go swagger https://goswagger.io/install.html"

scan-swagger: download-swagger-ui
	swagger generate spec -o $(SWAGGER_UI_FOLDER)/swagger.yaml --scan-models

serve-swagger:
	swagger serve -F=swagger $(SWAGGER_UI_FOLDER)/swagger.yaml

scan-serve-swagger: check-swagger scan-swagger serve-swagger

download-extract-ui:
	curl -L -o swagger-ui.tar.gz https://github.com/swagger-api/swagger-ui/archive/refs/tags/v4.1.3.tar.gz
	mkdir -p swagger-ui-bundle
	tar -xzf swagger-ui.tar.gz -C swagger-ui-bundle --strip-components 1
	mkdir -p $(SWAGGER_UI_FOLDER)
	mv swagger-ui-bundle/dist/* $(SWAGGER_UI_FOLDER)

download-swagger-ui:
  ifeq ($(wildcard $(SWAGGER_UI_FOLDER)),) # only create if does not exist
		$(info downloading an extracting swagger-ui)
		@make download-extract-ui
		sed 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/swagger.yaml/' $(SWAGGER_UI_FOLDER)/index.html > $(SWAGGER_UI_FOLDER)/index_temp.html
		mv $(SWAGGER_UI_FOLDER)/index_temp.html $(SWAGGER_UI_FOLDER)/index.html
  endif
	@make cleanup-download-swagger-ui

cleanup-download-swagger-ui:
  ifneq ($(wildcard swagger-ui.tar.gz),) # only delte if exists
		rm -rf swagger-ui.tar.gz
  endif
  ifneq ($(wildcard swagger-ui-bundle),)  # only delte if exists
		rm -rf swagger-ui-bundle
  endif

create-settings-private:
  ifeq ($(wildcard $(SETTINGS_PRIVATE)),) # only create if does not exist
		@touch $(SETTINGS_PRIVATE)
		@echo "WEATHER_BIT_API_KEY:" > $(SETTINGS_PRIVATE)
		$(info Created file: $(SETTINGS_PRIVATE))
  endif

lint:
	./git-hooks/pre-push.sh 
	forbidigo -set_exit_status ./...
	staticcheck ./...
