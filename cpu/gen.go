//-----------------------------------------------------------------------------
/*

RISC-V Code Generation

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

func genDecode(ins string) *insDecode {
	return nil
}

func addSet(set []string) []*insDecode {
	decodes := make([]*insDecode, len(set))
	for i := range set {
		decodes[i] = genDecode(set[i])
	}
	return decodes
}

//-----------------------------------------------------------------------------

type InsSet struct {
	name    string
	decodes []*insDecode
}

func NewInsSet(name string) *InsSet {
	return &InsSet{
		name: name,
	}
}

func (is *InsSet) Add(set []string) {
}

func (is *InsSet) GenerateDisassembler() string {
	return ""
}

//-----------------------------------------------------------------------------
