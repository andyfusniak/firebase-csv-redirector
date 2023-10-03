OUTPUT_DIR=./bin

all: redirector

redirector:
	@CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/redirector -ldflags "-X 'main.version=${VERSION}'" ./cmd/redirector/main.go

clean:
	-@rm -r $(OUTPUT_DIR)/* 2> /dev/null || true

.PHONY: clean
