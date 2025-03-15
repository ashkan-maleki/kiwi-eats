# ---- Project Settings ----

# ---- User-Defined Variables ----
# Default to "user" if not provided
SVN ?= user

# ---- Derived Variables ----
SERVICE_NAME ?= $(SVN)-service
PROJECT_MODULE = github.com/ashkan-maleki/kiwi-eats


PROTO_PATH_OUT = internal/$(SVN)/pb
PROTO_PACKAGE = $(PROJECT_MODULE)/$(PROTO_PATH_OUT)
PROTO_PATH = api/proto/
PROTO_PATH_SVN = $(PROTO_PATH)/$(SVN)

# ---- OS Detection ----
ifeq ($(OS),Windows_NT)
	MKDIR_P = powershell -Command "mkdir -Force '$(PROTO_PATH_OUT)'"
else
	MKDIR_P = mkdir -p "$(PROTO_PATH_OUT)"
endif



# ---- Tools ----
PROTOC = protoc
PROTOC_GEN_GO = protoc-gen-go
PROTOC_GEN_GRPC_GATEWAY = protoc-gen-grpc-gateway

# ---- Targets ----
.PHONY: help
help:  ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


tidy: ## Get all the dependencies
	@echo "Getting all the dependencies..."
	go mod tidy
	go mod vendor
	@echo "Done"

build: ## Build the Go binary for a service
	@echo "Building $(SERVICE_NAME)..."
	go build -o bin/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

test: ## Run all tests
	@echo "Testing..."
	go test -v ./...

grpc: ## Delete later
	protoc -I api/proto \
      --go_out=internal/user/pb --go_opt=module=github.com/ashkan-maleki/kiwi-eats/internal/user/pb \
      --go-grpc_out=internal/user/pb --go-grpc_opt=module=github.com/ashkan-maleki/kiwi-eats/internal/user/pb \
      --grpc-gateway_out=internal/user/pb --grpc-gateway_opt=module=github.com/ashkan-maleki/kiwi-eats/internal/user/pb \
      api/proto/user/auth.proto

generate-proto: ## Generate gRPC and gateway code from .proto files
	@echo "Creating output directory $(PROTO_PATH_OUT)..."
	$(MKDIR_P)
	@echo "Generating protobuf code for $(SVN)..."
	protoc -I $(PROTO_PATH) \
          --go_out=$(PROTO_PATH_OUT) --go_opt=module=$(PROTO_PACKAGE) \
          --go-grpc_out=$(PROTO_PATH_OUT) --go-grpc_opt=module=$(PROTO_PACKAGE) \
          --grpc-gateway_out=$(PROTO_PATH_OUT) --grpc-gateway_opt=module=$(PROTO_PACKAGE) \
          $(PROTO_PATH_SVN)/*.proto

gen-pb: generate-proto ## Generate gRPC and gateway code from .proto files
gen: gen-pb ## Generate all codes that should be generated

generate-proto1: ## Generate gRPC and gateway code from .proto files
	@echo "Generating protobuf code..."
	$(PROTOC) --go_out=. --go_opt=module=$(PROJECT_MODULE) --go-grpc_out=. --go-grpc_opt=module=$(PROJECT_MODULE) --grpc-gateway_out=. --grpc-gateway_opt=module=$(PROJECT_MODULE) 	-I=$(PROTO_PATH) -I=third_party $(PROTO_PATH)/*.proto

dockerize: ## Build a Docker image for the service
	@echo "Building Docker image for $(SERVICE_NAME)..."
	docker build -t kiwi-eats/$(SERVICE_NAME):latest -f deployments/Dockerfile .

run: ## Run the service locally
	@echo "Starting $(SERVICE_NAME)..."
	go run ./cmd/$(SERVICE_NAME)

install-tools: ## Install required Go tools (protoc-gen-go, etc.)
	@echo "Installing tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest