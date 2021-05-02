package gate

import "google.golang.org/grpc/credentials"

func getCert() (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile("bundlemc.io/cert.pem", "")
	if err != nil {
		return nil, err
	}
	return creds, nil
}
