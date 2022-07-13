# Relative root dir ("."|".."|"../.."|…)
_ROOT_DIR := $(patsubst ./%,%,$(patsubst %/.manala/Makefile,%,./$(filter %.manala/Makefile,$(MAKEFILE_LIST))))
# Is current dir root ? (""|"1")
_ROOT := $(if $(filter .,$(_ROOT_DIR)),1)
# Relative current dir ("."|"foo"|"foo/bar"|…)
_DIR := $(patsubst ./%,%,.$(patsubst $(realpath $(CURDIR)/$(_ROOT_DIR))%,%,$(CURDIR)))

_BIN_NAME := supateam

include ../../.make/text.mk
include ../../.make/help.mk

.DEFAULT_GOAL := help

.PHONY: setup
HELP += $(call help,setup, Setup the development environment)
setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: test
HELP += $(call help,test,  Run tests)
test: export APP_ENV_PATH = $(shell pwd)/.env.test
test:
	go test -v ./...

.PHONY: build
HELP += $(call help,build, Build the binary)
build: swagger
	@go mod download && \
 	go build -o $(_BIN_NAME) main.go && \
    chmod +x $(_BIN_NAME) && \
    echo "$(_BIN_NAME) has been built."

.PHONY: clean
HELP += $(call help,clean, Cleanup environment)
clean:
	@rm -f $(_BIN_NAME)

.PHONY: swagger
HELP += $(call help,test,  Run tests)
swagger: setup
	swag init --parseDependency