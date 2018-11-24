package main

import (
	"encoding/json"
	"fmt"
)

type Parent struct {
	Father string
	Mother string
}

type People struct {
	Name string
	Age int
	parent *Parent
}

func StructToJsonFunc()  {
	p := &People{"xiaodong", 18, &Parent{"a", "b"}}

	jsonBytes, err := json.Marshal(p)

	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("result: ", string(jsonBytes))
}

func main() {
	StructToJsonFunc()
}


