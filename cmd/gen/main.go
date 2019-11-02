//-----------------------------------------------------------------------------
/*

RISC-V Disassembler/Emulator Code Generation

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	"github.com/deadsy/riscv/cpu"
)

//-----------------------------------------------------------------------------

func main() {

	// create the ISA
	isa := cpu.NewISA("rv32g")
	err := isa.Add(cpu.ISArv32i, cpu.ISArv32m, cpu.ISArv32a, cpu.ISArv32f, cpu.ISArv32d)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	s := isa.GenLinearDecoder()
	fmt.Printf("%s\n", s)
	os.Exit(0)
}

//-----------------------------------------------------------------------------
