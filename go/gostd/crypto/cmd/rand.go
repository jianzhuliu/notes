package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	n := 10
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		fmt.Println("rand.Read() =====err:", err)
		os.Exit(2)
	}

	fmt.Println(b)
	fmt.Println("hex.EncodeToString() ----- ", hex.EncodeToString(b))
}
