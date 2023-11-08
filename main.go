package main

import "fmt"
import "os"
import "strings"
import "path/filepath"

func main() {
	fmt.Printf("*go-quarter*\n")

	listFile := os.Args[1]
	files := readListFile(listFile)

	bytes := readFiles(files)
	input := inputBytes{bytes, 0}

	Key := func(m *machine) {
		c := input.getChar()
		if isTrue(m.readValue(m.echoEnabledP)) {
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
}

func readListFile(listFile string) []string {
	bs, err := os.ReadFile(listFile)
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(listFile)
	var acc []string
	for _, line := range strings.Split(string(bs), "\n") {
		words := strings.Split(line, "#")
		filename := strings.TrimSpace(words[0])
		if len(filename) > 0 {
			prefixed := filepath.Join(dir, filename)
			acc = append(acc, prefixed)
		}
	}
	return acc
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
