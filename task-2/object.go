package task2

import (
	"fmt"
	"math"
)

type Shape interface {
	// 计算面积
	Area() float64
	// 计算周长
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (e *Employee) PrintInfo() {
	fmt.Println("姓名:", e.Name, "年龄:", e.Age, "员工ID:", e.EmployeeID)
}

func RunObject() {
	var it Shape = &Rectangle{Width: 10, Height: 20}
	var circle Shape = &Circle{Radius: 10}
	fmt.Println("Rectangle Area:", it.Area())
	fmt.Println("Rectangle Perimeter:", it.Perimeter())
	fmt.Println("Circle Area:", circle.Area())
	fmt.Println("Circle Perimeter:", circle.Perimeter())

	employee := &Employee{Person: Person{Name: "张三", Age: 20}, EmployeeID: 123456}
	employee.PrintInfo()
}
