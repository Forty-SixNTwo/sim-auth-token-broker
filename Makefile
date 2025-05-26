.PHONY: help

help:
	@echo "Usage: make [target]"
	@echo "Available targets in sub-makefiles: dev, test, build, deploy"

default:
	@$(MAKE) help

include makefiles/*.mk
