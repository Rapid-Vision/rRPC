SCHEMA := test.rrpc
ROOT   := ..
RRPC := $(ROOT)/rRPC

# It is easier to always rebuild everything
.PHONY: all $(RRPC) go-server py-client ts-client openapi clean

all: go-server go-client py-client ts-client openapi

go-server: $(SCHEMA) $(RRPC)
	$(RRPC) server -o ./go_server -f $(SCHEMA)

py-client: $(SCHEMA) $(RRPC)
	$(RRPC) client -o ./py_client -f $(SCHEMA)

go-client: $(SCHEMA) $(RRPC)
	$(RRPC) client --lang go -o ./go_client -f $(SCHEMA)

ts-client: $(SCHEMA) $(RRPC)
	$(RRPC) client --lang ts -o ./ts_client -f $(SCHEMA)
	$(RRPC) client --lang ts --ts-zod --pkg rpcclientzod -o ./ts_client -f $(SCHEMA)

openapi: $(SCHEMA) $(RRPC)
	$(RRPC) openapi -o . -f $(SCHEMA)

$(RRPC):
	cd $(ROOT) && go build

clean:
	rm -rf go_server/rpcserver go_client/rpcclient py_client/rpcclient openapi.json
