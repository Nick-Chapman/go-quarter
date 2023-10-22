package main

import "fmt"

func main() {
	fmt.Printf("*go-quarter*\n")
	input := readFile("../quarter-forth/f/quarter.q")
	m := newMachine(&input)
	m.installQuarterPrim('^',makePrim("key",key))
	m.installQuarterPrim('.',makePrim("emit",emit))
	m.run()
	fmt.Printf("*go-quarter*DONE\n")
}

func key(m *machine) {
	c := m.getChar()
	//fmt.Printf("key: %c\n", c)
	m.push(valueOfChar(c))
}

func emit(m *machine) {
	v := m.pop()
	c := charOfValue(v)
	//fmt.Printf("emit: %c\n", c)
	fmt.Printf("%c", c)
}
