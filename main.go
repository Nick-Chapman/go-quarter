package main

import "fmt"
import "os"

func main() {
	fmt.Printf("*go-quarter*\n")

	input := readFile("../quarter-forth/f/quarter.q")

	Key := func(m *machine) {
		c := input.getChar()
		//fmt.Printf("key: %c\n", c) //echo
		fmt.Printf("%c", c) //echo
		m.push(valueOfChar(c))
	}

	SetTabEntry := func(m *machine) {
		c := input.getChar()
		m.dt[c] = m.here()
	}

	m := newMachine(Key, Dispatch)

	m.installQuarterPrim('\n', "", Nop)
	m.installQuarterPrim(' ', "", Nop)
	m.installQuarterPrim('!', "", Store)
	m.installQuarterPrim('*', "", Mul)
	m.installQuarterPrim('+', "", Add)
	m.installQuarterPrim(',', "", Comma)
	m.installQuarterPrim('-', "", Minus)
	m.installQuarterPrim('.', "", Emit)
	m.installQuarterPrim('0', "", Zero)
	m.installQuarterPrim('1', "", One)
	m.installQuarterPrim(':', "", SetTabEntry)
	m.installQuarterPrim(';', "", RetComma)
	m.installQuarterPrim('<', "", LessThan)
	m.installQuarterPrim('=', "", Equal)
	m.installQuarterPrim('>', "", CompileComma)
	m.installQuarterPrim('?', "", Dispatch)
	m.installQuarterPrim('@', "", Fetch)
	m.installQuarterPrim('A', "", CrashOnlyDuringStartup)
	m.installQuarterPrim('B', "", Branch0)
	m.installQuarterPrim('C', "", C_Fetch)
	m.installQuarterPrim('D', "", Dup)
	m.installQuarterPrim('E', "", EntryComma)
	m.installQuarterPrim('G', "", XtToNext)
	m.installQuarterPrim('H', "", HerePointer)
	m.installQuarterPrim('I', "", IsImmediate)
	m.installQuarterPrim('J', "", Jump)
	m.installQuarterPrim('L', "", Lit)
	m.installQuarterPrim('M', "", CR)
	m.installQuarterPrim('N', "", XtToName)
	m.installQuarterPrim('O', "", Over)
	m.installQuarterPrim('P', "", Drop)
	m.installQuarterPrim('V', "", Execute)
	m.installQuarterPrim('W', "", Swap)
	m.installQuarterPrim('X', "", Exit)
	m.installQuarterPrim('Y', "", IsHidden)
	m.installQuarterPrim('Z', "", Latest)
	m.installQuarterPrim('^', "", Key)
	m.installQuarterPrim('`', "", C_Comma)

	m.run()
	fmt.Printf("\n*DONE*\n")
	m.see()
}

func readFile(filename string) inputBytes {
	bs, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return inputBytes{bs, 0}
}

type inputBytes struct {
	bs []byte
	n  int
}

func (x *inputBytes) getChar() char {
	if x.n == len(x.bs) {
		return 0
	}
	n := x.n
	c := x.bs[n]
	x.n = n + 1
	return char(c)
}

func Add(m *machine) {
	panic("Add")
}

func Branch0(m *machine) {
	panic("Branch0")
}

func C_Comma(m *machine) {
	//panic("C_Comma")
	c := charOfValue(m.pop())
	m.comma(c)
}

func C_Fetch(m *machine) {
	panic("C_Fetch")
}

func Comma(m *machine) {
	v := m.pop()
	m.comma(v)
}

func CompileComma(m *machine) {
	v := m.pop()
	m.comma(call{addrOfValue(v)})
}

func CR(m *machine) {
	fmt.Printf("\n")
}

func CrashOnlyDuringStartup(m *machine) {
	panic("CrashOnlyDuringStartup")
}

func Dispatch(m *machine) {
	c := charOfValue(m.pop())
	a := m.lookupDisaptch(c)
	m.push(valueOfAddr(a))
}

func Drop(m *machine) {
	panic("Drop")
}

func Dup(m *machine) {
	panic("Dup")
}

func Emit(m *machine) {
	c := charOfValue(m.pop())
	//fmt.Printf("Emit: %v '%c'\n", c, c)
	fmt.Printf("%c", c)
}

func EntryComma(m *machine) {
	panic("EntryComma")
}

func Equal(m *machine) {
	panic("Equal")
}

func Execute(m *machine) {
	panic("Execute")
}

func Exit(m *machine) {
	panic("Exit")
}

func Fetch(m *machine) {
	a := addrOfValue(m.pop())
	slot := m.lookupMem(a)
	m.push(slot.toLiteral())
}

func HerePointer(m *machine) {
	m.push(valueOfAddr(m.hereP))
}

func IsHidden(m *machine) {
	panic("IsHidden")
}

func IsImmediate(m *machine) {
	panic("IsImmediate")
}

func Jump(m *machine) {
	panic("Jump")
}

func Latest(m *machine) {
	panic("Latest")
}

func LessThan(m *machine) {
	panic("LessThan")
}

func Lit(m *machine) {
	a := addrOfValue(m.rsPop())
	slot := m.lookupMem(a)
	m.push(slot.toLiteral())
	valueOfAddr(a.offset(2))
}

func Minus(m *machine) {
	panic("Minus")
}

func Mul(m *machine) {
	panic("Mul")
}

func Nop(m *machine) {
	//nothing
}

func One(m *machine) {
	panic("One")
}

func Over(m *machine) {
	panic("Over")
}

func RetComma(m *machine) {
	m.comma(ret{})
}

func SetTabEntry(m *machine) {
	panic("SetTabEntry")
}

func Store(m *machine) {
	panic("Store")
}

func Swap(m *machine) {
	panic("Swap")
}

func XtToName(m *machine) {
	panic("XtToName")
}

func XtToNext(m *machine) {
	panic("XtToNext")
}

func Zero(m *machine) {
	m.push(value{0})
}
