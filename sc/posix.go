//-----------------------------------------------------------------------------
/*

RISC-V POSIX System Calls

See:

linux/include/uapi/asm-generic/unistd.h
glibc/sysdeps/unix/sysv/linux/riscv/sysdep.h

glibc passes the syscall number in a7.
Syscall arguments are passed in a0..a6.
The syscall return value is passed in a0.

*/
//-----------------------------------------------------------------------------

package sc

import "github.com/deadsy/riscv/rv"

//-----------------------------------------------------------------------------

type scEntry struct {
	name string
	sc32 rv.ScFunc32
	sc64 rv.ScFunc64
}

var scTable = map[uint]scEntry{
	57: {"close", sc32_close, sc64_close},
	80: {"fstat", sc32_fstat, sc64_fstat},
	93: {"exit", sc32_exit, sc64_exit},
}

func scLookup(n uint) *scEntry {
	if s, ok := scTable[n]; ok {
		return &s
	}
	return nil
}

//-----------------------------------------------------------------------------
// 32-bit system calls

func sc32_close(m *rv.RV32) {
	m.SetBreak()
}

func sc32_fstat(m *rv.RV32) {
	m.SetBreak()
}

func sc32_exit(m *rv.RV32) {
	m.SetExit(0)
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
	m.SetExit(0)
}

//-----------------------------------------------------------------------------

type Posix struct {
}

func NewPosix() *Posix {
	return &Posix{}
}

func (sc *Posix) Lookup32(m *rv.RV32) rv.ScFunc32 {
	n := uint(m.X[rv.RegA7])
	s := scLookup(n)
	if s != nil {
		return s.sc32
	}
	return nil
}

func (sc *Posix) Lookup64(m *rv.RV64) rv.ScFunc64 {
	n := uint(m.X[rv.RegA7])
	s := scLookup(n)
	if s != nil {
		return s.sc64
	}
	return nil
}

//-----------------------------------------------------------------------------
