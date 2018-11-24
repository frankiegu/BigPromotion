package main

import (
	"fmt"
	"math/rand"
)

func main()  {

	fmt.Println("random generator 1:")
	fmt.Println(random_generator_1())

	//fmt.Println("random generator 2:")
	//rand_service_handler := random_generator_2()
	//fmt.Printf("%d\n", <-rand_service_handler)
	//fmt.Printf("%d\n", <-rand_service_handler)
	//fmt.Printf("%d\n", <-rand_service_handler)

	//multiplexing
	fmt.Println("random generator 3:")
	rand_service_handler := random_generator_3()
	fmt.Println("%d\n", <-rand_service_handler)
	fmt.Println("%d\n", <-rand_service_handler)
	fmt.Println("%d\n", <-rand_service_handler)



}
func random_generator_3() chan int {
	//create two random data generation service
	rand_service_handler_1 := random_generator_2()
	rand_service_handler_2 := random_generator_2()

	//create chan
	out := make(chan int)
	//create coroutine
	go func() {
		for {
			out <- <-rand_service_handler_1
		}
	}()

	go func() {
		for {
			out <- <-rand_service_handler_2
		}
	}()

	return out
}

func random_generator_2() chan int {
	//create chan
	out := make(chan int)
	//create coroutine
	go func() {
		//
		for {
			out <- rand.Int()
		}
	}()
	return out

}
func random_generator_1() int {
	return rand.Int()
}
