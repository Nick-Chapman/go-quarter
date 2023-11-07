package main

import "fmt"

func (m *machine) setupPrims(Key, SetTabEntry native) {

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

	m.installPrim("/2", BitShiftRight)
	m.installPrim("/mod", DivMod)
	m.installPrim(">r", ToReturnStack)
	m.installPrim("KEY",KEY)
	m.installPrim("as-num", Nop)
	m.installPrim("branch", Branch)
	m.installPrim("c!", C_Store)
	m.installPrim("cls",Cls)
	m.installPrim("crash", Crash)
	m.installPrim("echo-enabled", EchoEnabled)
	m.installPrim("echo-off", EchoOff)
	m.installPrim("echo-on", EchoOn)
	m.installPrim("fx", Fx)
	m.installPrim("get-key", GetKey)
	m.installPrim("hidden^", FlipHidden)
	m.installPrim("immediate^", FlipImmediate)
	m.installPrim("key?", KeyNonBlocking)
	m.installPrim("mode", Mode)
	m.installPrim("r>", FromReturnStack)
	m.installPrim("read-char-col",ReadCharCol)
	m.installPrim("rsp", ReturnStackPointer)
	m.installPrim("rsp0", ReturnStackPointerBase)
	m.installPrim("set-cursor-position",SetCursorPosition)
	m.installPrim("set-cursor-shape",SetCursorShape)
	m.installPrim("set-key",SetKey)
	m.installPrim("sp", Sp)
	m.installPrim("sp0", Sp0)
	m.installPrim("startup-is-complete", StartupIsComplete)
	m.installPrim("time", Time)
	m.installPrim("write-char-col",WriteCharCol)
	m.installPrim("xor", Xor)
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
		m.rsPush(valueOfAddr(a.offset(2)))
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
	if !m.startupComplete {
		panic("CrashOnlyDuringStartup")
	}
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
	if int16(v1.i) < int16(v2.i) {
		m.push(value{65535})
	} else {
		m.push(value{0})
	}
}

func Lit(m *machine) {
	a := addrOfValue(m.rsPop())
	slot := m.lookupMem(a)
	m.push(slot.toLiteral())
	m.rsPush(valueOfAddr(a.offset(2)))
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
	b := m.rsPop()
	a := m.pop()
	m.rsPush(a)
	m.rsPush(b)
}

func FromReturnStack(m *machine) {
	b := m.rsPop()
	a := m.rsPop()
	m.push(a)
	m.rsPush(b)
}

func DivMod(m *machine) {
	v2 := m.pop()
	v1 := m.pop()
	m.push(value{v1.i % v2.i})
	m.push(value{v1.i / v2.i})
}

func KeyNonBlocking(m *machine) {
	m.push(value{0})
	fmt.Println("{KeyNonBlocking}")
}

func C_Store(m *machine) {
	a := addrOfValue(m.pop())
	c := charOfValue(m.pop())
	m.mem[a] = c
}

func BitShiftRight(m *machine) {
	v := m.pop()
	m.push(value{v.i / 2})
}

func Sp(m *machine) {
	m.push(valueOfAddr(m.psPointer))
}

func Sp0(m *machine) {
	m.push(valueOfAddr(addr{psBase}))
}

func ReturnStackPointer(m *machine) {
	panic("ReturnStackPointer")
}

func ReturnStackPointerBase(m *machine) {
	panic("ReturnStackPointerBase")
}

func GetKey(m *machine) {
	m.push(value{10000}) // TODO: NO!!
}

func Time(m *machine) {
	m.push(value{123})
	m.push(value{456})
	fmt.Println("{Time}")
}

func StartupIsComplete(m *machine) {
	m.startupComplete = true
}

func EchoOn(m *machine) {
	m.mem[m.echoEnabledP] = valueOfBool(true)
}

func EchoOff(m *machine) {
	m.mem[m.echoEnabledP] = valueOfBool(false)
}

func EchoEnabled(m *machine) {
	m.push(valueOfAddr(m.echoEnabledP))
}

func SetCursorShape(m *machine) {
	m.pop()
	fmt.Println("{SetCursorShape}")
}

func SetCursorPosition(m *machine) {
	m.pop()
	fmt.Println("{SetCursorPosition}")
}

func ReadCharCol(m *machine) {
	m.push(value{1})
	m.push(value{2})
	fmt.Println("{ReadCharCol}")
}

func WriteCharCol(m *machine) {
	m.pop()
	m.pop()
	fmt.Println("{WriteCharCol}")
}

func Cls(m *machine) {
	fmt.Println("{Cls}")
}

func KEY(m *machine) {
	panic("KEY")
}

func SetKey(m *machine) {
	m.pop()
	fmt.Println("{SetKey}")
}

func Fx(m *machine) {
	m.pop()
	m.pop()
	m.pop()
	fmt.Println("{Fx}")
}

func Mode(m *machine) {
	panic("Mode")
}
