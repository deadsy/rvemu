//-----------------------------------------------------------------------------
/*

RISC-V ISA Decodes

This code dumps the constant bit patterns for the ISA so I can use them
for instruction encoding in the RISC-V debugger.

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	"github.com/deadsy/riscv/rv"
)

//-----------------------------------------------------------------------------

func decode() error {
	// 32-bit ISA
	isa := rv.NewISA(0)
	err := isa.Add(rv.ISArv32g)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", isa.DecodeConstants())
	return nil
}

//-----------------------------------------------------------------------------

func main() {
	err := decode()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

//-----------------------------------------------------------------------------
