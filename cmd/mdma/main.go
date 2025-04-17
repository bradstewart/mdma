package main

import "fmt"
import "github.com/bradstewart/mdma/internal/hello"

func main() {

  person := hello.NewPerson("Jason", 42)

  fmt.Println("Hello DJ")
  fmt.Println(person.Greet())
}
