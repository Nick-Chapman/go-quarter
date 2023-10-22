package main

import "fmt"
import "os"

func main() {
	fmt.Printf("*go-quarter*\n")
	input := readFile("../quarter-forth/f/quarter.q")
	key := func(m *machine) {
		c := input.getChar()
		//fmt.Printf("%c", c) //echo
		m.push(valueOfChar(c))
	}
	m := newMachine()
	m.installQuarterPrim('^', makePrim("key", key))
	m.installQuarterPrim('.', makePrim("emit", emit))
	m.run()
	fmt.Printf("\n*DONE*\n")
	m.see()
}

func readFile(filename string) inputBytes {
	bs, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return inputBytes{bs[0:100], 0}
}

type inputBytes struct {
	bs []byte
	n  int
}

func (x *inputBytes) getChar() byte {
	if x.n == len(x.bs) {
		return 0
	}
	n := x.n
	c := x.bs[n]
	x.n = n + 1
	return c
}

func emit(m *machine) {
	v := m.pop()
	c := charOfValue(v)
	//fmt.Printf("emit: %v '%c'\n", c, c)
	fmt.Printf("%c", c)
}
