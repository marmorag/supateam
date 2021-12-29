# Relative root dir ("."|".."|"../.."|…)
_ROOT_DIR := $(patsubst ./%,%,$(patsubst %/.manala/Makefile,%,./$(filter %.manala/Makefile,$(MAKEFILE_LIST))))
# Is current dir root ? (""|"1")
_ROOT := $(if $(filter .,$(_ROOT_DIR)),1)
# Relative current dir ("."|"foo"|"foo/bar"|…)
_DIR := $(patsubst ./%,%,.$(patsubst $(realpath $(CURDIR)/$(_ROOT_DIR))%,%,$(CURDIR)))

include .make/text.mk
include .make/help.mk

.DEFAULT_GOAL := help

.PHONY: setup
HELP += $(call help,setup, Setup the development environment)
setup:
	$(MAKE) -C packages/api setup
	git config core.hooksPath .githooks