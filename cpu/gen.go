//-----------------------------------------------------------------------------
/*

RISC-V Disassembler/Emulator Code Generation

*/
//-----------------------------------------------------------------------------

package cpu

import (
	"fmt"
	"strings"
)

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

var knownDecodes = map[string]decodeType{
	"imm[31:12]_rd_7b":                       decodeNone,
	"imm[20|10:1|11|19:12]_rd_7b":            decodeNone,
	"imm[11:0]_rs1_3b_rd_7b":                 decodeNone,
	"imm[12|10:5]_rs2_rs1_3b_imm[4:1|11]_7b": decodeNone,
	"imm[11:5]_rs2_rs1_3b_imm[4:0]_7b":       decodeNone,
	"7b_shamt5_rs1_3b_rd_7b":                 decodeNone,
	"7b_rs2_rs1_3b_rd_7b":                    decodeNone,
	"4b_pred_succ_5b_3b_5b_7b":               decodeNone,
	"4b_4b_4b_5b_3b_5b_7b":                   decodeNone,
	"12b_5b_3b_5b_7b":                        decodeNone,
	"csr_rs1_3b_rd_7b":                       decodeNone,
	"csr_zimm_3b_rd_7b":                      decodeNone,
	"5b_aq_rl_5b_rs1_3b_rd_7b":               decodeNone,
	"5b_aq_rl_rs2_rs1_3b_rd_7b":              decodeNone,
	"rs3_2b_rs2_rs1_rm_rd_7b":                decodeNone,
	"7b_rs2_rs1_rm_rd_7b":                    decodeNone,
	"7b_5b_rs1_rm_rd_7b":                     decodeNone,
	"7b_5b_rs1_3b_rd_7b":                     decodeNone,
	"6b_shamt6_rs1_3b_rd_7b":                 decodeNone,
}

// getDecode returns the decode type for the instruction.
func getDecode(s string) (decodeType, error) {
	if t, ok := knownDecodes[s]; ok {
		return t, nil
	}
	return decodeNone, fmt.Errorf("decode signature not recognised \"%s\"", s)
}

//-----------------------------------------------------------------------------

// getInfo converts an instruction description string into an instruction information structure.
func getInfo(ins string, module ISAModule) (*insInfo, error) {
	parts := strings.Split(ins, " ")
	n := len(parts)
	if n <= 0 {
		return nil, fmt.Errorf("bad instruction string \"%s\"", ins)
	}

	d := insInfo{
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
		return nil, fmt.Errorf("bit length != 32 \"%s\"", ins)
	}
	d.val, d.mask = bits2vm(bits)

	// instruction decode
	t, err := getDecode(strings.Join(s1, "_"))
	if err != nil {
		return nil, err
	}
	d.decode = t

	return &d, nil
}

//-----------------------------------------------------------------------------

// ISA is an instruction set
type ISA struct {
	name string
	ins  []*insInfo
}

// NewISA creates an empty instruction set.
func NewISA(name string) *ISA {
	return &ISA{
		name: name,
		ins:  make([]*insInfo, 0),
	}
}

// Add adds an instruction to the instruction set.
func (isa *ISA) Add(defn string, module ISAModule) error {
	d, err := getInfo(defn, module)
	if err != nil {
		return err
	}
	is.ins = append(is.ins, d)
	return nil
}

//-----------------------------------------------------------------------------

// GenDecoder generates a decoder for a set of instructions.
func (is *InsSet) GenDecoder(name string) string {

	mask := uint32(0xffffffff)
	for i := range is.ins {
		mask &= is.ins[i].mask
	}

	sets := make(map[uint32][]*insInfo)

	for i := range is.ins {
		val := mask & is.ins[i].val
		sets[val] = append(sets[val], is.ins[i])
	}

	for k, v := range sets {
		fmt.Printf("%08x: %d\n", k, len(v))
	}

	s := make([]string, len(is.ins))
	for i, d := range is.ins {
		s[i] = fmt.Sprintf("%s %08x %08x %s", vm2bits(d.val, d.mask), d.val, d.mask, d.mneumonic)
	}
	return strings.Join(s, "\n")

}

//-----------------------------------------------------------------------------
