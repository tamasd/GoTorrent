package util

import (
	"crypto/rand"
	"crypto/sha1"
	"io"
	"log"
	"math/big"
)

func Hash(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	sum := h.Sum(nil)

	return string(sum)
}

func GeneratePeerID() string {
	h := sha1.New()

	n, err := rand.Int(rand.Reader, big.NewInt(4096))
	if err != nil {
		log.Print(err)
		return ""
	}

	io.CopyN(h, rand.Reader, n.Int64())

	sum := h.Sum(nil)

	return string(sum)
}
