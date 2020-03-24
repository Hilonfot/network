.PHONY: build clean tool lint help

FILE = "protofile"

default:
	go run main.go

race:
	go run -race main.go

scandir:
	python3 scandir.py

deploy:
	git pull
	git log -1

proto:
	cd ./$(FILE); protoc --go_out=plugins=grpc:. *.proto;  cd ..

help:
	@echo "用法：make command"
	@echo "	default		运行程序	"
	@echo "	race		竞态分析"
	@echo "	proto		生成protobuf.go文件"
	@echo "	scandir		扫描文件夹用于go mod replace"
	@echo "	deploy		代码部署 git pull && git log -1"


