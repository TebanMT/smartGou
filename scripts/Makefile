ifeq ($(OS),Windows_NT)
    ROOT_DIR := $(shell cygpath -w $(shell cd .. && pwd))
else
    ROOT_DIR := $(shell cd .. && pwd)
endif

build-lambda-function:
	cd src/modules/$(MODULE_NAME)/interfaces/lambda_handlers && \
	GOOS=linux GOARCH=amd64 go build -o $(ROOT_DIR)/bin/$(FUNC_NAME)/bootstrap $(FUNC_NAME).go && \
	mkdir -p $(ROOT_DIR)/bin/$(FUNC_NAME) && \
	cd $(ROOT_DIR)/bin/$(FUNC_NAME) && Compress-Archive -Path bootstrap -DestinationPath bootstrap.zip

deploy:
	cd infraestructure/cdk && npm run deploy

