package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

// NewTlsConfig 创建TLS配置
func NewTlsConfig(certPEM string) (*tls.Config, error) {
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM([]byte(certPEM))
	if !ok {
		return nil, fmt.Errorf("parse root certificate failed")
	}

	return &tls.Config{
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}, nil
}
