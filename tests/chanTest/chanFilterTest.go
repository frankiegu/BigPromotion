package main

import "fmt"

func main() {
	//result array
	var numbers[] int
	//channel size is 1
	ch := make(chan int)
	//the Generate goroutime use the channel: ch, ch is used once
	go Generate(ch)
	for i:=0; i< 1000; i++ {
		prime := <-ch
		numbers = append(numbers, prime)
		fmt.Println(prime, "\n")
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
		fmt.Println(<-ch1)
	}

	fmt.Printf("len=%d cap=%d slice=%v\n", len(numbers), cap(numbers), numbers)

}
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in
		fmt.Println("Filter: in pop:", i)
		if i%prime != 0 {
			fmt.Println("Filter:out add:", i)
			out <- i
		}
	}
}
// chan<- int means can only send data to channel
// <-chan int means can only receive data from channel
func Generate(ch chan<- int) {
	for i:= 2;;i++ {
		fmt.Println("Generate: i:", i)
		ch <- i
	}
}