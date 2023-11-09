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

	locate := map[string]uint16{

		// TODO: from beeb-asm
		"dispatch":                  0x12DD,
		"bye":                       0x1310,
		"crash":                     0x131F,
		"startup-is-complete":       0x1364,
		"crash-only-during-startup": 0x1389,
		"sp":                        0x1397,
		"sp0":                       0x13AD,
		"rsp":                       0x13C1,
		"rsp0":                      0x13DF,
		"as-num":                    0x13F6,
		"dup":                       0x1400,
		"swap":                      0x1415,
		"drop":                      0x1430,
		"over":                      0x143D,
		">r":                        0x1450,
		"r>":                        0x146D,
		"0":                         0x1489,
		"1":                         0x1499,
		"xor":                       0x14AD,
		"/2":                        0x14C4,
		"+":                         0x14D6,
		"-":                         0x14ED,
		"*":                         0x1504,
		"/mod":                      0x1540,
		"<":                         0x156F,
		"=":                         0x1592,
		"@":                         0x15B3,
		"!":                         0x15CF,
		"c@":                        0x15F0,
		"c!":                        0x1609,
		"here-pointer":              0x162E,
		",":                         0x1640,
		"c,":                        0x1663,
		"lit":                       0x167B,
		"execute":                   0x16AF,
		"jump":                      0x16C8,
		"exit":                      0x16D7,
		"0branch":                   0x16E7,
		"branch":                    0x1728,
		"ret,":                      0x1755,
		"compile,":                  0x1770,
		"xt->name":                  0x17A5,
		"xt->next":                  0x17CA,
		"immediate?":                0x17F1,
		"hidden?":                   0x1816,
		"immediate^":                0x183E,
		"hidden^":                   0x1860,
		"entry,":                    0x1881,
		"latest":                    0x18D8,
		"key":                       0x18EE,
		"set-key":                   0x190C,
		"get-key":                   0x1926,
		"echo-enabled":              0x1945,
		"echo-off":                  0x195E,
		"echo-on":                   0x1971,
		"emit":                      0x1981,
		"cr":                        0x1991,
		"time":                      0x199E,
		"key?":                      0x19B7,
		"fx":                        0x19F2,
		"mode":                      0x1A16,
	}

	m := newMachine(locate, Key, Dispatch)
	m.setupPrims(Key, SetTabEntry)

	here_start := addr{0x1b24} // TODO: from beeb-asm
	m.run(here_start)
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
