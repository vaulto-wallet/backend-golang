package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GenerateRSAKey(password []byte) ([]byte, []byte) {
	passwordHash := sha256.Sum256([]byte(password))

	privateRSAKey, err := rsa.GenerateKey(rand.Reader, 2048)

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateRSAKey)

	privateKeyEncryptedPEM, err := x509.EncryptPEMBlock(
		rand.Reader,
		"RSA PRIVATE KEY",
		privateKeyBytes,
		passwordHash[:],
		x509.PEMCipherAES256)

	if err != nil {
		fmt.Println("Error creating encrypted PEM:", err)
		return nil, nil
	}

	publicKeyPEM := x509.MarshalPKCS1PublicKey(privateRSAKey.Public().(*rsa.PublicKey))

	publicKeyPEMencoded := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyPEM,
		})

	privateKeyPEMencoded := pem.EncodeToMemory(privateKeyEncryptedPEM)

	return privateKeyPEMencoded, publicKeyPEMencoded
}

func EncryptWithRSA(publicKey []byte, ciphertext []byte) ([]byte, error) {
	publicPEM, _ := pem.Decode(publicKey)

	restoredRSAPublicKey, err := x509.ParsePKCS1PublicKey(publicPEM.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, restoredRSAPublicKey, ciphertext)

}

func DecryptWithRSA(password []byte, privateKey []byte, ciphertext []byte) ([]byte, error) {
	passwordHash := sha256.Sum256(password)
	encryptedPEM, _ := pem.Decode(privateKey)

	rawPEM, err := x509.DecryptPEMBlock(encryptedPEM, passwordHash[:])
	if err != nil {
		return nil, err
	}

	restoredRSAKey, err := x509.ParsePKCS1PrivateKey(rawPEM)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, restoredRSAKey, ciphertext)
}
