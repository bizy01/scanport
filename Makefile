.PHONY: default dev build_img

default: dev

# 本地环境
RELEASE_DOWNLOAD_ADDR = ""

PUB_DIR = dist
BUILD_DIR = dist

BIN = scanport
NAME = scanport
ENTRY = main.go

DEV_ARCHS = "linux/amd64|windows/amd64|darwin/amd64"
DEFAULT_ARCHS = "all"

VERSION := $(shell git describe --always --tags)
DATE := $(shell date -u +'%Y-%m-%d %H:%M:%S')
GOVERSION := $(shell go version)
COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMITER := $(shell git log -1 --pretty=format:'%an')
UPLOADER:= $(shell hostname)/${USER}/${COMMITER}

all:release  dev

define GIT_INFO
//nolint
package git
const (
	BuildAt string="$(DATE)"
	Version string="$(VERSION)"
	Golang string="$(GOVERSION)"
	Commit string="$(COMMIT)"
	Branch string="$(BRANCH)"
	Uploader string="$(UPLOADER)"
);
endef
export GIT_INFO

# build
define build
	@echo "===== $(BIN) ===="
	@rm -rf $(PUB_DIR)/$(1)/*
	@mkdir -p $(BUILD_DIR) $(PUB_DIR)
	@mkdir -p git
	@echo "$$GIT_INFO" > git/git.go
	@CGO_ENABLED=0 go run cmd/build/build.go -main $(ENTRY) -binary $(BIN) -name $(NAME) -build-dir $(BUILD_DIR) -archs $(1)
	@tree -Csh -L 3 $(BUILD_DIR)
endef

# pub
define pub
	@echo "publish $(NAME) ..."
	@GO111MODULE=off go run cmd/build/build.go -pub-dir $(PUB_DIR) -name $(NAME) -env $(1) -download-addr $(2)
	@tree -Csh -L 3 $(PUB_DIR)
endef

lint:
	@golangci-lint run | tee lint.err # https://golangci-lint.run/usage/install/#local-installation

vet:
	@go vet ./...

dev:
	$(call build, $(DEV_ARCHS))

release:
	$(call build, $(DEFAULT_ARCHS))

pub: release
	$(call pub, release, $(RELEASE_DOWNLOAD_ADDR))

build_img:dev
	@docker build -t scanport:$(VERSION) .

pub_image:build_img
	@docker tag scanport:$(VERSION) bizy01/scanport:latest
	@docker push bizy01/scanport:latest

clean:
	rm -rf dist/*

clean_img:
	docker rmi -rf scanport:**
