package hello
import "fmt"

// lowercase is unexported (private)
type info struct {
  name string
  age int8
}

// this is exported
type Person struct {
  Info info
  // Person.Info.name and age are private
}

func NewPerson(name string, age int8) Person {
  return Person{
    Info: info{
      name: name,
      age: age,
    },
  }
}

func (p Person) Greet() string {
  return fmt.Sprintf("Hello, %s. You are %d", p.Info.name, p.Info.age)
}
