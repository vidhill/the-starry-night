ROOT_PATH=cmd/webapp/main.go
SETTINGS_PRIVATE=settings_private.yaml
SWAGGER_UI_FOLDER=swagger-ui
SHELL := /bin/bash # mac quirk, need to declare which shell to use

UNIT_TESTS=$(shell go list ./... | grep -v /integration)
INTEGRATION_TESTS=$(shell go list ./... | grep /integration)

UNIT_TEST_OUTPUT_FILE=.testCoverage.txt

default: pre-build swagger.download-ui swagger.scan
	go build $(ROOT_PATH)

pre-build:
	rm -rf main

start:
	go run $(ROOT_PATH)

dev:
	air

test:
   ifneq (, $(shell richgo version))
		richgo test $(UNIT_TESTS) -coverprofile $(UNIT_TEST_OUTPUT_FILE) -covermode=atomic
   else
		go test $(UNIT_TESTS) -coverprofile $(UNIT_TEST_OUTPUT_FILE) -covermode=atomic
   endif
	
test.html-report: test
	go tool cover -html=$(UNIT_TEST_OUTPUT_FILE)

test.html-report-save:
# - save the html coverage report to disk
	go tool cover -o $(JUNIT_FILE_LOCATION)/test-coverage.html -html=$(UNIT_TEST_OUTPUT_FILE)

test.ci:
	gotestsum --packages="$(UNIT_TESTS)" --junitfile $(JUNIT_FILE_LOCATION)/gotestsum-report.xml --  -coverprofile $(UNIT_TEST_OUTPUT_FILE) -covermode=atomic
	@make test.html-report-save

test.integration:
  ifneq (, $(shell richgo version))
		richgo test -v -count=1 $(INTEGRATION_TESTS)
  else
		go test -v -count=1 $(INTEGRATION_TESTS)
  endif

setup-git-hooks:
	$(info Setting up git hooks)
	@printf '#!/bin/sh \nmake pre-push-hook' > .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push

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

pre-push-hook: lint test

lint:
	golangci-lint run

# 
# Check are dependencies installed
# 

check.dependencies: check.swagger check.golangci-lint check.forbidigo check.staticcheck check.air

check.swagger:
   ifeq (, $(shell which swagger))
		$(error swagger is not installed, Please install go swagger https://goswagger.io/install.html)
   endif

check.air:
   ifeq (, $(shell which air))
		$(error air is not installed, Please install run "go install github.com/cosmtrek/air@v1.27.10")
   endif

check.golangci-lint:
   ifeq (, $(shell which golangci-lint))
		$(error golangci-lint is not installed, Please install see https://golangci-lint.run/usage/install/#local-installation)
   endif