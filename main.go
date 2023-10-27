package main

import "fmt"
import "os"

func main() {
	fmt.Printf("*go-quarter*\n")

	prefix := "../quarter-forth/f/"

	bytes := readFiles([]string{
		prefix + "quarter.q",
		prefix + "forth.f",
		prefix + "tools.f",
		prefix + "regression.f",
		prefix + "examples.f",
		prefix + "primes.f",
		//prefix + "snake.f",
		//prefix + "buffer.f",
		prefix + "start.f",
	})
	input := inputBytes{bytes, 0}

	Key := func(m *machine) {
		c := input.getChar()
		if m.echoOn {
			fmt.Printf("%c", c)
		}
		m.push(valueOfChar(c))
	}

	SetTabEntry := func(m *machine) {
		c := input.getChar()
		m.dt[c] = m.here()
	}

	m := newMachine(Key, Dispatch)
	m.setupPrims(Key, SetTabEntry)
	m.run()
	fmt.Printf("\n*DONE*\n")
	m.see()
}

func readFiles(files []string) []byte {
	var acc []byte
	for _, file := range files {
		bs, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		acc = append(acc, bs...)
	}
	return acc
}

type inputBytes struct {
	bs []byte
	n  int
}

func (x *inputBytes) getChar() char {
	if x.n == len(x.bs) {
		fmt.Printf("*EOF*\n")
		os.Exit(0)
	}
	n := x.n
	c := x.bs[n]
	x.n = n + 1
	return char(c)
}
