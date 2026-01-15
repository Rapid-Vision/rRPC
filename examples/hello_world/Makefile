SCHEMA := hello.rrpc
ROOT   := ../..

ifneq ("$(wildcard $(ROOT)/go.mod)","")
RRPC := $(ROOT)/rRPC
NEED_GO_BUILD := 1
else
RRPC := rRPC
NEED_GO_BUILD :=
endif

.PHONY: all server client clean go-build

all: server client

ifeq ($(NEED_GO_BUILD),1)
server: go-build
client: go-build
endif

server:
	$(RRPC) server -o . -f $(SCHEMA)

client:
	$(RRPC) client -o . -f $(SCHEMA)

go-build:
	cd $(ROOT) && go build

clean:
	rm -rf rpcserver rpc_client
