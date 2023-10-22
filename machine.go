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
	name addr
}

type kdxLoop struct {
	key      native
	dispatch native
}

type machine struct {
	halt      addr
	kdx       addr
	hereP     addr
	dt        map[char]addr
	mem       map[addr]slot
	psPointer addr
	rsPointer addr
	steps     uint
}

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
		dt:        make(map[char]addr),
		mem:       mem,
		psPointer: addr{51000},
		rsPointer: addr{61000},
		steps:     0,
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
	p := &primitive{name, action}
	a := m.installPrim(p)
	m.dt[c] = a
}

func (m *machine) installPrim(prim *primitive) addr {
	// TODO: write name & entry to allow dictionary lookup
	// for now we will just write the native-slot code
	a := addrOfValue(m.mem[m.hereP].toLiteral())
	m.comma(prim)
	m.comma(ret{})
	return a
}

func (m *machine) here() addr {
	return addrOfValue(m.mem[m.hereP].toLiteral())
}

func (m *machine) setHere(a addr) {
	m.mem[m.hereP] = valueOfAddr(a)
}

func (m *machine) comma(s slot) {
	a := m.here()
	//fmt.Printf("comma: %v = %s\n", a, s.viewSlot())
	m.mem[a] = s
	m.setHere(a.next())
}

func (m *machine) lookupDispatch(c char) addr {
	//fmt.Printf("lookupDispatch: %v '%c'\n", c, c)
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

func (a addr) next() addr {
	return a.offset(1)
}

func (a addr) offset(n int) addr {
	return addr{a.u + uint16(n)}
}

func isZero(v value) bool {
	return v.i == 0
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
	return a.next()
}

func (call call) executeSlot(m *machine, a addr) addr {
	m.rsPush(valueOfAddr(a.next()))
	return call.addr
}

func (ret) executeSlot(m *machine, a addr) addr {
	return addrOfValue(m.rsPop())
}

func (value) executeSlot(*machine, addr) addr {
	panic("value/execute")
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

func (v value) toLiteral() value {
	return v
}

func (c char) toLiteral() value {
	return valueOfChar(c)
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

func (entry) viewSlot() string {
	return "entry"
}

func (v value) viewSlot() string {
	return fmt.Sprintf("value: %v", v)
}

func (c char) viewSlot() string {
	return fmt.Sprintf("char: %v", c)
}

// String...

func (addr addr) String() string {
	return fmt.Sprintf("[%d]",addr.u)
}

func (c char) String() string {
	return fmt.Sprintf("%d='%c'",c,c)
}

