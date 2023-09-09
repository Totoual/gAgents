# Makefile

# Define variables
PROTOC = protoc
PROTO_SRC_DIR = ./proto
GO_OUT_DIR = ./generated

# Define the list of proto files
PROTO_FILES = $(wildcard $(PROTO_SRC_DIR)/*.proto)

# Define the list of generated Go files
GO_FILES = $(patsubst $(PROTO_SRC_DIR)/%.proto,$(GO_OUT_DIR)/%.pb.go,$(PROTO_FILES))

# Define the target to generate Go files from proto files
$(GO_OUT_DIR)/%.pb.go: $(PROTO_SRC_DIR)/%.proto
	protoc --go_out=./generated --go-grpc_out=./generated $<

# Define a target to generate all Go files
generate: $(GO_FILES)

# Define a target to clean generated files
clean:
	rm -rf $(GO_OUT_DIR)

.PHONY: generate clean
