//-----------------------------------------------------------------------------
/*

RISC-V Disassembler

*/
//-----------------------------------------------------------------------------

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/deadsy/riscv/mem"
	"github.com/deadsy/riscv/rv"
)

//-----------------------------------------------------------------------------

// loadDump loads an objdump output file to memory.
func loadDump(m *mem.Memory, filename string) error {

	// get the file contents
	x, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	line := strings.Split(string(x), "\n")
	for i := range line {
		field := strings.Fields(line[i])
		n := len(field)

		if n <= 1 {
			continue
		}

		var adr uint64
		if n >= 2 {
			// address
			s := strings.Trim(field[0], ":")
			adr, err = strconv.ParseUint(s, 16, 32)
			if err != nil {
				return fmt.Errorf("error at line %d: %s", i+1, err)
			}
		}

		if n == 2 {
			// address + symbol
			s := strings.Trim(field[1], "<>:")
			m.AddSymbol(uint32(adr), s)
		}

		if n > 2 {
			insLength := len(field[1])
			if insLength == 4 {
				// 16-bit instruction
				ins, err := strconv.ParseUint(field[1], 16, 16)
				if err != nil {
					return fmt.Errorf("error at line %d: %s", i+1, err)
				}
				m.Wr16(uint32(adr), uint16(ins))
			} else if insLength == 8 {
				// 32-bit instruction
				ins, err := strconv.ParseUint(field[1], 16, 32)
				if err != nil {
					return fmt.Errorf("error at line %d: %s", i+1, err)
				}
				m.Wr32(uint32(adr), uint32(ins))
			} else {
				return fmt.Errorf("unrecognised instruction length at line %d", i+1)
			}
		}
	}

	return nil
}

//-----------------------------------------------------------------------------

func main() {

	// command line flags
	fname := flag.String("f", "dump.txt", "dump file to load")
	flag.Parse()

	// create the memory
	m := mem.NewMemory(0, 1<<20, false)
	// load the memory
	err := loadDump(m, *fname)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// create the ISA
	isa := rv.NewISA("rv32g")
	err = isa.Add(rv.ISArv32i, rv.ISArv32m, rv.ISArv32a, rv.ISArv32f, rv.ISArv32d, rv.ISArv32c)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// create the CPU
	cpu := rv.NewRV32(isa, m)
	adr := uint32(0)

	// Disassemble
	for true {
		da := cpu.Disassemble(adr)
		fmt.Printf("%s\n", da.String())
		adr += uint32(da.Length)
		if da.Assembly == "?" {
			break
		}
	}

	os.Exit(0)
}

//-----------------------------------------------------------------------------
