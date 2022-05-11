PROJECT="TiCheck"
VERSION=0.1.0

default:
	echo "Welcome to ${PROJECT} v${VERSION}"

fmt:
	echo "Formatting..."
	@clang-format -i *.c *.h

install:
	@echo "Installing ${PROJECT}..."
	@./install.sh
	@echo "Install Done."

build:
	@echo "Building ${PROJECT}..."
	@cd ./cmd/ticheck-server && go build -o ../../bin/ticheck-server
	@cd ../../
	@cd ./web && npm install && npm run build
	@echo "Build Done."

test:
	echo "Testing ${PROJECT}..."
	echo "Test Done."

.PHONY: default install test