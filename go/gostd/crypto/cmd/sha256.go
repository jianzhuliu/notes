package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {
	h1 := sha256.New()
	fmt.Println("sha256 ==  Size() === ", h1.Size())

	//第一种方式
	str := "Hello World"
	data := sha256.Sum256([]byte(str))
	fmt.Printf("%x\n", data)

	//转换为字符串
	sha256_str := hex.EncodeToString(data[:])
	fmt.Println("hex.EncodeToString() ----- ", sha256_str)

	//第二种方式
	h2 := sha256.New()
	io.WriteString(h2, str)
	fmt.Printf("%x\n", h2.Sum(nil))

	//第三种方式
	h3 := sha256.New()
	h3.Write([]byte(str))
	fmt.Printf("%x\n", h3.Sum(nil))

	//第四种方式
	f, err := os.Open("file.txt")
	if err != nil {
		fmt.Println("os.Open() ---err:", err)
		os.Exit(2)
	}
	defer f.Close()

	h4 := sha256.New()
	if _, err := io.Copy(h4, f); err != nil {
		fmt.Println("io.Copy() ---err:", err)
		os.Exit(2)
	}

	fmt.Printf("%x\n", h4.Sum(nil))
}
