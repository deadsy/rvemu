//-----------------------------------------------------------------------------
/*

Memory Symbols

*/
//-----------------------------------------------------------------------------

package mem

import "fmt"

//-----------------------------------------------------------------------------

// Symbol is a memory symbol.
type Symbol struct {
	Name string
	Addr uint
	Size uint
}

// sort symbols by address
type symbolByAddr []*Symbol

func (a symbolByAddr) Len() int           { return len(a) }
func (a symbolByAddr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a symbolByAddr) Less(i, j int) bool { return a[i].Addr < a[j].Addr }

//-----------------------------------------------------------------------------
// Symbol Functions

// SymbolByAddress returns a symbol for the memory address.
func (m *Memory) SymbolByAddress(adr uint) *Symbol {
	return m.symByAddr[adr]
}

// SymbolByName returns the symbol for a symbol name.
func (m *Memory) SymbolByName(s string) *Symbol {
	return m.symByName[s]
}

// SymbolGetAddress returns the symbol address for a symbol name.
func (m *Memory) SymbolGetAddress(s string) (uint, error) {
	symbol := m.symByName[s]
	if symbol == nil {
		return 0, fmt.Errorf("%s not found", s)
	}
	return symbol.Addr, nil
}

// AddSymbol adds a symbol to the symbol table.
func (m *Memory) AddSymbol(s string, adr, size uint) error {
	if m.findByAddr(adr, size) != nil {
		symbol := Symbol{s, adr, size}
		m.symByAddr[adr] = &symbol
		m.symByName[s] = &symbol
		return nil
	}
	return fmt.Errorf("%s is not in a memory region", s)
}

//-----------------------------------------------------------------------------
