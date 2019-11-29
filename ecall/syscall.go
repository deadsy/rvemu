//-----------------------------------------------------------------------------
/*

RISC-V Linux System Calls

See:

linux/include/uapi/asm-generic/unistd.h
glibc/sysdeps/unix/sysv/linux/riscv/sysdep.h

glibc passes the syscall number in a7.
Syscall arguments are passed in a0..a6.
The syscall return value is passed in a0.

*/
//-----------------------------------------------------------------------------

package ecall

import "github.com/deadsy/riscv/rv"

//-----------------------------------------------------------------------------
// 32-bit system calls

func sc32_close(m *rv.RV32) {
	m.SetBreak()
}

func sc32_fstat(m *rv.RV32) {
	m.SetBreak()
}

func sc32_exit(m *rv.RV32) {
	m.Exit(0)
}

//-----------------------------------------------------------------------------
// 64-bit system calls

func sc64_close(m *rv.RV64) {
	m.SetBreak()
}

func sc64_fstat(m *rv.RV64) {
	m.SetBreak()
}

func sc64_exit(m *rv.RV64) {
	m.Exit(0)
}

//-----------------------------------------------------------------------------

type scFunc32 func(m *rv.RV32)
type scFunc64 func(m *rv.RV64)

type scEntry struct {
	name string
	sc32 scFunc32
	sc64 scFunc64
}

var scTable = map[uint]scEntry{
	57: {"close", sc32_close, sc64_close},
	80: {"fstat", sc32_fstat, sc64_fstat},
	93: {"exit", sc32_exit, sc64_exit},
}

func scLookup(n uint) *scEntry {
	if e, ok := scTable[n]; ok {
		return &e
	}
	return nil
}

//-----------------------------------------------------------------------------

// Syscall is a syscall ecall object.
type Syscall struct {
}

// NewSyscall returns a syscall ecall object.
func NewSyscall() *Syscall {
	return &Syscall{}
}

// Call32 is a 32-bit ecall.
func (sc *Syscall) Call32(m *rv.RV32) {
	n := uint(m.X[rv.RegA7])
	e := scLookup(n)
	if e != nil {
		e.sc32(m)
	}
}

// Call64 is a 64-bit ecall.
func (sc *Syscall) Call64(m *rv.RV64) {
	n := uint(m.X[rv.RegA7])
	e := scLookup(n)
	if e != nil {
		e.sc64(m)
	}
}

//-----------------------------------------------------------------------------
