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

// Rounding modes.
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

// FCSR fflags bits.
const (
	fflagsNX = C.FFLAG_INEXACT     // Inexact
	fflagsUF = C.FFLAG_UNDERFLOW   // Underflow
	fflagsOF = C.FFLAG_OVERFLOW    // Overflow
	fflagsDZ = C.FFLAG_DIVIDE_ZERO // Divide by Zero
	fflagsNV = C.FFLAG_INVALID_OP  // Invalid Operation
)

//-----------------------------------------------------------------------------

const mask30to0 = (1 << 31) - 1
const f32SignMask = 1 << 31
const f64SignMask = 1 << 63

// neg32 changes the sign of a float32
func neg32(a uint32) uint32 {
	return a ^ f32SignMask
}

// neg64 changes the sign of a float64
func neg64(a uint64) uint64 {
	return a ^ f64SignMask
}

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

// feq_s returns a == b
func feq_s(a, b uint32, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.eq_quiet_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

// flt_s return a < b
func flt_s(a, b uint32, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.lt_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

// fle_s returns a <= b
func fle_s(a, b uint32, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.le_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

// feq_d returns a == b
func feq_d(a, b uint64, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.eq_quiet_sf64(C.sfloat64(a), C.sfloat64(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

// flt_d return a < b
func flt_d(a, b uint64, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.lt_sf64(C.sfloat64(a), C.sfloat64(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

// fle_d returns a <= b
func fle_d(a, b uint64, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.le_sf64(C.sfloat64(a), C.sfloat64(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

//-----------------------------------------------------------------------------

// fcvt_s_w converts an int32 to a 32-bit float
func fcvt_s_w(a int32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_i32_sf32(C.int32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_d_w converts an int32 to a 64-bit float
func fcvt_d_w(a int32, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.cvt_i32_sf64(C.int32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_s_wu converts a uint32 to a 32-bit float
func fcvt_s_wu(a uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_u32_sf32(C.uint32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_d_wu converts a uint32 to a 64-bit float
func fcvt_d_wu(a uint32, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.cvt_u32_sf64(C.uint32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_w_s converts a 32-bit float to a int32
func fcvt_w_s(a uint32, rm uint, s *csr.State) (int32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := int32(C.cvt_sf32_i32(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_wu_s converts a 32-bit float to a uint32
func fcvt_wu_s(a uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_sf32_u32(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_s_d converts a 64-bit float to a 32-bit float
func fcvt_s_d(a uint64, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_sf64_sf32(C.sfloat64(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fcvt_d_s converts a 32-bit float to a 64-bit float
func fcvt_d_s(a uint32, s *csr.State) uint64 {
	var flags C.uint32_t
	x := uint64(C.cvt_sf32_sf64(C.sfloat32(a), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

//-----------------------------------------------------------------------------

// fadd_s adds two 32-bit floats
func fadd_s(a, b uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.add_sf32(C.sfloat32(a), C.sfloat32(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fadd_d adds two 64-bit floats
func fadd_d(a, b uint64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.add_sf64(C.sfloat64(a), C.sfloat64(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fsub_s subtracts two 32-bit floats
func fsub_s(a, b uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.sub_sf32(C.sfloat32(a), C.sfloat32(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fsub_d subtracts two 64-bit floats
func fsub_d(a, b uint64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.sub_sf64(C.sfloat64(a), C.sfloat64(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

//-----------------------------------------------------------------------------

// fmul_s multiplies two 32-bit floats
func fmul_s(a, b uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.mul_sf32(C.sfloat32(a), C.sfloat32(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fmul_d multiplies two 64-bit floats
func fmul_d(a, b uint64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.mul_sf64(C.sfloat64(a), C.sfloat64(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fdiv_s divides two 32-bit floats
func fdiv_s(a, b uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.div_sf32(C.sfloat32(a), C.sfloat32(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fdiv_d divides two 64-bit floats
func fdiv_d(a, b uint64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.div_sf64(C.sfloat64(a), C.sfloat64(b), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

//-----------------------------------------------------------------------------

// fmin_s returns the minimum of two 32-bit floats
func fmin_s(a, b uint32, s *csr.State) uint32 {
	var flags C.uint32_t
	x := uint32(C.min_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

// fmax_s returns the maximum of two 32-bit floats
func fmax_s(a, b uint32, s *csr.State) uint32 {
	var flags C.uint32_t
	x := uint32(C.max_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x
}

//-----------------------------------------------------------------------------

// fclass_s returns the class of a 32-bit float
func fclass_s(a uint32) uint32 {
	return uint32(C.fclass_sf32(C.sfloat32(a)))
}

// fclass_d returns the class of a 64-bit float
func fclass_d(a uint64) uint32 {
	return uint32(C.fclass_sf64(C.sfloat64(a)))
}

//-----------------------------------------------------------------------------

// fsqrt_s returns the square root of a 32-bit float
func fsqrt_s(a uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.sqrt_sf32(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

// fsqrt_d returns the square root of a 64-bit float
func fsqrt_d(a uint64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.sqrt_sf64(C.sfloat64(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

//-----------------------------------------------------------------------------

// fmadd_s returns the fused-multiply-add of 32-bit floats
func fmadd_s(a, b, c uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.fma_sf32(C.sfloat32(a), C.sfloat32(b), C.sfloat32(c), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint(flags))
	return x, nil
}

//-----------------------------------------------------------------------------
