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
	{0, 0, "illegal"},
	{0, 0xffffffff, "illegal"},
	{0, 0x800005b7, "lui a1,0x80000"},
	{0, 0xdeadc7b7, "lui a5,0xdeadc"},
	{0, 0x00000097, "auipc ra,0x0"},
	{0, 0x00002297, "auipc t0,0x2"},
	{0, 0x00000017, "auipc zero,0x0"},
	{0, 0x00000297, "auipc t0,0x0"},
	{0x44, 0x0100006f, "j 54"},
	{0, 0x00008067, "ret"},
	{0, 0x000080e7, "jalr ra"},
	{0, 0xffe88fe7, "jalr t6,-2(a7)"},
	{0, 0x00390f67, "jalr t5,3(s2)"},
	{0, 0x001780e7, "jalr 1(a5)"},
	{0, 0xfff70167, "jalr sp,-1(a4)"},
	{0, 0x00040467, "jalr s0,s0"},
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
	{0, 0x00301073, "fscsr zero"},
	{0, 0xfea42623, "sw a0,-20(s0)"},
	{0, 0xfeb42423, "sw a1,-24(s0)"},
	{0, 0xfe842783, "lw a5,-24(s0)"},
	{0, 0xfec42783, "lw a5,-20(s0)"},
	{0, 0x00000073, "ecall"},
	{0, 0x30529073, "csrw mtvec,t0"},
	{0, 0x3b029073, "csrw pmpaddr0,t0"},
	{0, 0x34029073, "csrw mscratch,t0"},
	{0, 0x340290f3, "csrrw ra,mscratch,t0"},
	{0, 0x34029173, "csrrw sp,mscratch,t0"},
	{0, 0x340291f3, "csrrw gp,mscratch,t0"},
	{0, 0x34029273, "csrrw tp,mscratch,t0"},
	{0, 0x340012f3, "csrrw t0,mscratch,zero"},
	{0, 0x18005073, "csrwi satp,0"},
	{0, 0x3400f0f3, "csrrci ra,mscratch,1"},
	{0, 0x34007173, "csrrci sp,mscratch,0"},
	{0, 0x340ff1f3, "csrrci gp,mscratch,31"},
	{0, 0x34087273, "csrrci tp,mscratch,16"},
	{0, 0x3407f2f3, "csrrci t0,mscratch,15"},
	{0, 0x34202f73, "csrr t5,mcause"},
	{0, 0xf1402573, "csrr a0,mhartid"},
	{0, 0x3400b0f3, "csrrc ra,mscratch,ra"},
	{0, 0x34013173, "csrrc sp,mscratch,sp"},
	{0, 0x3401b1f3, "csrrc gp,mscratch,gp"},
	{0, 0x34023273, "csrrc tp,mscratch,tp"},
	{0, 0x3402b2f3, "csrrc t0,mscratch,t0"},
	{0, 0x34003cf3, "csrrc s9,mscratch,zero"},
	{0, 0x30200073, "mret"},
	{0, 0x0ff0000f, "fence"},
	{0, 0x010fa033, "slt zero,t6,a6"},
	{0, 0x00ff20b3, "slt ra,t5,a5"},
	{0, 0x00cda233, "slt tp,s11,a2"},
	{0, 0x800b2493, "slti s1,s6,-2048"},
	{0, 0x7ff9a613, "slti a2,s3,2047"},
	{0, 0x000fa013, "slti zero,t6,0"},
	{0, 0x000fb013, "sltiu zero,t6,0"},
	{0, 0x800d3293, "sltiu t0,s10,-2048"},
	{0, 0x80053a93, "sltiu s5,a0,-2048"},
	{0, 0x80173893, "sltiu a7,a4,-2047"},
	{0, 0x40605033, "sra zero,zero,t1"},
	{0, 0x407b54b3, "sra s1,s6,t2"},
	{0, 0x00bd52b3, "srl t0,s10,a1"},
	{0, 0x010fd033, "srl zero,t6,a6"},
	{0, 0x00cdb233, "sltu tp,s11,a2"},
	{0, 0x006ab533, "sltu a0,s5,t1"},
	{0, 0x00301073, "fscsr zero"},
	{0, 0x0030d073, "csrwi fcsr,1"},
	{0, 0x00302573, "frcsr a0"},
	{0, 0x00102573, "frflags a0"},
	{0, 0x001015f3, "fsflags a1,zero"},
	{0, 0x00215573, "fsrmi a0,2"},
	{0, 0x00127573, "csrrci a0,fflags,4"},
	{0, 0xeee59583, "lh a1,-274(a1)"},
	{0, 0x00a29a23, "sh a0,20(t0)"},
	{0, 0x00000013, "nop"},
	{0, 0x0000100f, "fence.i"},
	{0, 0x12000073, "sfence.vma"},
}

var rv32mTest = []daTest{
	{0, 0x025535b3, "mulhu a1,a0,t0"},
	{0, 0x036484b3, "mul s1,s1,s6"},
	{0, 0x03519ab3, "mulh s5,gp,s5"},
	{0, 0x039d1cb3, "mulh s9,s10,s9"},
	{0, 0x03089833, "mulh a6,a7,a6"},
}

var rv32aTest = []daTest{
	{0, 0x100526af, "lr.w a3,(a0)"},
	{0, 0x18c526af, "sc.w a3,a2,(a0)"},
	{0, 0x18e5272f, "sc.w a4,a4,(a0)"},
	{0, 0x0c55232f, "amoswap.w t1,t0,(a0)"},
	{0, 0x60b6a72f, "amoand.w a4,a1,(a3)"},
	{0, 0x00b6a72f, "amoadd.w a4,a1,(a3)"},
	{0, 0xe0b6a72f, "amomaxu.w a4,a1,(a3)"},
}

var rv32fTest = []daTest{
	{0, 0xfec42707, "flw fa4,-20(s0)"},
	{0, 0xfe842787, "flw fa5,-24(s0)"},
	{0, 0x10f777d3, "fmul.s fa5,fa4,fa5"},
	{0, 0x18f777d3, "fdiv.s fa5,fa4,fa5"},
	{0, 0x181071d3, "fdiv.s ft3,ft0,ft1"},
	{0, 0xf0000053, "fmv.w.x ft0,zero"},
	{0, 0xe0078553, "fmv.x.w a0,fa5"},
	{0, 0x20208053, "fsgnj.s ft0,ft1,ft2"},
	{0, 0xd0007053, "fcvt.s.w ft0,zero"},
	{0, 0x1000f0d3, "fmul.s ft1,ft1,ft0"},
	{0, 0xa0101553, "flt.s a0,ft0,ft1"},
	{0, 0xa0102553, "feq.s a0,ft0,ft1"},
	{0, 0xc0001553, "fcvt.w.s a0,ft0,rtz"},
	{0, 0xc000f0d3, "fcvt.w.s ra,ft1"},
	{0, 0x001071d3, "fadd.s ft3,ft0,ft1"},
	{0, 0xe0051553, "fclass.s a0,fa0"},
	{0, 0x580071d3, "fsqrt.s ft3,ft0"},
	{0, 0x101071c3, "fmadd.s ft3,ft0,ft1,ft2"},
	{0, 0x101071cf, "fnmadd.s ft3,ft0,ft1,ft2"},
	{0, 0x101071c7, "fmsub.s ft3,ft0,ft1,ft2"},
	{0, 0x281001d3, "fmin.s ft3,ft0,ft1"},
	{0, 0x281011d3, "fmax.s ft3,ft0,ft1"},
}

var rv32dTest = []daTest{
	{0, 0x0005b787, "fld fa5,0(a1)"},
	{0, 0x72a7f7c3, "fmadd.d fa5,fa5,fa0,fa4"},
	{0, 0xfef63c27, "fsd fa5,-8(a2)"},
	{0, 0x00853087, "fld ft1,8(a0)"},
	{0, 0x121071d3, "fmul.d ft3,ft0,ft1"},
	{0, 0x0a1071d3, "fsub.d ft3,ft0,ft1"},
	{0, 0x021071d3, "fadd.d ft3,ft0,ft1"},
	{0, 0x00353027, "fsd ft3,0(a0)"},
	{0, 0xe2051553, "fclass.d a0,fa0"},
	{0, 0xa2102553, "feq.d a0,ft0,ft1"},
	{0, 0xa2100553, "fle.d a0,ft0,ft1"},
	{0, 0xa2101553, "flt.d a0,ft0,ft1"},
	{0, 0x1a1071d3, "fdiv.d ft3,ft0,ft1"},
	{0, 0x5a0071d3, "fsqrt.d ft3,ft0"},
	{0, 0xd2058053, "fcvt.d.w ft0,a1"},
	{0, 0xd2158053, "fcvt.d.wu ft0,a1"},
	{0, 0xd2158053, "fcvt.d.wu ft0,a1"},
	{0, 0x40117153, "fcvt.s.d ft2,ft2"},
	{0, 0x42010153, "fcvt.d.s ft2,ft2"},
	{0, 0x2a1001d3, "fmin.d ft3,ft0,ft1"},
	{0, 0x2a1011d3, "fmax.d ft3,ft0,ft1"},
	{0, 0xc2001553, "fcvt.w.d a0,ft0,rtz"},
	{0, 0xc2101553, "fcvt.wu.d a0,ft0,rtz"},
	{0, 0xc210f0d3, "fcvt.wu.d ra,ft1"},
}

var rv32cTest = []daTest{
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
	{0x8000003e, 0xbfe5, "j 80000036"},
	{0x80000044, 0xe101, "bnez a0,80000044"},
	{0, 0x8f02, "jr t5"},
	{0, 0x52fd, "li t0,-1"},
	{0, 0x8cc9, "or s1,s1,a0"},
	{0, 0x8ca9, "xor s1,s1,a0"},
	{0, 0x8c89, "sub s1,s1,a0"},
	{0, 0x983d, "andi s0,s0,-17"},
	{0, 0x0000, "illegal"},
	{0, 0x0505, "addi a0,a0,1"},
	{0, 0xc1c8, "sw a0,4(a1)"},
	{0, 0x8431, "srai s0,s0,0xc"},
	{0, 0x8031, "srli s0,s0,0xc"},
}

var rv32cOnlyTest = []daTest{
	{0x358, 0x3d7d, "jal ra,216"},
}

var rv32fcTest = []daTest{
	{0, 0x7654, "flw fa3,44(a2)"},
	{0, 0xfedc, "fsw fa5,60(a3)"},
}

var rv32dcTest = []daTest{
	{0, 0x3210, "fld fa2,32(a2)"},
	{0, 0xba98, "fsd fa4,48(a3)"},
}

//-----------------------------------------------------------------------------
// rv64

var rv64iTest = []daTest{
	{0, 0x008a16bb, "sllw a3,s4,s0"},
	{0, 0x00e6163b, "sllw a2,a2,a4"},
	{0, 0x41978cbb, "subw s9,a5,s9"},
	{0, 0x40c989bb, "subw s3,s3,a2"},
	{0, 0x419509bb, "subw s3,a0,s9"},
	{0, 0x409786bb, "subw a3,a5,s1"},
	{0, 0x014786bb, "addw a3,a5,s4"},
	{0, 0x00d60bbb, "addw s7,a2,a3"},
	{0, 0x0087d79b, "srliw a5,a5,0x8"},
	{0, 0x0107d79b, "srliw a5,a5,0x10"},
	{0, 0x01f6d49b, "srliw s1,a3,0x1f"},
	{0, 0x0016169b, "slliw a3,a2,0x1"},
	{0, 0x0015959b, "slliw a1,a1,0x1"},
	{0, 0x4014d49b, "sraiw s1,s1,0x1"},
	{0, 0x4026571b, "sraiw a4,a2,0x2"},
	{0, 0x4027d69b, "sraiw a3,a5,0x2"},
	{0, 0x00859693, "slli a3,a1,0x8"},
	{0, 0x02059693, "slli a3,a1,0x20"},
	{0, 0x00279693, "slli a3,a5,0x2"},
	{0, 0x03279713, "slli a4,a5,0x32"},
	{0, 0x02051513, "slli a0,a0,0x20"},
	{0, 0x0094d793, "srli a5,s1,0x9"},
	{0, 0x0064d513, "srli a0,s1,0x6"},
	{0, 0x4037d493, "srai s1,a5,0x3"},
	{0, 0x40395913, "srai s2,s2,0x3"},
	{0, 0x413a59bb, "sraw s3,s4,s3"},
	{0, 0x41be5dbb, "sraw s11,t3,s11"},
	{0, 0xff80e703, "lwu a4,-8(ra)"},
	{0, 0x0002e103, "lwu sp,0(t0)"},
}

var rv64mTest = []daTest{
	{0, 0x02f777bb, "remuw a5,a4,a5"},
	{0, 0x02f757bb, "divuw a5,a4,a5"},
	{0, 0x031908bb, "mulw a7,s2,a7"},
	{0, 0x03be0dbb, "mulw s11,t3,s11"},
}

var rv64aTest = []daTest{
	{0, 0x08b6b72f, "amoswap.d a4,a1,(a3)"},
	{0, 0x00b6b72f, "amoadd.d a4,a1,(a3)"},
	{0, 0x20b6b72f, "amoxor.d a4,a1,(a3)"},
	{0, 0x60b6b72f, "amoand.d a4,a1,(a3)"},
	{0, 0x40b6b72f, "amoor.d a4,a1,(a3)"},
	{0, 0x80b6b72f, "amomin.d a4,a1,(a3)"},
	{0, 0xa0b6b72f, "amomax.d a4,a1,(a3)"},
	{0, 0xc0b6b72f, "amominu.d a4,a1,(a3)"},
	{0, 0xe0b6b72f, "amomaxu.d a4,a1,(a3)"},
}

var rv64fTest = []daTest{}

var rv64dTest = []daTest{
	{0, 0xe2018553, "fmv.x.d a0,ft3"},
	{0, 0xf2050553, "fmv.d.x fa0,a0"},
	{0, 0xd2257053, "fcvt.d.l ft0,a0"},
	{0, 0xd2357053, "fcvt.d.lu ft0,a0"},
	{0, 0xc2201553, "fcvt.l.d a0,ft0,rtz"},
	{0, 0xc2301553, "fcvt.lu.d a0,ft0,rtz"},
	{0, 0xc230f0d3, "fcvt.lu.d ra,ft1"},
}

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
	{0, 0xe822, "sd s0,16(sp)"},
	{0, 0xe04a, "sd s2,0(sp)"},
	{0, 0xec06, "sd ra,24(sp)"},
	{0, 0xe426, "sd s1,8(sp)"},
}

//-----------------------------------------------------------------------------

func testSet(module []ISAModule, tests []daTest) error {
	isa := NewISA(0)
	err := isa.Add(module)
	if err != nil {
		return err
	}
	for _, v := range tests {
		da := isa.daInstruction(v.pc, v.ins)
		if v.da != da {
			return fmt.Errorf("ins %08x \"%s\" (expected) \"%s\" (actual)", v.ins, v.da, da)
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
	rv32Tests = append(rv32Tests, rv32cOnlyTest...)
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
		{[]ISAModule{ISArv32cOnly}, rv32cOnlyTest},
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
