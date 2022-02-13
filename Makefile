SERVER_DIR := dispatcher
GRPC_FILES := $(SERVER_DIR)/dispatcher.pb.go        \
	      $(SERVER_DIR)/dispatcher_grpc.pb.go

CERTS_DIR ?= certs
CERTS := $(CERTS_DIR)/ca-key.pem     \
	 $(CERTS_DIR)/ca-cert.pem    \
	 $(CERTS_DIR)/server-key.pem \
	 $(CERTS_DIR)/server-req.pem \
	 $(CERTS_DIR)/server-cert.pem

BINARY := vinit

default: $(GRPC_FILES) $(CERTS) $(BINARY)

$(SERVER_DIR):
	mkdir -p $@

$(SERVER_DIR)/%.pb.go $(SERVER_DIR)/%_grpc.pb.go: protos/%.proto | $(SERVER_DIR)
	protoc -I protos/ $< --go_out=module=github.com/vinyl-linux/vinit/dispatcher:$(SERVER_DIR) --go-grpc_out=module=github.com/vinyl-linux/vinit/dispatcher:$(SERVER_DIR)

$(CERTS_DIR) $(GENERATED_DIR):
	mkdir -p $@

$(CERTS): | $(CERTS_DIR)
	(cd $(CERTS_DIR) && ../scripts/gen-cert)


$(BINARY): $(GRPC_FILES) $(CERTS) *.go go.mod go.sum
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $@
