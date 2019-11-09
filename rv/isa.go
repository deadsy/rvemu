//-----------------------------------------------------------------------------
/*

RISC-V ISA Definition

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

// daFunc is the disassembler function
type daFunc func(name string, pc uint32, ins uint) string

// emuFunc is the emulator function
type emuFunc func(m *RV, ins uint)

// insMeta is instruction meta-data determined at runtime
type insMeta struct {
	name      string     // instruction mneumonic
	val, mask uint       // value and mask of fixed bits in the instruction
	dt        decodeType // decode type
}

// insDefn is the base definition of an instruction
type insDefn struct {
	defn string  // instruction definition string (from the standard)
	meta insMeta // meta data determined at runtime
	da   daFunc  // disassembly function
	emu  emuFunc // emulation function
}

// ISAModule is a set (module) of RISC-V instructions.
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
		{"imm[31:12] rd 0110111 LUI", insMeta{}, daTypeUa, emuNone},                         // U
		{"imm[31:12] rd 0010111 AUIPC", insMeta{}, daTypeUa, emuNone},                       // U
		{"imm[20|10:1|11|19:12] rd 1101111 JAL", insMeta{}, daTypeJa, emuNone},              // J
		{"imm[11:0] rs1 000 rd 1100111 JALR", insMeta{}, daTypeIe, emuNone},                 // I
		{"imm[12|10:5] rs2 rs1 000 imm[4:1|11] 1100011 BEQ", insMeta{}, daTypeBa, emuNone},  // B
		{"imm[12|10:5] rs2 rs1 001 imm[4:1|11] 1100011 BNE", insMeta{}, daTypeBa, emuNone},  // B
		{"imm[12|10:5] rs2 rs1 100 imm[4:1|11] 1100011 BLT", insMeta{}, daTypeBa, emuNone},  // B
		{"imm[12|10:5] rs2 rs1 101 imm[4:1|11] 1100011 BGE", insMeta{}, daTypeBa, emuNone},  // B
		{"imm[12|10:5] rs2 rs1 110 imm[4:1|11] 1100011 BLTU", insMeta{}, daTypeBa, emuNone}, // B
		{"imm[12|10:5] rs2 rs1 111 imm[4:1|11] 1100011 BGEU", insMeta{}, daTypeBa, emuNone}, // B
		{"imm[11:0] rs1 000 rd 0000011 LB", insMeta{}, daTypeIc, emuNone},                   // I
		{"imm[11:0] rs1 001 rd 0000011 LH", insMeta{}, daTypeIc, emuNone},                   // I
		{"imm[11:0] rs1 010 rd 0000011 LW", insMeta{}, daTypeIc, emuNone},                   // I
		{"imm[11:0] rs1 100 rd 0000011 LBU", insMeta{}, daTypeIc, emuNone},                  // I
		{"imm[11:0] rs1 101 rd 0000011 LHU", insMeta{}, daTypeIc, emuNone},                  // I
		{"imm[11:5] rs2 rs1 000 imm[4:0] 0100011 SB", insMeta{}, daTypeSa, emuNone},         // S
		{"imm[11:5] rs2 rs1 001 imm[4:0] 0100011 SH", insMeta{}, daTypeSa, emuNone},         // S
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100011 SW", insMeta{}, daTypeSa, emuNone},         // S
		{"imm[11:0] rs1 000 rd 0010011 ADDI", insMeta{}, daTypeIb, emuNone},                 // I
		{"imm[11:0] rs1 010 rd 0010011 SLTI", insMeta{}, daTypeIa, emuNone},                 // I
		{"imm[11:0] rs1 011 rd 0010011 SLTIU", insMeta{}, daTypeIa, emuNone},                // I
		{"imm[11:0] rs1 100 rd 0010011 XORI", insMeta{}, daTypeIf, emuNone},                 // I
		{"imm[11:0] rs1 110 rd 0010011 ORI", insMeta{}, daTypeIa, emuNone},                  // I
		{"imm[11:0] rs1 111 rd 0010011 ANDI", insMeta{}, daTypeIa, emuNone},                 // I
		{"0000000 shamt5 rs1 001 rd 0010011 SLLI", insMeta{}, daTypeId, emuNone},            // I
		{"0000000 shamt5 rs1 101 rd 0010011 SRLI", insMeta{}, daTypeId, emuNone},            // I
		{"0100000 shamt5 rs1 101 rd 0010011 SRAI", insMeta{}, daTypeId, emuNone},            // I
		{"0000000 rs2 rs1 000 rd 0110011 ADD", insMeta{}, daTypeRa, emuNone},                // R
		{"0100000 rs2 rs1 000 rd 0110011 SUB", insMeta{}, daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 001 rd 0110011 SLL", insMeta{}, daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 010 rd 0110011 SLT", insMeta{}, daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 011 rd 0110011 SLTU", insMeta{}, daTypeRa, emuNone},               // R
		{"0000000 rs2 rs1 100 rd 0110011 XOR", insMeta{}, daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 101 rd 0110011 SRL", insMeta{}, daTypeRa, emuNone},                // R
		{"0100000 rs2 rs1 101 rd 0110011 SRA", insMeta{}, daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 110 rd 0110011 OR", insMeta{}, daTypeRa, emuNone},                 // R
		{"0000000 rs2 rs1 111 rd 0110011 AND", insMeta{}, daTypeRa, emuNone},                // R
		{"0000 pred succ 00000 000 00000 0001111 FENCE", insMeta{}, daNone, emuNone},        // I
		{"0000 0000 0000 00000 001 00000 0001111 FENCE.I", insMeta{}, daNone, emuNone},      // I
		{"000000000000 00000 000 00000 1110011 ECALL", insMeta{}, daNone, emuNone},          // I
		{"000000000001 00000 000 00000 1110011 EBREAK", insMeta{}, daNone, emuNone},         // I
		{"csr rs1 001 rd 1110011 CSRRW", insMeta{}, daNone, emuNone},                        // I
		{"csr rs1 010 rd 1110011 CSRRS", insMeta{}, daNone, emuNone},                        // I
		{"csr rs1 011 rd 1110011 CSRRC", insMeta{}, daNone, emuNone},                        // I
		{"csr zimm 101 rd 1110011 CSRRWI", insMeta{}, daNone, emuNone},                      // I
		{"csr zimm 110 rd 1110011 CSRRSI", insMeta{}, daNone, emuNone},                      // I
		{"csr zimm 111 rd 1110011 CSRRCI", insMeta{}, daNone, emuNone},                      // I
	},
}

// ISArv32m Integer Multiplication and Division
var ISArv32m = ISAModule{
	name: "rv32m",
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0110011 MUL", insMeta{}, daTypeRa, emuNone},    // R
		{"0000001 rs2 rs1 001 rd 0110011 MULH", insMeta{}, daTypeRa, emuNone},   // R
		{"0000001 rs2 rs1 010 rd 0110011 MULHSU", insMeta{}, daTypeRa, emuNone}, // R
		{"0000001 rs2 rs1 011 rd 0110011 MULHU", insMeta{}, daTypeRa, emuNone},  // R
		{"0000001 rs2 rs1 100 rd 0110011 DIV", insMeta{}, daTypeRa, emuNone},    // R
		{"0000001 rs2 rs1 101 rd 0110011 DIVU", insMeta{}, daTypeRa, emuNone},   // R
		{"0000001 rs2 rs1 110 rd 0110011 REM", insMeta{}, daTypeRa, emuNone},    // R
		{"0000001 rs2 rs1 111 rd 0110011 REMU", insMeta{}, daTypeRa, emuNone},   // R
	},
}

// ISArv32a Atomics
var ISArv32a = ISAModule{
	name: "rv32a",
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 010 rd 0101111 LR.W", insMeta{}, daTypeRb, emuNone},    // R
		{"00011 aq rl rs2 rs1 010 rd 0101111 SC.W", insMeta{}, daTypeRb, emuNone},      // R
		{"00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W", insMeta{}, daTypeRb, emuNone}, // R
		{"00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W", insMeta{}, daTypeRb, emuNone},  // R
		{"00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W", insMeta{}, daTypeRb, emuNone},  // R
		{"01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W", insMeta{}, daTypeRb, emuNone},  // R
		{"01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W", insMeta{}, daTypeRb, emuNone},   // R
		{"10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W", insMeta{}, daTypeRb, emuNone},  // R
		{"10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W", insMeta{}, daTypeRb, emuNone},  // R
		{"11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W", insMeta{}, daTypeRb, emuNone}, // R
		{"11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W", insMeta{}, daTypeRb, emuNone}, // R
	},
}

// ISArv32f Single-Precision Floating-Point
var ISArv32f = ISAModule{
	name: "rv32f",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 010 rd 0000111 FLW", insMeta{}, daTypeIa, emuNone},           // I
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100111 FSW", insMeta{}, daTypeSb, emuNone}, // S
		{"rs3 00 rs2 rs1 rm rd 1000011 FMADD.S", insMeta{}, daNone, emuNone},         // R4
		{"rs3 00 rs2 rs1 rm rd 1000111 FMSUB.S", insMeta{}, daNone, emuNone},         // R4
		{"rs3 00 rs2 rs1 rm rd 1001011 FNMSUB.S", insMeta{}, daNone, emuNone},        // R4
		{"rs3 00 rs2 rs1 rm rd 1001111 FNMADD.S", insMeta{}, daNone, emuNone},        // R4
		{"0000000 rs2 rs1 rm rd 1010011 FADD.S", insMeta{}, daTypeRa, emuNone},       // R
		{"0000100 rs2 rs1 rm rd 1010011 FSUB.S", insMeta{}, daTypeRa, emuNone},       // R
		{"0001000 rs2 rs1 rm rd 1010011 FMUL.S", insMeta{}, daTypeRa, emuNone},       // R
		{"0001100 rs2 rs1 rm rd 1010011 FDIV.S", insMeta{}, daTypeRa, emuNone},       // R
		{"0101100 00000 rs1 rm rd 1010011 FSQRT.S", insMeta{}, daNone, emuNone},      // R
		{"0010000 rs2 rs1 000 rd 1010011 FSGNJ.S", insMeta{}, daTypeRa, emuNone},     // R
		{"0010000 rs2 rs1 001 rd 1010011 FSGNJN.S", insMeta{}, daTypeRa, emuNone},    // R
		{"0010000 rs2 rs1 010 rd 1010011 FSGNJX.S", insMeta{}, daTypeRa, emuNone},    // R
		{"0010100 rs2 rs1 000 rd 1010011 FMIN.S", insMeta{}, daTypeRa, emuNone},      // R
		{"0010100 rs2 rs1 001 rd 1010011 FMAX.S", insMeta{}, daTypeRa, emuNone},      // R
		{"1100000 00000 rs1 rm rd 1010011 FCVT.W.S", insMeta{}, daNone, emuNone},     // R
		{"1100000 00001 rs1 rm rd 1010011 FCVT.WU.S", insMeta{}, daNone, emuNone},    // R
		{"1110000 00000 rs1 000 rd 1010011 FMV.X.W", insMeta{}, daNone, emuNone},     // R
		{"1010000 rs2 rs1 010 rd 1010011 FEQ.S", insMeta{}, daTypeRa, emuNone},       // R
		{"1010000 rs2 rs1 001 rd 1010011 FLT.S", insMeta{}, daTypeRa, emuNone},       // R
		{"1010000 rs2 rs1 000 rd 1010011 FLE.S", insMeta{}, daTypeRa, emuNone},       // R
		{"1110000 00000 rs1 001 rd 1010011 FCLASS.S", insMeta{}, daNone, emuNone},    // R
		{"1101000 00000 rs1 rm rd 1010011 FCVT.S.W", insMeta{}, daNone, emuNone},     // R
		{"1101000 00001 rs1 rm rd 1010011 FCVT.S.WU", insMeta{}, daNone, emuNone},    // R
		{"1111000 00000 rs1 000 rd 1010011 FMV.W.X", insMeta{}, daNone, emuNone},     // R
	},
}

// ISArv32d Double-Precision Floating-Point
var ISArv32d = ISAModule{
	name: "rv32d",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 011 rd 0000111 FLD", insMeta{}, daTypeIg, emuNone},           // I
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100111 FSD", insMeta{}, daTypeSb, emuNone}, // S
		{"rs3 01 rs2 rs1 rm rd 1000011 FMADD.D", insMeta{}, daTypeR4a, emuNone},      // R4
		{"rs3 01 rs2 rs1 rm rd 1000111 FMSUB.D", insMeta{}, daTypeR4a, emuNone},      // R4
		{"rs3 01 rs2 rs1 rm rd 1001011 FNMSUB.D", insMeta{}, daTypeR4a, emuNone},     // R4
		{"rs3 01 rs2 rs1 rm rd 1001111 FNMADD.D", insMeta{}, daTypeR4a, emuNone},     // R4
		{"0000001 rs2 rs1 rm rd 1010011 FADD.D", insMeta{}, daNone, emuNone},         // R
		{"0000101 rs2 rs1 rm rd 1010011 FSUB.D", insMeta{}, daNone, emuNone},         // R
		{"0001001 rs2 rs1 rm rd 1010011 FMUL.D", insMeta{}, daNone, emuNone},         // R
		{"0001101 rs2 rs1 rm rd 1010011 FDIV.D", insMeta{}, daNone, emuNone},         // R
		{"0101101 00000 rs1 rm rd 1010011 FSQRT.D", insMeta{}, daNone, emuNone},      // R
		{"0010001 rs2 rs1 000 rd 1010011 FSGNJ.D", insMeta{}, daNone, emuNone},       // R
		{"0010001 rs2 rs1 001 rd 1010011 FSGNJN.D", insMeta{}, daNone, emuNone},      // R
		{"0010001 rs2 rs1 010 rd 1010011 FSGNJX.D", insMeta{}, daNone, emuNone},      // R
		{"0010101 rs2 rs1 000 rd 1010011 FMIN.D", insMeta{}, daNone, emuNone},        // R
		{"0010101 rs2 rs1 001 rd 1010011 FMAX.D", insMeta{}, daNone, emuNone},        // R
		{"0100000 00001 rs1 rm rd 1010011 FCVT.S.D", insMeta{}, daNone, emuNone},     // R
		{"0100001 00000 rs1 rm rd 1010011 FCVT.D.S", insMeta{}, daNone, emuNone},     // R
		{"1010001 rs2 rs1 010 rd 1010011 FEQ.D", insMeta{}, daNone, emuNone},         // R
		{"1010001 rs2 rs1 001 rd 1010011 FLT.D", insMeta{}, daNone, emuNone},         // R
		{"1010001 rs2 rs1 000 rd 1010011 FLE.D", insMeta{}, daNone, emuNone},         // R
		{"1110001 00000 rs1 001 rd 1010011 FCLASS.D", insMeta{}, daNone, emuNone},    // R
		{"1100001 00000 rs1 rm rd 1010011 FCVT.W.D", insMeta{}, daNone, emuNone},     // R
		{"1100001 00001 rs1 rm rd 1010011 FCVT.WU.D", insMeta{}, daNone, emuNone},    // R
		{"1101001 00000 rs1 rm rd 1010011 FCVT.D.W", insMeta{}, daNone, emuNone},     // R
		{"1101001 00001 rs1 rm rd 1010011 FCVT.D.WU", insMeta{}, daNone, emuNone},    // R
	},
}

// ISArv32c Compressed
var ISArv32c = ISAModule{
	name: "rv32c",
	ilen: 16,
	defn: []insDefn{
		{"000 00000000 000 00 C.ILLEGAL", insMeta{}, daTypeCIWa, emuNone},                 // CIW (Quadrant 0)
		{"000 nzuimm[5:4|9:6|2|3] rd0 00 C.ADDI4SPN", insMeta{}, daNone, emuNone},         // CIW
		{"001 uimm[5:3] rs10 uimm[7:6] rd0 00 C.FLD", insMeta{}, daNone, emuNone},         // CL
		{"010 uimm[5:3] rs10 uimm[2|6] rd0 00 C.LW", insMeta{}, daNone, emuNone},          // CL
		{"011 uimm[5:3] rs10 uimm[2|6] rd0 00 C.FLW", insMeta{}, daNone, emuNone},         // CL
		{"101 uimm[5:3] rs10 uimm[7:6] rs20 00 C.FSD", insMeta{}, daNone, emuNone},        // CS
		{"110 uimm[5:3] rs10 uimm[2|6] rs20 00 C.SW", insMeta{}, daNone, emuNone},         // CS
		{"111 uimm[5:3] rs10 uimm[2|6] rs20 00 C.FSW", insMeta{}, daNone, emuNone},        // CS
		{"000 nzimm[5] 00000 nzimm[4:0] 01 C.NOP", insMeta{}, daNone, emuNone},            // CI (Quadrant 1)
		{"000 nzimm[5] rs1/rd!=0 nzimm[4:0] 01 C.ADDI", insMeta{}, daNone, emuNone},       // CI
		{"001 imm[11|4|9:8|10|6|7|3:1|5] 01 C.JAL", insMeta{}, daNone, emuNone},           // CJ
		{"010 imm[5] rd!=0 imm[4:0] 01 C.LI", insMeta{}, daTypeCIa, emuNone},              // CI
		{"011 nzimm[9] 00010 nzimm[4|6|8:7|5] 01 C.ADDI16SP", insMeta{}, daNone, emuNone}, // CI
		{"011 nzimm[17] rd!={0,2} nzimm[16:12] 01 C.LUI", insMeta{}, daNone, emuNone},     // CI
		{"100 nzuimm[5] 00 rs10/rd0 nzuimm[4:0] 01 C.SRLI", insMeta{}, daNone, emuNone},   // CI
		{"100 nzuimm[5] 01 rs10/rd0 nzuimm[4:0] 01 C.SRAI", insMeta{}, daNone, emuNone},   // CI
		{"100 imm[5] 10 rs10/rd0 imm[4:0] 01 C.ANDI", insMeta{}, daNone, emuNone},         // CI
		{"100 0 11 rs10/rd0 00 rs20 01 C.SUB", insMeta{}, daNone, emuNone},                // CR
		{"100 0 11 rs10/rd0 01 rs20 01 C.XOR", insMeta{}, daNone, emuNone},                // CR
		{"100 0 11 rs10/rd0 10 rs20 01 C.OR", insMeta{}, daNone, emuNone},                 // CR
		{"100 0 11 rs10/rd0 11 rs20 01 C.AND", insMeta{}, daNone, emuNone},                // CR
		{"101 imm[11|4|9:8|10|6|7|3:1|5] 01 C.J", insMeta{}, daNone, emuNone},             // CJ
		{"110 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BEQZ", insMeta{}, daNone, emuNone},      // CB
		{"111 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BNEZ", insMeta{}, daNone, emuNone},      // CB
		{"000 nzuimm[5] rs1/rd!=0 nzuimm[4:0] 10 C.SLLI", insMeta{}, daNone, emuNone},     // CI (Quadrant 2)
		{"000 0 rs1/rd!=0 00000 10 C.SLLI64", insMeta{}, daNone, emuNone},                 // CI
		{"001 uimm[5] rd uimm[4:3|8:6] 10 C.FLDSP", insMeta{}, daNone, emuNone},           // CSS
		{"010 uimm[5] rd!=0 uimm[4:2|7:6] 10 C.LWSP", insMeta{}, daNone, emuNone},         // CSS
		{"011 uimm[5] rd uimm[4:2|7:6] 10 C.FLWSP", insMeta{}, daNone, emuNone},           // CSS
		{"100 0 rs1!=0 00000 10 C.JR", insMeta{}, daTypeCJa, emuNone},                     // CJ
		{"100 0 rd!=0 rs2!=0 10 C.MV", insMeta{}, daNone, emuNone},                        // CR
		{"100 1 00000 00000 10 C.EBREAK", insMeta{}, daNone, emuNone},                     // CI
		{"100 1 rs1!=0 00000 10 C.JALR", insMeta{}, daNone, emuNone},                      // CJ
		{"100 1 rs1/rd!=0 rs2!=0 10 C.ADD", insMeta{}, daNone, emuNone},                   // CR
		{"101 uimm[5:3|8:6] rs2 10 C.FSDSP", insMeta{}, daNone, emuNone},                  // CSS
		{"110 uimm[5:2|7:6] rs2 10 C.SWSP", insMeta{}, daNone, emuNone},                   // CSS
		{"111 uimm[5:2|7:6] rs2 10 C.FSWSP", insMeta{}, daNone, emuNone},                  // CSS
	},
}

//-----------------------------------------------------------------------------
// RV64 instructions (+ RV32)

// ISArv64i Integer
var ISArv64i = ISAModule{
	name: "rv64i",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 110 rd 0000011 LWU", insMeta{}, daTypeIa, emuNone},          // I
		{"imm[11:0] rs1 011 rd 0000011 LD", insMeta{}, daTypeIa, emuNone},           // I
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD", insMeta{}, daTypeSa, emuNone}, // S
		{"000000 shamt6 rs1 001 rd 0010011 SLLI", insMeta{}, daNone, emuNone},       // I
		{"000000 shamt6 rs1 101 rd 0010011 SRLI", insMeta{}, daNone, emuNone},       // I
		{"010000 shamt6 rs1 101 rd 0010011 SRAI", insMeta{}, daNone, emuNone},       // I
		{"imm[11:0] rs1 000 rd 0011011 ADDIW", insMeta{}, daTypeIa, emuNone},        // I
		{"0000000 shamt5 rs1 001 rd 0011011 SLLIW", insMeta{}, daNone, emuNone},     // I
		{"0000000 shamt5 rs1 101 rd 0011011 SRLIW", insMeta{}, daNone, emuNone},     // I
		{"0100000 shamt5 rs1 101 rd 0011011 SRAIW", insMeta{}, daNone, emuNone},     // I
		{"0000000 rs2 rs1 000 rd 0111011 ADDW", insMeta{}, daNone, emuNone},         // R
		{"0100000 rs2 rs1 000 rd 0111011 SUBW", insMeta{}, daNone, emuNone},         // R
		{"0000000 rs2 rs1 001 rd 0111011 SLLW", insMeta{}, daNone, emuNone},         // R
		{"0000000 rs2 rs1 101 rd 0111011 SRLW", insMeta{}, daNone, emuNone},         // R
		{"0100000 rs2 rs1 101 rd 0111011 SRAW", insMeta{}, daNone, emuNone},         // R
	},
}

// ISArv64m Integer Multiplication and Division
var ISArv64m = ISAModule{
	name: "rv64m",
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0111011 MULW", insMeta{}, daNone, emuNone},  // R
		{"0000001 rs2 rs1 100 rd 0111011 DIVW", insMeta{}, daNone, emuNone},  // R
		{"0000001 rs2 rs1 101 rd 0111011 DIVUW", insMeta{}, daNone, emuNone}, // R
		{"0000001 rs2 rs1 110 rd 0111011 REMW", insMeta{}, daNone, emuNone},  // R
		{"0000001 rs2 rs1 111 rd 0111011 REMUW", insMeta{}, daNone, emuNone}, // R
	},
}

// ISArv64a Atomics
var ISArv64a = ISAModule{
	name: "rv64a",
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 011 rd 0101111 LR.D", insMeta{}, daNone, emuNone},    // R
		{"00011 aq rl rs2 rs1 011 rd 0101111 SC.D", insMeta{}, daNone, emuNone},      // R
		{"00001 aq rl rs2 rs1 011 rd 0101111 AMOSWAP.D", insMeta{}, daNone, emuNone}, // R
		{"00000 aq rl rs2 rs1 011 rd 0101111 AMOADD.D", insMeta{}, daNone, emuNone},  // R
		{"00100 aq rl rs2 rs1 011 rd 0101111 AMOXOR.D", insMeta{}, daNone, emuNone},  // R
		{"01100 aq rl rs2 rs1 011 rd 0101111 AMOAND.D", insMeta{}, daNone, emuNone},  // R
		{"01000 aq rl rs2 rs1 011 rd 0101111 AMOOR.D", insMeta{}, daNone, emuNone},   // R
		{"10000 aq rl rs2 rs1 011 rd 0101111 AMOMIN.D", insMeta{}, daNone, emuNone},  // R
		{"10100 aq rl rs2 rs1 011 rd 0101111 AMOMAX.D", insMeta{}, daNone, emuNone},  // R
		{"11000 aq rl rs2 rs1 011 rd 0101111 AMOMINU.D", insMeta{}, daNone, emuNone}, // R
		{"11100 aq rl rs2 rs1 011 rd 0101111 AMOMAXU.D", insMeta{}, daNone, emuNone}, // R
	},
}

// ISArv64f Single-Precision Floating-Point
var ISArv64f = ISAModule{
	name: "rv64f",
	ilen: 32,
	defn: []insDefn{
		{"1100000 00010 rs1 rm rd 1010011 FCVT.L.S", insMeta{}, daNone, emuNone},  // R
		{"1100000 00011 rs1 rm rd 1010011 FCVT.LU.S", insMeta{}, daNone, emuNone}, // R
		{"1101000 00010 rs1 rm rd 1010011 FCVT.S.L", insMeta{}, daNone, emuNone},  // R
		{"1101000 00011 rs1 rm rd 1010011 FCVT.S.LU", insMeta{}, daNone, emuNone}, // R
	},
}

// ISArv64d Double-Precision Floating-Point
var ISArv64d = ISAModule{
	name: "rv64d",
	ilen: 32,
	defn: []insDefn{
		{"1100001 00010 rs1 rm rd 1010011 FCVT.L.D", insMeta{}, daNone, emuNone},  // R
		{"1100001 00011 rs1 rm rd 1010011 FCVT.LU.D", insMeta{}, daNone, emuNone}, // R
		{"1110001 00000 rs1 000 rd 1010011 FMV.X.D", insMeta{}, daNone, emuNone},  // R
		{"1101001 00010 rs1 rm rd 1010011 FCVT.D.L", insMeta{}, daNone, emuNone},  // R
		{"1101001 00011 rs1 rm rd 1010011 FCVT.D.LU", insMeta{}, daNone, emuNone}, // R
		{"1111001 00000 rs1 000 rd 1010011 FMV.D.X", insMeta{}, daNone, emuNone},  // R
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
