//-----------------------------------------------------------------------------
/*

RISC-V Disassembler/Emulator Code Generation

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------

/*

// ISAModule is the numeric identifier of an ISA sub-module.
type ISAModuleID uint32

// Identifiers for ISA sub-modules.
const (
	IDv32i ISAModule = (1 << iota) // Integer
	IDrv32m                         // Integer Multiplication and Division
	IDrv32a                         // Atomics
	IDrv32f                         // Single-Precision Floating-Point
	IDrv32d                         // Double-Precision Floating-Point
	IDrv64i                         // Integer
	IDrv64m                         // Integer Multiplication and Division
	IDrv64a                         // Atomics
	IDrv64f                         // Single-Precision Floating-Point
	IDrv64d                         // Double-Precision Floating-Point
)

*/

//-----------------------------------------------------------------------------

// bits2vm converts a bit pattern to a value and mask.
func bits2vm(s string) (uint32, uint32) {
	var v uint32
	var m uint32
	for _, c := range s {
		v <<= 1
		m <<= 1
		if c == '0' || c == '1' {
			m++
			if c == '1' {
				v++
			}
		}
	}
	return v, m
}

// vm2bits converts a value and mask into a bit pattern.
func vm2bits(v, m uint32) string {
	s := make([]rune, 32)
	mask := uint32(1 << 31)
	for i := range s {
		if m&mask == 0 {
			s[i] = '.'
		} else {
			if v&mask != 0 {
				s[i] = '1'
			} else {
				s[i] = '0'
			}
		}
		mask >>= 1
	}
	return string(s)
}

// dontCare returns n don't care characters.
func dontCare(n int) string {
	s := make([]string, n)
	for i := range s {
		s[i] = "."
	}
	return strings.Join(s, "")
}

//-----------------------------------------------------------------------------

// isBits returns if a string contains only 0 and 1
func isBits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c != '0' && c != '1' {
			return false
		}
	}
	return true
}

var knownFields = map[string]int{
	"imm[31:12]":            20,
	"imm[20|10:1|11|19:12]": 20,
	"imm[11:0]":             12,
	"imm[12|10:5]":          7,
	"imm[4:1|11]":           5,
	"imm[11:5]":             7,
	"imm[4:0]":              5,
	"shamt5":                5,
	"shamt6":                6,
	"pred":                  4,
	"succ":                  4,
	"csr":                   12,
	"zimm":                  5,
	"rd":                    5,
	"rs1":                   5,
	"rs2":                   5,
	"rs3":                   5,
	"rm":                    3,
	"aq":                    1,
	"rl":                    1,
}

// isField returns the length of an instruction field.
func isField(s string) (int, error) {
	if n, ok := knownFields[s]; ok {
		return n, nil
	}
	return 0, fmt.Errorf("field not recognised \"%s\"", s)
}

//-----------------------------------------------------------------------------

type daFunc func(mneumonic string, adr, ins uint32) (string, string)

type decoders struct {
	da daFunc
}

var knownDecodes = map[string]decoders{
	"imm[31:12]_rd_7b":                       decoders{daNone},
	"imm[20|10:1|11|19:12]_rd_7b":            decoders{daNone},
	"imm[11:0]_rs1_3b_rd_7b":                 decoders{daNone},
	"imm[12|10:5]_rs2_rs1_3b_imm[4:1|11]_7b": decoders{daNone},
	"imm[11:5]_rs2_rs1_3b_imm[4:0]_7b":       decoders{daNone},
	"7b_shamt5_rs1_3b_rd_7b":                 decoders{daNone},
	"7b_rs2_rs1_3b_rd_7b":                    decoders{daNone},
	"4b_pred_succ_5b_3b_5b_7b":               decoders{daNone},
	"4b_4b_4b_5b_3b_5b_7b":                   decoders{daNone},
	"12b_5b_3b_5b_7b":                        decoders{daNone},
	"csr_rs1_3b_rd_7b":                       decoders{daNone},
	"csr_zimm_3b_rd_7b":                      decoders{daNone},
	"5b_aq_rl_5b_rs1_3b_rd_7b":               decoders{daNone},
	"5b_aq_rl_rs2_rs1_3b_rd_7b":              decoders{daNone},
	"rs3_2b_rs2_rs1_rm_rd_7b":                decoders{daNone},
	"7b_rs2_rs1_rm_rd_7b":                    decoders{daNone},
	"7b_5b_rs1_rm_rd_7b":                     decoders{daNone},
	"7b_5b_rs1_3b_rd_7b":                     decoders{daNone},
	"6b_shamt6_rs1_3b_rd_7b":                 decoders{daNone},
}

// getDecode returns the decode type for the instruction.
func getDecode(s string) (*decoders, error) {
	if decode, ok := knownDecodes[s]; ok {
		return &decode, nil
	}
	return nil, fmt.Errorf("decode signature not recognised \"%s\"", s)
}

//-----------------------------------------------------------------------------

// parseDefn parses an instruction definition string.
func parseDefn(defn string, module string) (*instructionInfo, error) {
	parts := strings.Split(defn, " ")
	n := len(parts)
	if n <= 0 {
		return nil, fmt.Errorf("bad instruction definition string \"%s\"", defn)
	}

	ii := instructionInfo{
		mneumonic: strings.ToLower(parts[n-1]),
		module:    module,
	}

	// remove the mneumonic from the end
	parts = parts[0 : n-1]

	s0 := make([]string, 0) // bit pattern
	s1 := make([]string, 0) // decode signature

	for _, x := range parts {
		if isBits(x) {
			s0 = append(s0, fmt.Sprintf("%s", x))
			s1 = append(s1, fmt.Sprintf("%db", len(x)))
		} else {
			n, err := isField(x)
			if err == nil {
				s0 = append(s0, dontCare(n))
				s1 = append(s1, x)
			} else {
				return nil, err
			}
		}
	}

	// instruction value and mask
	bits := strings.Join(s0, "")
	if len(bits) != 32 {
		return nil, fmt.Errorf("bit length != 32 \"%s\"", defn)
	}
	ii.val, ii.mask = bits2vm(bits)

	// instruction decode
	decode, err := getDecode(strings.Join(s1, "_"))
	if err != nil {
		return nil, err
	}
	ii.decode = decode

	return &ii, nil
}

//-----------------------------------------------------------------------------

// GenDecoder generates a decoder for a set of instructions.
func (isa *ISA) GenDecoder(name string) string {

	mask := uint32(0xffffffff)
	for _, ii := range isa.instruction {
		mask &= ii.mask
	}

	sets := make(map[uint32][]*instructionInfo)

	for _, ii := range isa.instruction {
		val := mask & ii.val
		sets[val] = append(sets[val], &ii)
	}

	for k, v := range sets {
		fmt.Printf("%08x: %d\n", k, len(v))
	}

	s := make([]string, len(isa.instruction))
	for i, ii := range isa.instruction {
		s[i] = fmt.Sprintf("%s %08x %08x %s", vm2bits(ii.val, ii.mask), ii.val, ii.mask, ii.mneumonic)
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
