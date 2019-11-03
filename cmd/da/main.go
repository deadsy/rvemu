//-----------------------------------------------------------------------------
/*

RISC-V Disassembler

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/deadsy/riscv/rv"
)

//-----------------------------------------------------------------------------

var symtab = rv.SymbolTable{
	0x00000000: "nybble",
	0x00000048: ".L2",
	0x00000054: ".L3",
	0x00000064: "hex8",
	0x000000f4: "hex16",
	0x0000016c: "hex32",
	0x000001e0: "itoa",
	0x00000218: ".L11",
	0x00000220: ".L13",
	0x000002ac: ".L14",
	0x000002cc: ".L16",
	0x0000032c: ".L15",
	0x00000350: "main",
}

var code = []uint32{
	0xfe010113, // addi	sp,sp,-32
	0x00812e23, // sw	s0,28(sp)
	0x02010413, // addi	s0,sp,32
	0x00050793, // mv	a5,a0
	0xfef407a3, // sb	a5,-17(s0)
	0xfef44783, // lbu	a5,-17(s0)
	0x00f7f793, // andi	a5,a5,15
	0xfef407a3, // sb	a5,-17(s0)
	0xfef44703, // lbu	a4,-17(s0)
	0x00900793, // li	a5,9
	0x02e7f063, // bgeu	a5,a4,48 <.L2>
	0xfef44703, // lbu	a4,-17(s0)
	0x00f00793, // li	a5,15
	0x00e7ea63, // bltu	a5,a4,48 <.L2>
	0xfef44783, // lbu	a5,-17(s0)
	0x05778793, // addi	a5,a5,87
	0x0ff7f793, // andi	a5,a5,255
	0x0100006f, // j	54 <.L3>
	0xfef44783, // lbu	a5,-17(s0)
	0x03078793, // addi	a5,a5,48
	0x0ff7f793, // andi	a5,a5,255
	0x00078513, // mv	a0,a5
	0x01c12403, // lw	s0,28(sp)
	0x02010113, // addi	sp,sp,32
	0x00008067, // ret
	0xfe010113, // addi	sp,sp,-32
	0x00112e23, // sw	ra,28(sp)
	0x00812c23, // sw	s0,24(sp)
	0x00912a23, // sw	s1,20(sp)
	0x02010413, // addi	s0,sp,32
	0xfea42623, // sw	a0,-20(s0)
	0x00058793, // mv	a5,a1
	0xfef405a3, // sb	a5,-21(s0)
	0xfeb44783, // lbu	a5,-21(s0)
	0x0047d793, // srli	a5,a5,0x4
	0x0ff7f793, // andi	a5,a5,255
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 94 <hex8+0x30>
	0x00050793, // mv	a5,a0
	0x00078713, // mv	a4,a5
	0xfec42783, // lw	a5,-20(s0)
	0x00e78023, // sb	a4,0(a5)
	0xfec42783, // lw	a5,-20(s0)
	0x00178493, // addi	s1,a5,1
	0xfeb44783, // lbu	a5,-21(s0)
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # bc <hex8+0x58>
	0x00050793, // mv	a5,a0
	0x00f48023, // sb	a5,0(s1)
	0xfec42783, // lw	a5,-20(s0)
	0x00278793, // addi	a5,a5,2
	0x00078023, // sb	zero,0(a5)
	0xfec42783, // lw	a5,-20(s0)
	0x00078513, // mv	a0,a5
	0x01c12083, // lw	ra,28(sp)
	0x01812403, // lw	s0,24(sp)
	0x01412483, // lw	s1,20(sp)
	0x02010113, // addi	sp,sp,32
	0x00008067, // ret
	0xfe010113, // addi	sp,sp,-32
	0x00112e23, // sw	ra,28(sp)
	0x00812c23, // sw	s0,24(sp)
	0x02010413, // addi	s0,sp,32
	0xfea42623, // sw	a0,-20(s0)
	0x00058793, // mv	a5,a1
	0xfef41523, // sh	a5,-22(s0)
	0xfea45783, // lhu	a5,-22(s0)
	0x0087d793, // srli	a5,a5,0x8
	0x01079793, // slli	a5,a5,0x10
	0x0107d793, // srli	a5,a5,0x10
	0x0ff7f793, // andi	a5,a5,255
	0x00078593, // mv	a1,a5
	0xfec42503, // lw	a0,-20(s0)
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 12c <hex16+0x38>
	0xfec42783, // lw	a5,-20(s0)
	0x00278793, // addi	a5,a5,2
	0xfea45703, // lhu	a4,-22(s0)
	0x0ff77713, // andi	a4,a4,255
	0x00070593, // mv	a1,a4
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 14c <hex16+0x58>
	0xfec42783, // lw	a5,-20(s0)
	0x00078513, // mv	a0,a5
	0x01c12083, // lw	ra,28(sp)
	0x01812403, // lw	s0,24(sp)
	0x02010113, // addi	sp,sp,32
	0x00008067, // ret
	0xfe010113, // addi	sp,sp,-32
	0x00112e23, // sw	ra,28(sp)
	0x00812c23, // sw	s0,24(sp)
	0x02010413, // addi	s0,sp,32
	0xfea42623, // sw	a0,-20(s0)
	0xfeb42423, // sw	a1,-24(s0)
	0xfe842783, // lw	a5,-24(s0)
	0x0107d793, // srli	a5,a5,0x10
	0x01079793, // slli	a5,a5,0x10
	0x0107d793, // srli	a5,a5,0x10
	0x00078593, // mv	a1,a5
	0xfec42503, // lw	a0,-20(s0)
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 19c <hex32+0x30>
	0xfec42783, // lw	a5,-20(s0)
	0x00478793, // addi	a5,a5,4
	0xfe842703, // lw	a4,-24(s0)
	0x01071713, // slli	a4,a4,0x10
	0x01075713, // srli	a4,a4,0x10
	0x00070593, // mv	a1,a4
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 1c0 <hex32+0x54>
	0xfec42783, // lw	a5,-20(s0)
	0x00078513, // mv	a0,a5
	0x01c12083, // lw	ra,28(sp)
	0x01812403, // lw	s0,24(sp)
	0x02010113, // addi	sp,sp,32
	0x00008067, // ret
	0xfd010113, // addi	sp,sp,-48
	0x02112623, // sw	ra,44(sp)
	0x02812423, // sw	s0,40(sp)
	0x03010413, // addi	s0,sp,48
	0xfca42e23, // sw	a0,-36(s0)
	0xfcb42c23, // sw	a1,-40(s0)
	0xfe042423, // sw	zero,-24(s0)
	0xfe042223, // sw	zero,-28(s0)
	0xfd842783, // lw	a5,-40(s0)
	0x0007da63, // bgez	a5,218 <.L11>
	0xfd842783, // lw	a5,-40(s0)
	0x40f007b3, // neg	a5,a5
	0xfef42623, // sw	a5,-20(s0)
	0x00c0006f, // j	220 <.L13>
	0xfd842783, // lw	a5,-40(s0)
	0xfef42623, // sw	a5,-20(s0)
	0xfec42783, // lw	a5,-20(s0)
	0x00a00593, // li	a1,10
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 22c <.L13+0xc>
	0x00050793, // mv	a5,a0
	0x0ff7f713, // andi	a4,a5,255
	0xfe842783, // lw	a5,-24(s0)
	0x00178693, // addi	a3,a5,1
	0xfed42423, // sw	a3,-24(s0)
	0x00078693, // mv	a3,a5
	0xfdc42783, // lw	a5,-36(s0)
	0x00d787b3, // add	a5,a5,a3
	0x03070713, // addi	a4,a4,48
	0x0ff77713, // andi	a4,a4,255
	0x00e78023, // sb	a4,0(a5)
	0xfec42783, // lw	a5,-20(s0)
	0x00a00593, // li	a1,10
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 26c <.L13+0x4c>
	0x00050793, // mv	a5,a0
	0xfef42623, // sw	a5,-20(s0)
	0xfec42783, // lw	a5,-20(s0)
	0xfa0790e3, // bnez	a5,220 <.L13>
	0xfd842783, // lw	a5,-40(s0)
	0x0207d263, // bgez	a5,2ac <.L14>
	0xfe842783, // lw	a5,-24(s0)
	0x00178713, // addi	a4,a5,1
	0xfee42423, // sw	a4,-24(s0)
	0x00078713, // mv	a4,a5
	0xfdc42783, // lw	a5,-36(s0)
	0x00e787b3, // add	a5,a5,a4
	0x02d00713, // li	a4,45
	0x00e78023, // sb	a4,0(a5)
	0xfe842783, // lw	a5,-24(s0)
	0xfdc42703, // lw	a4,-36(s0)
	0x00f707b3, // add	a5,a4,a5
	0x00078023, // sb	zero,0(a5)
	0xfe842783, // lw	a5,-24(s0)
	0xfff78793, // addi	a5,a5,-1
	0xfef42423, // sw	a5,-24(s0)
	0x0640006f, // j	32c <.L15>
	0xfe442783, // lw	a5,-28(s0)
	0xfdc42703, // lw	a4,-36(s0)
	0x00f707b3, // add	a5,a4,a5
	0x0007c783, // lbu	a5,0(a5)
	0xfef401a3, // sb	a5,-29(s0)
	0xfe842783, // lw	a5,-24(s0)
	0xfdc42703, // lw	a4,-36(s0)
	0x00f70733, // add	a4,a4,a5
	0xfe442783, // lw	a5,-28(s0)
	0x00178693, // addi	a3,a5,1
	0xfed42223, // sw	a3,-28(s0)
	0x00078693, // mv	a3,a5
	0xfdc42783, // lw	a5,-36(s0)
	0x00d787b3, // add	a5,a5,a3
	0x00074703, // lbu	a4,0(a4)
	0x00e78023, // sb	a4,0(a5)
	0xfe842783, // lw	a5,-24(s0)
	0xfff78713, // addi	a4,a5,-1
	0xfee42423, // sw	a4,-24(s0)
	0x00078713, // mv	a4,a5
	0xfdc42783, // lw	a5,-36(s0)
	0x00e787b3, // add	a5,a5,a4
	0xfe344703, // lbu	a4,-29(s0)
	0x00e78023, // sb	a4,0(a5)
	0xfe442703, // lw	a4,-28(s0)
	0xfe842783, // lw	a5,-24(s0)
	0xf8f74ce3, // blt	a4,a5,2cc <.L16>
	0xfdc42783, // lw	a5,-36(s0)
	0x00078513, // mv	a0,a5
	0x02c12083, // lw	ra,44(sp)
	0x02812403, // lw	s0,40(sp)
	0x03010113, // addi	sp,sp,48
	0x00008067, // ret
	0xfd010113, // addi	sp,sp,-48
	0x02112623, // sw	ra,44(sp)
	0x02812423, // sw	s0,40(sp)
	0x03010413, // addi	s0,sp,48
	0xfd040793, // addi	a5,s0,-48
	0x00000593, // li	a1,0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 36c <main+0x1c>
	0x00050793, // mv	a5,a0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 37c <main+0x2c>
	0xfd040793, // addi	a5,s0,-48
	0x4d200593, // li	a1,1234
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 390 <main+0x40>
	0x00050793, // mv	a5,a0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 3a0 <main+0x50>
	0xfd040793, // addi	a5,s0,-48
	0xb2e00593, // li	a1,-1234
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 3b4 <main+0x64>
	0x00050793, // mv	a5,a0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 3c4 <main+0x74>
	0xfd040713, // addi	a4,s0,-48
	0x800007b7, // lui	a5,0x80000
	0xfff7c593, // not	a1,a5
	0x00070513, // mv	a0,a4
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 3dc <main+0x8c>
	0x00050793, // mv	a5,a0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 3ec <main+0x9c>
	0xfd040793, // addi	a5,s0,-48
	0x800005b7, // lui	a1,0x80000
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 400 <main+0xb0>
	0x00050793, // mv	a5,a0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 410 <main+0xc0>
	0xfd040793, // addi	a5,s0,-48
	0x0ab00593, // li	a1,171
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 424 <main+0xd4>
	0x00050793, // mv	a5,a0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 434 <main+0xe4>
	0xfd040713, // addi	a4,s0,-48
	0x0000b7b7, // lui	a5,0xb
	0xbcd78593, // addi	a1,a5,-1075 # abcd <main+0xa87d>
	0x00070513, // mv	a0,a4
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 44c <main+0xfc>
	0x00050793, // mv	a5,a0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 45c <main+0x10c>
	0xfd040713, // addi	a4,s0,-48
	0xdeadc7b7, // lui	a5,0xdeadc
	0xeef78593, // addi	a1,a5,-273 # deadbeef <main+0xdeadbb9f>
	0x00070513, // mv	a0,a4
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 474 <main+0x124>
	0x00050793, // mv	a5,a0
	0x00078513, // mv	a0,a5
	0x00000097, // auipc	ra,0x0
	0x000080e7, // jalr	ra # 484 <main+0x134>
	0x00000793, // li	a5,0
	0x00078513, // mv	a0,a5
	0x02c12083, // lw	ra,44(sp)
	0x02812403, // lw	s0,40(sp)
	0x03010113, // addi	sp,sp,48
	0x00008067, // ret
}

//-----------------------------------------------------------------------------

type memory struct {
	base uint32
	mem  []uint32
}

func (m *memory) Read32(adr uint32) uint32 {
	adr -= m.base
	if adr&3 != 0 {
		panic(fmt.Sprintf("mis-aligned 32 bit read @ %08x", adr))
	}
	return m.mem[adr>>2]
}

func (m *memory) Write32(adr uint32, val uint32) {
	// nop
}

//-----------------------------------------------------------------------------

func main() {

	// create the ISA
	isa := rv.NewISA("rv32g")
	err := isa.Add(rv.ISArv32i, rv.ISArv32m, rv.ISArv32a, rv.ISArv32f, rv.ISArv32d)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	err = isa.GenDecoders()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// create the memory
	adr := uint32(0)
	m := &memory{
		base: adr,
		mem:  code,
	}

	// create the CPU
	cpu, err := rv.NewRV(rv.VariantRV32, isa, m)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// Disassemble
	size := len(code) * int(unsafe.Sizeof(code[0]))
	for size > 0 {
		da := cpu.Disassemble(adr, symtab)
		fmt.Printf("%s\n", da.String())
		size -= da.N
		adr += uint32(da.N)
	}

	os.Exit(0)
}

//-----------------------------------------------------------------------------
