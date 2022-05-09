package imes

import (
	"fmt"
)

//Middleware struct to hold wails runtime for all middleware implementations
type Middleware struct {
}

func (s *Middleware) OpenFile(Hash string) bool {
	fmt.Println("OpenFile")
	return true
}

func (s *Middleware) OpenLink(Link string) bool {
	fmt.Println("OpenLink")
	return true
}

func (s *Middleware) OpenLog() bool {
	fmt.Println("OpenLog")
	return true
}

func (s *Middleware) OpenFolder(Hash string) bool {
	fmt.Println("OpenFolder")
	return true
}
