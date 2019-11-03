//-----------------------------------------------------------------------------
/*

RISC-V Disassembler/Emulator Code Generation

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	"github.com/deadsy/riscv/rv"
)

//-----------------------------------------------------------------------------

func main() {

	// create the ISA
	isa := rv.NewISA("rv32g")
	err := isa.Add(rv.ISArv32i, rv.ISArv32m, rv.ISArv32a, rv.ISArv32f, rv.ISArv32d)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

//-----------------------------------------------------------------------------
