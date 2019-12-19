//-----------------------------------------------------------------------------
/*

RISC-V Floating Point Routines

These is glue to the C-based softfp library that does all the real work.
See: https://bellard.org/softfp/

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

const upper32 = uint64(((1 << 32) - 1) << 32)
const mask30to0 = (1 << 31) - 1
const mask62to0 = (1 << 63) - 1
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
		x, _ := s.Rd(csr.FRM)
		rm = uint(x)
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
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

// flt_s return a < b
func flt_s(a, b uint32, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.lt_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

// fle_s returns a <= b
func fle_s(a, b uint32, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.le_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

// feq_d returns a == b
func feq_d(a, b uint64, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.eq_quiet_sf64(C.sfloat64(a), C.sfloat64(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

// flt_d return a < b
func flt_d(a, b uint64, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.lt_sf64(C.sfloat64(a), C.sfloat64(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

// fle_d returns a <= b
func fle_d(a, b uint64, s *csr.State) uint {
	var flags C.uint32_t
	x := uint(C.le_sf64(C.sfloat64(a), C.sfloat64(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

//-----------------------------------------------------------------------------
// fcvt to {s,d} from {d,s}

// fcvt_s_d converts to float32 from float64
func fcvt_s_d(a uint64, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_sf64_sf32(C.sfloat64(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_d_s converts to float64 from float32
func fcvt_d_s(a uint32, s *csr.State) uint64 {
	var flags C.uint32_t
	x := uint64(C.cvt_sf32_sf64(C.sfloat32(a), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

//-----------------------------------------------------------------------------
// fcvt to {s,d} from {w,wu,l,lu}

// fcvt_s_w converts to float32 from int32
func fcvt_s_w(a int32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_i32_sf32(C.int32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_s_wu converts to float32 from uint32
func fcvt_s_wu(a uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_u32_sf32(C.uint32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_s_l converts to float32 from int64
func fcvt_s_l(a int64, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_i64_sf32(C.int64_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_s_lu converts to float32 from uint64
func fcvt_s_lu(a uint64, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_u64_sf32(C.uint64_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_d_w converts to float64 from int32
func fcvt_d_w(a int32, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.cvt_i32_sf64(C.int32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_d_wu converts to float64 from uint32
func fcvt_d_wu(a uint32, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.cvt_u32_sf64(C.uint32_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_d_l converts to float64 from int64
func fcvt_d_l(a int64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.cvt_i64_sf64(C.int64_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_d_lu converts to float64 from uint64
func fcvt_d_lu(a uint64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.cvt_u64_sf64(C.uint64_t(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

//-----------------------------------------------------------------------------
// fcvt to {w,wu,l,lu} from {s,d}

// fcvt_w_s converts to int32 from float32.
func fcvt_w_s(a uint32, rm uint, s *csr.State) (int32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := int32(C.cvt_sf32_i32(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_wu_s converts to uint32 from float32
func fcvt_wu_s(a uint32, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_sf32_u32(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_l_s converts to int64 from float32.
func fcvt_l_s(a uint32, rm uint, s *csr.State) (int64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := int64(C.cvt_sf32_i64(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_lu_s converts to uint64 from float32
func fcvt_lu_s(a uint32, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.cvt_sf32_u64(C.sfloat32(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_w_d converts to int32 from float64
func fcvt_w_d(a uint64, rm uint, s *csr.State) (int32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := int32(C.cvt_sf64_i32(C.sfloat64(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_wu_d converts to uint32 from float64
func fcvt_wu_d(a uint64, rm uint, s *csr.State) (uint32, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint32(C.cvt_sf64_u32(C.sfloat64(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_l_d converts to int64 from float64
func fcvt_l_d(a uint64, rm uint, s *csr.State) (int64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := int64(C.cvt_sf64_i64(C.sfloat64(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fcvt_lu_d converts to uint64 from float64
func fcvt_lu_d(a uint64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.cvt_sf64_u64(C.sfloat64(a), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

//-----------------------------------------------------------------------------

// fmin_s returns the minimum of two 32-bit floats
func fmin_s(a, b uint32, s *csr.State) uint32 {
	var flags C.uint32_t
	x := uint32(C.min_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

// fmin_d returns the minimum of two 64-bit floats
func fmin_d(a, b uint64, s *csr.State) uint64 {
	var flags C.uint32_t
	x := uint64(C.min_sf64(C.sfloat64(a), C.sfloat64(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

// fmax_s returns the maximum of two 32-bit floats
func fmax_s(a, b uint32, s *csr.State) uint32 {
	var flags C.uint32_t
	x := uint32(C.max_sf32(C.sfloat32(a), C.sfloat32(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

// fmax_d returns the maximum of two 64-bit floats
func fmax_d(a, b uint64, s *csr.State) uint64 {
	var flags C.uint32_t
	x := uint64(C.max_sf64(C.sfloat64(a), C.sfloat64(b), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x
}

//-----------------------------------------------------------------------------

// fclass_s returns the class of a 32-bit float
func fclass_s(a uint32) uint {
	return uint(C.fclass_sf32(C.sfloat32(a)))
}

// fclass_d returns the class of a 64-bit float
func fclass_d(a uint64) uint {
	return uint(C.fclass_sf64(C.sfloat64(a)))
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
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
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

// fmadd_d returns the fused-multiply-add of 64-bit floats
func fmadd_d(a, b, c uint64, rm uint, s *csr.State) (uint64, error) {
	rm, err := getRoundingMode(rm, s)
	if err != nil {
		return 0, err
	}
	var flags C.uint32_t
	x := uint64(C.fma_sf64(C.sfloat64(a), C.sfloat64(b), C.sfloat64(c), C.RoundingModeEnum(rm), &flags))
	s.Wr(csr.FFLAGS, uint64(flags))
	return x, nil
}

//-----------------------------------------------------------------------------
