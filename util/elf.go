//-----------------------------------------------------------------------------
/*

ELF File Utilities

*/
//-----------------------------------------------------------------------------

package util

import (
	"debug/elf"
	"fmt"
)

//-----------------------------------------------------------------------------

// GetELFClass return the elf file class (32 or 64 bit)
func GetELFClass(filename string) (elf.Class, error) {

	f, err := elf.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("%s %s", filename, err)
	}
	defer f.Close()

	if f.Machine != elf.EM_RISCV {
		return 0, fmt.Errorf("%s is not a RISC-V ELF file", filename)
	}

	if f.Type != elf.ET_EXEC {
		return 0, fmt.Errorf("%s is not an executable ELF file", filename)
	}

	return f.Class, nil
}

//-----------------------------------------------------------------------------
