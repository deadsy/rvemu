//-----------------------------------------------------------------------------
/*

RISC-V Compliance Testing

https://github.com/riscv/riscv-compliance

The complicance tests exit the test by setting status in the gp register
and using the ecall instruction. At that point the non-emulation code can
read the output from memory and compare it with the known good signature
value for the test.

*/
//-----------------------------------------------------------------------------

package ecall

import "github.com/deadsy/riscv/rv"

//-----------------------------------------------------------------------------

// Compliance is a compliance ecall object.
type Compliance struct {
}

// NewCompliance returns a compliance ecall object.
func NewCompliance() *Compliance {
	return &Compliance{}
}

// Call32 is a 32-bit ecall.
func (c *Compliance) Call32(m *rv.RV32) {
	m.Exit(m.X[rv.RegGp])
}

// Call64 is a 64-bit ecall.
func (c *Compliance) Call64(m *rv.RV64) {
	m.Exit(m.X[rv.RegGp])
}

//-----------------------------------------------------------------------------
