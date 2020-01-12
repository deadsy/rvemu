//-----------------------------------------------------------------------------
/*

RISC-V Emulation Errors

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"strings"

	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/mem"
)

//-----------------------------------------------------------------------------

// Emulation error types.
const (
	ErrIllegal = (1 << iota) // illegal instruction
	ErrMemory                // memory exception
	ErrEcall                 // ecall exception
	ErrEbreak                // ebreak exception
	ErrCSR                   // CSR exception
	ErrTodo                  // unimplemented instruction
	ErrStuck                 // stuck program counter
	//ErrExit                  // exit from emulation
)

// Error is a general emulation error.
type Error struct {
	Type uint   // error type
	alen uint   // address length
	ins  uint   // illegal instruction value
	pc   uint64 // program counter at which error occurrred
	err  error  // sub error
}

func (e *Error) Error() string {
	pcStr := ""
	if e.alen == 32 {
		pcStr = fmt.Sprintf("%08x", e.pc)
	} else {
		pcStr = fmt.Sprintf("%016x", e.pc)
	}
	switch e.Type {
	case ErrIllegal:
		return "illegal instruction at PC " + pcStr
	case ErrMemory:
		return fmt.Sprintf("memory exception at PC %s, %s", pcStr, e.err)
	case ErrEcall:
		return "ecall exception at PC " + pcStr
	case ErrEbreak:
		return "ebreak exception at PC " + pcStr
	case ErrCSR:
		return fmt.Sprintf("csr exception at PC %s, %s", pcStr, e.err)
	//case ErrExit:
	//	return "exit at PC " + pcStr
	case ErrTodo:
		return "unimplemented instruction at PC " + pcStr
	case ErrStuck:
		return "stuck at PC " + pcStr
	}
	return "unknown exception at PC " + pcStr
}

// GetMemError returns a memory error from the general CPU error.
func (e *Error) GetMemError() *mem.Error {
	if e.Type != ErrMemory {
		return nil
	}
	return e.err.(*mem.Error)
}

// GetCSRError returns a CSR error from the general CPU error.
func (e *Error) GetCSRError() *csr.Error {
	if e.Type != ErrCSR {
		return nil
	}
	return e.err.(*csr.Error)
}

// errIllegal returns the error for an illegal instruction exception.
func (m *RV) errIllegal(ins uint) error {
	return &Error{
		Type: ErrIllegal,
		ins:  ins,
		alen: m.xlen,
		pc:   m.PC,
	}
}

// errEcall returns the error for an environment call exception.
func (m *RV) errEcall() error {
	return &Error{
		Type: ErrEcall,
		alen: m.xlen,
		pc:   m.PC,
	}
}

// errEbreak returns the error for an environment break exception.
func (m *RV) errEbreak() error {
	return &Error{
		Type: ErrEbreak,
		alen: m.xlen,
		pc:   m.PC,
	}
}

// errMemory returns the error for a memory exception.
func (m *RV) errMemory(err error) error {
	return &Error{
		Type: ErrMemory,
		alen: m.xlen,
		pc:   m.PC,
		err:  err,
	}
}

// errCSR returns the error for CSR access exception.
func (m *RV) errCSR(err error, ins uint) error {
	return &Error{
		Type: ErrCSR,
		ins:  ins,
		alen: m.xlen,
		pc:   m.PC,
		err:  err,
	}
}

func (m *RV) errStuckPC() error {
	return &Error{
		Type: ErrStuck,
		alen: m.xlen,
		pc:   m.PC,
	}
}

func (m *RV) errTodo() error {
	return &Error{
		Type: ErrTodo,
		alen: m.xlen,
		pc:   m.PC,
	}
}

//-----------------------------------------------------------------------------
// error buffer - record emulation errors as they occur.
// Some of these will be handled as exceptions, but it's useful to look at
// them so you can get an idea of what is happening.

type errBuffer struct {
	buf    []*Error
	rd, wr uint
	size   uint
}

func newErrBuffer(size uint) *errBuffer {
	return &errBuffer{
		buf:  make([]*Error, size),
		size: size,
	}
}

func (eb *errBuffer) write(e *Error) {
	eb.buf[eb.wr] = e
	eb.wr = (eb.wr + 1) % eb.size
	if eb.wr == eb.rd {
		// chase the read index forward
		eb.rd = (eb.rd + 1) % eb.size
	}
}

func (eb *errBuffer) reset() {
	eb.rd = 0
	eb.wr = 0
}

// DisplayErrorBuffer returns a string of the errors in the error buffer.
func (m *RV) DisplayErrorBuffer() string {
	eb := m.err
	s := []string{}
	rd := eb.rd
	for rd != eb.wr {
		s = append(s, eb.buf[eb.rd].Error())
		rd = (rd + 1) % eb.size
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
