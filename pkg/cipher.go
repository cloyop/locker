package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

func md(k []byte) string {
	rs := md5.Sum(k)
	return hex.EncodeToString(rs[:])
}
func MustCipherText(keyPhrase, value []byte) []byte {
	gcm := gcmInstance(keyPhrase)
	nonce := make([]byte, gcm.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, []byte(value), nil)
}
func gcmInstance(keyPhrase []byte) cipher.AEAD {
	hashedPhrase := md(keyPhrase)
	aesBlock, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}
	return gcm
}
func MustUnCipher(keyPhrase, ciphered []byte) ([]byte, error) {
	gcm := gcmInstance(keyPhrase)
	nonceSize := gcm.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]
	return gcm.Open(nil, nonce, cipheredText, nil)
}
