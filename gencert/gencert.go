package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"Manning Publications Co."},
		OrganizationalUnit: []string{"Books"},
		CommonName:         "Go Web Programming",
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, _ := os.Create("key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}

/*
	How HTTPS Works

	HTTPS verifies whether a website is authenticated/certified or not.
	An authenticated website has a unique personal certificate purchased
	from a Certification Authority (CA).

	CAs are trusted companies (GoDaddy, GeoTrust, VeriSign, etc.) who
	provide digital certificates to websites.

	1.	The website owner generates a public key and a private key. He gives
		a Certificate Signing Request (CSR) file and his public key to the CA

	2.	The CA creates a personal certificate based on the CSR, inlcuding domain
		name, owner name, expiry date, and serial#. The CA then adds its own digital
		signature (encrypted text) to the certificate. This digital signature is
		encrypted via the CA's own private key. The CA then encrypts the whole
		certificate with the public key of the website, and sends it back to
		the website owner

	3.	The website owner uses its private key to decrypt the certificate and
		installs it on the website

	4.	Note: the encrypted text is the digital signature of the CA. That text
		is encrypted by the private key of the CA and can only be decrypted by
		a public key of the CA. This means that only that particular CA could
		have encrypted that text--as evident by the fact that you are able to
		decrypt it using that CA's public key

	5.	When you visit www.google.com, Google's server sends you its public key
		and certificate, which was signed by a trusted CA (GeoTrust)
	
	6.	The browser now has to verify the authenticity of the certificate, i.e.
		whether it's actually from GeoTrust or not. Browsers come with a pre-
		installed list of public keys from all the major CAs. It picks the public
		key of GeoTrust and tries to decrypt the digital signature of the certificate

	7.	If the browser successfully decrypts GeoTrust's signature, it means
		that only GeoTrust could have encrypted it to begin with. The browser
		decides that it can trust the website

	8.	Once the certificate is validated, the browser creates a Session Key
		and makes two copies of it. This key can encrypt as well as decrypt
		the data

	9.	The browser encrypts the Session Key (plus other request data) using
		Google's public key, and sends it back to Google's server

	10.	Google's server decrypts the encrypted data using its private key, thereby
		obtaining the session key and other request data

	11.	At this point, both the server and browser have the Session Key, but no one
		else has it. This key can then be used for back-and-forth communication
		between the browser and the server--but only for the duration of that session	

*/
