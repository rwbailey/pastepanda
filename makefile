certs:
	mkdir -p tls \
	&& cd tls \
	&& go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

build:
	go build -o /tmp/web ./cmd/web/ \
	&& cp -r ./tls /tmp/

run:
	cd /tmp && ./web

phony: certs build run
