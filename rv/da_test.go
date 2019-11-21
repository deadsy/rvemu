//-----------------------------------------------------------------------------
/*

RISC-V Disassembler Testing

*/
//-----------------------------------------------------------------------------

package rv

import (
	"fmt"
	"testing"
)

//-----------------------------------------------------------------------------

type daTest struct {
	pc  uint   // program counter
	ins uint   // instruction code
	da  string // expected disassembly
}

//-----------------------------------------------------------------------------
// rv32

var rv32iTest = []daTest{
	{0, 0, "?"},
	{0, 0xffffffff, "?"},
	{0, 0x800005b7, "lui a1,0x80000"},
	{0, 0xdeadc7b7, "lui a5,0xdeadc"},
	{0, 0x00000097, "auipc ra,0x0"},
	{0x44, 0x0100006f, "j 54"},
	{0, 0x00008067, "ret"},
	{0, 0x000080e7, "jalr ra"},
	{0x28, 0x02e7f063, "bgeu a5,a4,48"},
	{0x34, 0x00e7ea63, "bltu a5,a4,48"},
	{0x334, 0xf8f74ce3, "blt a4,a5,2cc"},
	{0x204, 0x0007da63, "bgez a5,218"},
	{0, 0x02050463, "beqz a0,28"},
	{0x20, 0x01185a63, "bge a6,a7,34"},
	{0x24, 0xfea614e3, "bne a2,a0,c"},
	{0x4, 0x06b69e63, "bne a3,a1,80"},
	{0xc, 0xfe069ae3, "bnez a3,0"},
	{0, 0xfef44783, "lbu a5,-17(s0)"},
	{0, 0xfef44703, "lbu a4,-17(s0)"},
	{0, 0x01c12403, "lw s0,28(sp)"},
	{0, 0xffc62883, "lw a7,-4(a2)"},
	{0, 0xfef407a3, "sb a5,-17(s0)"},
	{0, 0x00812e23, "sw s0,28(sp)"},
	{0, 0xfe010113, "addi sp,sp,-32"},
	{0, 0x02010413, "addi s0,sp,32"},
	{0, 0xeef78593, "addi a1,a5,-273"},
	{0, 0x00050793, "mv a5,a0"},
	{0, 0x00078513, "mv a0,a5"},
	{0, 0x00f00793, "li a5,15"},
	{0, 0x00f7f793, "andi a5,a5,15"},
	{0, 0x0ff7f793, "andi a5,a5,255"},
	{0, 0xfff7c593, "not a1,a5"},
	{0, 0x01079793, "slli a5,a5,0x10"},
	{0, 0x0107d793, "srli a5,a5,0x10"},
	{0, 0x00a60533, "add a0,a2,a0"},
	{0, 0x40f007b3, "neg a5,a5"},
	{0x40, 0x00301073, "fscsr zero"},
	{0xce, 0xf1402573, "csrr a0,mhartid"},
	{0x114, 0xfea42623, "sw a0,-20(s0)"},
	{0x118, 0xfeb42423, "sw a1,-24(s0)"},
	{0x11c, 0xfe842783, "lw a5,-24(s0)"},
	{0xfe, 0xfec42783, "lw a5,-20(s0)"},
	{0, 0x00000073, "ecall"},
}

var rv32mTest = []daTest{
	{0, 0x025535b3, "mulhu a1,a0,t0"},
}

var rv32aTest = []daTest{
	{0, 0x100526af, "lr.w a3,(a0)"},
	{0, 0x18c526af, "sc.w a3,a2,(a0)"},
	{0, 0x0c55232f, "amoswap.w t1,t0,(a0)"},
}

var rv32fTest = []daTest{
	{0x14, 0xfec42707, "flw fa4,-20(s0)"},
	{0x18, 0xfe842787, "flw fa5,-24(s0)"},
	{0x1c, 0x10f777d3, "fmul.s fa5,fa4,fa5"},
	{0x20, 0xe0078553, "fmv.x.w a0,fa5"},
	{0x4c, 0x18f777d3, "fdiv.s fa5,fa4,fa5"},
	{0x44, 0xf0000053, "fmv.w.x ft0,zero"},
}

var rv32dTest = []daTest{
	{0, 0x0005b787, "fld fa5,0(a1)"},
	{0, 0x72a7f7c3, "fmadd.d fa5,fa5,fa0,fa4"},
	{0, 0xfef63c27, "fsd fa5,-8(a2)"},
}

var rv32cTest = []daTest{
	{0, 0, "?"},
	{0, 0x4705, "li a4,1"},
	{0, 0x8082, "ret"},
	{0, 0xce06, "sw ra,28(sp)"},
	{0, 0xcc22, "sw s0,24(sp)"},
	{0, 0xca26, "sw s1,20(sp)"},
	{0, 0x40f2, "lw ra,28(sp)"},
	{0, 0x4462, "lw s0,24(sp)"},
	{0, 0x44d2, "lw s1,20(sp)"},
	{0, 0x6145, "addi sp,sp,48"},
	{0, 0x1800, "addi s0,sp,48"},
	{0, 0x1101, "addi sp,sp,-32"},
	{0, 0x873e, "mv a4,a5"},
	{0, 0x8391, "srli a5,a5,0x4"},
	{0, 0x0742, "slli a4,a4,0x10"},
	{0, 0x8bbd, "andi a5,a5,15"},
	{0, 0x97b6, "add a5,a5,a3"},
	{0x186, 0xa029, "j 190"},
	{0, 0x67ad, "lui a5,0xb"},
	{0x1d0, 0xf3e1, "bnez a5,190"},
	{0x358, 0x3d7d, "jal ra,216"},
	{0, 0x0001, "nop"},
	{0, 0x8e09, "sub a2,a2,a0"},
	{0, 0xc30c, "sw a1,0(a4)"},
	{0, 0xc74c, "sw a1,12(a4)"},
	{0, 0xc78c, "sw a1,8(a5)"},
	{0, 0x4398, "lw a4,0(a5)"},
	{0, 0x40d4, "lw a3,4(s1)"},
	{0, 0x4388, "lw a0,0(a5)"},
	{0, 0x9682, "jalr a3"},
	{0, 0x9782, "jalr a5"},
}

var rv32fcTest = []daTest{}
var rv32dcTest = []daTest{}

//-----------------------------------------------------------------------------
// rv64

var rv64iTest = []daTest{}
var rv64mTest = []daTest{}
var rv64aTest = []daTest{}
var rv64fTest = []daTest{}
var rv64dTest = []daTest{}

var rv64cTest = []daTest{
	{0, 0xe30c, "sd a1,0(a4)"},
	{0, 0xe70c, "sd a1,8(a4)"},
	{0, 0x6398, "ld a4,0(a5)"},
	{0, 0x60a6, "ld ra,72(sp)"},
	{0, 0x6406, "ld s0,64(sp)"},
	{0, 0x74e2, "ld s1,56(sp)"},
	{0, 0x2705, "addiw a4,a4,1"},
	{0, 0x347d, "addiw s0,s0,-1"},
	{0, 0x37fd, "addiw a5,a5,-1"},
}

//-----------------------------------------------------------------------------

func testSet(module []ISAModule, tests []daTest) error {
	isa := NewISA()
	err := isa.Add(module)
	if err != nil {
		return err
	}
	for _, v := range tests {
		da := isa.daInstruction(v.pc, v.ins)
		if v.da != da {
			return fmt.Errorf("ins %08x \"%s\" (expected) \"%s\" (actual)\n", v.ins, v.da, da)
		}
	}
	return nil
}

//-----------------------------------------------------------------------------

func Test_Disassembly(t *testing.T) {

	rv32Tests := make([]daTest, 0)
	rv32Tests = append(rv32Tests, rv32iTest...)
	rv32Tests = append(rv32Tests, rv32mTest...)
	rv32Tests = append(rv32Tests, rv32aTest...)
	rv32Tests = append(rv32Tests, rv32fTest...)
	rv32Tests = append(rv32Tests, rv32dTest...)
	rv32Tests = append(rv32Tests, rv32cTest...)
	rv32Tests = append(rv32Tests, rv32fcTest...)
	rv32Tests = append(rv32Tests, rv32dcTest...)

	rv64Tests := make([]daTest, 0)
	rv64Tests = append(rv64Tests, rv32iTest...)
	rv64Tests = append(rv64Tests, rv32mTest...)
	rv64Tests = append(rv64Tests, rv32aTest...)
	rv64Tests = append(rv64Tests, rv32fTest...)
	rv64Tests = append(rv64Tests, rv32dTest...)
	rv64Tests = append(rv64Tests, rv32cTest...)
	rv64Tests = append(rv64Tests, rv32dcTest...)
	rv64Tests = append(rv64Tests, rv64iTest...)
	rv64Tests = append(rv64Tests, rv64mTest...)
	rv64Tests = append(rv64Tests, rv64aTest...)
	rv64Tests = append(rv64Tests, rv64fTest...)
	rv64Tests = append(rv64Tests, rv64dTest...)
	rv64Tests = append(rv64Tests, rv64cTest...)

	testCases := []struct {
		module []ISAModule
		tests  []daTest
	}{
		// rv32
		{[]ISAModule{ISArv32i}, rv32iTest},
		{[]ISAModule{ISArv32m}, rv32mTest},
		{[]ISAModule{ISArv32a}, rv32aTest},
		{[]ISAModule{ISArv32f}, rv32fTest},
		{[]ISAModule{ISArv32d}, rv32dTest},
		{[]ISAModule{ISArv32c}, rv32cTest},
		{[]ISAModule{ISArv32fc}, rv32fcTest},
		{[]ISAModule{ISArv32dc}, rv32dcTest},
		// rv64
		{[]ISAModule{ISArv64i}, rv64iTest},
		{[]ISAModule{ISArv64m}, rv64mTest},
		{[]ISAModule{ISArv64a}, rv64aTest},
		{[]ISAModule{ISArv64f}, rv64fTest},
		{[]ISAModule{ISArv64d}, rv64dTest},
		{[]ISAModule{ISArv64c}, rv64cTest},
		// together
		{ISArv32gc, rv32Tests},
		{ISArv64gc, rv64Tests},
	}
	for _, v := range testCases {
		err := testSet(v.module, v.tests)
		if err != nil {
			fmt.Printf("%s\n", err)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------
