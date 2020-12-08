package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {
	h1 := md5.New()
	fmt.Println("md5 ==  Size() === ", h1.Size())

	//第一种方式
	str := "Hello World"
	data := md5.Sum([]byte(str))
	fmt.Printf("%x\n", data)

	//转换为字符串
	md5_str := hex.EncodeToString(data[:])
	fmt.Println("hex.EncodeToString() ----- ", md5_str)

	//第二种方式
	h2 := md5.New()
	io.WriteString(h2, str)
	fmt.Printf("%x\n", h2.Sum(nil))

	//第三种方式
	h3 := md5.New()
	h3.Write([]byte(str))
	fmt.Printf("%x\n", h3.Sum(nil))

	//第四种方式
	f, err := os.Open("file.txt")
	if err != nil {
		fmt.Println("os.Open() ---err:", err)
		os.Exit(2)
	}
	defer f.Close()

	h4 := md5.New()
	if _, err := io.Copy(h4, f); err != nil {
		fmt.Println("io.Copy() ---err:", err)
		os.Exit(2)
	}

	fmt.Printf("%x\n", h4.Sum(nil))
}
