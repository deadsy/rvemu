//-----------------------------------------------------------------------------
/*

ELF File Handling

*/
//-----------------------------------------------------------------------------

package mem

import (
	"debug/elf"
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------

func (m *Memory) loadSymbols(f *elf.File) string {
	st, err := f.Symbols()
	if err != nil {
		return fmt.Sprintf("can't load symbols")
	}
	n := 0
	for i := range st {
		var err error
		switch elf.ST_TYPE(st[i].Info) {
		case elf.STT_FUNC:
			err = m.AddSymbol(st[i].Name, uint(st[i].Value), uint(st[i].Size))
			if err == nil {
				n++
			}
		}
	}
	return fmt.Sprintf("loaded %d symbols", n)
}

//-----------------------------------------------------------------------------

func makeSection(f *elf.File, s *elf.Section) (*Section, string) {

	if s.Size == 0 {
		return nil, fmt.Sprintf("%s (0 bytes)", s.Name)
	}

	// work out the memory attribute
	attr := AttrR
	if s.Flags&elf.SHF_WRITE != 0 {
		attr |= AttrW
	}
	if s.Flags&elf.SHF_EXECINSTR != 0 {
		attr |= AttrX
	}

	// create the memory section
	ms := NewSection(s.Name, uint(s.Addr), uint(s.Size), attr)

	// read the section data from the ELF file
	data, err := s.Data()
	if err != nil {
		return nil, fmt.Sprintf("can't read section %s (%s)", s.Name, err)
	}

	// write the data to the memory section
	for i, v := range data {
		ms.Wr8(uint(s.Addr)+uint(i), v)
	}

	end := s.Addr + s.Size - 1
	return ms, fmt.Sprintf("%-16s %08x-%08x (%d bytes)", s.Name, s.Addr, end, s.Size)
}

//-----------------------------------------------------------------------------

// LoadELF loads an ELF file to memory.
func (m *Memory) LoadELF(filename string, class elf.Class) (string, error) {

	f, err := elf.Open(filename)
	if err != nil {
		return "", fmt.Errorf("%s %s", filename, err)
	}

	defer f.Close()

	if f.Machine != elf.EM_RISCV {
		return "", fmt.Errorf("%s is not a RISC-V ELF file", filename)
	}

	if f.Class != class {
		return "", fmt.Errorf("%s is not an %s file", filename, class)
	}

	if f.Type != elf.ET_EXEC {
		return "", fmt.Errorf("%s is not an executable ELF file", filename)
	}

	s := make([]string, 0)

	// load the sections
	for _, fs := range f.Sections {
		if fs.Flags&elf.SHF_ALLOC != 0 {
			ms, status := makeSection(f, fs)
			if ms != nil {
				m.Add(ms)
			}
			s = append(s, status)
		}
	}

	// set the program entry point
	m.Entry = f.Entry
	s = append(s, fmt.Sprintf("%-16s %08x", "entry point", m.Entry))

	// load the symbols
	s = append(s, m.loadSymbols(f))

	return strings.Join(s, "\n"), nil
}

//-----------------------------------------------------------------------------
