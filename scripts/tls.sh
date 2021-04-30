#!/bin/bash
#!/bin/bash
rm tls/*.pem
openssl genrsa -out tls/server-key.pem 2048 
openssl req -nodes -new -x509 -sha256 -days 1825 -config tls/cert.conf -extensions 'req_ext' -key tls/server-key.pem -out tls/server-cert.pem