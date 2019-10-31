//-----------------------------------------------------------------------------
/*

RISC-V Code Generation

*/
//-----------------------------------------------------------------------------

package cpu

import (
	"errors"
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------

func parseBits(s string) (int, int, error) {
	n := len(s)
	if n == 0 {
		return 0, 0, errors.New("zero length bit string")
	}
	val := 0
	for _, c := range s {
		if c != '0' && c != '1' {
			return 0, 0, errors.New("not a bit string")
		}
		val <<= 1
		if c == '1' {
			val++
		}
	}
	return val, n, nil
}

var fields = map[string]int{
	"imm[31:12]":            20,
	"imm[20|10:1|11|19:12]": 20,
	"imm[11:0]":             12,
	"imm[12|10:5]":          7,
	"imm[4:1|11]":           5,
	"imm[11:5]":             7,
	"imm[4:0]":              5,
	"shamt":                 5,
	"pred":                  4,
	"succ":                  4,
	"csr":                   12,
	"zimm":                  5,
	"rd":                    4,
	"rs1":                   5,
	"rs2":                   5,
}

func parseField(s string) (int, error) {
	if n, ok := fields[s]; ok {
		return n, nil
	}
	return 0, fmt.Errorf("field not recognised %s", s)
}

func dontCare(n int) string {
	s := make([]string, n)
	for i := range s {
		s[i] = "."
	}
	return strings.Join(s, "")
}

//-----------------------------------------------------------------------------

func genDecode(ins string) (*insDecode, error) {
	parts := strings.Split(ins, " ")
	n := len(parts)
	if n <= 0 {
		return nil, fmt.Errorf("bad instruction string %s\n", ins)
	}

	var d insDecode

	// get the mneumonic off the end
	d.mneumonic = strings.ToLower(parts[n-1])
	parts = parts[0 : n-1]

	s := make([]string, 0)
	for _, x := range parts {
		_, _, err := parseBits(x)
		if err == nil {
			s = append(s, fmt.Sprintf("%s", x))
		} else {
			n, err := parseField(x)
			if err == nil {
				s = append(s, dontCare(n))
			} else {
				return nil, err
			}
		}
	}

	d.blah = strings.Join(s, "")

	return &d, nil
}

//-----------------------------------------------------------------------------

type InsSet struct {
	name   string
	decode []*insDecode
}

func NewInsSet(name string) *InsSet {
	return &InsSet{
		name:   name,
		decode: make([]*insDecode, 0),
	}
}

func (is *InsSet) Add(set []string) error {
	for i := range set {
		d, err := genDecode(set[i])
		if err != nil {
			return err
		}
		is.decode = append(is.decode, d)
	}
	return nil
}

func (is *InsSet) GenerateDisassembler() string {
	s := make([]string, len(is.decode))
	for i, d := range is.decode {
		s[i] = fmt.Sprintf("\"%s\"\t\t\t%s", d.mneumonic, d.blah)
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
