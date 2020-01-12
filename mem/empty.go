//-----------------------------------------------------------------------------
/*

Empty Memory

This memory region is a backstop to empty areas in the memory map.
When accessed it returns default values and a memory error.

*/
//-----------------------------------------------------------------------------

package mem

//-----------------------------------------------------------------------------

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
	err := rdInsError(adr, m.attr, m.name)
	err.(*Error).Type |= ErrEmpty
	return 0, err
}

// Rd64 reads a 64-bit data value from memory.
func (m *empty) Rd64(adr uint) (uint64, error) {
	err := rdError(adr, m.attr, m.name, 8)
	err.(*Error).Type |= ErrEmpty
	return 0, err
}

// Rd32 reads a 32-bit data value from memory.
func (m *empty) Rd32(adr uint) (uint32, error) {
	err := rdError(adr, m.attr, m.name, 4)
	err.(*Error).Type |= ErrEmpty
	return 0, err
}

// Rd16 reads a 16-bit data value from memory.
func (m *empty) Rd16(adr uint) (uint16, error) {
	err := rdError(adr, m.attr, m.name, 2)
	err.(*Error).Type |= ErrEmpty
	return 0, err
}

// Rd8 reads an 8-bit data value from memory.
func (m *empty) Rd8(adr uint) (uint8, error) {
	err := rdError(adr, m.attr, m.name, 1)
	err.(*Error).Type |= ErrEmpty
	return 0, err
}

// Wr64 writes a 64-bit data value to memory.
func (m *empty) Wr64(adr uint, val uint64) error {
	err := wrError(adr, m.attr, m.name, 8)
	err.(*Error).Type |= ErrEmpty
	return err
}

// Wr32 writes a 32-bit data value to memory.
func (m *empty) Wr32(adr uint, val uint32) error {
	err := wrError(adr, m.attr, m.name, 4)
	err.(*Error).Type |= ErrEmpty
	return err
}

// Wr16 writes a 16-bit data value to memory.
func (m *empty) Wr16(adr uint, val uint16) error {
	err := wrError(adr, m.attr, m.name, 2)
	err.(*Error).Type |= ErrEmpty
	return err
}

// Wr8 writes an 8-bit data value to memory.
func (m *empty) Wr8(adr uint, val uint8) error {
	err := wrError(adr, m.attr, m.name, 1)
	err.(*Error).Type |= ErrEmpty
	return err
}

//-----------------------------------------------------------------------------
