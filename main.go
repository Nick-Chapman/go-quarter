package main

import "fmt"
import "os"

func main() {
	fmt.Printf("*go-quarter*\n")

	input := readFile("../quarter-forth/f/quarter.q")

	key := func(m *machine) {
		c := input.getChar()
		//fmt.Printf("key: %c\n", c) //echo
		m.push(valueOfChar(c))
	}

	m := newMachine(key, dispatch)

	m.installQuarterPrim('^', "key", key)
	m.installQuarterPrim('.', "emit", emit)
	m.installQuarterPrim('?', "dispatch", dispatch)
	//m.installQuarterPrim('V', "execute", execute)
	m.installQuarterPrim('M', "cr", cr)
	m.installQuarterPrim(10, "nop", nop)
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

func dispatch(m *machine) {
	//fmt.Printf("dispatch\n")
	c := charOfValue(m.pop())
	a := m.lookupDisaptch(c)
	m.push(valueOfAddr(a))
}

func emit(m *machine) {
	c := charOfValue(m.pop())
	fmt.Printf("emit: %v '%c'\n", c, c)
	//fmt.Printf("%c", c)
}

func cr(m *machine) {
	//fmt.Printf("{cr}\n")
	fmt.Printf("\n")
}

func nop(*machine) {
}
