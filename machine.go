package main

import "fmt"

type native func(*machine)

type primitive struct {
	name   string
	action native
}

type value struct {
	i int16
}

type char byte

type addr struct {
	u uint16
}

type ret struct {
}

type call struct {
	addr addr
}

type entry struct {
	name      addr
	next      addr
	hidden    bool
	immediate bool
}

type kdxLoop struct {
	key      native
	dispatch native
}

type machine struct {
	halt      addr
	kdx       addr
	hereP     addr
	latest    addr
	dt        map[char]addr
	mem       map[addr]slot
	psPointer addr
	rsPointer addr
	steps     uint
	echoOn	  bool
}

const psBase uint16 = 51000

func newMachine(key, dispatch native) *machine {
	mem := make(map[addr]slot)
	halt := addr{0}
	kdx := addr{1}
	hereP := addr{2}
	mem[kdx] = kdxLoop{key, dispatch}
	mem[hereP] = valueOfAddr(addr{100})
	return &machine{
		halt:      halt,
		kdx:       kdx,
		hereP:     hereP,
		latest:    addr{0},
		dt:        make(map[char]addr),
		mem:       mem,
		psPointer: addr{psBase},
		rsPointer: addr{61000},
		steps:     0,
		echoOn:	   false,
	}
}

func (m *machine) run() { // the inner interpreter (aka trampoline!)
	var a addr = m.kdx
	for {
		if a == m.halt {
			break
		}
		m.tick()
		//m.see()
		//fmt.Printf("addr=%v\n",a)
		slot := m.lookupMem(a)
		a = slot.executeSlot(m, a)
	}
}

func (m *machine) tick() {
	m.steps++
}

func (m *machine) see() {
	here := m.mem[m.hereP]
	fmt.Printf("steps: %v, here: %v, ps: %v, rs: %v\n",
		m.steps, here, m.psPointer, m.rsPointer)
}

func (m *machine) installQuarterPrim(c char, name string, action native) {
	a := m.installPrim(name, action)
	m.dt[c] = a
}

func (m *machine) installPrim(nameString string, action native) addr {
	prim := &primitive{nameString, action}
	name := m.here()
	m.commaString(prim.name)
	m.comma(entry{name, m.latest, false, false})
	xt := m.here()
	m.comma(prim)
	m.comma(ret{})
	m.latest = xt
	return xt
}

func AsEntry(slot slot) entry {
	entry, ok := slot.(entry)
	if !ok {
		panic("AsEntry/non-entry")
	}
	return entry
}

func AsChar(slot slot) char {
	char, ok := slot.(char)
	if !ok {
		panic("AsChar/non-char")
	}
	return char
}

func (m *machine) here() addr {
	return addrOfValue(m.mem[m.hereP].toLiteral())
}

func (m *machine) setHere(a addr) {
	m.mem[m.hereP] = valueOfAddr(a)
}

func (m *machine) commaString(s string) {
	for i := 0; i < len(s); i++ {
		c := char(s[i])
		m.comma(c)
	}
	m.comma(char(0))
}

func (m *machine) comma(s slot) {
	a := m.here()
	//fmt.Printf("comma: %v = %s\n", a, s.viewSlot())
	m.mem[a] = s
	m.setHere(a.offset(s.sizeSlot()))
}

func (m *machine) lookupDispatch(c char) addr {
	if c == 0 {
		return m.halt
	}
	addr, ok := m.dt[c]
	if !ok {
		panic(fmt.Sprintf("lookupDispatch: %v", c))
	}
	return addr
}

func (m *machine) lookupMem(a addr) slot {
	slot, ok := m.mem[a]
	if !ok {
		panic(fmt.Sprintf("lookupMem: %v", a))
	}
	return slot
}

func (m *machine) push(v value) {
	m.psPointer = m.psPointer.offset(-2)
	m.mem[m.psPointer] = v
}

func (m *machine) pop() value {
	slot := m.lookupMem(m.psPointer)
	m.psPointer = m.psPointer.offset(2)
	return slot.toLiteral()
}

func (m *machine) rsPush(v value) {
	m.rsPointer = m.rsPointer.offset(-2)
	m.mem[m.rsPointer] = v
}

func (m *machine) rsPop() value {
	slot := m.lookupMem(m.rsPointer)
	m.rsPointer = m.rsPointer.offset(2)
	return slot.toLiteral()
}

func (a addr) offset(n int) addr {
	return addr{a.u + uint16(n)}
}

func isZero(v value) bool {
	return v.i == 0
}

func valueOfBool(b bool) value {
	if b {
		return value{-1}
	} else {
		return value{0}
	}
}

func valueOfChar(c char) value {
	return value{int16(c)}
}

func valueOfAddr(a addr) value {
	return value{int16(a.u)}
}

func charOfValue(v value) char {
	return char(v.i % 256)
}

func addrOfValue(v value) addr {
	return addr{uint16(v.i)}
}

type slot interface {
	executeSlot(*machine, addr) addr
	toLiteral() value
	viewSlot() string
	sizeSlot() int
}

// executeSlot...

func (x kdxLoop) executeSlot(m *machine, a addr) addr {
	x.key(m)
	x.dispatch(m)
	m.rsPush(valueOfAddr(m.kdx))
	return addrOfValue(m.pop())
}

func (p primitive) executeSlot(m *machine, a addr) addr {
	//fmt.Printf("* %s\n", p.name)
	p.action(m)
	return a.offset(1)
}

func (call call) executeSlot(m *machine, a addr) addr {
	m.rsPush(valueOfAddr(a.offset(3)))
	return call.addr
}

func (ret) executeSlot(m *machine, a addr) addr {
	return addrOfValue(m.rsPop())
}

func (v value) executeSlot(*machine, addr) addr {
	panic(fmt.Sprintf("value/execute: %v", v))
}

func (char) executeSlot(*machine, addr) addr {
	panic("char/execute")
}

func (entry) executeSlot(*machine, addr) addr {
	panic("entry/execute")
}

// toLiteral...

func (kdxLoop) toLiteral() value {
	panic("kdxLoop/toLiteral")
}

func (primitive) toLiteral() value {
	panic("primitive/toLiteral")
}

func (call) toLiteral() value {
	panic("call/toLiteral")
}

func (ret) toLiteral() value {
	panic("ret/toLiteral")
}

func (entry) toLiteral() value {
	panic("entry/toLiteral")
}

func (c char) toLiteral() value {
	panic("char/toLiteral")
}

func (v value) toLiteral() value {
	return v
}

// viewSlot...

func (kdxLoop) viewSlot() string {
	return "KDXL"
}

func (p primitive) viewSlot() string {
	return fmt.Sprintf("prim: %s", p.name)
}

func (call call) viewSlot() string {
	return fmt.Sprintf("call: %v", call.addr)
}

func (ret) viewSlot() string {
	return "ret"
}

func (e entry) viewSlot() string {
	return fmt.Sprintf("entry: name=%v, next=%v", e.name, e.next)
}

func (v value) viewSlot() string {
	return fmt.Sprintf("value: %v", v)
}

func (c char) viewSlot() string {
	return fmt.Sprintf("char: %v", c)
}



func (kdxLoop) sizeSlot() int {
	return 100
}

func (p primitive) sizeSlot() int {
	return 1
}

func (call call) sizeSlot() int {
	return 3
}

func (ret) sizeSlot() int {
	return 1
}

func (e entry) sizeSlot() int {
	return 1
}

func (v value) sizeSlot() int {
	return 2
}

func (c char) sizeSlot() int {
	return 1
}


// String...

func (addr addr) String() string {
	return fmt.Sprintf("[%d]", addr.u)
}

func (c char) String() string {
	return fmt.Sprintf("%d='%c'", c, c)
}
