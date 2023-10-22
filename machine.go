package main

import "fmt"

type primitive struct {
	name   string
	action func(*machine)
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
}

func newMachine() *machine {
	dt := make(map[byte]addr)
	mem := make(map[addr]slot)
	here := addr{100}
	psPointer := addr{50000}
	return &machine{0, dt, mem, here, psPointer}
}

func (m *machine) tick() {
	m.steps++
}

func (m *machine) see() {
	fmt.Printf("machine: steps = %v, here = %v, psPointer= %v\n",
		m.steps, m.here, m.psPointer)
}

func (m *machine) installQuarterPrim(c byte, p *primitive) {
	a := m.installPrim(p)
	m.dt[c] = a
}

func (m *machine) installPrim(prim *primitive) addr {
	// TODO: write name & entry to allow dictionary lookup
	// for now we will just write the native-slot code
	a := m.here
	slot := prim
	comma(m, slot)
	return a
}

func comma(m *machine, s slot) {
	a := m.here
	m.mem[a] = s
	m.here = a.next()
}

func (m *machine) lookupDisaptch(c byte) addr {
	addr, ok := m.dt[c]
	if !ok {
		panic(fmt.Sprintf("lookupDisaptch: %c", c))
	}
	return addr
}

func (m *machine) lookupMem(a addr) slot {
	slot, ok := m.mem[a]
	//fmt.Println("lookupMem",a,"->",slot,ok)
	if !ok {
		panic(fmt.Sprintf("lookupMem: %a", a))
	}
	return slot
}

func (m *machine) run() {
	for {
		m.tick()
		//m.see()
		m.executeQ('^')
		if isZero(m.top()) {
			break
		}
		m.executeQ('.')
	}
}

func (m *machine) executeQ(c byte) {
	a := m.lookupDisaptch(c)
	a.execute(m)
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

func (m *machine) top() value {
	a := m.psPointer
	slot := m.lookupMem(a)
	return slot.toLiteral()
}

type addr struct {
	i uint16
}

func (a addr) next() addr {
	return a.offset(1)
}

func (a addr) offset(n int) addr {
	return addr{a.i + uint16(n)}
}

func (a addr) execute(m *machine) {
	s := m.lookupMem(a)
	s.execute(m)
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

func charOfValue(v value) byte {
	return byte(v.i % 256)
}

type slot interface { // primitive or value
	execute(*machine)
	toLiteral() value
}

func (n primitive) execute(m *machine) {
	n.action(m)
}

func (n primitive) toLiteral() value {
	panic("primitive/toLiteral")
}

func (value) execute(*machine) {
	panic("slotLiteral/execute")
}

func (v value) toLiteral() value {
	return v
}
