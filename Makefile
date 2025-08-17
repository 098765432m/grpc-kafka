PROTO_SRC = services/common/protobuf
PROTO_GO_OUT = services/common/
PROTO_FILES := $(wildcard $(PROTO_SRC)/*.proto)

SQLC_CONFIG_FILES = $(wildcard services/*/sqlc.yaml)

MAIN_FILES = $(wildcard services/*/main.go)
API_GATEWAY_MAIN_FILE = $(wildcard services/api-gateway/*/main.go)
MAIN_FILES_WITHOUT_API_GATEWAY = ${filter-out ${API_GATEWAY_MAIN_FILE}, ${MAIN_FILES}}

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

run-dev:
	@echo ">> Running Development <<"
	@$(foreach main, $(MAIN_FILES_WITHOUT_API_GATEWAY), \
		(cd $(dir $(main)) && go run main.go &); \
	)
	@sleep 3
	@(cd $(dir $(API_GATEWAY_MAIN_FILE)) && go run main.go &)
	wait

docker-compose-up:
	@echo ">> Start Docker Compose <<"
	@cd services && sudo docker-compose up -d

docker-compose-down:
	@echo ">> Stop Docker Compose <<"
	@cd services && sudo docker-compose down -v --rmi all

docker-compose-reset:
	@echo ">> Reset Docker Compose <<"
	@cd services && sudo docker-compose down
	@$(MAKE) docker-compose-up

clean:
	@echo ">> Cleaning generated files <<"
	@find $(PROTO_GO_OUT) -name "*.pb.go" -type f -delete
	
kill-main:
	@pkill -x main

