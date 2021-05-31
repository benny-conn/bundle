package db

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/logger"
	"github.com/johanbrandhorst/certify"
	"github.com/johanbrandhorst/certify/issuers/vault"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type RSA struct {
	bits int
}

func (r RSA) Generate() (crypto.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, r.bits)
}

func RunServer() error {
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	mode := os.Getenv("MODE")
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return err
	}
	var creds credentials.TransportCredentials
	if mode == "PROD" {
		creds, err = vaultCert()
	} else {
		creds, err = credentials.NewServerTLSFromFile("out/grpc/service.pem", "out/grpc/service.key")
	}
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	api.RegisterUsersServiceServer(grpcServer, newUsersServer())
	api.RegisterPluginsServiceServer(grpcServer, newPluginsServer())
	api.RegisterReadmeServiceServer(grpcServer, newReadmesServer())
	api.RegisterChangelogServiceServer(grpcServer, newChangelogServer())
	logger.InfoLog.Printf("Started Database Server on :%v", port)

	grpcServer.Serve(lis)
	return nil
}

// TODO make this work with kubernetes... after learning kubernetes
func vaultCert() (credentials.TransportCredentials, error) {
	b, err := ioutil.ReadFile("out/grpc/ca.cert")
	if err != nil {
		return nil, fmt.Errorf("vaultCert: problem with input file")
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("vaultCert: failed to append certificates")
	}
	issuer := &vault.Issuer{
		URL: &url.URL{
			Scheme: "https",
			Host:   "vault:8200",
		},
		TLSConfig: &tls.Config{
			RootCAs: cp,
		},
		AuthMethod: vault.ConstantToken(os.Getenv("TOKEN")),
		Role:       "api",
	}
	cfg := certify.CertConfig{
		SubjectAlternativeNames: []string{"localhost"},
		IPSubjectAlternativeNames: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
		KeyGenerator: RSA{bits: 2048},
	}
	c := &certify.Certify{
		CommonName:  "localhost",
		Issuer:      issuer,
		Cache:       certify.DirCache("/data/vault"),
		CertConfig:  &cfg,
		RenewBefore: 24 * time.Hour,
	}
	tlsConfig := &tls.Config{
		GetCertificate: c.GetCertificate,
	}
	return credentials.NewTLS(tlsConfig), nil
}
