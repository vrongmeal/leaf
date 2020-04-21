build:
	@./scripts/build.sh

format:
	@./scripts/format.sh

help:
	@echo "Leaf Makefile: make [<command>]"
	@echo "Available commands:"
	@echo "build   -- build leaf into ./build"
	@echo "format  -- format code"
	@echo "help    -- print this message"
	@echo "install -- install development tools"
	@echo "lint    -- lint code for mistakes"

install:
	@./scripts/install.sh

lint:
	@./scripts/lint.sh

.PHONY: build format help install lint
