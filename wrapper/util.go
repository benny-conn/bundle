package wrapper

import (
	"google.golang.org/grpc/credentials"
)

func GetCert() (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile("tls/server-cert.pem", "")
	if err != nil {
		return nil, err
	}
	return creds, nil
}
