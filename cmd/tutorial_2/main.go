package main

import (
	"fmt"
	"errors"
)

func main () {
	var printValue string = "Hello World"	
	printMe(printValue)

	numerator := 10
	denominator := 2
	var result, reminder, err = intDivision(numerator, denominator)

	switch {
		case err != nil:
			fmt.Printf(err.Error() + "\n")
		case reminder == 0:
			fmt.Printf("The result of the integer division is %v\n", result)
		default:
			fmt.Printf("The result of the integer division is %v with remainder %v\n", result, reminder)
	}
}

func printMe(printValue string) {
	fmt.Println(printValue)
}

func intDivision(numerator int, denominator int) (int, int, error) {
	var err error
	if denominator==0{
		err = errors.New("Cannot devide by Zero")
		return 0,0,err
	}
	var result int = numerator / denominator 
	var remainder int = numerator%denominator
	return result, remainder, err
}