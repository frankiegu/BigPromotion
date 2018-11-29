package main

func main() {
	Any(2)
	Any("666")
}

func Any(v interface{}) {
	if v2, ok := v.(string); ok {
		println(v2)
	}else if v3, ok := v.(int); ok {
		println(v3)
	}
}
