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

func Test_RV32G(t *testing.T) {

	isa := NewISA()
	err := isa.Add(ISArv32g)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.Error("FAIL")
	}

	daTest := []struct {
		adr uint32 // program counter
		ins uint   // instruction code
		da  string // expected disassembly
	}{
		{0, 0, "?"},
		{0, 0xffffffff, "?"},
		// rv32i
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
		// rv32m
		{0, 0x025535b3, "mulhu a1,a0,t0"},
		// rv32a
		{0, 0x100526af, "lr.w a3,(a0)"},
		{0, 0x18c526af, "sc.w a3,a2,(a0)"},
		{0, 0x0c55232f, "amoswap.w t1,t0,(a0)"},
		// rv32f
		{0x14, 0xfec42707, "flw fa4,-20(s0)"},
		{0x18, 0xfe842787, "flw fa5,-24(s0)"},
		{0x1c, 0x10f777d3, "fmul.s fa5,fa4,fa5"},
		{0x20, 0xe0078553, "fmv.x.w a0,fa5"},
		{0x4c, 0x18f777d3, "fdiv.s fa5,fa4,fa5"},
		{0x44, 0xf0000053, "fmv.w.x ft0,zero"},
		// rv32d
		{0, 0x0005b787, "fld fa5,0(a1)"},
		{0, 0x72a7f7c3, "fmadd.d fa5,fa5,fa0,fa4"},
		{0, 0xfef63c27, "fsd fa5,-8(a2)"},
	}

	for _, v := range daTest {
		da := isa.daInstruction(v.adr, v.ins)
		if v.da != da {
			fmt.Printf("ins %08x \"%s\" (expected) \"%s\" (actual)\n", v.ins, v.da, da)
			t.Error("FAIL")
		}
	}

}

//-----------------------------------------------------------------------------

func Test_RV32C(t *testing.T) {

	isa := NewISA()
	err := isa.Add([]ISAModule{ISArv32c})
	if err != nil {
		fmt.Printf("%s\n", err)
		t.Error("FAIL")
	}

	daTest := []struct {
		adr uint32 // program counter
		ins uint   // instruction code
		da  string // expected disassembly
	}{
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
	}

	for _, v := range daTest {
		da := isa.daInstruction(v.adr, v.ins)
		if v.da != da {
			fmt.Printf("ins %04x \"%s\" (expected) \"%s\" (actual)\n", v.ins, v.da, da)
			t.Error("FAIL")
		}
	}

}

//-----------------------------------------------------------------------------
