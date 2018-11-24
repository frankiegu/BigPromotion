package main

import "fmt"

type query struct {
	//param channel
	sql chan string
	//result channel
	result chan string
}

func execQuery(q query) {
	go func() {
		sql := <-q.sql
		q.result <- "get " + sql
	}()
}

func main() {
	q := query{make(chan string, 1), make(chan string, 2)}
	execQuery(q)

	q.sql <- "select * from table"
	fmt.Println(<-q.result)
}
