.PHONY: install uninstall generate check test-agents health

install:
	./install.sh --global

uninstall:
	./uninstall.sh --global

generate:
	bash scripts/generate-configs.sh

check:
	bash scripts/check-parity.sh

test-agents:
	bash scripts/test-agents.sh

health:
	bash scripts/health-check.sh --verbose
