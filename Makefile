
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
	rm -r out/
	mkdir out
	certstrap init --common-name "Bundle" 
	certstrap request-cert -domain "localhost,bundlemc.io,*.bundlemc.io,web,db,mongo,gate,repo" -ip "127.0.0.1,::1" --common-name "server"
	certstrap sign server --CA Bundle
	certstrap request-cert -domain "localhost,bundlemc.io,*.bundlemc.io,web,db,mongo,gate,repo" -ip "127.0.0.1,::1" --common-name "client"
	certstrap sign client --CA Bundle
	rm tls/*.cert
	rm tls/*.key
	rm tls/*.srl
	rm tls/*.csr
	rm tls/*.pem    
	openssl genrsa -out tls/ca.key 4096
	openssl req -new -x509 -key tls/ca.key -sha256 -subj "/C=US/ST=California/O=BundleMC" -days 1825 -out tls/ca.cert
	openssl genrsa -out tls/service.key 4096
	openssl req -new -key tls/service.key -out tls/service.csr -config tls/cert.conf
	openssl x509 -req -in tls/service.csr -CA tls/ca.cert -CAkey tls/ca.key -CAcreateserial \
		-out tls/service.pem -days 1825 -sha256 -extfile tls/cert.conf -extensions req_ext


.PHONY: clean help
clean:
	scripts/kill.sh


