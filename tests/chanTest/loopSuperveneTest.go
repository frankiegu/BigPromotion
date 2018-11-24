package main

import "fmt"


func doSomething(i int, xi int) {
	fmt.Printf("i=%d, xi=%d\n", i, xi)
}

func main() {
	data := []int{1, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3, 2, 3, 1, 2, 3}
	N := len(data)
	sem := make(chan int, N)
	for i, xi := range data {
		go func() {
			doSomething(i, xi)
			sem <- 0
		}()
	}

	fmt.Println("loop supervene start ... ")

	for i := 0; i< N; i++ {
		<-sem
	}

	fmt.Println("end")

}