package main

import (
	"fmt"
)



func main() {
	//ARRAYS AND SLICES

	// var intArr [3]int32
	// intArr := [...]int32{1,2,3} // this is another way of initialiting the array
	var intArr [3]int32 = [3]int32{1,2,3} // this is another way of initialiting the array

	fmt.Println(intArr[0])
	// fmt.Println(intArr[1:3])
	
	// fmt.Println(&intArr[0])
	// fmt.Println(&intArr[1])
	// fmt.Println(&intArr[2])

	var intSlice []int32 = []int32{4,5,6}
	fmt.Printf("The length is %x with capaticy %x\n", len(intSlice), cap(intSlice))
	intSlice = append(intSlice, 7)
	// fmt.Println(intSlice)
	fmt.Printf("The length is %x with capaticy %x\n", len(intSlice), cap(intSlice))


	var intSlice2 []int32 = []int32{8,9}
	intSlice = append(intSlice, intSlice2...)
	fmt.Println(intSlice)

	var intSlice3 []int32 = make([]int32, 3, 10)
	fmt.Println(intSlice3)

	// MAPS
	var myMap map[string]uint8 = make(map[string]uint8)
	fmt.Println(myMap)
	
	var myMap2 = map[string]uint8{"Adam":23, "Sarah":45}
	fmt.Println(myMap2["Adam"])
	var age,ok = myMap2["Adam"] // the parameter OK in this case is a boolean, true if "Adam" exists or false if it doesnt

	// delete(myMap2, "Adam")
	if ok {
		fmt.Printf("The age is %v\n", age)
	} else {
		fmt.Printf("Invalid Name\n")
	}

	for name, age := range myMap2 {
		fmt.Printf("Name: %v, Age: %v\n", name, age)
	}

	for i , v := range intArr { 
		fmt.Printf("Index: %v, Value: %v\n", i , v)
	}

	for i:=0; i < 10; i++ {
		fmt.Println(i)
	}

}