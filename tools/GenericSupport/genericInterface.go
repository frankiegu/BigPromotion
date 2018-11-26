package main

import "fmt"

func max(first interface{}, rest ...interface{}) interface{} {
	max := first
	for _, v := range rest {
		switch v := v.(type) {
		case int:
			if v > max.(int) {
				max = v
			}
		case float64:
			if v > max.(float64) {
				max = v
			}
		}
	}
	return max
}

func main() {
	fmt.Println(max(1, 2, 3, 4, 5))
	fmt.Println(max(2.1, 3.4, 6.7, 9.1))
}
