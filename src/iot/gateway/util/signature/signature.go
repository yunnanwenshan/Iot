package signature

import (
	"crypto/rsa"
	"fmt"
	"crypto"
	"crypto/rand"
	"encoding/pem"
	"crypto/x509"
	"github.com/pkg/errors"
)

type Type int64

const (
	PKCS1 Type = iota
	PKCS8
)

var CipherInstance Cipher

func init() {
	//-----BEGIN PRIVATE KEY-----, -----END PRIVATE KEY-----
	//-----BEGIN PUBLIC KEY-----, -----END PUBLIC KEY-----
	client, err := NewDefault(`-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAMrYUqi7GFB6vrlC
BaMCqG5XccgBmyLM3iX/kjGigKGxaPzwlwIKOGLAIYgX5IKXy9YjUJALtAJQWOT0
CWlXccih/Hw4B/haOEgQSL5k49Qk1xlGY6Fz7ILrMvp7cRrdOORtkOAXPlSYKOYo
VRuI9MWq3+L2iuB9gGhLHGs/bwbJAgMBAAECgYBKBRRsxBFEVPZCDiiWaoLh+QDp
NkTRNycdgJxthlogJugj3PuN4ALhbjEOQ4G8cf4M/0gHuG2Qppc5vR+uFB3NrT3C
uyif6hjx/fWwloOwAICc2cQh4NiqPhLklq6KEyt2tOjS0pmSxrjg8+v9kOCH8kru
cV0KIrAQyS8/tEziXQJBAPMdhFl58ifz+Oy41NiMGvxLlcv0xuDxwOfxsUc9gk7l
xzYIVGltsgjnnb/eSvVfI6hcTl3Dha0VdC/1Ciepi8cCQQDVmG+iDzTOm9QtKuXl
UPEWeBS8PrVaw5MmnxrAJt7rloAyu2frz8kKFUIl5Jh/I0Gd8h5nhZSZsZ6+vaa9
3jjvAkEAtxvnQEFB62eteBZqccNs29POOnTdijVr1wbKQF8Kk4Qre/3gHhw5+M0C
mq3CBXen8rm7aJHIUCoVfb1w7ZicpwJAao05mxmE2VCJHuMYfjXLns7WYTXTGG0Z
2hliqdp6OAIC/8vXQp6MBpimP+ryW/IFiLpAipnrkGQ38aUAKhVSRwJAW00a52I2
U+Q1aVqzB5FtzDhVst+LruBVZ1DKT++Dv7wzITHErMJmfz7fXgWoxeg9tDEfJKu3
r1+94qaF7XV/qw==
-----END PRIVATE KEY-----`, `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDK2FKouxhQer65QgWjAqhuV3HI
AZsizN4l/5IxooChsWj88JcCCjhiwCGIF+SCl8vWI1CQC7QCUFjk9AlpV3HIofx8
OAf4WjhIEEi+ZOPUJNcZRmOhc+yC6zL6e3Ea3TjkbZDgFz5UmCjmKFUbiPTFqt/i
9orgfYBoSxxrP28GyQIDAQAB
-----END PUBLIC KEY-----`)

	if err != nil {
		fmt.Println(err)
	}

	CipherInstance = client
}

type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	Verify(src []byte, sign []byte, hash crypto.Hash) error
}

type pkcsClient struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (this *pkcsClient) Encrypt(plaintext []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, this.publicKey, plaintext)
}
func (this *pkcsClient) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, this.privateKey, ciphertext)
}

func (this *pkcsClient) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, this.privateKey, hash, hashed)
}

func (this *pkcsClient) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(this.publicKey, hash, hashed, sign)
}

//默认客户端，pkcs8私钥格式，pem编码
func NewDefault(privateKey, publicKey string) (Cipher, error) {
	blockPri, _ := pem.Decode([]byte(privateKey))
	if blockPri == nil {
		return nil, errors.New("private key error")
	}

	blockPub, _ := pem.Decode([]byte(publicKey))
	if blockPub == nil {
		return nil, errors.New("public key error")
	}

	return New(blockPri.Bytes, blockPub.Bytes, PKCS8)
}

func New(privateKey, publicKey []byte, privateKeyType Type) (Cipher, error) {

	priKey, err := genPriKey(privateKey, privateKeyType)
	if err != nil {
		return nil, err
	}
	pubKey, err := genPubKey(publicKey)
	if err != nil {
		return nil, err
	}
	return &pkcsClient{privateKey: priKey, publicKey: pubKey}, nil
}

func genPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func genPriKey(privateKey []byte, privateKeyType Type) (*rsa.PrivateKey, error) {
	var priKey *rsa.PrivateKey
	var err error
	switch privateKeyType {
	case PKCS1:
		{
			priKey, err = x509.ParsePKCS1PrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
		}
	case PKCS8:
		{
			prkI, err := x509.ParsePKCS8PrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
			priKey = prkI.(*rsa.PrivateKey)
		}
	default:
		{
			return nil, errors.New("unsupport private key type")
		}
	}
	return priKey, nil
}


