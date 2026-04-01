BINARY_NAME := cyblog
BINARY_PATH := ./bin/${BINARY_NAME}
MAIN_FILE_DIR := ./

#init:
#	go install github.com/google/wire/cmd/wire@latest

#wire进行依赖注入生成代码
wire:
	wire ${MAIN_FILE_DIR}

#构建二进制程序
build:
	go build -o ${BINARY_PATH} ${MAIN_FILE_DIR}


rebuild: wire build



swag:
	swag init

# 编译至 Linux AMD64 平台
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_PATH} ${MAIN_FILE_DIR}
