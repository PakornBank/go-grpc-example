package security

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/PakornBank/go-grpc-example/gateway/internal/config"
	"google.golang.org/grpc/credentials"
)

func NewCredentials(cfg *config.Config) credentials.TransportCredentials {
	caCert, err := os.ReadFile(cfg.CACertPath)
	if err != nil {
		log.Fatalf("Failed to load CA cert: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(cfg.ClientCertPath, cfg.ClientKeyPath)
	if err != nil {
		log.Fatalf("Failed to load client cert: %v", err)
	}

	return credentials.NewTLS(&tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            certPool,
		InsecureSkipVerify: false,
	})
}
