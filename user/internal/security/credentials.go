package security

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/PakornBank/go-grpc-example/user/internal/config"
	"google.golang.org/grpc/credentials"
)

func NewCredentials(cfg *config.Config) credentials.TransportCredentials {
	caCert, err := os.ReadFile(cfg.CACertPath)
	if err != nil {
		log.Fatalf("Failed to load CA cert: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	serverCert, err := tls.LoadX509KeyPair(cfg.ServerCertPath, cfg.ServerKeyPath)
	if err != nil {
		log.Fatalf("Failed to load server cert and key: %v", err)
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})
}
