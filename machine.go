package main

import "fmt"

type native func(*machine)

type primitive struct {
	name   string
	action native
}

type value struct {
	i uint16
}

type char byte

type addr struct {
	u uint16
}

type rts struct {
}

type jsr struct {
}

type flags struct {
	hidden    bool
	immediate bool
}

type kdxLoop struct {
	key      native
	dispatch native
}

type machine struct {
	locate          map[string]uint16
	halt            addr
	kdx             addr
	hereP           addr
	echoEnabledP    addr
	latest          addr
	dt              map[char]addr
	mem             map[addr]slot
	psPointer       addr
	rsPointer       addr
	steps           uint
	startupComplete bool
}

const psBase uint16 = 0x90

func newMachine(locate map[string]uint16, key, dispatch native) *machine {

	mem := make(map[addr]slot)
	halt := addr{0}
	kdx := addr{1}
	hereP := addr{2}
	echoEnabledP := addr{4}
	m := machine{
		locate:          locate,
		halt:            halt,
		kdx:             kdx,
		hereP:           hereP,
		echoEnabledP:    echoEnabledP,
		latest:          addr{0},
		dt:              make(map[char]addr),
		mem:             mem,
		psPointer:       addr{psBase},
		rsPointer:       addr{61000},
		steps:           0,
		startupComplete: false,
	}
	m.mem[kdx] = kdxLoop{key, dispatch}
	m.setHere(addr{100}) // for 2x Nop & SetTabEntry
	m.setValue(echoEnabledP, valueOfBool(false))
	return &m
}

func (m *machine) setValue(a addr, v value) {
	m.mem[a] = char(v.i % 256)
	m.mem[a.offset(1)] = char(v.i / 256)
}

func (m *machine) run(here_start addr) { // the inner interpreter (aka trampoline!)
	m.setHere(here_start)
	var a addr = m.kdx
	for {
		if a == m.halt {
			break
		}
		m.tick()
		slot := m.lookupMem(a)
		a = slot.executeSlot(m, a)
	}
}

func (m *machine) tick() {
	m.steps++
}

func (m *machine) installQuarterOnly(c char, name string, action native) {
	xt := m.here()
	prim := &primitive{name, action}
	m.comma(prim)
	m.dt[c] = xt
}

func (m *machine) installQuarterPrim(c char, name string, action native) {
	xt := m.installPrim(name, action)
	m.dt[c] = xt
}

func (m *machine) installPrim(name string, action native) addr {
	at, ok := m.locate[name]
	if !ok {
		panic(fmt.Sprintf("installPrim: %v", name))
	}
	m.setHere(addr{at}.offset(-(len(name) + 6))) // null + 5 for dict entry
	prim := &primitive{name, action}
	nameP := m.here()
	m.commaString(name)
	m.commaValue(valueOfAddr(nameP))
	m.commaValue(valueOfAddr(m.latest))
	m.comma(flags{false, false})
	xt := m.here()
	m.comma(prim)
	m.latest = xt
	return xt
}

func AsFlags(slot slot) flags {
	flags, ok := slot.(flags)
	if !ok {
		panic("AsFlags/non-flags")
	}
	return flags
}

func (m *machine) here() addr {
	return addrOfValue(m.readValue(m.hereP))
}

func (m *machine) setHere(a addr) {
	m.setValue(m.hereP, valueOfAddr(a))
}

func (m *machine) commaString(s string) {
	for i := 0; i < len(s); i++ {
		c := char(s[i])
		m.comma(c)
	}
	m.comma(char(0))
}

func (m *machine) commaValue(v value) {
	m.comma(char(v.i % 256))
	m.comma(char(v.i / 256))
}

func (m *machine) comma(s slot) {
	a := m.here()
	//fmt.Printf("comma: %v = %s\n", a, s.viewSlot())
	m.mem[a] = s
	m.setHere(a.offset(1))
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

func (m *machine) readValue(a addr) value {
	lo := uint16(m.lookupMem(a).toChar())
	hi := uint16(m.lookupMem(a.offset(1)).toChar())
	return value{lo + 256*hi}
}

func (m *machine) push(v value) {
	m.psPointer = m.psPointer.offset(-2)
	m.setValue(m.psPointer, v)
}

func (m *machine) pop() value {
	v := m.readValue(m.psPointer)
	m.psPointer = m.psPointer.offset(2)
	return v
}

func (m *machine) rsPush(v value) {
	m.rsPointer = m.rsPointer.offset(-2)
	m.setValue(m.rsPointer, v)
}

func (m *machine) rsPop() value {
	v := m.readValue(m.rsPointer)
	m.rsPointer = m.rsPointer.offset(2)
	return v
}

func (a addr) offset(n int) addr {
	return addr{a.u + uint16(n)}
}

func isZero(v value) bool {
	return v.i == 0
}

func isTrue(v value) bool {
	return v.i != 0
}

func valueOfBool(b bool) value {
	if b {
		return value{65535}
	} else {
		return value{0}
	}
}

func valueOfChar(c char) value {
	return value{uint16(c)}
}

func valueOfAddr(a addr) value {
	return value{uint16(a.u)}
}

func charOfValue(v value) char {
	return char(v.i % 256)
}

func addrOfValue(v value) addr {
	return addr{uint16(v.i)}
}

type slot interface {
	executeSlot(*machine, addr) addr
	viewSlot() string
	toChar() char
}

// executeSlot...

func (x kdxLoop) executeSlot(m *machine, a addr) addr {
	x.key(m)
	x.dispatch(m)
	m.rsPush(valueOfAddr(m.kdx))
	return addrOfValue(m.pop())
}

func (p primitive) executeSlot(m *machine, a addr) addr {
	p.action(m)
	return addrOfValue(m.rsPop())
}

func (jsr) executeSlot(m *machine, a addr) addr {
	m.rsPush(valueOfAddr(a.offset(3)))
	return addrOfValue(m.readValue(a.offset(1)))
}

func (rts) executeSlot(m *machine, a addr) addr {
	return addrOfValue(m.rsPop())
}

func (c char) executeSlot(m *machine, a addr) addr {
	panic("char/execute")
}

func (flags) executeSlot(*machine, addr) addr {
	panic("flags/execute")
}

// toChar...

const jsrOpcode = 0x20 // 0xe8
const rtsOpcode = 0x60 // 0xc3

func (kdxLoop) toChar() char {
	panic("kdxLoop/toChar")
}

func (primitive) toChar() char {
	panic("primitive/toChar")
}

func (jsr) toChar() char {
	return char(jsrOpcode)
}

func (rts) toChar() char {
	return char(rtsOpcode)
}

func (f flags) toChar() char {
	acc := 0
	if f.immediate {
		acc += 0x40
	}
	if f.hidden {
		acc += 0x80
	}
	return char(acc)
}

func (c char) toChar() char {
	return c
}

// viewSlot...

func (kdxLoop) viewSlot() string {
	return "KDXL"
}

func (p primitive) viewSlot() string {
	return fmt.Sprintf("prim: %s", p.name)
}

func (jsr) viewSlot() string {
	return "jsr"
}

func (rts) viewSlot() string {
	return "rts"
}

func (e flags) viewSlot() string {
	return fmt.Sprintf("flags: immediate=%v, hidden=%v", e.immediate, e.hidden)
}

func (c char) viewSlot() string {
	return fmt.Sprintf("char: %v", c)
}

// String...

func (addr addr) String() string {
	return fmt.Sprintf("[%x]", addr.u)
}

func (c char) String() string {
	if c >= 32 && c <= 126 {
		return fmt.Sprintf("%d='%c'", c, c)
	} else {
		return fmt.Sprintf("%d", c)
	}
}
