FROM uber/prototool:1.8.1

WORKDIR /bin
COPY todo.proto todo.proto
RUN protoc --go_out=plugins=grpc:. -I. todo.proto