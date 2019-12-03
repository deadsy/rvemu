//-----------------------------------------------------------------------------
/*

RISC-V Floating Point Routines

*/
//-----------------------------------------------------------------------------

package rv

import (
	"errors"
	"math"

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

// rounding modes
const (
	frmRNE = 0 // Round to Nearest, ties to Even
	frmRTZ = 1 // Round towards Zero
	frmRDN = 2 // Round Down (towards -inf)
	frmRUP = 3 // Round Up (towards +inf)
	frmRRM = 4 // Round to Nearest, ties to Max Magnitude
	frmDYN = 7 // Use the value in the FRM csr
)

// rounding mode names
var rmName = [8]string{
	"rne", "rtz", "rdn", "rup", "rrm", "rm5", "rm6", "dyn",
}

// FCSR fflags bits.
const (
	fflagsNX = (1 << iota) // Inexact
	fflagsUF               // Underflow
	fflagsOF               // Overflow
	fflagsDZ               // Divide by Zero
	fflagsNV               // Invalid Operation
)

const fflagsALL = fflagsNX | fflagsUF | fflagsOF | fflagsDZ | fflagsNV

//-----------------------------------------------------------------------------

// feq_s (NV)
func feq_s(f1, f2 float32, s *csr.State) uint {
	if f1 == f2 {
		return 1
	}
	return 0
}

// flt_s (NV)
func flt_s(f1, f2 float32, s *csr.State) uint {
	if f1 < f2 {
		return 1
	}
	return 0
}

// fle_s (NV)
func fle_s(f1, f2 float32, s *csr.State) uint {
	if f1 <= f2 {
		return 1
	}
	return 0
}

//-----------------------------------------------------------------------------

// fcvt_s_w (NX)
func fcvt_s_w(x int32, s *csr.State) float32 {
	return float32(x)
}

// fcvt_s_wu (NX)
func fcvt_s_wu(x uint32, s *csr.State) float32 {
	return float32(x)
}

//-----------------------------------------------------------------------------

// fcvt_w_s converts a float32 to an int32 (NV, NX)
func fcvt_w_s(f float32, rm uint, s *csr.State) (int32, error) {

	// with dynamic rounding rm = FRM
	if rm == frmDYN {
		rm, _ = s.Rd(csr.FRM)
	}

	// Clear the FCSR flags
	s.Clr(csr.FFLAGS, fflagsALL)

	if math.IsNaN(float64(f)) || f > float32(math.MaxInt32) {
		s.Set(csr.FFLAGS, fflagsNV)
		return math.MaxInt32, nil
	}

	if f < float32(math.MinInt32) {
		s.Set(csr.FFLAGS, fflagsNV)
		return math.MinInt32, nil
	}

	var x int32
	switch rm {
	case frmRNE:
		x = int32(math.RoundToEven(float64(f)))
	case frmRTZ:
		x = int32(math.Trunc(float64(f)))
	case frmRDN:
	case frmRUP:
	case frmRRM:
	default:
		return 0, errors.New("illegal")
	}

	if f != float32(x) {
		s.Set(csr.FFLAGS, fflagsNX)
	}

	return x, nil
}

// fcvt_wu_s converts a float32 to an uint32 (NV, NX)
func fcvt_wu_s(f float32, rm uint, s *csr.State) (uint32, error) {

	// with dynamic rounding rm = FRM
	if rm == frmDYN {
		rm, _ = s.Rd(csr.FRM)
	}

	// Clear the FCSR flags
	s.Clr(csr.FFLAGS, fflagsALL)

	if math.IsNaN(float64(f)) || f > float32(math.MaxUint32) {
		s.Set(csr.FFLAGS, fflagsNV)
		return math.MaxUint32, nil
	}

	if f <= -1 {
		s.Set(csr.FFLAGS, fflagsNV)
		return 0, nil
	}

	var x uint32
	switch rm {
	case frmRNE:
		x = uint32(math.RoundToEven(float64(f)))
	case frmRTZ:
		x = uint32(math.Trunc(float64(f)))
	case frmRDN:
	case frmRUP:
	case frmRRM:
	default:
		return 0, errors.New("illegal")
	}

	if f != float32(x) {
		s.Set(csr.FFLAGS, fflagsNX)
	}

	return x, nil
}

//-----------------------------------------------------------------------------

// fadd_s adds two 32-bit floats (NV, OF, UF, NX)
func fadd_s(f1, f2 float32, rm uint, s *csr.State) (float32, error) {

	// with dynamic rounding rm = FRM
	if rm == frmDYN {
		rm, _ = s.Rd(csr.FRM)
	}

	// Clear the FCSR flags
	s.Clr(csr.FFLAGS, fflagsALL)

	var x float32
	switch rm {
	case frmRNE:
		x = f1 + f2
	case frmRTZ:
		x = f1 + f2
	case frmRDN:
	case frmRUP:
	case frmRRM:
	default:
		return 0, errors.New("illegal")
	}

	if ((x - f2) != f1) || ((x - f1) != f2) {
		s.Set(csr.FFLAGS, fflagsNX)
	}

	return x, nil
}

// fsub_s subtracts two 32-bit floats (NV, OF, UF, NX)
func fsub_s(f1, f2 float32, rm uint, s *csr.State) (float32, error) {

	// with dynamic rounding rm = FRM
	if rm == frmDYN {
		rm, _ = s.Rd(csr.FRM)
	}

	// Clear the FCSR flags
	s.Clr(csr.FFLAGS, fflagsALL)

	var x float32
	switch rm {
	case frmRNE:
		x = f1 - f2
	case frmRTZ:
		x = f1 - f2
	case frmRDN:
	case frmRUP:
	case frmRRM:
	default:
		return 0, errors.New("illegal")
	}

	if ((x + f2) != f1) || ((f1 - x) != f2) {
		s.Set(csr.FFLAGS, fflagsNX)
	}

	return x, nil
}

//-----------------------------------------------------------------------------

// fmul_s multiplies two 32-bit floats (NV, OF, UF, NX)
func fmul_s(f1, f2 float32, rm uint, s *csr.State) (float32, error) {

	// with dynamic rounding rm = FRM
	if rm == frmDYN {
		rm, _ = s.Rd(csr.FRM)
	}

	// Clear the FCSR flags
	s.Clr(csr.FFLAGS, fflagsALL)

	var x float32
	switch rm {
	case frmRNE:
		x = f1 * f2
	case frmRTZ:
		x = f1 * f2
	case frmRDN:
	case frmRUP:
	case frmRRM:
	default:
		return 0, errors.New("illegal")
	}

	return x, nil
}

//-----------------------------------------------------------------------------
