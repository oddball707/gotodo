# Set base image
FROM golang:1.13 AS builder

# Project setup
WORKDIR /bin

# Force the go compiler to use modules
ENV GO111MODULE=on

# Grab the code
COPY . .

# Fetch go modules
RUN go mod download

# Run tests
RUN go test ./... -test.v

# Build the source
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

FROM scratch

# Copy the binary from builder
COPY --from=builder /bin/main /main

# Run
ENTRYPOINT ["/main"]

EXPOSE 8080
EXPOSE 8090
