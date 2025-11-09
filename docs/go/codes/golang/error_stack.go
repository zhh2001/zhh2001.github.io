package main

import (
	"fmt"
	"github.com/pkg/errors"
)

type Student struct {
	Name string
	Age  int
}

func (s *Student) SetName(name string) error {
	if name == "" {
		return errors.New("name is empty")
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
		return nil, errors.Wrap(err, "set name failed")
	}
	return stu, nil
}

func main() {
	_, err := NewStudent()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}
