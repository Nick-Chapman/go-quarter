package main

import "fmt"

type input interface {
	getChar() byte
}

type primitive struct {
	name string
	action func(*machine)
}

func makePrim(n string, f func(*machine)) *primitive {
	return &primitive{n,f}
}

type machine struct {
	dt map[byte]addr
	mem map[addr]slot
	here addr
	input input
	psPointer addr
}

func newMachine(input input) *machine {
	dt := make(map[byte]addr)
	mem := make(map[addr]slot)
	here := addr{100}
	psPointer := addr{50000}
	return &machine{dt,mem,here,input,psPointer}
}

func (m *machine) installQuarterPrim(c byte, p *primitive) {
	a := m.installPrim(p)
	m.dt[c] = a
}

func (m *machine) installPrim(p *primitive) addr {
	// TODO: write name & entry to allow dictionary lookup
	// for now we will just write the native-slot code
	a := m.here
	slot := native{p.action}
	comma(m,slot)
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
		panic(fmt.Sprintf("lookupDisaptch: %c",c))
	}
	return addr
}

func (m *machine) lookupMem(a addr) slot {
	slot, ok := m.mem[a]
	if !ok {
		panic(fmt.Sprintf("lookupMem: %a",a))
	}
	return slot
}

func (m *machine) run() {
	for {
		m.executeQ('^')
		m.executeQ('.')
	}
}

func (m *machine) executeQ(c byte) {
	a := m.lookupDisaptch(c)
	a.execute(m)
}


func (m *machine) getChar() byte {
	return m.input.getChar()
}

func (m *machine) push(v value) {
	a := m.psPointer
	m.mem[a] = slotLiteral{v}
}

func (m *machine) pop() value {
	a := m.psPointer
	slot := m.lookupMem(a)
	return slot.toLiteral()
}

type value struct {
	i int16
}

func valueOfChar(c byte) value {
	return value{int16(c)}
}

func charOfValue(v value) byte {
	return byte(v.i % 256)
}


type addr struct {
	i uint16
}

func (a addr) next() addr {
	return addr{a.i+1}
}

func (a addr) execute(m *machine) {
	s := m.lookupMem(a)
	s.execute(m)
}

type slot interface {
	execute(*machine)
	toLiteral() value
}

type native struct {
	action func(*machine)
}

func (n native) execute(m *machine) {
	n.action(m)
}

func (n native) toLiteral() value {
	panic("native/toLiteral")
}

type slotLiteral struct {
	literal value
}

func (slotLiteral) execute(*machine) {
	panic("slotLiteral/execute")
}

func (s slotLiteral) toLiteral() value {
	return s.literal
}
