CONFIG_SCHEMA = config-schema.json

GO_CONFIG = cli/config/dto.go
GO_SRC_FILES = $(shell find cli -name '*.go') \
							 $(wildcard *.go) $(GO_CONFIG)

all: cli

cli: $(GO_SRC_FILES)
	go build

go-config: $(GO_CONFIG)

clean: clean-go-config

clean-go-config:
	rm $(GO_CONFIG)

$(GO_CONFIG): $(CONFIG_SCHEMA)
	dir=$$(mktemp -d) && \
	cp $< $$dir/config && \
	go-jsonschema -p config $$dir/config > $@; \
	rm $$dir -r

.PHONY: clean clean-go-config
