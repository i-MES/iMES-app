package testset

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

type Parser struct {
	filePath string
}

func (fp *Parser) ParsePython(tgidbase int, file string) []TestGroup {
	fp.filePath = file
	var err error
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	validClass := regexp.MustCompile(`^class\ *(.*):`)
	validFunc := regexp.MustCompile(`\ *def (test_.*)\(`)
	tgs := make([]TestGroup, 0)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			fmt.Println("EoF of", file)
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}

		cname := validClass.FindStringSubmatch(line)
		if len(cname) > 1 {
			fmt.Println("Match: ", cname[1])
			tgs = append(tgs,
				TestGroup{tgidbase + 1, cname[1], file, make([]TestItem, 0)})
		}
		fname := validFunc.FindStringSubmatch(line)
		if len(fname) > 1 {
			fmt.Println("Match: ", fname[1])
			_l := len(tgs) - 1
			tgs[_l].TestItems = append(tgs[_l].TestItems,
				TestItem{fname[1], file, fname[1], 0})
		}
	}
	return tgs
}

func (fp *Parser) ParseGolang(file string) error {
	return nil
}
