SCHEMA := test.rrpc
ROOT   := ..
RRPC := $(ROOT)/rRPC

.PHONY: all $(RRPC) go-server py-client openapi clean

all: go-server go-client py-client openapi

go-server: $(SCHEMA) $(RRPC)
	$(RRPC) server -o ./go_server -f $(SCHEMA)

py-client: $(SCHEMA) $(RRPC)
	$(RRPC) client -o ./py_client -f $(SCHEMA)

go-client: $(SCHEMA) $(RRPC)
	$(RRPC) client --lang go -o ./go_client -f $(SCHEMA)

openapi: $(SCHEMA) $(RRPC)
	$(RRPC) openapi -o . -f $(SCHEMA)

$(RRPC):
	cd $(ROOT) && go build

clean:
	rm -rf go_server/rpcserver go_client/rpc_client py_client/rpc_client openapi.json
