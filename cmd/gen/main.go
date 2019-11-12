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
	isa := rv.NewISA()
	err := isa.Add(rv.ISArv32g)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

//-----------------------------------------------------------------------------
