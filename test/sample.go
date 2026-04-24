package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string {
	return fmt.Sprintf("hello, %s", p.Name)
}

func Add(a, b int) int {
	return a + b
}

func main() {
	p := Person{Name: "Alice", Age: 30}
	fmt.Println(p.Greet())
	fmt.Println(Add(2, 3))
}
