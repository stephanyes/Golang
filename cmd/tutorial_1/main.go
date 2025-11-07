package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	var intNum int64
	intNum = 10
	fmt.Println(intNum)

	var floatNum float64
	floatNum = 10.5
	fmt.Println(floatNum)

	var stringNum string
	stringNum = "10"
	fmt.Println(stringNum)

	var uintNum uint64
	uintNum = 10
	fmt.Println(uintNum)

	var boolNum bool
	boolNum = true
	fmt.Println(boolNum)

	var complexNum complex128
	complexNum = 10 + 10i
	fmt.Println(complexNum)

	var myRune rune = 'a'
	fmt.Println(myRune)

	var rune2 rune
	fmt.Println(rune2)

	// in GO instead of using the var keyword we can use :=
	newVar := "text"
}
