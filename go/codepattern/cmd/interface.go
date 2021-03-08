package main

import (
	"codepattern"
	"fmt"
)

func main() {
	r := codepattern.Rect{Width: 3.0, Height: 4.0}
	c := codepattern.Circle{Radius: 2.0}

	s := []codepattern.Shape{&r, &c}

	for _, sh := range s {
		fmt.Printf("%#v\n", sh)
		fmt.Println("Area() == ", sh.Area())
		fmt.Println("Perimeter() == ", sh.Perimeter())
		fmt.Println("---------------------")
	}
}
