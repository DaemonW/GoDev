# 指定执行shell
SHELL=/bin/bash

#工程名
PROJECT_NAME=daemonw

##当前工作目录作为GOPATH
GOPATH=$(shell pwd)
##系统设置的GOPATH
SYS_GOPATH=$(shell echo $$GOPATH)

#Go命令
GO=$(GOROOT)/bin/go
GODEP=$(SYS_GOPATH)/bin/dep
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get



##自定义参数
##可执行文件名
BINARY_NAME=

##源码路径,需要包括main包
WORK_DIR=./src/$(PROJECT_NAME)




# build and test
all: build test

checkenv:
ifneq ("$(BINARY_NAME)","")
	echo $(BINARY_NAME)
	BINARY_NAME+=_
endif

#init a empty project structure
init:
	mkdir pkg bin src
	mkdir $(WORK_DIR)
	touch $(WORK_DIR)/main.go
#default linux
build:checkenv
ifeq ("release", "$(type)")
	##release版去除调试信息
	$(GOBUILD) -ldflags "-s -w" -o $(BINARY_NAME)release -v -tags=jsoniter $(WORK_DIR)/.
	##使用upx压缩可执行文件,减小size,效果很不错
	upx $(BINARY_NAME)release
else
	$(GOBUILD) -o $(BINARY_NAME)debug -v -tags=jsoniter $(WORK_DIR)/.
endif

test:
	$(GOTEST) -v $(WORK_DIR)/.

clean:

	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
ifeq ("release", "$(type)")
	./$(BINARY_NAME)release
else
	./$(BINARY_NAME)debug
endif
depinstall:
	cd $(WORK_DIR) && $(GODEP) ensure
depupdate:
	cd $(WORK_DIR) && $(GODEP) ensure -update
depinit:
	export GOPATH=$(SYS_GOPATH)
	cd $(WORK_DIR) && $(GODEP) init -gopath=true

# Cross compilation
build-linux:checkenv
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	$(GOBUILD) -o $(BINARY_NAME) -v -tags=jsoniter $(WORK_DIR)/.