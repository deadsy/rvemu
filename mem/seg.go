//-----------------------------------------------------------------------------
/*

Memory Segments

*/
//-----------------------------------------------------------------------------

package mem

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

//-----------------------------------------------------------------------------

// Error is a memory acccess error.
type Error struct {
	n    uint   // bitmap of memory errors
	addr uint   // memory address causing the error
	name string // section name for the address
}

// Memory error bits.
const (
	ErrAlign = 1 << iota // misaligned read/write
	ErrRead              // can't read this memory
	ErrWrite             // can't write this memory
	ErrExec              // can't read instructions from this memory
	ErrBreak             // break on memory access
)

func (e *Error) Error() string {
	s := make([]string, 0)
	if e.n&ErrAlign != 0 {
		s = append(s, "align")
	}
	if e.n&ErrRead != 0 {
		s = append(s, "read")
	}
	if e.n&ErrWrite != 0 {
		s = append(s, "write")
	}
	if e.n&ErrExec != 0 {
		s = append(s, "exec")
	}
	if e.n&ErrBreak != 0 {
		s = append(s, "break")
	}
	errStr := strings.Join(s, ",")
	return fmt.Sprintf("%s @ %08x (%s)", errStr, e.addr, e.name)
}

//-----------------------------------------------------------------------------

// Attribute is a bit mask of memory access attributes.
type Attribute uint

// Attribute values.
const (
	AttrR Attribute = 1 << iota // read
	AttrW                       // write
	AttrX                       // execute
	AttrM                       // misaligned access
)

// AttrRW = read/write
const AttrRW = AttrR | AttrW

// AttrRW = read/write/misaligned
const AttrRWM = AttrR | AttrW | AttrM

// AttrRX = read/execute
const AttrRX = AttrR | AttrX

// AttrRWX = read/write/execute
const AttrRWX = AttrR | AttrW | AttrX

func (a Attribute) String() string {
	s := make([]string, 4)
	for i := range s {
		s[i] = "-"
	}
	if a&AttrR != 0 {
		s[0] = "r"
	}
	if a&AttrW != 0 {
		s[1] = "w"
	}
	if a&AttrX != 0 {
		s[2] = "x"
	}
	if a&AttrM != 0 {
		s[3] = "m"
	}
	return strings.Join(s, "")
}

//-----------------------------------------------------------------------------
// memory access errors

func wrError(addr uint, attr Attribute, name string, align uint) error {
	var n uint
	if attr&AttrW == 0 {
		n |= ErrWrite
	}
	if (attr&AttrM == 0) && (addr&(align-1) != 0) {
		n |= ErrAlign
	}
	if n != 0 {
		return &Error{n, addr, name}
	}
	return nil
}

func rdError(addr uint, attr Attribute, name string, align uint) error {
	var n uint
	if attr&AttrR == 0 {
		n |= ErrRead
	}
	if (attr&AttrM == 0) && (addr&(align-1) != 0) {
		n |= ErrAlign
	}
	if n != 0 {
		return &Error{n, addr, name}
	}
	return nil
}

func rdInsError(addr uint, attr Attribute, name string) error {
	// rv32c has mixed 32/16 bit instruction streams so
	// we allow 32-bit reads on 2 byte address boundaries.
	err := rdError(addr, attr, name, 2)
	if err != nil && attr&AttrX == 0 {
		err.(*Error).n |= ErrExec
	}
	return err
}

//-----------------------------------------------------------------------------

// RegionInfo contains information for the memory region.
type RegionInfo struct {
	name       string
	start, end uint
	attr       Attribute
}

// sort regions by start address
type regionByStart []*RegionInfo

func (a regionByStart) Len() int           { return len(a) }
func (a regionByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a regionByStart) Less(i, j int) bool { return a[i].start < a[j].start }

// Region is an interface to a contiguous region of memory.
type Region interface {
	Info() *RegionInfo
	SetAttr(attr Attribute)
	RdIns(adr uint) (uint, error)
	Rd64(adr uint) (uint64, error)
	Rd32(adr uint) (uint32, error)
	Rd16(adr uint) (uint16, error)
	Rd8(adr uint) (uint8, error)
	Wr64(adr uint, val uint64) error
	Wr32(adr uint, val uint32) error
	Wr16(adr uint, val uint16) error
	Wr8(adr uint, val uint8) error
	In(adr, size uint) bool
}

//-----------------------------------------------------------------------------

// Section is a contiguous region of real memory.
type Section struct {
	name       string    // section name
	attr       Attribute // bitmask of attributes
	start, end uint      // address range
	mem        []uint8   // memory array
}

// NewSection allocates and returns a memory chunk.
func NewSection(name string, start, size uint, attr Attribute) *Section {
	// allocate the memory and set it to all ones
	mem := make([]uint8, size)
	for i := range mem {
		mem[i] = 0xff
	}
	return &Section{
		name:  name,
		attr:  attr,
		start: start,
		end:   start + size - 1,
		mem:   mem,
	}
}

// SetAttr sets the attributes for this memory section.
func (m *Section) SetAttr(attr Attribute) {
	m.attr = attr
}

// Info returns the information for this region.
func (m *Section) Info() *RegionInfo {
	return &RegionInfo{
		name:  m.name,
		start: m.start,
		end:   m.end,
		attr:  m.attr,
	}
}

// In returns true if the adr, size is entirely within the memory chunk.
func (m *Section) In(adr, size uint) bool {
	end := adr + size - 1
	return (adr >= m.start) && (end <= m.end)
}

// RdIns reads a 32-bit instruction from memory.
func (m *Section) RdIns(adr uint) (uint, error) {
	return uint(binary.LittleEndian.Uint32(m.mem[adr-m.start:])), rdInsError(adr, m.attr, m.name)
}

// Rd64 reads a 64-bit data value from memory.
func (m *Section) Rd64(adr uint) (uint64, error) {
	return binary.LittleEndian.Uint64(m.mem[adr-m.start:]), rdError(adr, m.attr, m.name, 8)
}

// Rd32 reads a 32-bit data value from memory.
func (m *Section) Rd32(adr uint) (uint32, error) {
	return binary.LittleEndian.Uint32(m.mem[adr-m.start:]), rdError(adr, m.attr, m.name, 4)
}

// Rd16 reads a 16-bit data value from memory.
func (m *Section) Rd16(adr uint) (uint16, error) {
	return binary.LittleEndian.Uint16(m.mem[adr-m.start:]), rdError(adr, m.attr, m.name, 2)
}

// Rd8 reads an 8-bit data value from memory.
func (m *Section) Rd8(adr uint) (uint8, error) {
	return m.mem[adr-m.start], rdError(adr, m.attr, m.name, 1)
}

// Wr64 writes a 64-bit data value to memory.
func (m *Section) Wr64(adr uint, val uint64) error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint64(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, m.name, 8)
}

// Wr32 writes a 32-bit data value to memory.
func (m *Section) Wr32(adr uint, val uint32) error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint32(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, m.name, 4)
}

// Wr16 writes a 16-bit data value to memory.
func (m *Section) Wr16(adr uint, val uint16) error {
	if m.attr&AttrW != 0 {
		binary.LittleEndian.PutUint16(m.mem[adr-m.start:], val)
	}
	return wrError(adr, m.attr, m.name, 2)
}

// Wr8 writes an 8-bit data value to memory.
func (m *Section) Wr8(adr uint, val uint8) error {
	if m.attr&AttrW != 0 {
		m.mem[adr-m.start] = val
	}
	return wrError(adr, m.attr, m.name, 1)
}

//-----------------------------------------------------------------------------
// If a memory access does not correspond to a defined memory Region the
// empty memory region will be used.

// empty memory region.
type empty struct {
	attr Attribute // bitmask of attributes
	name string
}

// newEmpty allocates and returns the empty memory region.
func newEmpty(attr Attribute) *empty {
	return &empty{
		attr: attr,
		name: "empty",
	}
}

// SetAttr sets the attributes for the empty region.
func (m *empty) SetAttr(attr Attribute) {
	m.attr = attr
}

// Info returns the information for the empty region.
func (m *empty) Info() *RegionInfo {
	return &RegionInfo{
		name: m.name,
		attr: m.attr,
	}
}

// In returns true if the adr, size is entirely within the empty region.
func (m *empty) In(adr, size uint) bool {
	return true
}

// RdIns reads a 32-bit instruction from memory.
func (m *empty) RdIns(adr uint) (uint, error) {
	return math.MaxUint32, rdInsError(adr, m.attr, m.name)
}

// Rd64 reads a 64-bit data value from memory.
func (m *empty) Rd64(adr uint) (uint64, error) {
	return math.MaxUint64, rdError(adr, m.attr, m.name, 8)
}

// Rd32 reads a 32-bit data value from memory.
func (m *empty) Rd32(adr uint) (uint32, error) {
	return math.MaxUint32, rdError(adr, m.attr, m.name, 4)
}

// Rd16 reads a 16-bit data value from memory.
func (m *empty) Rd16(adr uint) (uint16, error) {
	return math.MaxUint16, rdError(adr, m.attr, m.name, 2)
}

// Rd8 reads an 8-bit data value from memory.
func (m *empty) Rd8(adr uint) (uint8, error) {
	return math.MaxUint8, rdError(adr, m.attr, m.name, 1)
}

// Wr64 writes a 64-bit data value to memory.
func (m *empty) Wr64(adr uint, val uint64) error {
	return wrError(adr, m.attr, m.name, 8)
}

// Wr32 writes a 32-bit data value to memory.
func (m *empty) Wr32(adr uint, val uint32) error {
	return wrError(adr, m.attr, m.name, 4)
}

// Wr16 writes a 16-bit data value to memory.
func (m *empty) Wr16(adr uint, val uint16) error {
	return wrError(adr, m.attr, m.name, 2)
}

// Wr8 writes an 8-bit data value to memory.
func (m *empty) Wr8(adr uint, val uint8) error {
	return wrError(adr, m.attr, m.name, 1)
}

//-----------------------------------------------------------------------------
