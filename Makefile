CONTAINER_NAME=builder
CONTAINER_TAG=1.0.0
PROTO_PATH=`pwd`/proto

.PHONY: build
build:
	docker build -t $(CONTAINER_NAME):$(CONTAINER_TAG) .

.PHONY: local
local:
	GOPROXY=direct CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main ./main.go

.PHONY: mocks
mocks:
	mockery -all -inpkg -case snake

.PHONY: proto
proto:
	docker run --rm -v$(PROTO_PATH):$(PROTO_PATH) -w$(PROTO_PATH) jaegertracing/protobuf:latest --proto_path=$(PROTO_PATH) \
    --go_out=$(PROTO_PATH) -I/usr/include/github.com/gogo/protobuf $(PROTO_PATH)/todo.proto
# 	docker build -t protogen ./proto
# 	docker run --name protogen protogen
# 	docker cp protogen:/bin/todo.pb.go ./proto
# 	docker rm protogen