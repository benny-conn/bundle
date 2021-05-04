package gate

import "google.golang.org/grpc/credentials"

func getCert() (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile("out/grpc/ca.cert", "")
	if err != nil {
		return nil, err
	}
	return creds, nil
}
