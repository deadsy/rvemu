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
			m.AddSymbol(s, uint(adr), 0)
		}

		if n > 2 {
			insLength := len(field[1])
			if insLength == 4 {
				// 16-bit instruction
				ins, err := strconv.ParseUint(field[1], 16, 16)
				if err != nil {
					return fmt.Errorf("error at line %d: %s", i+1, err)
				}
				m.Wr16(uint(adr), uint16(ins))
			} else if insLength == 8 {
				// 32-bit instruction
				ins, err := strconv.ParseUint(field[1], 16, 32)
				if err != nil {
					return fmt.Errorf("error at line %d: %s", i+1, err)
				}
				m.Wr32(uint(adr), uint32(ins))
			} else {
				return fmt.Errorf("unrecognised instruction length at line %d", i+1)
			}
		}

		// get the reference disassembly
		if n > 2 {
			s := make([]string, 0)
			for j := range field[2:] {
				x := field[2+j]
				if x[0] == '#' || x[0] == '<' {
					break
				}
				s = append(s, x)
			}
			m.AddDisassembly(strings.Join(s, " "), uint(adr))
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
	m := mem.NewMem32()
	m.Add(mem.NewChunk(0, 64<<10, mem.AttrRW))

	// load the memory
	err := loadDump(m, *fname)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// create the ISA
	isa := rv.NewISA()
	err = isa.Add(rv.ISArv32gc)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// Disassemble
	adr := uint(0)
	for true {
		da := isa.Disassemble(m, adr)
		if da.Assembly == "?" {
			break
		}
		if da.Assembly == m.Disassembly(adr) {
			fmt.Printf("%s\n", da.String())
		} else {
			fmt.Printf("%s should be: \"%s\"\n", da.String(), m.Disassembly(adr))
		}
		adr += uint(da.Length)
	}

	os.Exit(0)
}

//-----------------------------------------------------------------------------
