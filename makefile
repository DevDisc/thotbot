# Params
SERVICE  ?= thot-bot
VERSION	 ?= $(shell cat version | head -n 1)
REVISION ?= $(shell git rev-parse --short HEAD)

.DEFAULT_GOAL=test

APP ?= src/cmd/bot.go

# Development tasks
COVERPROFILE=.cover.out
COVERDIR=.cover

dep:
	@go get ./...

run:
	@go run $(APP)

test:
	@go test -coverprofile=$(COVERPROFILE) ./...

test-report:
	@go test -coverprofile=$(COVERPROFILE) -json ./... | tee test-report.json

cover: test
	@mkdir -p $(COVERDIR)
	@go tool cover -html=$(COVERPROFILE) -o $(COVERDIR)/index.html
	@cd $(COVERDIR) && python -m SimpleHTTPServer 3000

# Build tasks

clean:
	@rm -rf bin

build: clean
	@echo Building verion: $(VERSION), revision: $(REVISION)
	@CGO_ENDABLED=0 go build -a --installsuffix cgo \
		-o bin/$(SERVICE) $(APP)
