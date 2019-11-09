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

// bits2vm converts a bit pattern to a value and mask.
func bits2vm(s string) (uint, uint) {
	var v uint
	var m uint
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
func vm2bits(v, m uint) string {
	s := make([]rune, 32)
	mask := uint(1 << 31)
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
	"imm[31:12]":                 20,
	"imm[20|10:1|11|19:12]":      20,
	"imm[11:0]":                  12,
	"imm[12|10:5]":               7,
	"imm[4:1|11]":                5,
	"imm[11:5]":                  7,
	"imm[4:0]":                   5,
	"shamt5":                     5,
	"shamt6":                     6,
	"pred":                       4,
	"succ":                       4,
	"csr":                        12,
	"zimm":                       5,
	"rd":                         5,
	"rs1":                        5,
	"rs2":                        5,
	"rs3":                        5,
	"rm":                         3,
	"aq":                         1,
	"rl":                         1,
	"imm[11|4|9:8|10|6|7|3:1|5]": 11,
	"imm[7:6|2:1|5]":             5,
	"imm[8|4:3]":                 3,
	"imm[5]":                     1,
	"uimm[4:3|8:6]":              5,
	"uimm[4:2|7:6]":              5,
	"uimm[5:3|8:6]":              6,
	"uimm[5:2|7:6]":              6,
	"uimm[5:3]":                  3,
	"uimm[7:6]":                  2,
	"uimm[2|6]":                  2,
	"uimm[5]":                    1,
	"nzimm[4|6|8:7|5]":           5,
	"nzimm[16:12]":               5,
	"nzimm[4:0]":                 5,
	"nzimm[5]":                   1,
	"nzimm[9]":                   1,
	"nzimm[17]":                  1,
	"nzuimm[5:4|9:6|2|3]":        8,
	"nzuimm[5]":                  1,
	"nzuimm[4:0]":                5,
	"rd0":                        3,
	"rs10":                       3,
	"rs20":                       3,
	"rs1/rd!=0":                  5,
	"rd!=0":                      5,
	"rs1!=0":                     5,
	"rs2!=0":                     5,
	"rd!={0,2}":                  5,
	"rs10/rd0":                   3,
}

// isField returns the length of an instruction field.
func isField(s string) (int, error) {
	if n, ok := knownFields[s]; ok {
		return n, nil
	}
	return 0, fmt.Errorf("field not recognised \"%s\"", s)
}

//-----------------------------------------------------------------------------

type decodeType int

const (
	decodeTypeNone = iota // unknown
	decodeTypeR           // Register
	decodeTypeI           // Immediate
	decodeTypeS           // Store
	decodeTypeB           // Branch
	decodeTypeU           // Upper Immediate
	decodeTypeJ           // Jump
	decodeTypeR4          // Register(4)
	decodeTypeCR          // Compressed Register
	decodeTypeCI          // Compressed Immediate
	decodeTypeCSS         // Compressed Stack-relative Store
	decodeTypeCIW         // Compressed Wide Immediate
	decodeTypeCL          // Compressed Load
	decodeTypeCS          // Compressed Store
	decodeTypeCB          // Compressed Branch
	decodeTypeCJ          // Compressed Jump
)

var knownDecodes = map[string]decodeType{
	"imm[31:12]_rd_7b":                        decodeTypeU,
	"imm[20|10:1|11|19:12]_rd_7b":             decodeTypeJ, // aka UJ
	"imm[11:5]_rs2_rs1_3b_imm[4:0]_7b":        decodeTypeS,
	"imm[12|10:5]_rs2_rs1_3b_imm[4:1|11]_7b":  decodeTypeB, // aka SB
	"7b_shamt5_rs1_3b_rd_7b":                  decodeTypeI,
	"6b_shamt6_rs1_3b_rd_7b":                  decodeTypeI,
	"imm[11:0]_rs1_3b_rd_7b":                  decodeTypeI,
	"csr_rs1_3b_rd_7b":                        decodeTypeI,
	"csr_zimm_3b_rd_7b":                       decodeTypeI,
	"4b_pred_succ_5b_3b_5b_7b":                decodeTypeI,
	"12b_5b_3b_5b_7b":                         decodeTypeI,
	"4b_4b_4b_5b_3b_5b_7b":                    decodeTypeI,
	"7b_rs2_rs1_3b_rd_7b":                     decodeTypeR,
	"7b_rs2_rs1_rm_rd_7b":                     decodeTypeR,
	"7b_5b_rs1_rm_rd_7b":                      decodeTypeR,
	"5b_aq_rl_5b_rs1_3b_rd_7b":                decodeTypeR,
	"5b_aq_rl_rs2_rs1_3b_rd_7b":               decodeTypeR,
	"7b_5b_rs1_3b_rd_7b":                      decodeTypeR,
	"rs3_2b_rs2_rs1_rm_rd_7b":                 decodeTypeR4,
	"3b_nzuimm[5:4|9:6|2|3]_rd0_2b":           decodeTypeCIW,
	"3b_8b_3b_2b":                             decodeTypeCIW,
	"3b_uimm[5:3]_rs10_uimm[7:6]_rd0_2b":      decodeTypeCL,
	"3b_uimm[5:3]_rs10_uimm[2|6]_rd0_2b":      decodeTypeCL,
	"3b_uimm[5:3]_rs10_uimm[7:6]_rs20_2b":     decodeTypeCS,
	"3b_uimm[5:3]_rs10_uimm[2|6]_rs20_2b":     decodeTypeCS,
	"3b_nzimm[5]_5b_nzimm[4:0]_2b":            decodeTypeCI,
	"3b_nzimm[5]_rs1/rd!=0_nzimm[4:0]_2b":     decodeTypeCI,
	"3b_imm[11|4|9:8|10|6|7|3:1|5]_2b":        decodeTypeCJ,
	"3b_imm[5]_rd!=0_imm[4:0]_2b":             decodeTypeCI,
	"3b_nzimm[9]_5b_nzimm[4|6|8:7|5]_2b":      decodeTypeCI,
	"3b_nzimm[17]_rd!={0,2}_nzimm[16:12]_2b":  decodeTypeCI,
	"3b_nzuimm[5]_2b_rs10/rd0_nzuimm[4:0]_2b": decodeTypeCI,
	"3b_imm[5]_2b_rs10/rd0_imm[4:0]_2b":       decodeTypeCI,
	"3b_1b_2b_rs10/rd0_2b_rs20_2b":            decodeTypeCR,
	"3b_imm[8|4:3]_rs10_imm[7:6|2:1|5]_2b":    decodeTypeCB,
	"3b_nzuimm[5]_rs1/rd!=0_nzuimm[4:0]_2b":   decodeTypeCI,
	"3b_1b_rs1/rd!=0_5b_2b":                   decodeTypeCI,
	"3b_uimm[5]_rd_uimm[4:3|8:6]_2b":          decodeTypeCSS,
	"3b_uimm[5]_rd!=0_uimm[4:2|7:6]_2b":       decodeTypeCSS,
	"3b_uimm[5]_rd_uimm[4:2|7:6]_2b":          decodeTypeCSS,
	"3b_1b_rs1!=0_5b_2b":                      decodeTypeCJ,
	"3b_1b_rd!=0_rs2!=0_2b":                   decodeTypeCR,
	"3b_1b_5b_5b_2b":                          decodeTypeCI,
	"3b_1b_rs1/rd!=0_rs2!=0_2b":               decodeTypeCR,
	"3b_uimm[5:3|8:6]_rs2_2b":                 decodeTypeCSS,
	"3b_uimm[5:2|7:6]_rs2_2b":                 decodeTypeCSS,
}

// getDecode returns the decode type for the instruction.
func getDecode(s string) (decodeType, error) {
	if t, ok := knownDecodes[s]; ok {
		return t, nil
	}
	return decodeTypeNone, fmt.Errorf("decode signature not recognised \"%s\"", s)
}

//-----------------------------------------------------------------------------

// parseDefn parses an instruction definition string and returns the meta-data.
func parseDefn(id *insDefn, ilen int) (*insMeta, error) {

	im := insMeta{
		defn: id,
		n:    ilen,
	}

	parts := strings.Split(id.defn, " ")
	n := len(parts)
	if n <= 0 {
		return nil, fmt.Errorf("bad instruction definition string \"%s\"", id.defn)
	}

	// mneumonic
	im.name = strings.ToLower(parts[n-1])
	// strip the "c." for the compressed instruction set
	im.name = strings.TrimPrefix(im.name, "c.")

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
	if len(bits) != ilen {
		return nil, fmt.Errorf("instruction length != %d \"%s\"", ilen, id.defn)
	}
	im.val, im.mask = bits2vm(bits)

	// set the decode type
	dt, err := getDecode(strings.Join(s1, "_"))
	if err != nil {
		return nil, err
	}
	im.dt = dt

	return &im, nil
}

//-----------------------------------------------------------------------------
