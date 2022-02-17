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
TEST_BINARY := vinit-test

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


# Because we need to do things like setuid on binaries, we need
# to run tests via sudo.
#
# In order to do this, we need to compile the tests first, and then
# run it with sudo.
#
# This avoids having to redownload dependencies and having to configure
# toolchains/ go env vars
$(TEST_BINARY): $(GRPC_FILES) $(CERTS) *.go go.mod go.sum
	-go test -covermode=count -o $@ -tags sudo > /dev/null

.PHONY: test
test: $(TEST_BINARY)
	sudo ./$< -test.coverprofile=count.out -test.v
	-go tool cover -html=count.out
