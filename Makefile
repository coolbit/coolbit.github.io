.PHONY: setup

setup:
	cd editor && go mod tidy
	cd editor/frontend && npm install
