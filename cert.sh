#!/bin/bash
rm -f out/*.crl
rm -f out/*.crt
rm -f out/*.key
rm -f out/*.csr
certstrap init --common-name "Bundle" 
certstrap request-cert -domain "localhost,bundlemc.io,*.bundlemc.io,web,db,mongo,gate,repo" -ip "127.0.0.1,::1" --common-name "server"
certstrap sign server --CA Bundle
certstrap request-cert -domain "localhost,bundlemc.io,*.bundlemc.io,web,db,mongo,gate,repo" -ip "127.0.0.1,::1" --common-name "client"
certstrap sign client --CA Bundle
rm -f out/grpc/*.cert
rm -f out/grpc/*.key
rm -f out/grpc/*.srl
rm -f out/grpc/*.csr
rm -f out/grpc/*.pem    
openssl genrsa -out out/grpc/ca.key 4096
openssl req -new -x509 -key out/grpc/ca.key -sha256 -subj "/C=US/ST=California/O=BundleMC" -days 3650 -out out/grpc/ca.cert
openssl genrsa -out out/grpc/service.key 4096
openssl req -new -key out/grpc/service.key -out out/grpc/service.csr -config out/grpc/cert.conf
openssl x509 -req -in out/grpc/service.csr -CA out/grpc/ca.cert -CAkey out/grpc/ca.key -CAcreateserial \
    -out out/grpc/service.pem -days 3650 -sha256 -extfile out/grpc/cert.conf -extensions req_ext
sudo rm -r /usr/share/ca-certificates/extra
sudo mkdir /usr/share/ca-certificates/extra
sudo cp out/Bundle.crt /usr/share/ca-certificates/extra
sudo dpkg-reconfigure ca-certificates