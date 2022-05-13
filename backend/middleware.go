package imes

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
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

func (s *Middleware) OpenGithub() {
	var err error
	url := "https://github.com/i-MES"
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

type TestItem struct {
	Id       int    `json:"id"` 
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Funcname string `json:"funcname"`
}

// Load testitems from a file
func (s *Middleware) LoadTestitems(path string) []TestItem {
	tis := make([]TestItem, 0)

	tis = append(tis,
		TestItem{1, "MCU Test", "MCU Test...", "test_mcu"},
		TestItem{2, "Memory Test", "Memory Test...", "test_memory"},
		TestItem{3, "Network Test", "Network Test...", "test_network"},
	)
	return tis
}

var counter = 0

func (s *Middleware) AddCounter() int {
	counter += 1
	fmt.Println(counter)
	return counter
}

func (s *Middleware) LoadCounter() int {
	return counter
}
