PROTO_IN_DIR=./pkg/api/v1/
PROTO_OUT_DIR=./pkg/api/v1/

all: install build-docker gen-cert

lint:
	golangci-lint run --exclude 'unused'

test:
	go test -v ./... --count=1

build-proto:
	protoc -I ${PROTO_IN_DIR} --go_out=$(PROTO_OUT_DIR) $(PROTO_IN_DIR)/v1.proto
	protoc -I ${PROTO_IN_DIR} --go-grpc_out=$(PROTO_OUT_DIR) $(PROTO_IN_DIR)/v1.proto

install:
	go install cmd/mdb-tutorial.go

build-docker:
	CGO_ENABLED=0 GOOS=linux go build -o ./build/mdb-tutorial ./cmd
	docker build --tag mdb-tutorial:1.0 ./build/

gen-cert:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./build/resources/tls.key -out ./build/resources/tls.cert -subj '/CN=nginx'
