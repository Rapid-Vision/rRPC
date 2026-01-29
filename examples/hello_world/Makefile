SCHEMA := hello.rrpc
ROOT   := ../..

ifneq ("$(wildcard $(ROOT)/go.mod)","")
RRPC := $(ROOT)/rRPC
NEED_GO_BUILD := 1
else
RRPC := rRPC
NEED_GO_BUILD :=
endif

.PHONY: all server client openapi clean go-build

all: server client openapi

ifeq ($(NEED_GO_BUILD),1)
server: go-build
client: go-build
openapi: go-build
endif

server: $(SCHEMA)
	$(RRPC) server -o ./go_server -f $(SCHEMA)

client: $(SCHEMA)
	$(RRPC) client -o ./py_client -f $(SCHEMA)

openapi: $(SCHEMA)
	$(RRPC) openapi -o . -f $(SCHEMA)

go-build:
	cd $(ROOT) && go build

clean:
	rm -rf ./go_server/rpcserver ./py_client/rpcclient openapi.json
