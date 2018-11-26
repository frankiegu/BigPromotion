package main

import (
	"fmt"
	"reflect"
)

func index(slice interface{}, v interface{}) int {
	fmt.Println("slice: ", slice)
	fmt.Println("v: ", v)
	//println("reflect slice: ", reflect.ValueOf(slice))

	if slice := reflect.ValueOf(slice); slice.Kind() == reflect.Slice {
		for i := 0; i< slice.Len(); i++ {
			if reflect.DeepEqual(v, slice.Index(i).Interface()) {
				return i
			}
		}
	}
	return -1
}

func main() {
	slice1 := []int{0, 4, 5, 7, 9}
	//fmt.Println(slice1, 7)
	res := index(slice1, 7)
	println("res: ", res)
}
