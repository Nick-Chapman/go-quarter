package main

import "fmt"

type native func(*machine)

type primitive struct {
	name   string
	action native
}

func makePrim(n string, f func(*machine)) *primitive {
	return &primitive{n, f}
}

type machine struct {
	steps     uint
	dt        map[byte]addr
	mem       map[addr]slot
	here      addr
	psPointer addr
	rsPointer addr
}

func newMachine(key, dispatch native) *machine {
	mem := make(map[addr]slot)
	mem[addr{0}] = kdxLoop{key, dispatch}
	return &machine{
		steps:     0,
		dt:        make(map[byte]addr),
		mem:       mem,
		here:      addr{100},
		psPointer: addr{50000},
		rsPointer: addr{60000},
	}
}

func (m *machine) run() { // the inner interpreter (aka trampoline!)
	var a addr = addr{0}
	for {
		m.tick()
		//m.see()
		//fmt.Printf("addr=%v\n",a)
		slot := m.lookupMem(a)
		a = slot.executeSlot(m, a.next())
	}
}

func (m *machine) tick() {
	m.steps++
}

func (m *machine) see() {
	fmt.Printf("steps: %v, here: %v, ps: %v, rs: %v\n",
		m.steps, m.here, m.psPointer, m.rsPointer)
}

func (m *machine) installQuarterPrim(c byte, name string, native native) {
	p := makePrim(name, native)
	a := m.installPrim(p)
	m.dt[c] = a
}

func (m *machine) installPrim(prim *primitive) addr {
	// TODO: write name & entry to allow dictionary lookup
	// for now we will just write the native-slot code
	a := m.here
	m.comma(prim.action)
	m.comma(ret{})
	return a
}

func (m *machine) comma(s slot) {
	a := m.here
	m.mem[a] = s
	m.here = a.next()
}

func (m *machine) lookupDisaptch(c byte) addr {
	addr, ok := m.dt[c]
	if !ok {
		panic(fmt.Sprintf("lookupDisaptch: %v '%c'", c, c))
	}
	return addr
}

func (m *machine) lookupMem(a addr) slot {
	slot, ok := m.mem[a]
	//fmt.Println("lookupMem",a,"->",slot,ok)
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

type addr struct {
	u uint16
}

func (a addr) next() addr {
	return a.offset(1)
}

func (a addr) offset(n int) addr {
	return addr{a.u + uint16(n)}
}

type value struct {
	i int16
}

func isZero(v value) bool {
	return v.i == 0
}

func valueOfChar(c byte) value {
	return value{int16(c)}
}

func valueOfAddr(a addr) value {
	return value{int16(a.u)}
}

func charOfValue(v value) byte {
	return byte(v.i % 256)
}

func addrOfValue(v value) addr {
	return addr{uint16(v.i)}
}

type slot interface {
	executeSlot(*machine, addr) addr
	toLiteral() value
}

type ret struct {
}

type call struct {
	addr addr
}

type kdxLoop struct {
	key      native
	dispatch native
}

// executeSlot...

func (x kdxLoop) executeSlot(m *machine, a addr) addr {
	x.key(m)
	x.dispatch(m)
	m.rsPush(valueOfAddr(addr{0}))
	return addrOfValue(m.pop())
}

func (native native) executeSlot(m *machine, a addr) addr {
	native(m)
	return a
}

func (call) executeSlot(m *machine, a addr) addr {
	panic("call/execute") // TODO: this will be needed
}

func (ret) executeSlot(m *machine, a addr) addr {
	return addrOfValue(m.rsPop())
}

func (value) executeSlot(*machine, addr) addr {
	panic("value/execute")
}

// toLiteral...

func (kdxLoop) toLiteral() value {
	panic("kdxLoop/toLiteral")
}

func (native) toLiteral() value {
	panic("native/toLiteral")
}

func (call) toLiteral() value {
	panic("call/toLiteral")
}

func (ret) toLiteral() value {
	panic("ret/toLiteral")
}

func (v value) toLiteral() value {
	return v
}
