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

func (m *Memory) loadSection(f *elf.File, name string) string {
	s := f.Section(name)
	if s == nil || s.Size == 0 {
		return fmt.Sprintf("%s (0 bytes)", name)
	}
	data, err := s.Data()
	if err != nil {
		return fmt.Sprintf("%s can't read section: %s", name, err)
	}
	for i, v := range data {
		m.Wr8(uint(s.Addr)+uint(i), v)
	}
	end := s.Addr + s.Size - 1
	return fmt.Sprintf("%s %08x-%08x (%d bytes)", name, s.Addr, end, s.Size)
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
	s = append(s, m.loadSymbols(f))
	s = append(s, m.loadSection(f, ".text"))
	s = append(s, m.loadSection(f, ".rodata"))
	s = append(s, m.loadSection(f, ".data"))

	return strings.Join(s, "\n"), nil
}

//-----------------------------------------------------------------------------
