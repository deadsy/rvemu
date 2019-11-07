//-----------------------------------------------------------------------------
/*

RISC-V ISA Definition

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

// daFunc is the disassembler function for a  16/32-bit instructions.
type daFunc func(name string, pc uint32, ins uint) string

type insDefn struct {
	defn string     // instruction definition string (from the standard)
	dt   decodeType // decode type
	da   daFunc
}

// ISAModule is a set of RISC-V instructions as described in the specification.
type ISAModule struct {
	name string    // name of module
	ilen int       // instruction length
	defn []insDefn // instruction definitions
}

//-----------------------------------------------------------------------------
// RV32 instructions

// ISArv32i Integer
var ISArv32i = ISAModule{
	name: "rv32i",
	ilen: 32,
	defn: []insDefn{
		{"imm[31:12] rd 0110111 LUI", decodeTypeU, daTypeUa},
		{"imm[31:12] rd 0010111 AUIPC", decodeTypeU, daTypeUa},
		{"imm[20|10:1|11|19:12] rd 1101111 JAL", decodeTypeJ, daTypeJa},
		{"imm[11:0] rs1 000 rd 1100111 JALR", decodeTypeI, daTypeIe},
		{"imm[12|10:5] rs2 rs1 000 imm[4:1|11] 1100011 BEQ", decodeTypeB, daTypeBa},
		{"imm[12|10:5] rs2 rs1 001 imm[4:1|11] 1100011 BNE", decodeTypeB, daTypeBa},
		{"imm[12|10:5] rs2 rs1 100 imm[4:1|11] 1100011 BLT", decodeTypeB, daTypeBa},
		{"imm[12|10:5] rs2 rs1 101 imm[4:1|11] 1100011 BGE", decodeTypeB, daTypeBa},
		{"imm[12|10:5] rs2 rs1 110 imm[4:1|11] 1100011 BLTU", decodeTypeB, daTypeBa},
		{"imm[12|10:5] rs2 rs1 111 imm[4:1|11] 1100011 BGEU", decodeTypeB, daTypeBa},
		{"imm[11:0] rs1 000 rd 0000011 LB", decodeTypeI, daTypeIc},
		{"imm[11:0] rs1 001 rd 0000011 LH", decodeTypeI, daTypeIc},
		{"imm[11:0] rs1 010 rd 0000011 LW", decodeTypeI, daTypeIc},
		{"imm[11:0] rs1 100 rd 0000011 LBU", decodeTypeI, daTypeIc},
		{"imm[11:0] rs1 101 rd 0000011 LHU", decodeTypeI, daTypeIc},
		{"imm[11:5] rs2 rs1 000 imm[4:0] 0100011 SB", decodeTypeS, daTypeSa},
		{"imm[11:5] rs2 rs1 001 imm[4:0] 0100011 SH", decodeTypeS, daTypeSa},
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100011 SW", decodeTypeS, daTypeSa},
		{"imm[11:0] rs1 000 rd 0010011 ADDI", decodeTypeI, daTypeIb},
		{"imm[11:0] rs1 010 rd 0010011 SLTI", decodeTypeI, daTypeIa},
		{"imm[11:0] rs1 011 rd 0010011 SLTIU", decodeTypeI, daTypeIa},
		{"imm[11:0] rs1 100 rd 0010011 XORI", decodeTypeI, daTypeIf},
		{"imm[11:0] rs1 110 rd 0010011 ORI", decodeTypeI, daTypeIa},
		{"imm[11:0] rs1 111 rd 0010011 ANDI", decodeTypeI, daTypeIa},
		{"0000000 shamt5 rs1 001 rd 0010011 SLLI", decodeTypeI, daTypeId},
		{"0000000 shamt5 rs1 101 rd 0010011 SRLI", decodeTypeI, daTypeId},
		{"0100000 shamt5 rs1 101 rd 0010011 SRAI", decodeTypeI, daTypeId},
		{"0000000 rs2 rs1 000 rd 0110011 ADD", decodeTypeR, daTypeRa},
		{"0100000 rs2 rs1 000 rd 0110011 SUB", decodeTypeR, daTypeRa},
		{"0000000 rs2 rs1 001 rd 0110011 SLL", decodeTypeR, daTypeRa},
		{"0000000 rs2 rs1 010 rd 0110011 SLT", decodeTypeR, daTypeRa},
		{"0000000 rs2 rs1 011 rd 0110011 SLTU", decodeTypeR, daTypeRa},
		{"0000000 rs2 rs1 100 rd 0110011 XOR", decodeTypeR, daTypeRa},
		{"0000000 rs2 rs1 101 rd 0110011 SRL", decodeTypeR, daTypeRa},
		{"0100000 rs2 rs1 101 rd 0110011 SRA", decodeTypeR, daTypeRa},
		{"0000000 rs2 rs1 110 rd 0110011 OR", decodeTypeR, daTypeRa},
		{"0000000 rs2 rs1 111 rd 0110011 AND", decodeTypeR, daTypeRa},
		{"0000 pred succ 00000 000 00000 0001111 FENCE", decodeTypeI, daNone},
		{"0000 0000 0000 00000 001 00000 0001111 FENCE.I", decodeTypeI, daNone},
		{"000000000000 00000 000 00000 1110011 ECALL", decodeTypeI, daNone},
		{"000000000001 00000 000 00000 1110011 EBREAK", decodeTypeI, daNone},
		{"csr rs1 001 rd 1110011 CSRRW", decodeTypeI, daNone},
		{"csr rs1 010 rd 1110011 CSRRS", decodeTypeI, daNone},
		{"csr rs1 011 rd 1110011 CSRRC", decodeTypeI, daNone},
		{"csr zimm 101 rd 1110011 CSRRWI", decodeTypeI, daNone},
		{"csr zimm 110 rd 1110011 CSRRSI", decodeTypeI, daNone},
		{"csr zimm 111 rd 1110011 CSRRCI", decodeTypeI, daNone},
	},
}

// ISArv32m Integer Multiplication and Division
var ISArv32m = ISAModule{
	name: "rv32m",
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0110011 MUL", decodeTypeR, daTypeRa},
		{"0000001 rs2 rs1 001 rd 0110011 MULH", decodeTypeR, daTypeRa},
		{"0000001 rs2 rs1 010 rd 0110011 MULHSU", decodeTypeR, daTypeRa},
		{"0000001 rs2 rs1 011 rd 0110011 MULHU", decodeTypeR, daTypeRa},
		{"0000001 rs2 rs1 100 rd 0110011 DIV", decodeTypeR, daTypeRa},
		{"0000001 rs2 rs1 101 rd 0110011 DIVU", decodeTypeR, daTypeRa},
		{"0000001 rs2 rs1 110 rd 0110011 REM", decodeTypeR, daTypeRa},
		{"0000001 rs2 rs1 111 rd 0110011 REMU", decodeTypeR, daTypeRa},
	},
}

// ISArv32a Atomics
var ISArv32a = ISAModule{
	name: "rv32a",
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 010 rd 0101111 LR.W", decodeTypeR, daTypeRb},
		{"00011 aq rl rs2 rs1 010 rd 0101111 SC.W", decodeTypeR, daTypeRb},
		{"00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W", decodeTypeR, daTypeRb},
		{"00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W", decodeTypeR, daTypeRb},
		{"00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W", decodeTypeR, daTypeRb},
		{"01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W", decodeTypeR, daTypeRb},
		{"01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W", decodeTypeR, daTypeRb},
		{"10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W", decodeTypeR, daTypeRb},
		{"10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W", decodeTypeR, daTypeRb},
		{"11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W", decodeTypeR, daTypeRb},
		{"11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W", decodeTypeR, daTypeRb},
	},
}

// ISArv32f Single-Precision Floating-Point
var ISArv32f = ISAModule{
	name: "rv32f",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 010 rd 0000111 FLW", decodeTypeI, daTypeIa},
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100111 FSW", decodeTypeS, daTypeSb},
		{"rs3 00 rs2 rs1 rm rd 1000011 FMADD.S", decodeTypeR4, daNone},
		{"rs3 00 rs2 rs1 rm rd 1000111 FMSUB.S", decodeTypeR4, daNone},
		{"rs3 00 rs2 rs1 rm rd 1001011 FNMSUB.S", decodeTypeR4, daNone},
		{"rs3 00 rs2 rs1 rm rd 1001111 FNMADD.S", decodeTypeR4, daNone},
		{"0000000 rs2 rs1 rm rd 1010011 FADD.S", decodeTypeR, daTypeRa},
		{"0000100 rs2 rs1 rm rd 1010011 FSUB.S", decodeTypeR, daTypeRa},
		{"0001000 rs2 rs1 rm rd 1010011 FMUL.S", decodeTypeR, daTypeRa},
		{"0001100 rs2 rs1 rm rd 1010011 FDIV.S", decodeTypeR, daTypeRa},
		{"0101100 00000 rs1 rm rd 1010011 FSQRT.S", decodeTypeR, daNone},
		{"0010000 rs2 rs1 000 rd 1010011 FSGNJ.S", decodeTypeR, daTypeRa},
		{"0010000 rs2 rs1 001 rd 1010011 FSGNJN.S", decodeTypeR, daTypeRa},
		{"0010000 rs2 rs1 010 rd 1010011 FSGNJX.S", decodeTypeR, daTypeRa},
		{"0010100 rs2 rs1 000 rd 1010011 FMIN.S", decodeTypeR, daTypeRa},
		{"0010100 rs2 rs1 001 rd 1010011 FMAX.S", decodeTypeR, daTypeRa},
		{"1100000 00000 rs1 rm rd 1010011 FCVT.W.S", decodeTypeR, daNone},
		{"1100000 00001 rs1 rm rd 1010011 FCVT.WU.S", decodeTypeR, daNone},
		{"1110000 00000 rs1 000 rd 1010011 FMV.X.W", decodeTypeR, daNone},
		{"1010000 rs2 rs1 010 rd 1010011 FEQ.S", decodeTypeR, daTypeRa},
		{"1010000 rs2 rs1 001 rd 1010011 FLT.S", decodeTypeR, daTypeRa},
		{"1010000 rs2 rs1 000 rd 1010011 FLE.S", decodeTypeR, daTypeRa},
		{"1110000 00000 rs1 001 rd 1010011 FCLASS.S", decodeTypeR, daNone},
		{"1101000 00000 rs1 rm rd 1010011 FCVT.S.W", decodeTypeR, daNone},
		{"1101000 00001 rs1 rm rd 1010011 FCVT.S.WU", decodeTypeR, daNone},
		{"1111000 00000 rs1 000 rd 1010011 FMV.W.X", decodeTypeR, daNone},
	},
}

// ISArv32d Double-Precision Floating-Point
var ISArv32d = ISAModule{
	name: "rv32d",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 011 rd 0000111 FLD", decodeTypeI, daTypeIg},
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100111 FSD", decodeTypeS, daTypeSb},
		{"rs3 01 rs2 rs1 rm rd 1000011 FMADD.D", decodeTypeR4, daTypeR4a},
		{"rs3 01 rs2 rs1 rm rd 1000111 FMSUB.D", decodeTypeR4, daTypeR4a},
		{"rs3 01 rs2 rs1 rm rd 1001011 FNMSUB.D", decodeTypeR4, daTypeR4a},
		{"rs3 01 rs2 rs1 rm rd 1001111 FNMADD.D", decodeTypeR4, daTypeR4a},
		{"0000001 rs2 rs1 rm rd 1010011 FADD.D", decodeTypeR, daNone},
		{"0000101 rs2 rs1 rm rd 1010011 FSUB.D", decodeTypeR, daNone},
		{"0001001 rs2 rs1 rm rd 1010011 FMUL.D", decodeTypeR, daNone},
		{"0001101 rs2 rs1 rm rd 1010011 FDIV.D", decodeTypeR, daNone},
		{"0101101 00000 rs1 rm rd 1010011 FSQRT.D", decodeTypeR, daNone},
		{"0010001 rs2 rs1 000 rd 1010011 FSGNJ.D", decodeTypeR, daNone},
		{"0010001 rs2 rs1 001 rd 1010011 FSGNJN.D", decodeTypeR, daNone},
		{"0010001 rs2 rs1 010 rd 1010011 FSGNJX.D", decodeTypeR, daNone},
		{"0010101 rs2 rs1 000 rd 1010011 FMIN.D", decodeTypeR, daNone},
		{"0010101 rs2 rs1 001 rd 1010011 FMAX.D", decodeTypeR, daNone},
		{"0100000 00001 rs1 rm rd 1010011 FCVT.S.D", decodeTypeR, daNone},
		{"0100001 00000 rs1 rm rd 1010011 FCVT.D.S", decodeTypeR, daNone},
		{"1010001 rs2 rs1 010 rd 1010011 FEQ.D", decodeTypeR, daNone},
		{"1010001 rs2 rs1 001 rd 1010011 FLT.D", decodeTypeR, daNone},
		{"1010001 rs2 rs1 000 rd 1010011 FLE.D", decodeTypeR, daNone},
		{"1110001 00000 rs1 001 rd 1010011 FCLASS.D", decodeTypeR, daNone},
		{"1100001 00000 rs1 rm rd 1010011 FCVT.W.D", decodeTypeR, daNone},
		{"1100001 00001 rs1 rm rd 1010011 FCVT.WU.D", decodeTypeR, daNone},
		{"1101001 00000 rs1 rm rd 1010011 FCVT.D.W", decodeTypeR, daNone},
		{"1101001 00001 rs1 rm rd 1010011 FCVT.D.WU", decodeTypeR, daNone},
	},
}

// ISArv32c Compressed
var ISArv32c = ISAModule{
	name: "rv32c",
	ilen: 16,
	defn: []insDefn{
		// Quadrant 0
		{"000 00000000 000 00 ILLEGAL", decodeTypeCIW, daTypeCIWa},
		{"000 nzuimm[5:4|9:6|2|3] rd0 00 C.ADDI4SPN", decodeTypeCIW, daNone},
		{"001 uimm[5:3] rs10 uimm[7:6] rd0 00 C.FLD", decodeTypeCL, daNone},
		{"010 uimm[5:3] rs10 uimm[2|6] rd0 00 C.LW", decodeTypeCL, daNone},
		{"011 uimm[5:3] rs10 uimm[2|6] rd0 00 C.FLW", decodeTypeCL, daNone},
		{"101 uimm[5:3] rs10 uimm[7:6] rs20 00 C.FSD", decodeTypeCS, daNone},
		{"110 uimm[5:3] rs10 uimm[2|6] rs20 00 C.SW", decodeTypeCS, daNone},
		{"111 uimm[5:3] rs10 uimm[2|6] rs20 00 C.FSW", decodeTypeCS, daNone},
		// Quadrant 1
		{"000 nzimm[5] 00000 nzimm[4:0] 01 C.NOP", decodeTypeCI, daNone},
		{"000 nzimm[5] rs1/rd!=0 nzimm[4:0] 01 C.ADDI", decodeTypeCI, daNone},
		{"001 imm[11|4|9:8|10|6|7|3:1|5] 01 C.JAL", decodeTypeCJ, daNone},
		{"010 imm[5] rd!=0 imm[4:0] 01 C.LI", decodeTypeCI, daTypeCIa},
		{"011 nzimm[9] 00010 nzimm[4|6|8:7|5] 01 C.ADDI16SP", decodeTypeCI, daNone},
		{"011 nzimm[17] rd!={0,2} nzimm[16:12] 01 C.LUI", decodeTypeCI, daNone},
		{"100 nzuimm[5] 00 rs10/rd0 nzuimm[4:0] 01 C.SRLI", decodeTypeCI, daNone},
		{"100 nzuimm[5] 01 rs10/rd0 nzuimm[4:0] 01 C.SRAI", decodeTypeCI, daNone},
		{"100 imm[5] 10 rs10/rd0 imm[4:0] 01 C.ANDI", decodeTypeCI, daNone},
		{"100 0 11 rs10/rd0 00 rs20 01 C.SUB", decodeTypeCR, daNone},
		{"100 0 11 rs10/rd0 01 rs20 01 C.XOR", decodeTypeCR, daNone},
		{"100 0 11 rs10/rd0 10 rs20 01 C.OR", decodeTypeCR, daNone},
		{"100 0 11 rs10/rd0 11 rs20 01 C.AND", decodeTypeCR, daNone},
		{"101 imm[11|4|9:8|10|6|7|3:1|5] 01 C.J", decodeTypeCJ, daNone},
		{"110 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BEQZ", decodeTypeCB, daNone},
		{"111 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BNEZ", decodeTypeCB, daNone},
		// Quadrant 2
		{"000 nzuimm[5] rs1/rd!=0 nzuimm[4:0] 10 C.SLLI", decodeTypeCI, daNone},
		{"000 0 rs1/rd!=0 00000 10 C.SLLI64", decodeTypeCI, daNone},
		{"001 uimm[5] rd uimm[4:3|8:6] 10 C.FLDSP", decodeTypeCSS, daNone},
		{"010 uimm[5] rd!=0 uimm[4:2|7:6] 10 C.LWSP", decodeTypeCSS, daNone},
		{"011 uimm[5] rd uimm[4:2|7:6] 10 C.FLWSP", decodeTypeCSS, daNone},
		{"100 0 rs1!=0 00000 10 C.JR", decodeTypeCJ, daTypeCJa},
		{"100 0 rd!=0 rs2!=0 10 C.MV", decodeTypeCR, daNone},
		{"100 1 00000 00000 10 C.EBREAK", decodeTypeCI, daNone},
		{"100 1 rs1!=0 00000 10 C.JALR", decodeTypeCJ, daNone},
		{"100 1 rs1/rd!=0 rs2!=0 10 C.ADD", decodeTypeCR, daNone},
		{"101 uimm[5:3|8:6] rs2 10 C.FSDSP", decodeTypeCSS, daNone},
		{"110 uimm[5:2|7:6] rs2 10 C.SWSP", decodeTypeCSS, daNone},
		{"111 uimm[5:2|7:6] rs2 10 C.FSWSP", decodeTypeCSS, daNone},
	},
}

//-----------------------------------------------------------------------------
// RV64 instructions (+ RV32)

// ISArv64i Integer
var ISArv64i = ISAModule{
	name: "rv64i",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 110 rd 0000011 LWU", decodeTypeI, daTypeIa},
		{"imm[11:0] rs1 011 rd 0000011 LD", decodeTypeI, daTypeIa},
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD", decodeTypeS, daTypeSa},
		{"000000 shamt6 rs1 001 rd 0010011 SLLI", decodeTypeI, daNone},
		{"000000 shamt6 rs1 101 rd 0010011 SRLI", decodeTypeI, daNone},
		{"010000 shamt6 rs1 101 rd 0010011 SRAI", decodeTypeI, daNone},
		{"imm[11:0] rs1 000 rd 0011011 ADDIW", decodeTypeI, daTypeIa},
		{"0000000 shamt5 rs1 001 rd 0011011 SLLIW", decodeTypeI, daNone},
		{"0000000 shamt5 rs1 101 rd 0011011 SRLIW", decodeTypeI, daNone},
		{"0100000 shamt5 rs1 101 rd 0011011 SRAIW", decodeTypeI, daNone},
		{"0000000 rs2 rs1 000 rd 0111011 ADDW", decodeTypeR, daNone},
		{"0100000 rs2 rs1 000 rd 0111011 SUBW", decodeTypeR, daNone},
		{"0000000 rs2 rs1 001 rd 0111011 SLLW", decodeTypeR, daNone},
		{"0000000 rs2 rs1 101 rd 0111011 SRLW", decodeTypeR, daNone},
		{"0100000 rs2 rs1 101 rd 0111011 SRAW", decodeTypeR, daNone},
	},
}

// ISArv64m Integer Multiplication and Division
var ISArv64m = ISAModule{
	name: "rv64m",
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0111011 MULW", decodeTypeR, daNone},
		{"0000001 rs2 rs1 100 rd 0111011 DIVW", decodeTypeR, daNone},
		{"0000001 rs2 rs1 101 rd 0111011 DIVUW", decodeTypeR, daNone},
		{"0000001 rs2 rs1 110 rd 0111011 REMW", decodeTypeR, daNone},
		{"0000001 rs2 rs1 111 rd 0111011 REMUW", decodeTypeR, daNone},
	},
}

// ISArv64a Atomics
var ISArv64a = ISAModule{
	name: "rv64a",
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 011 rd 0101111 LR.D", decodeTypeR, daNone},
		{"00011 aq rl rs2 rs1 011 rd 0101111 SC.D", decodeTypeR, daNone},
		{"00001 aq rl rs2 rs1 011 rd 0101111 AMOSWAP.D", decodeTypeR, daNone},
		{"00000 aq rl rs2 rs1 011 rd 0101111 AMOADD.D", decodeTypeR, daNone},
		{"00100 aq rl rs2 rs1 011 rd 0101111 AMOXOR.D", decodeTypeR, daNone},
		{"01100 aq rl rs2 rs1 011 rd 0101111 AMOAND.D", decodeTypeR, daNone},
		{"01000 aq rl rs2 rs1 011 rd 0101111 AMOOR.D", decodeTypeR, daNone},
		{"10000 aq rl rs2 rs1 011 rd 0101111 AMOMIN.D", decodeTypeR, daNone},
		{"10100 aq rl rs2 rs1 011 rd 0101111 AMOMAX.D", decodeTypeR, daNone},
		{"11000 aq rl rs2 rs1 011 rd 0101111 AMOMINU.D", decodeTypeR, daNone},
		{"11100 aq rl rs2 rs1 011 rd 0101111 AMOMAXU.D", decodeTypeR, daNone},
	},
}

// ISArv64f Single-Precision Floating-Point
var ISArv64f = ISAModule{
	name: "rv64f",
	ilen: 32,
	defn: []insDefn{
		{"1100000 00010 rs1 rm rd 1010011 FCVT.L.S", decodeTypeR, daNone},
		{"1100000 00011 rs1 rm rd 1010011 FCVT.LU.S", decodeTypeR, daNone},
		{"1101000 00010 rs1 rm rd 1010011 FCVT.S.L", decodeTypeR, daNone},
		{"1101000 00011 rs1 rm rd 1010011 FCVT.S.LU", decodeTypeR, daNone},
	},
}

// ISArv64d Double-Precision Floating-Point
var ISArv64d = ISAModule{
	name: "rv64d",
	ilen: 32,
	defn: []insDefn{
		{"1100001 00010 rs1 rm rd 1010011 FCVT.L.D", decodeTypeR, daNone},
		{"1100001 00011 rs1 rm rd 1010011 FCVT.LU.D", decodeTypeR, daNone},
		{"1110001 00000 rs1 000 rd 1010011 FMV.X.D", decodeTypeR, daNone},
		{"1101001 00010 rs1 rm rd 1010011 FCVT.D.L", decodeTypeR, daNone},
		{"1101001 00011 rs1 rm rd 1010011 FCVT.D.LU", decodeTypeR, daNone},
		{"1111001 00000 rs1 000 rd 1010011 FMV.D.X", decodeTypeR, daNone},
	},
}

//-----------------------------------------------------------------------------

// insInfo is instruction information.
type insInfo struct {
	name      string // instruction mneumonic
	val, mask uint   // value and mask of fixed bits in the instruction
	da        daFunc // disassembler
}

// ISA is an instruction set
type ISA struct {
	name        string     // the name of the ISA
	instruction []*insInfo // the set of instruction in the ISA
}

// NewISA creates an empty instruction set.
func NewISA(name string) *ISA {
	return &ISA{
		name:        name,
		instruction: make([]*insInfo, 0),
	}
}

// Add a ISA sub-module to the ISA.
func (isa *ISA) Add(module ...ISAModule) error {
	for i := range module {
		for _, id := range module[i].defn {
			ii, err := parseDefn(&id, module[i].ilen)
			if err != nil {
				return err
			}
			isa.instruction = append(isa.instruction, ii)
		}
	}
	return nil
}

// lookup returns the instruction information for an instruction.
func (isa *ISA) lookup(ins uint) *insInfo {
	for _, ii := range isa.instruction {
		if ins&ii.mask == ii.val {
			return ii
		}
	}
	return nil
}

//-----------------------------------------------------------------------------
