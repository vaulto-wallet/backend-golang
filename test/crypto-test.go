package main

import (
	"crypto/rsa"
	"fmt"
)
import "crypto/sha256"
import "crypto/rand"
import "crypto/x509"
import "encoding/pem"
import h "../helpers"

func main() {
	privateRSAKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA2048 Private Key", err)
		return
	}

	passwordHash := sha256.Sum256([]byte("Password"))

	fmt.Println("Password Hash :", passwordHash)

	cryptoPrivateKey := make([]byte, 32)
	rand.Read(cryptoPrivateKey)
	fmt.Println("Crypto PrivateKey", cryptoPrivateKey)

	fmt.Println("PrivateRSAKey", privateRSAKey.E, privateRSAKey.D.Bytes())

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateRSAKey)

	privateKeyEncryptedPEM, err := x509.EncryptPEMBlock(
		rand.Reader,
		"RSA PRIVATE KEY",
		privateKeyBytes,
		passwordHash[:],
		x509.PEMCipherAES256)

	if err != nil {
		fmt.Println("Error creating encrypted PEM:", err)
		return
	}

	publicKeyPEM := x509.MarshalPKCS1PublicKey(privateRSAKey.Public().(*rsa.PublicKey))

	publicKeyPEMencoded := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyPEM,
		})

	privateKeyPEMencoded := pem.EncodeToMemory(privateKeyEncryptedPEM)

	fmt.Println("PrivateRSAKey PEM\n", string(privateKeyPEMencoded))
	fmt.Println("PublicRSAKey PEM\n", string(publicKeyPEMencoded))

	//hasher := sha256.New()

	encryptedPEM, _ := pem.Decode(privateKeyPEMencoded)

	rawPEM, err := x509.DecryptPEMBlock(encryptedPEM, passwordHash[:])
	if err != nil {
		fmt.Println("Error creating decrypt PEM:", err)
		return
	}

	restoredRSAKey, err := x509.ParsePKCS1PrivateKey(rawPEM)
	if err != nil {
		fmt.Println("Error parse PKCS1 key:", err)
		return
	}

	restoredPublicKey := restoredRSAKey.Public()
	fmt.Println("Decrypted RSA Public Key", restoredPublicKey)

	password := []byte("password")
	passwordHash = sha256.Sum256(password)

	priv, pub := h.GenerateRSAKey(password)
	fmt.Println("Priv", priv)
	fmt.Println("Pub", pub)

	plain := []byte("Plaintextmessage")
	encryptedText, err := h.EncryptWithRSA(pub, plain)
	if err != nil {
		fmt.Println("Error encrypting", err)
		return
	}
	fmt.Println("Encrypted", encryptedText)
	decryptedText, err := h.DecryptWithRSA(password, priv, encryptedText)
	if err != nil {
		fmt.Println("Error decrypting", err)
		return
	}
	fmt.Println("Decrypted", string(decryptedText))
}
