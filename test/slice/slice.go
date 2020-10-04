package main

import (
	"fmt"
)

func main() {
	buff := make([]int, 10)

	n := 0
	for i := 0; i < 10 ; i++ {
		temp := []int{1,2,3,4,5,6,7,8,9,10}
		fmt.Println("buff len: ", len(buff))
		fmt.Println("buff cap: ", cap(buff))
		// buff = append(buff, temp[:len(temp)]...)
		buff = append(buff, 1)
		n += len(temp)
	}

	fmt.Printf("buff: %d\n", buff)
}