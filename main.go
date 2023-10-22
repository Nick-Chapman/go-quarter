package main

import "fmt"
import "os"

func main() {
	fmt.Printf("*go-quarter*\n")

	input := readFile("../quarter-forth/f/quarter.q")

	Key := func(m *machine) {
		c := input.getChar()
		//fmt.Printf("Key: %c\n", c) //echo
		fmt.Printf("%c", c) //echo
		m.push(valueOfChar(c))
	}

	SetTabEntry := func(m *machine) {
		c := input.getChar()
		//fmt.Printf("SetTabEntry: %c\n", c)
		m.dt[c] = m.here()
	}

	m := newMachine(Key, Dispatch)

	m.installQuarterPrim('\n', "NopNL", Nop)
	m.installQuarterPrim(' ', "NopSpace", Nop)
	m.installQuarterPrim('!', "Store", Store)
	m.installQuarterPrim(',', "Comma", Comma)
	m.installQuarterPrim('-', "Minus", Minus)
	m.installQuarterPrim('.', "Emit", Emit)
	m.installQuarterPrim('0', "Zero", Zero)
	m.installQuarterPrim(':', "SetTabEntry", SetTabEntry)
	m.installQuarterPrim(';', "RetComma", RetComma)
	m.installQuarterPrim('=', "Equal", Equal)
	m.installQuarterPrim('>', "CompileComma", CompileComma)
	m.installQuarterPrim('?', "Dispatch", Dispatch)
	m.installQuarterPrim('@', "Fetch", Fetch)
	m.installQuarterPrim('B', "Branch0", Branch0)
	m.installQuarterPrim('D', "Dup", Dup)
	m.installQuarterPrim('E', "EntryComma", EntryComma)
	m.installQuarterPrim('H', "HerePointer", HerePointer)
	m.installQuarterPrim('J', "Jump", Jump)
	m.installQuarterPrim('L', "Lit", Lit)
	m.installQuarterPrim('M', "CR", CR)
	m.installQuarterPrim('W', "Swap", Swap)
	m.installQuarterPrim('X', "Exit", Exit)
	m.installQuarterPrim('^', "Key", Key)
	m.installQuarterPrim('`', "C_Comma", C_Comma)

	// m.installQuarterPrim('*', "Mul", Mul)
	// m.installQuarterPrim('+', "Add", Add)
	// m.installQuarterPrim('1', "One", One)
	// m.installQuarterPrim('<', "LessThan", LessThan)
	// m.installQuarterPrim('A', "CrashOnlyDuringStartup", CrashOnlyDuringStartup)
	// m.installQuarterPrim('C', "C_Fetch", C_Fetch)
	// m.installQuarterPrim('G', "XtToNext", XtToNext)
	// m.installQuarterPrim('I', "IsImmediate", IsImmediate)
	// m.installQuarterPrim('N', "XtToName", XtToName)
	// m.installQuarterPrim('O', "Over", Over)
	// m.installQuarterPrim('P', "Drop", Drop)
	// m.installQuarterPrim('V', "Execute", Execute)
	// m.installQuarterPrim('Y', "IsHidden", IsHidden)
	// m.installQuarterPrim('Z', "Latest", Latest)

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
	a := addrOfValue(m.rsPop())
	v := m.pop()
	if isZero(v) {
		slot := m.lookupMem(a)
		n := int(slot.toLiteral().i)
		m.rsPush(valueOfAddr(a.offset(n)))
	} else {
		m.rsPush(valueOfAddr(a.next()))
	}
}

func C_Comma(m *machine) {
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
	a := m.lookupDispatch(c)
	m.push(valueOfAddr(a))
}

func Drop(m *machine) {
	panic("Drop")
}

func Dup(m *machine) {
	v := m.pop()
	m.push(v)
	m.push(v)
}

func Emit(m *machine) {
	c := charOfValue(m.pop())
	//fmt.Printf("Emit: %v\n", c)
	fmt.Printf("%c", c)
}

func EntryComma(m *machine) {
	v := m.pop()
	m.comma(entry{addrOfValue(v)})
}

func Equal(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	if v1.i == v2.i {
		m.push(value{-1})
	} else {
		m.push(value{0})
	}
}

func Execute(m *machine) {
	panic("Execute")
}

func Exit(m *machine) {
	m.rsPop()
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
	v := m.pop()
	m.rsPop()
	m.rsPush(v)
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
	m.rsPush(valueOfAddr(a.next())) //AGGGH, was 2 and forgot to rs-push
}

func Minus(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	m.push(value{v1.i - v2.i})
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

func Store(m *machine) {
	a := addrOfValue(m.pop())
	v := m.pop()
	m.mem[a] = v
}

func Swap(m *machine) {
	v1 := m.pop()
	v2 := m.pop()
	m.push(v1)
	m.push(v2)
}

func XtToName(m *machine) {
	panic("XtToName")
	/*(a := addrOfValue(m.pop()).offset(-1)
	fmt.Printf("XtToName: %v\n", a)
	slot := m.lookupMem(a)
	entry, ok := slot.(entry)
	if !ok {
		panic("XtToName/non-entry")
	}
	m.push(valueOfAddr(entry.name))*/
}

func XtToNext(m *machine) {
	panic("XtToNext")
}

func Zero(m *machine) {
	m.push(value{0})
}
