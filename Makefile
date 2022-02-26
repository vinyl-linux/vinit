SERVER_DIR := dispatcher
GRPC_FILES := $(SERVER_DIR)/dispatcher.pb.go        \
	      $(SERVER_DIR)/dispatcher_grpc.pb.go

CERTS_DIR ?= certs
CERTS := $(CERTS_DIR)/ca-key.pem     \
	 $(CERTS_DIR)/ca-cert.pem    \
	 $(CERTS_DIR)/server-key.pem \
	 $(CERTS_DIR)/server-req.pem \
	 $(CERTS_DIR)/server-cert.pem

PKG_DIR := vinit-x86_64

BINARY := vinit
CLIENT_BINARY := vinitctl
TEST_BINARY := vinit-test

BUILT_ON ?= $(shell date --rfc-3339=seconds | sed 's/ /T/')
BUILT_BY ?= $(shell whoami)
BUILD_REF ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

default: $(GRPC_FILES) $(CERTS) $(BINARY) $(CLIENT_BINARY)

$(SERVER_DIR) $(CERTS_DIR) $(GENERATED_DIR) $(PKG_DIR):
	mkdir -p $@

$(SERVER_DIR)/%.pb.go $(SERVER_DIR)/%_grpc.pb.go: protos/%.proto | $(SERVER_DIR)
	protoc -I protos/ $< --go_out=module=github.com/vinyl-linux/vinit/dispatcher:$(SERVER_DIR) --go-grpc_out=module=github.com/vinyl-linux/vinit/dispatcher:$(SERVER_DIR)

$(CERTS): | $(CERTS_DIR)
	(cd $(CERTS_DIR) && $(CURDIR)/scripts/gen-cert)


$(BINARY): $(GRPC_FILES) *.go go.mod go.sum
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.certDir=$(CERTS_DIR) -X main.ref=$(BUILD_REF) -X main.buildUser=$(BUILT_BY) -X main.builtOn=$(BUILT_ON)" -trimpath -o $@

$(CLIENT_BINARY): pkg = "github.com/vinyl-linux/vinit/client/cmd"
$(CLIENT_BINARY): client client/*.go client/**/*.go
	(cd $< && CGO_ENABLED=0 go build -ldflags="-s -w -X $(pkg).Ref=$(BUILD_REF) -X $(pkg).BuildUser=$(BUILT_BY) -X $(pkg).BuiltOn=$(BUILT_ON)" -trimpath -o ../$@)


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

.PHONY: package
package: $(BINARY) $(CLIENT_BINARY) | $(PKG_DIR) ./scripts
	install -m 0700 $(BINARY) $(PKG_DIR)
	install -m 0700 $(CLIENT_BINARY) $(PKG_DIR)

	cp -r scripts $(PKG_DIR)
	cp Makefile $(PKG_DIR)
