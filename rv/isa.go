//-----------------------------------------------------------------------------
/*

RISC-V ISA Definition

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

// daFunc is an instruction disassembly function
type daFunc func(name string, pc uint32, ins uint) string

// emuFunc is an instruction emulation function
type emuFunc32 func(m *RV32, ins uint)

// insDefn is the definition of an instruction
type insDefn struct {
	defn  string    // instruction definition string (from the standard)
	da    daFunc    // disassembly function
	emu32 emuFunc32 // 32-bit emulation function
}

// ISAModule is a set/module of RISC-V instructions.
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
		{"imm[31:12] rd 0110111 LUI", daTypeUa, emuNone},                         // U
		{"imm[31:12] rd 0010111 AUIPC", daTypeUa, emuNone},                       // U
		{"imm[20|10:1|11|19:12] rd 1101111 JAL", daTypeJa, emuJAL},               // J
		{"imm[11:0] rs1 000 rd 1100111 JALR", daTypeIe, emuNone},                 // I
		{"imm[12|10:5] rs2 rs1 000 imm[4:1|11] 1100011 BEQ", daTypeBa, emuNone},  // B
		{"imm[12|10:5] rs2 rs1 001 imm[4:1|11] 1100011 BNE", daTypeBa, emuNone},  // B
		{"imm[12|10:5] rs2 rs1 100 imm[4:1|11] 1100011 BLT", daTypeBa, emuNone},  // B
		{"imm[12|10:5] rs2 rs1 101 imm[4:1|11] 1100011 BGE", daTypeBa, emuNone},  // B
		{"imm[12|10:5] rs2 rs1 110 imm[4:1|11] 1100011 BLTU", daTypeBa, emuNone}, // B
		{"imm[12|10:5] rs2 rs1 111 imm[4:1|11] 1100011 BGEU", daTypeBa, emuNone}, // B
		{"imm[11:0] rs1 000 rd 0000011 LB", daTypeIc, emuNone},                   // I
		{"imm[11:0] rs1 001 rd 0000011 LH", daTypeIc, emuNone},                   // I
		{"imm[11:0] rs1 010 rd 0000011 LW", daTypeIc, emuNone},                   // I
		{"imm[11:0] rs1 100 rd 0000011 LBU", daTypeIc, emuNone},                  // I
		{"imm[11:0] rs1 101 rd 0000011 LHU", daTypeIc, emuNone},                  // I
		{"imm[11:5] rs2 rs1 000 imm[4:0] 0100011 SB", daTypeSa, emuNone},         // S
		{"imm[11:5] rs2 rs1 001 imm[4:0] 0100011 SH", daTypeSa, emuNone},         // S
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100011 SW", daTypeSa, emuNone},         // S
		{"imm[11:0] rs1 000 rd 0010011 ADDI", daTypeIb, emuNone},                 // I
		{"imm[11:0] rs1 010 rd 0010011 SLTI", daTypeIa, emuNone},                 // I
		{"imm[11:0] rs1 011 rd 0010011 SLTIU", daTypeIa, emuNone},                // I
		{"imm[11:0] rs1 100 rd 0010011 XORI", daTypeIf, emuNone},                 // I
		{"imm[11:0] rs1 110 rd 0010011 ORI", daTypeIa, emuNone},                  // I
		{"imm[11:0] rs1 111 rd 0010011 ANDI", daTypeIa, emuNone},                 // I
		{"0000000 shamt5 rs1 001 rd 0010011 SLLI", daTypeId, emuNone},            // I
		{"0000000 shamt5 rs1 101 rd 0010011 SRLI", daTypeId, emuNone},            // I
		{"0100000 shamt5 rs1 101 rd 0010011 SRAI", daTypeId, emuNone},            // I
		{"0000000 rs2 rs1 000 rd 0110011 ADD", daTypeRa, emuNone},                // R
		{"0100000 rs2 rs1 000 rd 0110011 SUB", daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 001 rd 0110011 SLL", daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 010 rd 0110011 SLT", daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 011 rd 0110011 SLTU", daTypeRa, emuNone},               // R
		{"0000000 rs2 rs1 100 rd 0110011 XOR", daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 101 rd 0110011 SRL", daTypeRa, emuNone},                // R
		{"0100000 rs2 rs1 101 rd 0110011 SRA", daTypeRa, emuNone},                // R
		{"0000000 rs2 rs1 110 rd 0110011 OR", daTypeRa, emuNone},                 // R
		{"0000000 rs2 rs1 111 rd 0110011 AND", daTypeRa, emuNone},                // R
		{"0000 pred succ 00000 000 00000 0001111 FENCE", daNone, emuNone},        // I
		{"0000 0000 0000 00000 001 00000 0001111 FENCE.I", daNone, emuNone},      // I
		{"000000000000 00000 000 00000 1110011 ECALL", daNone, emuNone},          // I
		{"000000000001 00000 000 00000 1110011 EBREAK", daNone, emuNone},         // I
		{"csr rs1 001 rd 1110011 CSRRW", daTypeIh, emuCSRRW},                     // I
		{"csr rs1 010 rd 1110011 CSRRS", daTypeIh, emuNone},                      // I
		{"csr rs1 011 rd 1110011 CSRRC", daTypeIh, emuNone},                      // I
		{"csr zimm 101 rd 1110011 CSRRWI", daNone, emuNone},                      // I
		{"csr zimm 110 rd 1110011 CSRRSI", daNone, emuNone},                      // I
		{"csr zimm 111 rd 1110011 CSRRCI", daNone, emuNone},                      // I
	},
}

// ISArv32m Integer Multiplication and Division
var ISArv32m = ISAModule{
	name: "rv32m",
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0110011 MUL", daTypeRa, emuNone},    // R
		{"0000001 rs2 rs1 001 rd 0110011 MULH", daTypeRa, emuNone},   // R
		{"0000001 rs2 rs1 010 rd 0110011 MULHSU", daTypeRa, emuNone}, // R
		{"0000001 rs2 rs1 011 rd 0110011 MULHU", daTypeRa, emuNone},  // R
		{"0000001 rs2 rs1 100 rd 0110011 DIV", daTypeRa, emuNone},    // R
		{"0000001 rs2 rs1 101 rd 0110011 DIVU", daTypeRa, emuNone},   // R
		{"0000001 rs2 rs1 110 rd 0110011 REM", daTypeRa, emuNone},    // R
		{"0000001 rs2 rs1 111 rd 0110011 REMU", daTypeRa, emuNone},   // R
	},
}

// ISArv32a Atomics
var ISArv32a = ISAModule{
	name: "rv32a",
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 010 rd 0101111 LR.W", daTypeRb, emuNone},    // R
		{"00011 aq rl rs2 rs1 010 rd 0101111 SC.W", daTypeRb, emuNone},      // R
		{"00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W", daTypeRb, emuNone}, // R
		{"00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W", daTypeRb, emuNone},  // R
		{"00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W", daTypeRb, emuNone},  // R
		{"01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W", daTypeRb, emuNone},  // R
		{"01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W", daTypeRb, emuNone},   // R
		{"10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W", daTypeRb, emuNone},  // R
		{"10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W", daTypeRb, emuNone},  // R
		{"11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W", daTypeRb, emuNone}, // R
		{"11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W", daTypeRb, emuNone}, // R
	},
}

// ISArv32f Single-Precision Floating-Point
var ISArv32f = ISAModule{
	name: "rv32f",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 010 rd 0000111 FLW", daTypeIg, emuNone},            // I
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100111 FSW", daTypeSb, emuNone},  // S
		{"rs3 00 rs2 rs1 rm rd 1000011 FMADD.S", daNone, emuNone},          // R4
		{"rs3 00 rs2 rs1 rm rd 1000111 FMSUB.S", daNone, emuNone},          // R4
		{"rs3 00 rs2 rs1 rm rd 1001011 FNMSUB.S", daNone, emuNone},         // R4
		{"rs3 00 rs2 rs1 rm rd 1001111 FNMADD.S", daNone, emuNone},         // R4
		{"0000000 rs2 rs1 rm rd 1010011 FADD.S", daTypeRc, emuNone},        // R
		{"0000100 rs2 rs1 rm rd 1010011 FSUB.S", daTypeRc, emuNone},        // R
		{"0001000 rs2 rs1 rm rd 1010011 FMUL.S", daTypeRc, emuNone},        // R
		{"0001100 rs2 rs1 rm rd 1010011 FDIV.S", daTypeRc, emuNone},        // R
		{"0101100 00000 rs1 rm rd 1010011 FSQRT.S", daNone, emuNone},       // R
		{"0010000 rs2 rs1 000 rd 1010011 FSGNJ.S", daTypeRc, emuNone},      // R
		{"0010000 rs2 rs1 001 rd 1010011 FSGNJN.S", daTypeRc, emuNone},     // R
		{"0010000 rs2 rs1 010 rd 1010011 FSGNJX.S", daTypeRc, emuNone},     // R
		{"0010100 rs2 rs1 000 rd 1010011 FMIN.S", daTypeRc, emuNone},       // R
		{"0010100 rs2 rs1 001 rd 1010011 FMAX.S", daTypeRc, emuNone},       // R
		{"1100000 00000 rs1 rm rd 1010011 FCVT.W.S", daNone, emuNone},      // R
		{"1100000 00001 rs1 rm rd 1010011 FCVT.WU.S", daNone, emuNone},     // R
		{"1110000 00000 rs1 000 rd 1010011 FMV.X.W", daTypeRd, emuNone},    // R
		{"1010000 rs2 rs1 010 rd 1010011 FEQ.S", daTypeRa, emuNone},        // R
		{"1010000 rs2 rs1 001 rd 1010011 FLT.S", daTypeRa, emuNone},        // R
		{"1010000 rs2 rs1 000 rd 1010011 FLE.S", daTypeRa, emuNone},        // R
		{"1110000 00000 rs1 001 rd 1010011 FCLASS.S", daNone, emuNone},     // R
		{"1101000 00000 rs1 rm rd 1010011 FCVT.S.W", daNone, emuNone},      // R
		{"1101000 00001 rs1 rm rd 1010011 FCVT.S.WU", daNone, emuNone},     // R
		{"1111000 00000 rs1 000 rd 1010011 FMV.W.X", daTypeRe, emuFMVxWxX}, // R
	},
}

// ISArv32d Double-Precision Floating-Point
var ISArv32d = ISAModule{
	name: "rv32d",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 011 rd 0000111 FLD", daTypeIg, emuNone},           // I
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100111 FSD", daTypeSb, emuNone}, // S
		{"rs3 01 rs2 rs1 rm rd 1000011 FMADD.D", daTypeR4a, emuNone},      // R4
		{"rs3 01 rs2 rs1 rm rd 1000111 FMSUB.D", daTypeR4a, emuNone},      // R4
		{"rs3 01 rs2 rs1 rm rd 1001011 FNMSUB.D", daTypeR4a, emuNone},     // R4
		{"rs3 01 rs2 rs1 rm rd 1001111 FNMADD.D", daTypeR4a, emuNone},     // R4
		{"0000001 rs2 rs1 rm rd 1010011 FADD.D", daNone, emuNone},         // R
		{"0000101 rs2 rs1 rm rd 1010011 FSUB.D", daNone, emuNone},         // R
		{"0001001 rs2 rs1 rm rd 1010011 FMUL.D", daNone, emuNone},         // R
		{"0001101 rs2 rs1 rm rd 1010011 FDIV.D", daNone, emuNone},         // R
		{"0101101 00000 rs1 rm rd 1010011 FSQRT.D", daNone, emuNone},      // R
		{"0010001 rs2 rs1 000 rd 1010011 FSGNJ.D", daNone, emuNone},       // R
		{"0010001 rs2 rs1 001 rd 1010011 FSGNJN.D", daNone, emuNone},      // R
		{"0010001 rs2 rs1 010 rd 1010011 FSGNJX.D", daNone, emuNone},      // R
		{"0010101 rs2 rs1 000 rd 1010011 FMIN.D", daNone, emuNone},        // R
		{"0010101 rs2 rs1 001 rd 1010011 FMAX.D", daNone, emuNone},        // R
		{"0100000 00001 rs1 rm rd 1010011 FCVT.S.D", daNone, emuNone},     // R
		{"0100001 00000 rs1 rm rd 1010011 FCVT.D.S", daNone, emuNone},     // R
		{"1010001 rs2 rs1 010 rd 1010011 FEQ.D", daNone, emuNone},         // R
		{"1010001 rs2 rs1 001 rd 1010011 FLT.D", daNone, emuNone},         // R
		{"1010001 rs2 rs1 000 rd 1010011 FLE.D", daNone, emuNone},         // R
		{"1110001 00000 rs1 001 rd 1010011 FCLASS.D", daNone, emuNone},    // R
		{"1100001 00000 rs1 rm rd 1010011 FCVT.W.D", daNone, emuNone},     // R
		{"1100001 00001 rs1 rm rd 1010011 FCVT.WU.D", daNone, emuNone},    // R
		{"1101001 00000 rs1 rm rd 1010011 FCVT.D.W", daNone, emuNone},     // R
		{"1101001 00001 rs1 rm rd 1010011 FCVT.D.WU", daNone, emuNone},    // R
	},
}

// ISArv32c Compressed
var ISArv32c = ISAModule{
	name: "rv32c",
	ilen: 16,
	defn: []insDefn{
		{"000 00000000 000 00 C.ILLEGAL", daTypeCIWa, emuNone},                    // CIW (Quadrant 0)
		{"000 nzuimm[5:4|9:6|2|3] rd0 00 C.ADDI4SPN", daTypeCIWb, emuNone},        // CIW
		{"001 uimm[5:3] rs10 uimm[7:6] rd0 00 C.FLD", daNone, emuNone},            // CL
		{"010 uimm[5:3] rs10 uimm[2|6] rd0 00 C.LW", daNone, emuNone},             // CL
		{"011 uimm[5:3] rs10 uimm[2|6] rd0 00 C.FLW", daNone, emuNone},            // CL
		{"101 uimm[5:3] rs10 uimm[7:6] rs20 00 C.FSD", daNone, emuNone},           // CS
		{"110 uimm[5:3] rs10 uimm[2|6] rs20 00 C.SW", daNone, emuNone},            // CS
		{"111 uimm[5:3] rs10 uimm[2|6] rs20 00 C.FSW", daNone, emuNone},           // CS
		{"000 nzimm[5] 00000 nzimm[4:0] 01 C.NOP", daNone, emuNone},               // CI (Quadrant 1)
		{"000 nzimm[5] rs1/rd!=0 nzimm[4:0] 01 C.ADDI", daTypeCIc, emuNone},       // CI
		{"001 imm[11|4|9:8|10|6|7|3:1|5] 01 C.JAL", daNone, emuNone},              // CJ
		{"010 imm[5] rd!=0 imm[4:0] 01 C.LI", daTypeCIa, emuCxLI},                 // CI
		{"011 nzimm[9] 00010 nzimm[4|6|8:7|5] 01 C.ADDI16SP", daTypeCIb, emuNone}, // CI
		{"011 nzimm[17] rd!={0,2} nzimm[16:12] 01 C.LUI", daTypeCIg, emuNone},     // CI
		{"100 nzuimm[5] 00 rs10/rd0 nzuimm[4:0] 01 C.SRLI", daTypeCId, emuNone},   // CI
		{"100 nzuimm[5] 01 rs10/rd0 nzuimm[4:0] 01 C.SRAI", daTypeCId, emuNone},   // CI
		{"100 imm[5] 10 rs10/rd0 imm[4:0] 01 C.ANDI", daTypeCIf, emuNone},         // CI
		{"100 0 11 rs10/rd0 00 rs20 01 C.SUB", daNone, emuNone},                   // CR
		{"100 0 11 rs10/rd0 01 rs20 01 C.XOR", daNone, emuNone},                   // CR
		{"100 0 11 rs10/rd0 10 rs20 01 C.OR", daNone, emuNone},                    // CR
		{"100 0 11 rs10/rd0 11 rs20 01 C.AND", daNone, emuNone},                   // CR
		{"101 imm[11|4|9:8|10|6|7|3:1|5] 01 C.J", daTypeCJb, emuNone},             // CJ
		{"110 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BEQZ", daTypeCBa, emuNone},      // CB
		{"111 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BNEZ", daTypeCBa, emuNone},      // CB
		{"000 nzuimm[5] rs1/rd!=0 nzuimm[4:0] 10 C.SLLI", daTypeCIe, emuNone},     // CI (Quadrant 2)
		{"000 0 rs1/rd!=0 00000 10 C.SLLI64", daNone, emuNone},                    // CI
		{"001 uimm[5] rd uimm[4:3|8:6] 10 C.FLDSP", daNone, emuNone},              // CSS
		{"010 uimm[5] rd!=0 uimm[4:2|7:6] 10 C.LWSP", daTypeCSSa, emuNone},        // CSS
		{"011 uimm[5] rd uimm[4:2|7:6] 10 C.FLWSP", daNone, emuNone},              // CSS
		{"100 0 rs1!=0 00000 10 C.JR", daTypeCJa, emuNone},                        // CJ
		{"100 0 rd!=0 rs2!=0 10 C.MV", daTypeCRa, emuNone},                        // CR
		{"100 1 00000 00000 10 C.EBREAK", daNone, emuNone},                        // CI
		{"100 1 rs1!=0 00000 10 C.JALR", daNone, emuNone},                         // CJ
		{"100 1 rs1/rd!=0 rs2!=0 10 C.ADD", daTypeCRb, emuNone},                   // CR
		{"101 uimm[5:3|8:6] rs2 10 C.FSDSP", daNone, emuNone},                     // CSS
		{"110 uimm[5:2|7:6] rs2 10 C.SWSP", daTypeCSSb, emuNone},                  // CSS
		{"111 uimm[5:2|7:6] rs2 10 C.FSWSP", daNone, emuNone},                     // CSS
	},
}

//-----------------------------------------------------------------------------
// RV64 instructions (+ RV32)

// ISArv64i Integer
var ISArv64i = ISAModule{
	name: "rv64i",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 110 rd 0000011 LWU", daTypeIa, emuNone},          // I
		{"imm[11:0] rs1 011 rd 0000011 LD", daTypeIa, emuNone},           // I
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD", daTypeSa, emuNone}, // S
		{"000000 shamt6 rs1 001 rd 0010011 SLLI", daNone, emuNone},       // I
		{"000000 shamt6 rs1 101 rd 0010011 SRLI", daNone, emuNone},       // I
		{"010000 shamt6 rs1 101 rd 0010011 SRAI", daNone, emuNone},       // I
		{"imm[11:0] rs1 000 rd 0011011 ADDIW", daTypeIa, emuNone},        // I
		{"0000000 shamt5 rs1 001 rd 0011011 SLLIW", daNone, emuNone},     // I
		{"0000000 shamt5 rs1 101 rd 0011011 SRLIW", daNone, emuNone},     // I
		{"0100000 shamt5 rs1 101 rd 0011011 SRAIW", daNone, emuNone},     // I
		{"0000000 rs2 rs1 000 rd 0111011 ADDW", daNone, emuNone},         // R
		{"0100000 rs2 rs1 000 rd 0111011 SUBW", daNone, emuNone},         // R
		{"0000000 rs2 rs1 001 rd 0111011 SLLW", daNone, emuNone},         // R
		{"0000000 rs2 rs1 101 rd 0111011 SRLW", daNone, emuNone},         // R
		{"0100000 rs2 rs1 101 rd 0111011 SRAW", daNone, emuNone},         // R
	},
}

// ISArv64m Integer Multiplication and Division
var ISArv64m = ISAModule{
	name: "rv64m",
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0111011 MULW", daNone, emuNone},  // R
		{"0000001 rs2 rs1 100 rd 0111011 DIVW", daNone, emuNone},  // R
		{"0000001 rs2 rs1 101 rd 0111011 DIVUW", daNone, emuNone}, // R
		{"0000001 rs2 rs1 110 rd 0111011 REMW", daNone, emuNone},  // R
		{"0000001 rs2 rs1 111 rd 0111011 REMUW", daNone, emuNone}, // R
	},
}

// ISArv64a Atomics
var ISArv64a = ISAModule{
	name: "rv64a",
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 011 rd 0101111 LR.D", daNone, emuNone},    // R
		{"00011 aq rl rs2 rs1 011 rd 0101111 SC.D", daNone, emuNone},      // R
		{"00001 aq rl rs2 rs1 011 rd 0101111 AMOSWAP.D", daNone, emuNone}, // R
		{"00000 aq rl rs2 rs1 011 rd 0101111 AMOADD.D", daNone, emuNone},  // R
		{"00100 aq rl rs2 rs1 011 rd 0101111 AMOXOR.D", daNone, emuNone},  // R
		{"01100 aq rl rs2 rs1 011 rd 0101111 AMOAND.D", daNone, emuNone},  // R
		{"01000 aq rl rs2 rs1 011 rd 0101111 AMOOR.D", daNone, emuNone},   // R
		{"10000 aq rl rs2 rs1 011 rd 0101111 AMOMIN.D", daNone, emuNone},  // R
		{"10100 aq rl rs2 rs1 011 rd 0101111 AMOMAX.D", daNone, emuNone},  // R
		{"11000 aq rl rs2 rs1 011 rd 0101111 AMOMINU.D", daNone, emuNone}, // R
		{"11100 aq rl rs2 rs1 011 rd 0101111 AMOMAXU.D", daNone, emuNone}, // R
	},
}

// ISArv64f Single-Precision Floating-Point
var ISArv64f = ISAModule{
	name: "rv64f",
	ilen: 32,
	defn: []insDefn{
		{"1100000 00010 rs1 rm rd 1010011 FCVT.L.S", daNone, emuNone},  // R
		{"1100000 00011 rs1 rm rd 1010011 FCVT.LU.S", daNone, emuNone}, // R
		{"1101000 00010 rs1 rm rd 1010011 FCVT.S.L", daNone, emuNone},  // R
		{"1101000 00011 rs1 rm rd 1010011 FCVT.S.LU", daNone, emuNone}, // R
	},
}

// ISArv64d Double-Precision Floating-Point
var ISArv64d = ISAModule{
	name: "rv64d",
	ilen: 32,
	defn: []insDefn{
		{"1100001 00010 rs1 rm rd 1010011 FCVT.L.D", daNone, emuNone},  // R
		{"1100001 00011 rs1 rm rd 1010011 FCVT.LU.D", daNone, emuNone}, // R
		{"1110001 00000 rs1 000 rd 1010011 FMV.X.D", daNone, emuNone},  // R
		{"1101001 00010 rs1 rm rd 1010011 FCVT.D.L", daNone, emuNone},  // R
		{"1101001 00011 rs1 rm rd 1010011 FCVT.D.LU", daNone, emuNone}, // R
		{"1111001 00000 rs1 000 rd 1010011 FMV.D.X", daNone, emuNone},  // R
	},
}

//-----------------------------------------------------------------------------
// pre-canned ISA module sets

// ISArv32g = RV32imafd
var ISArv32g = []ISAModule{
	ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d,
}

// ISArv32gc = RV32imafdc
var ISArv32gc = []ISAModule{
	ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d, ISArv32c,
}

// ISArv64g = RV64imafd
var ISArv64g = []ISAModule{
	ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d,
	ISArv64i, ISArv64m, ISArv64a, ISArv64f, ISArv64d,
}

// ISArv64gc = RV64imafdc
var ISArv64gc = []ISAModule{
	ISArv32i, ISArv32m, ISArv32a, ISArv32f, ISArv32d, ISArv32c,
	ISArv64i, ISArv64m, ISArv64a, ISArv64f, ISArv64d,
}

//-----------------------------------------------------------------------------

// insMeta is instruction meta-data determined at runtime
type insMeta struct {
	defn      *insDefn   // the instruction definition
	name      string     // instruction mneumonic
	n         int        // instruction bit length
	val, mask uint       // value and mask of fixed bits in the instruction
	dt        decodeType // decode type
}

// ISA is an instruction set
type ISA struct {
	ins16 []*insMeta // the set of 16-bit instructions in the ISA
	ins32 []*insMeta // the set of 32-bit instructions in the ISA
}

// NewISA creates an empty instruction set.
func NewISA() *ISA {
	return &ISA{
		ins16: make([]*insMeta, 0),
		ins32: make([]*insMeta, 0),
	}
}

// Add a ISA sub-module to the ISA.
func (isa *ISA) Add(module []ISAModule) error {
	for i := range module {
		for j := range module[i].defn {
			im, err := parseDefn(&module[i].defn[j], module[i].ilen)
			if err != nil {
				return err
			}
			if im.n == 16 {
				isa.ins16 = append(isa.ins16, im)
			} else {
				isa.ins32 = append(isa.ins32, im)
			}
		}
	}
	return nil
}

// lookup returns the instruction meta information for an instruction.
func (isa *ISA) lookup(ins uint) *insMeta {
	if ins&3 == 3 {
		// 32-bit instruction
		for _, im := range isa.ins32 {
			if ins&im.mask == im.val {
				return im
			}
		}
	} else {
		// 16-bit instruction
		for _, im := range isa.ins16 {
			if ins&im.mask == im.val {
				return im
			}
		}
	}
	return nil
}

//-----------------------------------------------------------------------------
