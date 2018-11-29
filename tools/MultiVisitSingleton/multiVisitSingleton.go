package main

import (
	"fmt"
	"myproject/util"
)

func main() {
	fmt.Println("begin")
	res := util.GetCodisUtilInstance()
	fmt.Println("res: ", res)
	for i:=0; i< 100; i++ {
		fmt.Println("res: ", res)
	}
}
