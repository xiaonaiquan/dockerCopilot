VERSION := $(shell echo "v1.0.2")
BUILD_DATE := $(shell date)
LDFLAGS := -X 'dockerCopilot/internal/config.Version=$(VERSION)' -X 'dockerCopilot/internal/config.BuildDate=$(BUILD_DATE)'

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dockerCopilot .

.PHONY: run
run:
	go run -ldflags="$(LDFLAGS)" .

.PHONY: dev
dev:
	go run -ldflags="-X 'dockerCopilot/internal/config.Version=dev' -X 'dockerCopilot/internal/config.BuildDate=$(BUILD_DATE)'" .