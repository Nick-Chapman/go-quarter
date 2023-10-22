package main

import "os"

func readFile(filename string) inputBytes {
	bs, err := os.ReadFile(filename)
	check(err)
	return inputBytes{bs,0}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type inputBytes struct {
	bs []byte
	n int
}

func (x *inputBytes) getChar() byte {
	// TODO: handle EOF better
	n := x.n
	c := x.bs[n]
	x.n = n+1
	return c
}
