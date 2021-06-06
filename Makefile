all: maas_cli_amd64

.PHONY: maas_cli_amd64
maas_cli_amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/maas ./cmd/cli

.PHONY: clean
clean:
	rm -rf build/
