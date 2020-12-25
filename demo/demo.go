package main

import (
	"fmt"
)

func main() {
	bbb := BBB{&AAA{A: "a", B: "b"}, CCC : "ccc"}
	bbb.A = "aaa"
	bbb.B = "bbb"

	fmt.Println("bbb", bbb.AAA)
}


type AAA struct {
	A string
	B string
}

type BBB struct {
	*AAA
	CCC string
}