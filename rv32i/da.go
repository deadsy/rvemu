//-----------------------------------------------------------------------------
/*

RISC-V RV32I Disassembler

*/
//-----------------------------------------------------------------------------

package rv32i

import (
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------

// Disassembly returns the result of the disassembler call.
type Disassembly struct {
	Symbol      string  // symbol for the address (if any)
	Instruction string  // instruction decode
	Comment     string  // useful comment
}

//-----------------------------------------------------------------------------

func Disassemble(ins uint32) *Disassembly {








}

//-----------------------------------------------------------------------------
