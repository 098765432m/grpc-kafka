PROTO_SRC = services/common/protobuf
PROTO_GO_OUT = services/common/
PROTO_FILES := $(wildcard $(PROTO_SRC)/*.proto)

SQLC_CONFIG_FILES = $(wildcard services/*/sqlc.yaml)
all:
	@echo ">> Running all targets <<"
	@$(MAKE) proto
	@$(MAKE) sqlc

proto:
	@echo ">> protoc all files <<"
	@$(foreach file,$(PROTO_FILES), \
		protoc --proto_path=$(PROTO_SRC) \
		       --go_out=$(PROTO_GO_OUT) \
		       --go-grpc_out=$(PROTO_GO_OUT) \
		       $(file);)
			   
sqlc:
	@echo ">> Running sqlc <<"
	@$(foreach config, $(SQLC_CONFIG_FILES), \
		(cd $(dir $(config)) && sqlc generate);)

docker-compose-up:
	@echo ">> Start Docker Compose <<"
	@cd services && sudo docker-compose up -d

docker-compose-down:
	@echo ">> Start Docker Compose <<"
	@cd services && sudo docker-compose down

clean:
	@echo ">> Cleaning generated files <<"
	@find $(PROTO_GO_OUT) -name "*.pb.go" -type f -delete
	
	
	

