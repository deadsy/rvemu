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

	isa := NewISA("rv32g")
	err := isa.Add(ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.Error("FAIL")
	}

	daTest := []struct {
		adr uint32
		ins uint32
		da  string
	}{
		{0, 0, "?"},
		{0, 0xffffffff, "?"},
		{0, 0xfe010113, "addi sp,sp,-32"},
		{0, 0x02010413, "addi s0,sp,32"},
		{0, 0x00050793, "mv a5,a0"},
		{0, 0x00f00793, "li a5,15"},
		{0, 0x00078513, "mv a0,a5"},
		{0, 0x00812e23, "sw s0,28(sp)"},
		{0, 0xfef407a3, "sb a5,-17(s0)"},
		{0, 0xfef44783, "lbu a5,-17(s0)"},
		{0, 0x00f7f793, "andi a5,a5,15"},
		{0, 0x0ff7f793, "andi a5,a5,255"},
		{0, 0xfef44703, "lbu a4,-17(s0)"},
		{0x28, 0x02e7f063, "bgeu a5,a4,48"},
		{0x34, 0x00e7ea63, "bltu a5,a4,48"},
		{0x334, 0xf8f74ce3, "blt a4,a5,2cc"},
		{0, 0x01c12403, "lw s0,28(sp)"},
		{0, 0x800005b7, "lui a1,0x80000"},
		{0, 0x00000097, "auipc ra,0x0"},
		{0x44, 0x0100006f, "j 54"},
		{0, 0x00008067, "ret"},
		{0, 0x000080e7, "jalr ra"},
		{0, 0x01079793, "slli a5,a5,0x10"},
		{0, 0x0107d793, "srli a5,a5,0x10"},
	}

	for _, v := range daTest {
		da, _ := isa.daInstruction(v.adr, v.ins)
		if v.da != da {
			fmt.Printf("ins %08x \"%s\" (expected) \"%s\" (actual)\n", v.ins, v.da, da)
			t.Error("FAIL")
		}
	}

}

//-----------------------------------------------------------------------------
