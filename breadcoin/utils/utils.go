package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
)

func Hash(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	val := h.Sum(nil)
	return val
}

func GenerateKeypair() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		panic(err)
	}
	return key
}

func Sign(privKey *rsa.PrivateKey, msg string) []byte {
	h := Hash(msg)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, h[:])
	if err != nil {
		panic(err)
	}
	return sig
}

func VerifySignature(privKey *rsa.PrivateKey, msg string, sig []byte) bool {
	h := Hash(msg)
	err := rsa.VerifyPKCS1v15(&privKey.PublicKey, crypto.SHA256, h[:], sig)
	if err != nil {
		return false
	}
	return true
}

func CalculateAddress(privKey *rsa.PrivateKey) string {
	stringKey, err := json.Marshal(&privKey.PublicKey)
	if err != nil {
		panic(err)
	}
	h := Hash(string(stringKey))
	stringOfHash := b64.StdEncoding.EncodeToString(h)
	return stringOfHash
}

func AddressMatchesKey(addr string, privKey *rsa.PrivateKey) bool {
	return addr == CalculateAddress(privKey)
}
