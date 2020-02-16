build:
	@./scripts/build.sh

format:
	@./scripts/format.sh

install:
	@./scripts/install.sh

lint:
	@./scripts/lint.sh

.PHONY: build format help install lint
