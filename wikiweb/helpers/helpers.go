package helpers

import (
	"crypto/rand"
	"math/big"
)

func randASCIIInt() int{
	//ascii range for printable characters
	min:=33
	max:=126

	bound := max-min

	nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(bound))) 

	num := int(nBig.Int64()) + min

	return num

}

// RandString returns a cryptographically random but printable string using ascii characters with encoding in [33,126]
func RandString(len int) string{
	b := make([]rune, len)

	for i := 0; i < len; i++ {
		b[i] = rune(randASCIIInt())
	}
	return string(b)
}