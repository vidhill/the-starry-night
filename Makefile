ROOT_PATH=cmd/webapp/main.go
SETTINGS_PRIVATE=settings_private.yaml
SWAGGER_UI_FOLDER=swagger-ui

default: pre-build swagger.download-ui swagger.scan
	go build $(ROOT_PATH)

pre-build:
	rm -rf main

start:
	go run $(ROOT_PATH)

dev:
	air

test:
	go test $(shell go list ./... | grep -v /integration) -coverprofile .testCoverage.txt

integration-test:
	go test $(shell go list ./... | grep /integration)

setup-git-hooks:
	cp git-hooks/pre-push.sh .git/hooks/pre-push

swagger.scan: check.swagger swagger.download-ui
	swagger generate spec -i swagger-base.yaml -o $(SWAGGER_UI_FOLDER)/swagger.yaml --scan-models

swagger.download-extract-ui:
	curl -L -o swagger-ui.tar.gz https://github.com/swagger-api/swagger-ui/archive/refs/tags/v4.1.3.tar.gz
	mkdir -p swagger-ui-bundle
	tar -xzf swagger-ui.tar.gz -C swagger-ui-bundle --strip-components 1
	mkdir -p $(SWAGGER_UI_FOLDER)
	mv swagger-ui-bundle/dist/* $(SWAGGER_UI_FOLDER)

swagger.download-ui:
  ifeq ($(wildcard $(SWAGGER_UI_FOLDER)),) # only create if does not exist
		$(info downloading and extracting swagger-ui)
		@make swagger.download-extract-ui
		sed 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/swagger.yaml/' $(SWAGGER_UI_FOLDER)/index.html > $(SWAGGER_UI_FOLDER)/index_temp.html
		mv $(SWAGGER_UI_FOLDER)/index_temp.html $(SWAGGER_UI_FOLDER)/index.html
  endif
	@make swagger.cleanup-download-ui

swagger.cleanup-download-ui:
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
	./bash_scripts/go-fmt-msg.sh
	forbidigo -set_exit_status ./...
	staticcheck ./...

# 
# Check are dependencies installed
# 

check.dependencies: check.swagger check.forbidigo check.staticcheck check.air

check.swagger:
   ifeq (, $(shell which swagger))
		$(error swagger is not installed, Please install go swagger https://goswagger.io/install.html)
   endif

check.forbidigo:
   ifeq (, $(shell which forbidigo))
		$(error forbidigo is not installed, Please install run "go install github.com/ashanbrown/forbidigo@v1.3.0")
   endif

check.staticcheck:
   ifeq (, $(shell which staticcheck))
		$(error staticcheck is not installed, Please install run "go install honnef.co/go/tools/cmd/staticcheck@2022.1.1")
   endif

check.air:
   ifeq (, $(shell which air))
		$(error air is not installed, Please install run "go install github.com/cosmtrek/air@v1.27.10")
   endif