package main

import "fmt"
import "os"

func main() {
	fmt.Printf("*go-quarter*\n")

	bs1, err := os.ReadFile("../quarter-forth/f/quarter.q")
	if err != nil {
		panic(err)
	}
	bs2, err := os.ReadFile("../quarter-forth/f/forth.f")
	if err != nil {
		panic(err)
	}
	bs := append(bs1, bs2...)
	input := inputBytes{bs, 0}

	Key := func(m *machine) {
		c := input.getChar()
		//fmt.Printf("Key: %v\n", c) //echo
		//fmt.Printf("%c", c) //echo
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
	m.installQuarterPrim('!', "!", Store)
	m.installQuarterPrim('*', "*", Mul)
	m.installQuarterPrim('+', "+", Add)
	m.installQuarterPrim(',', ",", Comma)
	m.installQuarterPrim('-', "-", Minus)
	m.installQuarterPrim('.', "emit", Emit)
	m.installQuarterPrim('0', "0", Zero)
	m.installQuarterPrim('1', "1", One)
	m.installQuarterPrim(':', "SetTabEntry", SetTabEntry)
	m.installQuarterPrim(';', "ret,", RetComma)
	m.installQuarterPrim('<', "<", LessThan)
	m.installQuarterPrim('=', "=", Equal)
	m.installQuarterPrim('>', "compile,", CompileComma)
	m.installQuarterPrim('?', "Dispatch", Dispatch)
	m.installQuarterPrim('@', "@", Fetch)
	m.installQuarterPrim('A', "crash-only-during-startup", CrashOnlyDuringStartup)
	m.installQuarterPrim('B', "0branch", Branch0)
	m.installQuarterPrim('C', "c@", C_Fetch)
	m.installQuarterPrim('D', "dup", Dup)
	m.installQuarterPrim('E', "entry,", EntryComma)
	m.installQuarterPrim('G', "xt->next", XtToNext)
	m.installQuarterPrim('H', "here-pointer", HerePointer)
	m.installQuarterPrim('I', "immediate?", IsImmediate)
	m.installQuarterPrim('J', "jump", Jump)
	m.installQuarterPrim('L', "lit", Lit)
	m.installQuarterPrim('M', "cr", CR)
	m.installQuarterPrim('N', "xt->name", XtToName)
	m.installQuarterPrim('O', "over", Over)
	m.installQuarterPrim('P', "drop", Drop)
	m.installQuarterPrim('V', "execute", Execute)
	m.installQuarterPrim('W', "swap", Swap)
	m.installQuarterPrim('X', "exit", Exit)
	m.installQuarterPrim('Y', "hidden?", IsHidden)
	m.installQuarterPrim('Z', "latest", Latest)
	m.installQuarterPrim('^', "key", Key)
	m.installQuarterPrim('`', "c,", C_Comma)

	m.installPrim("immediate^", FlipImmediate)
	m.installPrim("hidden^", FlipHidden)
	m.installPrim("branch", Branch)
	m.installPrim("xor", Xor)
	m.installPrim("crash", Crash)
	m.installPrim(">r", ToReturnStack)
	m.installPrim("r>", FromReturnStack)
	m.installPrim("/mod", DivMod)
	m.installPrim("key?", KeyNonBlocking)
	m.installPrim("c!", C_Store)

	m.run()
	fmt.Printf("\n*DONE*\n")
	m.see()
}

type inputBytes struct {
	bs []byte
	n  int
}

func (x *inputBytes) getChar() char {
	if x.n == len(x.bs) {
		panic("EOF")
	}
	n := x.n
	c := x.bs[n]
	x.n = n + 1
	return char(c)
}

func Add(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	m.push(value{v1.i + v2.i})
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
	a := addrOfValue(m.pop())
	slot := m.lookupMem(a)
	char := AsChar(slot)
	m.push(valueOfChar(char))
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
	m.pop()
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
	m.comma(entry{addrOfValue(v), m.latest, false, false})
	m.latest = m.here()
}

func Equal(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	m.push(valueOfBool(v1.i == v2.i))
}

func Execute(m *machine) {
	v := m.pop()
	m.rsPush(v)
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
	a := addrOfValue(m.pop()).offset(-1)
	slot := m.lookupMem(a)
	entry := AsEntry(slot)
	b := entry.hidden
	m.push(valueOfBool(b))
}

func IsImmediate(m *machine) {
	a := addrOfValue(m.pop()).offset(-1)
	slot := m.lookupMem(a)
	entry := AsEntry(slot)
	b := entry.immediate
	m.push(valueOfBool(b))
}

func Jump(m *machine) {
	v := m.pop()
	m.rsPop()
	m.rsPush(v)
}

func Latest(m *machine) {
	m.push(valueOfAddr(m.latest))
}

func LessThan(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	if v1.i < v2.i {
		m.push(value{-1})
	} else {
		m.push(value{0})
	}
}

func Lit(m *machine) {
	a := addrOfValue(m.rsPop())
	slot := m.lookupMem(a)
	m.push(slot.toLiteral())
	m.rsPush(valueOfAddr(a.next()))
}

func Minus(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	m.push(value{v1.i - v2.i})
}

func Mul(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	m.push(value{v1.i * v2.i})
}

func Nop(m *machine) {
	//nothing
}

func One(m *machine) {
	m.push(value{1})
}

func Over(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	m.push(v1)
	m.push(v2)
	m.push(v1)
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
	a := addrOfValue(m.pop()).offset(-1)
	slot := m.lookupMem(a)
	entry := AsEntry(slot)
	m.push(valueOfAddr(entry.name))
}

func XtToNext(m *machine) {
	a := addrOfValue(m.pop()).offset(-1)
	slot := m.lookupMem(a)
	entry := AsEntry(slot)
	m.push(valueOfAddr(entry.next))
}

func Zero(m *machine) {
	m.push(value{0})
}

func FlipImmediate(m *machine) {
	a := addrOfValue(m.pop()).offset(-1)
	slot := m.lookupMem(a)
	e := AsEntry(slot)
	m.mem[a] = entry{e.name, e.next, e.hidden, !e.immediate}
}

func FlipHidden(m *machine) {
	a := addrOfValue(m.pop()).offset(-1)
	slot := m.lookupMem(a)
	e := AsEntry(slot)
	m.mem[a] = entry{e.name, e.next, !e.hidden, e.immediate}
}

func Branch(m *machine) {
	a := addrOfValue(m.rsPop())
	slot := m.lookupMem(a)
	n := int(slot.toLiteral().i)
	m.rsPush(valueOfAddr(a.offset(n)))
}

func Xor(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	m.push(value{v1.i ^ v2.i})
}

func Crash(m *machine) {
	panic("Crash")
}

func ToReturnStack(m *machine) {
	panic("ToReturnStack")
}

func FromReturnStack(m *machine) {
	panic("FromReturnStack")
}

func DivMod(m *machine) {
	panic("DivMod")
}

func KeyNonBlocking(m *machine) {
	panic("KeyNonBlocking")
}

func C_Store(m *machine) {
	panic("C_Store")
}
