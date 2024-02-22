DIR = $(shell pwd)/app
GOPATH := $(shell go env GOPATH)

IDL_PATH = $(shell pwd)/idl
PORT ?= 10001

SERVICES := im user
service = $(word 1, $@)

node = 0

BIN = $(shell pwd)/bin

.PHONY: proto
proto:
	@for dir in $(IDL_PATH)/*; do \
		for file in $$dir/*.proto; do \
			protoc -I $$dir $$file --proto_path=.:$(GOPATH)/src:../ --go-grpc_out=../ --go_out=../; \
		done; \
		find $$dir -type f -name '*.pb.go' -exec protoc-go-inject-tag -input {} -remove_tag_comment \;; \
	done

.PHONY: $(SERVICES)
$(SERVICES):
	go build -o $(BIN)/$(service)_$(PORT) -ldflags="-X 'main.Port=$(PORT)'" $(shell pwd)/cmd/$(service)
	$(BIN)/$(service)_$(PORT)

.PHONY: gateway
gateway:
	go build -o $(BIN)/$(service) $(shell pwd)/cmd/$(service)
	$(BIN)/$(service)

.PHONY: env-up
env-up:
	docker-compose up -d

.PHONY: env-down
env-down:
	docker-compose down