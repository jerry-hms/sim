DIR = $(shell pwd)/app

IDL_PATH = $(shell pwd)/idl

SERVICES := gateway im user
service = $(word 1, $@)

node = 0

BIN = $(shell pwd)/bin

.PHONY: proto
proto:
	@for file in $(IDL_PATH)/*.proto; do \
		protoc -I $(IDL_PATH) $$file --go-grpc_out=$(IDL_PATH) --go_out=$(IDL_PATH); \
	done
	@find $(IDL_PATH)/pb/* -type f -name '*.pb.go' -exec protoc-go-inject-tag -input {} \;

.PHONY: $(SERVICES)
$(SERVICES):
	go build -o $(BIN)/$(service) $(shell pwd)/cmd/$(service)
	$(BIN)/$(service)

.PHONY: env-up
env-up:
	docker-compose up -d

.PHONY: env-down
env-down:
	docker-compose down