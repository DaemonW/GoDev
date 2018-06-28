# 指定执行shell
SHELL=/bin/bash

#Go命令
GOCMD=$(GOROOT)/bin/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

##工作目录
WORKDIR=$(shell pwd)
## $$为系统环境变量, $为脚本内变量
SYS_GOPATH=$(shell echo -n $$GOPATH)
##重新定义GOPATH环境变量,追加工作目录
export GOPATH=$(SYS_GOPATH):$(WORKDIR)


##自定义参数
##可执行文件名
BINARY_NAME=GoDev

##源码路径,需要包括main包
SRC_FOLDER=./src/daemonw


all: test build

#default linux
build:
ifeq ("release", "$(type)")
	##release版去除调试信息
	$(GOBUILD) -ldflags "-s -w" -o $(BINARY_NAME) -v -tags=jsoniter $(SRC_FOLDER)/.
	##使用upx压缩可执行文件,效果客观
	upx $(BINARY_NAME)
else
	$(GOBUILD) -o $(BINARY_NAME) -v -tags=jsoniter $(SRC_FOLDER)/.
endif

test:
	$(GOTEST) -v $(SRC_FOLDER)/.

clean:

	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)
deps:
	$(GOGET) -u github.com/gin-gonic/gin
	$(GOGET) -u github.com/gin-gonic/gin/binding
	$(GOGET) -u golang.org/x/crypto/acme/autocert
	$(GOGET) -u github.com/jmoiron/sqlx
	$(GOGET) -u github.com/lib/pq
	$(GOGET) -u github.com/dgrijalva/jwt-go
	$(GOGET) -u github.com/koding/multiconfig
	$(GOGET) -u github.com/go-redis/redis
	$(GOGET) -u github.com/rs/zerolog
	$(GOGET) -u golang.org/x/time/rate
	$(GOGET) -u gopkg.in/gomail.v2

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	$(GOBUILD) -o $(BINARY_NAME) -v -tags=jsoniter $(SRC_FOLDER)/.