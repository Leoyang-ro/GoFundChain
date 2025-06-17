package main

import "fmt"

/*加减乘除*/
func Add(a, b int) int {
    return a + b
}

func Subtract(a, b int) int {
	return a - b
}
func Multiply(a, b int) int {
	return a * b
}
func Divide(a, b int) int {
	if b == 0 {
		return 0 // Avoid division by zero
	}
	return a / b
}

/* 增删改查
go语言都有什么数据类型
 */
func Create() {
}

func Read() {
}

func Update() 
	a := Add(1, 2)
}

func Delete() {