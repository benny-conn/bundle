
version = 1.0


all: docker-run

docker-build:
	go mod vendor
	docker-compose build 

docker-run: docker-build
	docker-compose up -d

dev:
	scripts/run.sh

cert:
	scripts/tls.sh
	openssl genrsa -out tls/ca.key 4096
	openssl req -new -x509 -key tls/ca.key -sha256 -subj "/C=US/ST=CA/O=BundleMC" -days 1825 -out tls/ca.cert
	openssl genrsa -out tls/service.key 4096
	openssl req -new -key tls/service.key -out tls/service.csr -config tls/cert.conf
	openssl x509 -req -in tls/service.csr -CA tls/ca.cert -CAkey tls/ca.key -CAcreateserial \
		-out tls/service.pem -days 1825 -sha256 -extfile tls/cert.conf -extensions req_ext


.PHONY: clean help
clean:
	scripts/kill.sh


