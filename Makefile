# 指定执行shell
SHELL=/bin/bash

#工程名
PROJECT_NAME=daemonw

##当前工作目录作为GOPATH
GOPATH=$(shell pwd)

#Go命令
GO=go
GODEP=dep
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get

## 0/1
CGO_ENABLED=0
## linux/windows
GOOS=linux
## amd64/in3686
GOARCH=in386

##自定义参数
##可执行文件名
BINARY_NAME=test

##源码路径
WORK_DIR=./src/$(PROJECT_NAME)

main:build

#init a empty project structure
init:
	mkdir pkg bin src
	mkdir $(WORK_DIR)
	touch $(WORK_DIR)/main.go
#default linux
build:
ifeq ("release", "$(type)")
##release版去除调试信息
	source ./setenv.sh $(arch) && $(GOBUILD) -ldflags "-s -w" -o $(BINARY_NAME) -v -tags=jsoniter $(WORK_DIR)/.
##使用upx压缩可执行文件,减小size,效果很不错
	upx $(BINARY_NAME)
else
	source ./setenv.sh $(arch) && $(GOBUILD) -o $(BINARY_NAME) -v -tags=jsoniter $(WORK_DIR)/.
endif

test:
	$(GOTEST) -v $(WORK_DIR)/.

clean:

	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)
depinstall:
	cd $(WORK_DIR) && $(GODEP) ensure
depupdate:
	cd $(WORK_DIR) && $(GODEP) ensure -update
depinit:
##初始化依赖时设置为系统默认GOPATH
	@export GOPATH=$(shell echo $$GOPATH)
	cd $(WORK_DIR) && $(GODEP) init -gopath=true
