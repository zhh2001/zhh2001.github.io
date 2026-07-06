package main

import (
	"errors"
	"fmt"
	"runtime/debug"
)

var errEmptyName = errors.New("name is empty")

type Student struct {
	Name string
	Age  int
}

func (s *Student) SetName(name string) error {
	if name == "" {
		return errEmptyName
	}
	s.Name = name
	return nil
}

func NewStudent() (*Student, error) {
	stu := &Student{
		Age: 18,
	}
	err := stu.SetName("")
	if err != nil {
		return nil, fmt.Errorf("set name: %w", err)
	}
	return stu, nil
}

func main() {
	_, err := NewStudent()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		fmt.Printf("is empty name: %t\n", errors.Is(err, errEmptyName))
		fmt.Printf("stack:\n%s", debug.Stack())
	}
}
