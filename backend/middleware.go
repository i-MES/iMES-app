package imes

import (
	"fmt"
)

//Middleware struct to hold wails runtime for all middleware implementations
type Middleware struct {
}

func (s *Middleware) OpenFile(Hash string) {
	fmt.Println("OpenFile")
}

func (s *Middleware) OpenLink(Link string) {
	fmt.Println("OpenLink")

}

func (s *Middleware) OpenLog() {
	fmt.Println("OpenLog")

}

func (s *Middleware) OpenFolder(Hash string) {
	fmt.Println("OpenFolder")

}
