//-----------------------------------------------------------------------------
/*

RISC-V Floating Point Routines

*/
//-----------------------------------------------------------------------------

package rv

import (
	"errors"

	"github.com/deadsy/riscv/csr"
)

//-----------------------------------------------------------------------------

/*
#cgo linux LDFLAGS: -L../softfp -lsoftfp
#cgo linux CFLAGS: -I../softfp
#include "softfp.h"
*/
import "C"

//-----------------------------------------------------------------------------

const mask30to0 = (1 << 31) - 1
const mask31 = (1 << 31)

//-----------------------------------------------------------------------------

// Rounding modes. These are the RISC-V values and match the softfp library.
const (
	frmRNE = 0 // Round to Nearest, ties to Even
	frmRTZ = 1 // Round towards Zero
	frmRDN = 2 // Round Down (towards -inf)
	frmRUP = 3 // Round Up (towards +inf)
	frmRRM = 4 // Round to Nearest, ties to Max Magnitude
	frmDYN = 7 // Use the value in the FRM csr
)

// Rounding mode names.
var rmName = [8]string{
	"rne", "rtz", "rdn", "rup", "rrm", "rm5", "rm6", "dyn",
}

// FCSR fflags bits. These are the RISC-V values and match the softfp library.
const (
	fflagsNX = (1 << iota) // Inexact
	fflagsUF               // Underflow
	fflagsOF               // Overflow
	fflagsDZ               // Divide by Zero
	fflagsNV               // Invalid Operation
)

const fflagsALL = fflagsNX | fflagsUF | fflagsOF | fflagsDZ | fflagsNV

//-----------------------------------------------------------------------------

func getRoundingMode(rm uint, s *csr.State) (uint, error) {
	// with dynamic rounding rm = FRM
	if rm == frmDYN {
		rm, _ = s.Rd(csr.FRM)
	}
	if rm > frmRRM {
		return 0, errors.New("illegal")
	}
	return rm, nil
}

//-----------------------------------------------------------------------------

// feq_s (NV)
func feq_s(a, b uint32, s *csr.State) uint {
	var flags C.uint
	x := uint(C.eq_quiet_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

// flt_s (NV)
func flt_s(a, b uint32, s *csr.State) uint {
	var flags C.uint
	x := uint(C.lt_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

// fle_s (NV)
func fle_s(a, b uint32, s *csr.State) uint {
	var flags C.uint
	x := uint(C.le_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

//-----------------------------------------------------------------------------

// fcvt_s_w (NX)
func fcvt_s_w(a int32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint
	x := uint32(C.cvt_i32_sf32(C.int32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_s_wu (NX)
func fcvt_s_wu(a uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint
	x := uint32(C.cvt_u32_sf32(C.uint32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

//-----------------------------------------------------------------------------

// fcvt_w_s converts a float32 to an int32 (NV, NX)
func fcvt_w_s(a uint32, rm uint, s *csr.State) (int32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint
	x := int32(C.cvt_sf32_i32(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_wu_s converts a float32 to an uint32 (NV, NX)
func fcvt_wu_s(a uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint
	x := uint32(C.cvt_sf32_u32(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

//-----------------------------------------------------------------------------

// fadd_s adds two 32-bit floats (NV, OF, UF, NX)
func fadd_s(a, b uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint
	x := uint32(C.add_sf32(C.sfloat32(a), C.sfloat32(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fsub_s subtracts two 32-bit floats (NV, OF, UF, NX)
func fsub_s(a, b uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint
	x := uint32(C.sub_sf32(C.sfloat32(a), C.sfloat32(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

//-----------------------------------------------------------------------------

// fmul_s multiplies two 32-bit floats (NV, OF, UF, NX)
func fmul_s(a, b uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint
	x := uint32(C.mul_sf32(C.sfloat32(a), C.sfloat32(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

//-----------------------------------------------------------------------------
