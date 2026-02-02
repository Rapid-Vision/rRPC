SCHEMA := hello.rrpc
ROOT   := ../..

ifneq ("$(wildcard $(ROOT)/go.mod)","")
RRPC := $(ROOT)/rRPC
NEED_GO_BUILD := 1
else
RRPC := rRPC
NEED_GO_BUILD :=
endif

.PHONY: all go_server py_server py_client openapi clean go-build

all: go_server py_server py_client openapi

ifeq ($(NEED_GO_BUILD),1)
go_server: go-build
py_server: go-build
py_client: go-build
openapi: go-build
endif

go_server: $(SCHEMA)
	$(RRPC) server --lang go -o ./go_server -f $(SCHEMA)

py_server: $(SCHEMA)
	$(RRPC) server --lang py -o ./py_server -f $(SCHEMA)

py_client: $(SCHEMA)
	$(RRPC) client -o ./py_client -f $(SCHEMA)

openapi: $(SCHEMA)
	$(RRPC) openapi -o . -f $(SCHEMA)

go-build:
	cd $(ROOT) && go build

clean:
	rm -rf ./go_server/rpcserver ./py_server/rpcserver ./py_client/rpcclient openapi.json
