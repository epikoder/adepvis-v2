package crypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"reflect"
)

const alphabet = "./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var (
	bcEncoding = base64.NewEncoding(alphabet)
)

func KeyGen() (k *ecdsa.PrivateKey, err error) {
	if k, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		return nil, err
	}
	return k, nil
}

func PemKeyPair(key *ecdsa.PrivateKey) (privateKeyPEM []byte, publicKeyPEM []byte, err error) {
	var der []byte
	if der, err = x509.MarshalECPrivateKey(key); err != nil {
		return nil, nil, err
	}

	privateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	})

	if der, err = x509.MarshalPKIXPublicKey(key.Public()); err != nil {
		return nil, nil, err
	}

	publicKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: der,
	})
	return
}

func Base64Encode(src string, i interface{}) (string, error) {
	if i != nil && reflect.TypeOf(i).Kind() != reflect.Int {
		return "", fmt.Errorf("int required got: %s", reflect.TypeOf(i).Kind().String())
	}
	n := bcEncoding.EncodedLen(len(src))
	dst := make([]byte, n)
	bcEncoding.Encode(dst, []byte(src))
	for dst[n-1] == '=' {
		n--
	}
	if i != nil && i.(int) <= len(dst) {
		return string(dst[:i.(int)]), nil
	}
	return string(dst[:n]), nil
}

func Base64Decode(s string, i interface{}) (string, error) {
	if i != nil && reflect.TypeOf(i).Kind() != reflect.Int {
		return "", fmt.Errorf("int required got: %s", reflect.TypeOf(i).Kind().String())
	}
	var bcEncoding = base64.NewEncoding(alphabet)
	src := []byte(s)
	numOfEquals := 4 - (len(src) % 4)
	for i := 0; i < numOfEquals; i++ {
		src = append(src, '=')
	}

	dst := make([]byte, bcEncoding.DecodedLen(len(src)))
	n, err := bcEncoding.Decode(dst, []byte(src))
	if err != nil {
		return "", err
	}
	return string(dst[:n]), nil
}
