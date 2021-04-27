package main

import "fmt"

type A struct {a int}
type B struct {c, b int}

type C struct {A; B}

type engine interface {
	start()
	stop()
}

type Car struct {
	wheelCount int
	engine
}

func (self *Car) numberOfWheels() int {
	return self.wheelCount
}

func (self *Car) start() {
	fmt.Println("start engine! Buzz~ Buzz~ Buzz~")
}

func (self *Car) stop() {
	fmt.Println("Ops~ The engine is stopped.")
}

type Mercedes struct {
	Car
}

func (m *Mercedes) sayHiToMerkel() {
	fmt.Println("Mercedes Bens ! I will get it.")
}

func main() {
	c := Car{
		wheelCount: 4,
	}

	fmt.Println(c.numberOfWheels())

	c.start()
	c.stop()

	m := Mercedes{}
	m.sayHiToMerkel()
}
