package testdata

import "fmt"

type MyInterface interface {
	Print()
}

type MyStruct1 struct {
	X int
}

func (s MyStruct1) Print() {
	fmt.Println("MyStruct1")
}

type MyStruct2 struct {
	Y string
}

func (s MyStruct2) Print() {
	fmt.Println("MyStruct2")
}

type MyStruct3 struct {
	Z float64
}

func (s MyStruct3) Print() {
	fmt.Println("MyStruct3")
}
