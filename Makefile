.PHONY: all clean
# 被编译的文件
BUILDFILE=main.go
# 编译后的静态链接库文件名称
TARGETNAME=Q-wind
# GOOS为目标主机系统 
# mac os : "darwin"
# linux  : "linux"
# windows: "windows"
GOOS=linux
GOOSMAC=linux
# GOARCH为目标主机CPU架构, 默认为amd64 
GOARCH=amd64

VER=$(shell sh ./version/ver.sh)
all: format test build clean print build_mac

test:
	go test -v .

format:
	gofmt -w .

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o ~/Documents/$(TARGETNAME) $(BUILDFILE)



print:
	echo $(VER)

clean:
	go clean -i