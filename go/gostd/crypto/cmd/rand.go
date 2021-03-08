package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

var idChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func main() {
	n := 10
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		fmt.Println("rand.Read() =====err:", err)
		os.Exit(2)
	}

	fmt.Println(b)
	fmt.Println("hex.EncodeToString() ----- ", hex.EncodeToString(b))

	fmt.Println("=====================")

	//生成随机字母数字
	rStr := make([]byte, n)
	for i, v := range b {
		index := v % byte(len(idChars))
		rStr[i] = idChars[index]
	}

	fmt.Println(string(rStr))

}
