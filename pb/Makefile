.PHONY: all

all:
	@protoc --proto_path=. --go_out=plugins=grpc:. tags.proto
	@echo "Protofiles rebuild done"
