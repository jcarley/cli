package commands

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"
)

// VerifyChain takes a chain and ensures it is a full chain optionally ensuring
// the given private key matches that chain.
func VerifyChain(chainPath string, privateKeyPath string, hostname string, selfSigned bool) {
	if _, err := os.Stat(chainPath); os.IsNotExist(err) {
		fmt.Printf("A cert does not exist at path '%s'\n", chainPath)
		os.Exit(1)
	}
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		fmt.Printf("A private key does not exist at path '%s'\n", chainPath)
		os.Exit(1)
	}
	cert, err := tls.LoadX509KeyPair(chainPath, privateKeyPath)
	if err != nil {
		panic(err)
	}
	if !selfSigned {
		x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			panic(err)
		}
		certPool := x509.NewCertPool()
		for i := 1; i < len(cert.Certificate); i++ {
			c, err := x509.ParseCertificate(cert.Certificate[i])
			if err != nil {
				panic(err)
			}
			certPool.AddCert(c)
		}

		// verify we can make a chain from cert down through the intermediates to a root
		if _, err := x509Cert.Verify(x509.VerifyOptions{
			Intermediates: certPool,
		}); err != nil {
			panic("Failed to verify certificate chain: " + err.Error())
		}
		// verify the cert we pulled out matches the hostname specified
		if err := x509Cert.VerifyHostname(hostname); err != nil {
			panic("Certificate hostname mismatch: " + err.Error())
		}
		outputCertInfo(x509Cert)
		warnOnExpired(x509Cert)
	}
	fmt.Println("Certificate chain and key are valid")
}

func outputCertInfo(cert *x509.Certificate) {
	fmt.Printf("Issued by: %s\n", cert.Issuer.CommonName)
	fmt.Printf("Subject: %s\n", cert.Subject.CommonName)
	switch cert.SignatureAlgorithm {
	case x509.UnknownSignatureAlgorithm:
		fmt.Println("Signature Algorithm: Unknown")
	case x509.MD2WithRSA:
		fmt.Println("Signature Algorithm: MD2 with RSA")
	case x509.MD5WithRSA:
		fmt.Println("Signature Algorithm: MD5 with RSA")
	case x509.SHA1WithRSA:
		fmt.Println("Signature Algorithm: SHA 1 with RSA")
	case x509.SHA256WithRSA:
		fmt.Println("Signature Algorithm: SHA 256 with RSA")
	case x509.SHA384WithRSA:
		fmt.Println("Signature Algorithm: SHA 384 with RSA")
	case x509.SHA512WithRSA:
		fmt.Println("Signature Algorithm: SHA 512 with RSA")
	case x509.DSAWithSHA1:
		fmt.Println("Signature Algorithm: DSA with SHA 1")
	case x509.DSAWithSHA256:
		fmt.Println("Signature Algorithm: DSA with SHA 256")
	case x509.ECDSAWithSHA1:
		fmt.Println("Signature Algorithm: ECDSA with SHA 1")
	case x509.ECDSAWithSHA256:
		fmt.Println("Signature Algorithm: ECDSA with SHA 256")
	case x509.ECDSAWithSHA384:
		fmt.Println("Signature Algorithm: ECDSA with SHA 384")
	case x509.ECDSAWithSHA512:
		fmt.Println("Signature Algorithm: ECDSA with SHA 512")
	}
	switch cert.PublicKeyAlgorithm {
	case x509.UnknownPublicKeyAlgorithm:
	case x509.RSA:
		fmt.Println("Public Key Algorithm: RSA")
		publicKey := cert.PublicKey.(*rsa.PublicKey)
		fmt.Printf("Key Size: %d\n", publicKey.N.BitLen())
	case x509.DSA:
		fmt.Println("Public Key Algorithm: DSA")
		publicKey := cert.PublicKey.(*dsa.PublicKey)
		fmt.Printf("Key Size: %d\n", publicKey.Y.BitLen())
	case x509.ECDSA:
		fmt.Println("Public Key Algorithm: ECDSA")
		publicKey := cert.PublicKey.(*ecdsa.PublicKey)
		fmt.Printf("Key Size: %d\n", publicKey.Y.BitLen())
	}
	fmt.Printf("Not Valid Before: %s\n", cert.NotBefore.Local().String())
	fmt.Printf("Not Valid After: %s\n", cert.NotAfter.Local().String())
	fmt.Println()
}

func warnOnExpired(cert *x509.Certificate) {
	if time.Since(cert.NotBefore) < 0 {
		fmt.Println("WARNING! This certificate is not yet valid!")
	}
	if time.Since(cert.NotAfter) > 0 {
		fmt.Println("WARNING! This certificate is expired!")
	}
}
