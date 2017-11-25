package certgen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"time"
	"os"
)

func CreateNewCert(host string) (*tls.Certificate, error) {
	publicBlock, privateBlock, err := createNewCert(host)
	if err != nil {
		return nil, err
	}
	certificate, err := tls.X509KeyPair(pem.EncodeToMemory(publicBlock), pem.EncodeToMemory(privateBlock))
	if err != nil {
		return nil, err
	}
	return &certificate, nil
}

func CreateNewKeyPair(host string) {
	publicBlock, privateBlock, err := createNewCert(host)
	certOut, err := os.Create("crt.pem")
	if err != nil {
		log.Fatalf("failed to open crt.pem for writing: %s", err)
		return
	}

	pem.Encode(certOut, publicBlock)
	certOut.Close()
	log.Print("written crt.pem\n")

	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("failed to open key.pem for writing:", err)
		return
	}
	pem.Encode(keyOut, privateBlock)
	keyOut.Close()
	log.Print("written key.pem\n")
}

func createNewCert(host string) (*pem.Block, *pem.Block, error) {
	if host == "" {
		host = "HTTP Shell Server"
	}
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
		return nil, nil, err
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
		return nil, nil, err
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: host,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
		return nil, nil, err
	}
	publicBlock := pem.Block{Type: "CERTIFICATE", Bytes: derBytes}
	privateBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	return &publicBlock, &privateBlock, nil
}
