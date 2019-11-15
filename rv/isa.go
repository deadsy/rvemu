//-----------------------------------------------------------------------------
/*

RISC-V ISA Definition

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

// daFunc is an instruction disassembly function
type daFunc func(name string, pc uint32, ins uint) string

// emuFunc32 is 32-bit emulation function
type emuFunc32 func(m *RV32, ins uint)

// emuFunc64 is a 64-bit emulation function
type emuFunc64 func(m *RV64, ins uint)

// insDefn is the definition of an instruction
type insDefn struct {
	defn  string    // instruction definition string (from the standard)
	da    daFunc    // disassembly function
	emu32 emuFunc32 // 32-bit emulation function
	emu64 emuFunc64 // 64-bit emulation function
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
		{"imm[31:12] rd 0110111 LUI", daTypeUa, emu32_LUI, emu64_LUI},                            // U
		{"imm[31:12] rd 0010111 AUIPC", daTypeUa, emu32_AUIPC, emu64_AUIPC},                      // U
		{"imm[20|10:1|11|19:12] rd 1101111 JAL", daTypeJa, emu32_JAL, emu64_JAL},                 // J
		{"imm[11:0] rs1 000 rd 1100111 JALR", daTypeIe, emu32_JALR, emu64_JALR},                  // I
		{"imm[12|10:5] rs2 rs1 000 imm[4:1|11] 1100011 BEQ", daTypeBa, emu32_BEQ, emu64_BEQ},     // B
		{"imm[12|10:5] rs2 rs1 001 imm[4:1|11] 1100011 BNE", daTypeBa, emu32_BNE, emu64_BNE},     // B
		{"imm[12|10:5] rs2 rs1 100 imm[4:1|11] 1100011 BLT", daTypeBa, emu32_BLT, emu64_BLT},     // B
		{"imm[12|10:5] rs2 rs1 101 imm[4:1|11] 1100011 BGE", daTypeBa, emu32_BGE, emu64_BGE},     // B
		{"imm[12|10:5] rs2 rs1 110 imm[4:1|11] 1100011 BLTU", daTypeBa, emu32_BLTU, emu64_BLTU},  // B
		{"imm[12|10:5] rs2 rs1 111 imm[4:1|11] 1100011 BGEU", daTypeBa, emu32_BGEU, emu64_BGEU},  // B
		{"imm[11:0] rs1 000 rd 0000011 LB", daTypeIc, emu32_LB, emu64_LB},                        // I
		{"imm[11:0] rs1 001 rd 0000011 LH", daTypeIc, emu32_LH, emu64_LH},                        // I
		{"imm[11:0] rs1 010 rd 0000011 LW", daTypeIc, emu32_LW, emu64_LW},                        // I
		{"imm[11:0] rs1 100 rd 0000011 LBU", daTypeIc, emu32_LBU, emu64_LBU},                     // I
		{"imm[11:0] rs1 101 rd 0000011 LHU", daTypeIc, emu32_LHU, emu64_LHU},                     // I
		{"imm[11:5] rs2 rs1 000 imm[4:0] 0100011 SB", daTypeSa, emu32_SB, emu64_SB},              // S
		{"imm[11:5] rs2 rs1 001 imm[4:0] 0100011 SH", daTypeSa, emu32_SH, emu64_SH},              // S
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100011 SW", daTypeSa, emu32_SW, emu64_SW},              // S
		{"imm[11:0] rs1 000 rd 0010011 ADDI", daTypeIb, emu32_ADDI, emu64_ADDI},                  // I
		{"imm[11:0] rs1 010 rd 0010011 SLTI", daTypeIa, emu32_SLTI, emu64_SLTI},                  // I
		{"imm[11:0] rs1 011 rd 0010011 SLTIU", daTypeIa, emu32_SLTIU, emu64_SLTIU},               // I
		{"imm[11:0] rs1 100 rd 0010011 XORI", daTypeIf, emu32_XORI, emu64_XORI},                  // I
		{"imm[11:0] rs1 110 rd 0010011 ORI", daTypeIa, emu32_ORI, emu64_ORI},                     // I
		{"imm[11:0] rs1 111 rd 0010011 ANDI", daTypeIa, emu32_ANDI, emu64_ANDI},                  // I
		{"0000000 shamt5 rs1 001 rd 0010011 SLLI", daTypeId, emu32_SLLI, emu64_SLLI},             // I
		{"0000000 shamt5 rs1 101 rd 0010011 SRLI", daTypeId, emu32_SRLI, emu64_SRLI},             // I
		{"0100000 shamt5 rs1 101 rd 0010011 SRAI", daTypeId, emu32_SRAI, emu64_SRAI},             // I
		{"0000000 rs2 rs1 000 rd 0110011 ADD", daTypeRa, emu32_ADD, emu64_ADD},                   // R
		{"0100000 rs2 rs1 000 rd 0110011 SUB", daTypeRa, emu32_SUB, emu64_SUB},                   // R
		{"0000000 rs2 rs1 001 rd 0110011 SLL", daTypeRa, emu32_SLL, emu64_SLL},                   // R
		{"0000000 rs2 rs1 010 rd 0110011 SLT", daTypeRa, emu32_SLT, emu64_SLT},                   // R
		{"0000000 rs2 rs1 011 rd 0110011 SLTU", daTypeRa, emu32_SLTU, emu64_SLTU},                // R
		{"0000000 rs2 rs1 100 rd 0110011 XOR", daTypeRa, emu32_XOR, emu64_XOR},                   // R
		{"0000000 rs2 rs1 101 rd 0110011 SRL", daTypeRa, emu32_SRL, emu64_SRL},                   // R
		{"0100000 rs2 rs1 101 rd 0110011 SRA", daTypeRa, emu32_SRA, emu64_SRA},                   // R
		{"0000000 rs2 rs1 110 rd 0110011 OR", daTypeRa, emu32_OR, emu64_OR},                      // R
		{"0000000 rs2 rs1 111 rd 0110011 AND", daTypeRa, emu32_AND, emu64_AND},                   // R
		{"0000 pred succ 00000 000 00000 0001111 FENCE", daNone, emu32_FENCE, emu64_FENCE},       // I
		{"0000 0000 0000 00000 001 00000 0001111 FENCE.I", daNone, emu32_FENCE_I, emu64_FENCE_I}, // I
		{"000000000000 00000 000 00000 1110011 ECALL", daNone, emu32_ECALL, emu64_ECALL},         // I
		{"000000000001 00000 000 00000 1110011 EBREAK", daNone, emu32_EBREAK, emu64_EBREAK},      // I
		{"csr rs1 001 rd 1110011 CSRRW", daTypeIh, emu32_CSRRW, emu64_CSRRW},                     // I
		{"csr rs1 010 rd 1110011 CSRRS", daTypeIh, emu32_CSRRS, emu64_CSRRS},                     // I
		{"csr rs1 011 rd 1110011 CSRRC", daTypeIh, emu32_CSRRC, emu64_CSRRC},                     // I
		{"csr zimm 101 rd 1110011 CSRRWI", daNone, emu32_CSRRWI, emu64_CSRRWI},                   // I
		{"csr zimm 110 rd 1110011 CSRRSI", daNone, emu32_CSRRSI, emu64_CSRRSI},                   // I
		{"csr zimm 111 rd 1110011 CSRRCI", daNone, emu32_CSRRCI, emu64_CSRRCI},                   // I
	},
}

// ISArv32m Integer Multiplication and Division
var ISArv32m = ISAModule{
	name: "rv32m",
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0110011 MUL", daTypeRa, emu32_MUL, emu64_MUL},          // R
		{"0000001 rs2 rs1 001 rd 0110011 MULH", daTypeRa, emu32_MULH, emu64_MULH},       // R
		{"0000001 rs2 rs1 010 rd 0110011 MULHSU", daTypeRa, emu32_MULHSU, emu64_MULHSU}, // R
		{"0000001 rs2 rs1 011 rd 0110011 MULHU", daTypeRa, emu32_MULHU, emu64_MULHU},    // R
		{"0000001 rs2 rs1 100 rd 0110011 DIV", daTypeRa, emu32_DIV, emu64_DIV},          // R
		{"0000001 rs2 rs1 101 rd 0110011 DIVU", daTypeRa, emu32_DIVU, emu64_DIVU},       // R
		{"0000001 rs2 rs1 110 rd 0110011 REM", daTypeRa, emu32_REM, emu64_REM},          // R
		{"0000001 rs2 rs1 111 rd 0110011 REMU", daTypeRa, emu32_REMU, emu64_REMU},       // R
	},
}

// ISArv32a Atomics
var ISArv32a = ISAModule{
	name: "rv32a",
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 010 rd 0101111 LR.W", daTypeRb, emu32_LR_W, emu64_LR_W},              // R
		{"00011 aq rl rs2 rs1 010 rd 0101111 SC.W", daTypeRb, emu32_SC_W, emu64_SC_W},                // R
		{"00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W", daTypeRb, emu32_AMOSWAP_W, emu64_AMOSWAP_W}, // R
		{"00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W", daTypeRb, emu32_AMOADD_W, emu64_AMOADD_W},    // R
		{"00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W", daTypeRb, emu32_AMOXOR_W, emu64_AMOXOR_W},    // R
		{"01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W", daTypeRb, emu32_AMOAND_W, emu64_AMOAND_W},    // R
		{"01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W", daTypeRb, emu32_AMOOR_W, emu64_AMOOR_W},       // R
		{"10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W", daTypeRb, emu32_AMOMIN_W, emu64_AMOMIN_W},    // R
		{"10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W", daTypeRb, emu32_AMOMAX_W, emu64_AMOMAX_W},    // R
		{"11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W", daTypeRb, emu32_AMOMINU_W, emu64_AMOMINU_W}, // R
		{"11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W", daTypeRb, emu32_AMOMAXU_W, emu64_AMOMAXU_W}, // R
	},
}

// ISArv32f Single-Precision Floating-Point
var ISArv32f = ISAModule{
	name: "rv32f",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 010 rd 0000111 FLW", daTypeIg, emu32_FLW, emu64_FLW},                    // I
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100111 FSW", daTypeSb, emu32_FSW, emu64_FSW},          // S
		{"rs3 00 rs2 rs1 rm rd 1000011 FMADD.S", daNone, emu32_FMADD_S, emu64_FMADD_S},          // R4
		{"rs3 00 rs2 rs1 rm rd 1000111 FMSUB.S", daNone, emu32_FMSUB_S, emu64_FMSUB_S},          // R4
		{"rs3 00 rs2 rs1 rm rd 1001011 FNMSUB.S", daNone, emu32_FNMSUB_S, emu64_FNMSUB_S},       // R4
		{"rs3 00 rs2 rs1 rm rd 1001111 FNMADD.S", daNone, emu32_FNMADD_S, emu64_FNMADD_S},       // R4
		{"0000000 rs2 rs1 rm rd 1010011 FADD.S", daTypeRc, emu32_FADD_S, emu64_FADD_S},          // R
		{"0000100 rs2 rs1 rm rd 1010011 FSUB.S", daTypeRc, emu32_FSUB_S, emu64_FSUB_S},          // R
		{"0001000 rs2 rs1 rm rd 1010011 FMUL.S", daTypeRc, emu32_FMUL_S, emu64_FMUL_S},          // R
		{"0001100 rs2 rs1 rm rd 1010011 FDIV.S", daTypeRc, emu32_FDIV_S, emu64_FDIV_S},          // R
		{"0101100 00000 rs1 rm rd 1010011 FSQRT.S", daNone, emu32_FSQRT_S, emu64_FSQRT_S},       // R
		{"0010000 rs2 rs1 000 rd 1010011 FSGNJ.S", daTypeRc, emu32_FSGNJ_S, emu64_FSGNJ_S},      // R
		{"0010000 rs2 rs1 001 rd 1010011 FSGNJN.S", daTypeRc, emu32_FSGNJN_S, emu64_FSGNJN_S},   // R
		{"0010000 rs2 rs1 010 rd 1010011 FSGNJX.S", daTypeRc, emu32_FSGNJX_S, emu64_FSGNJX_S},   // R
		{"0010100 rs2 rs1 000 rd 1010011 FMIN.S", daTypeRc, emu32_FMIN_S, emu64_FMIN_S},         // R
		{"0010100 rs2 rs1 001 rd 1010011 FMAX.S", daTypeRc, emu32_FMAX_S, emu64_FMAX_S},         // R
		{"1100000 00000 rs1 rm rd 1010011 FCVT.W.S", daNone, emu32_FCVT_W_S, emu64_FCVT_W_S},    // R
		{"1100000 00001 rs1 rm rd 1010011 FCVT.WU.S", daNone, emu32_FCVT_WU_S, emu64_FCVT_WU_S}, // R
		{"1110000 00000 rs1 000 rd 1010011 FMV.X.W", daTypeRd, emu32_FMV_X_W, emu64_FMV_X_W},    // R
		{"1010000 rs2 rs1 010 rd 1010011 FEQ.S", daTypeRa, emu32_FEQ_S, emu64_FEQ_S},            // R
		{"1010000 rs2 rs1 001 rd 1010011 FLT.S", daTypeRa, emu32_FLT_S, emu64_FLT_S},            // R
		{"1010000 rs2 rs1 000 rd 1010011 FLE.S", daTypeRa, emu32_FLE_S, emu64_FLE_S},            // R
		{"1110000 00000 rs1 001 rd 1010011 FCLASS.S", daNone, emu32_FCLASS_S, emu64_FCLASS_S},   // R
		{"1101000 00000 rs1 rm rd 1010011 FCVT.S.W", daNone, emu32_FCVT_S_W, emu64_FCVT_S_W},    // R
		{"1101000 00001 rs1 rm rd 1010011 FCVT.S.WU", daNone, emu32_FCVT_S_WU, emu64_FCVT_S_WU}, // R
		{"1111000 00000 rs1 000 rd 1010011 FMV.W.X", daTypeRe, emu32_FMV_W_X, emu64_FMV_W_X},    // R
	},
}

// ISArv32d Double-Precision Floating-Point
var ISArv32d = ISAModule{
	name: "rv32d",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 011 rd 0000111 FLD", daTypeIg, emu32_FLD, emu64_FLD},                    // I
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100111 FSD", daTypeSb, emu32_FSD, emu64_FSD},          // S
		{"rs3 01 rs2 rs1 rm rd 1000011 FMADD.D", daTypeR4a, emu32_FMADD_D, emu64_FMADD_D},       // R4
		{"rs3 01 rs2 rs1 rm rd 1000111 FMSUB.D", daTypeR4a, emu32_FMSUB_D, emu64_FMSUB_D},       // R4
		{"rs3 01 rs2 rs1 rm rd 1001011 FNMSUB.D", daTypeR4a, emu32_FNMSUB_D, emu64_FNMSUB_D},    // R4
		{"rs3 01 rs2 rs1 rm rd 1001111 FNMADD.D", daTypeR4a, emu32_FNMADD_D, emu64_FNMADD_D},    // R4
		{"0000001 rs2 rs1 rm rd 1010011 FADD.D", daNone, emu32_FADD_D, emu64_FADD_D},            // R
		{"0000101 rs2 rs1 rm rd 1010011 FSUB.D", daNone, emu32_FSUB_D, emu64_FSUB_D},            // R
		{"0001001 rs2 rs1 rm rd 1010011 FMUL.D", daNone, emu32_FMUL_D, emu64_FMUL_D},            // R
		{"0001101 rs2 rs1 rm rd 1010011 FDIV.D", daNone, emu32_FDIV_D, emu64_FDIV_D},            // R
		{"0101101 00000 rs1 rm rd 1010011 FSQRT.D", daNone, emu32_FSQRT_D, emu64_FSQRT_D},       // R
		{"0010001 rs2 rs1 000 rd 1010011 FSGNJ.D", daNone, emu32_FSGNJ_D, emu64_FSGNJ_D},        // R
		{"0010001 rs2 rs1 001 rd 1010011 FSGNJN.D", daNone, emu32_FSGNJN_D, emu64_FSGNJN_D},     // R
		{"0010001 rs2 rs1 010 rd 1010011 FSGNJX.D", daNone, emu32_FSGNJX_D, emu64_FSGNJX_D},     // R
		{"0010101 rs2 rs1 000 rd 1010011 FMIN.D", daNone, emu32_FMIN_D, emu64_FMIN_D},           // R
		{"0010101 rs2 rs1 001 rd 1010011 FMAX.D", daNone, emu32_FMAX_D, emu64_FMAX_D},           // R
		{"0100000 00001 rs1 rm rd 1010011 FCVT.S.D", daNone, emu32_FCVT_S_D, emu64_FCVT_S_D},    // R
		{"0100001 00000 rs1 rm rd 1010011 FCVT.D.S", daNone, emu32_FCVT_D_S, emu64_FCVT_D_S},    // R
		{"1010001 rs2 rs1 010 rd 1010011 FEQ.D", daNone, emu32_FEQ_D, emu64_FEQ_D},              // R
		{"1010001 rs2 rs1 001 rd 1010011 FLT.D", daNone, emu32_FLT_D, emu64_FLT_D},              // R
		{"1010001 rs2 rs1 000 rd 1010011 FLE.D", daNone, emu32_FLE_D, emu64_FLE_D},              // R
		{"1110001 00000 rs1 001 rd 1010011 FCLASS.D", daNone, emu32_FCLASS_D, emu64_FCLASS_D},   // R
		{"1100001 00000 rs1 rm rd 1010011 FCVT.W.D", daNone, emu32_FCVT_W_D, emu64_FCVT_W_D},    // R
		{"1100001 00001 rs1 rm rd 1010011 FCVT.WU.D", daNone, emu32_FCVT_WU_D, emu64_FCVT_WU_D}, // R
		{"1101001 00000 rs1 rm rd 1010011 FCVT.D.W", daNone, emu32_FCVT_D_W, emu64_FCVT_D_W},    // R
		{"1101001 00001 rs1 rm rd 1010011 FCVT.D.WU", daNone, emu32_FCVT_D_WU, emu64_FCVT_D_WU}, // R
	},
}

// ISArv32c Compressed
var ISArv32c = ISAModule{
	name: "rv32c",
	ilen: 16,
	defn: []insDefn{
		{"000 00000000 000 00 C.ILLEGAL", daTypeCIWa, emu32_C_ILLEGAL, emu64_C_ILLEGAL},                      // CIW (Quadrant 0)
		{"000 nzuimm[5:4|9:6|2|3] rd0 00 C.ADDI4SPN", daTypeCIWb, emu32_C_ADDI4SPN, emu64_C_ADDI4SPN},        // CIW
		{"001 uimm[5:3] rs10 uimm[7:6] rd0 00 C.FLD", daNone, emu32_C_FLD, emu64_C_FLD},                      // CL
		{"010 uimm[5:3] rs10 uimm[2|6] rd0 00 C.LW", daNone, emu32_C_LW, emu64_C_LW},                         // CL
		{"011 uimm[5:3] rs10 uimm[2|6] rd0 00 C.FLW", daNone, emu32_C_FLW, emu64_C_FLW},                      // CL
		{"101 uimm[5:3] rs10 uimm[7:6] rs20 00 C.FSD", daNone, emu32_C_FSD, emu64_C_FSD},                     // CS
		{"110 uimm[5:3] rs10 uimm[2|6] rs20 00 C.SW", daNone, emu32_C_SW, emu64_C_SW},                        // CS
		{"111 uimm[5:3] rs10 uimm[2|6] rs20 00 C.FSW", daNone, emu32_C_FSW, emu64_C_FSW},                     // CS
		{"000 nzimm[5] 00000 nzimm[4:0] 01 C.NOP", daNone, emu32_C_NOP, emu64_C_NOP},                         // CI (Quadrant 1)
		{"000 nzimm[5] rs1/rd!=0 nzimm[4:0] 01 C.ADDI", daTypeCIc, emu32_C_ADDI, emu64_C_ADDI},               // CI
		{"001 imm[11|4|9:8|10|6|7|3:1|5] 01 C.JAL", daTypeCJc, emu32_C_JAL, emu64_C_JAL},                     // CJ
		{"010 imm[5] rd!=0 imm[4:0] 01 C.LI", daTypeCIa, emu32_C_LI, emu64_C_LI},                             // CI
		{"011 nzimm[9] 00010 nzimm[4|6|8:7|5] 01 C.ADDI16SP", daTypeCIb, emu32_C_ADDI16SP, emu64_C_ADDI16SP}, // CI
		{"011 nzimm[17] rd!={0,2} nzimm[16:12] 01 C.LUI", daTypeCIg, emu32_C_LUI, emu64_C_LUI},               // CI
		{"100 nzuimm[5] 00 rs10/rd0 nzuimm[4:0] 01 C.SRLI", daTypeCId, emu32_C_SRLI, emu64_C_SRLI},           // CI
		{"100 nzuimm[5] 01 rs10/rd0 nzuimm[4:0] 01 C.SRAI", daTypeCId, emu32_C_SRAI, emu64_C_SRAI},           // CI
		{"100 imm[5] 10 rs10/rd0 imm[4:0] 01 C.ANDI", daTypeCIf, emu32_C_ANDI, emu64_C_ANDI},                 // CI
		{"100 0 11 rs10/rd0 00 rs20 01 C.SUB", daNone, emu32_C_SUB, emu64_C_SUB},                             // CR
		{"100 0 11 rs10/rd0 01 rs20 01 C.XOR", daNone, emu32_C_XOR, emu64_C_XOR},                             // CR
		{"100 0 11 rs10/rd0 10 rs20 01 C.OR", daNone, emu32_C_OR, emu64_C_OR},                                // CR
		{"100 0 11 rs10/rd0 11 rs20 01 C.AND", daNone, emu32_C_AND, emu64_C_AND},                             // CR
		{"101 imm[11|4|9:8|10|6|7|3:1|5] 01 C.J", daTypeCJb, emu32_C_J, emu64_C_J},                           // CJ
		{"110 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BEQZ", daTypeCBa, emu32_C_BEQZ, emu64_C_BEQZ},              // CB
		{"111 imm[8|4:3] rs10 imm[7:6|2:1|5] 01 C.BNEZ", daTypeCBa, emu32_C_BNEZ, emu64_C_BNEZ},              // CB
		{"000 nzuimm[5] rs1/rd!=0 nzuimm[4:0] 10 C.SLLI", daTypeCIe, emu32_C_SLLI, emu64_C_SLLI},             // CI (Quadrant 2)
		{"000 0 rs1/rd!=0 00000 10 C.SLLI64", daNone, emu32_C_SLLI64, emu64_C_SLLI64},                        // CI
		{"001 uimm[5] rd uimm[4:3|8:6] 10 C.FLDSP", daNone, emu32_C_FLDSP, emu64_C_FLDSP},                    // CSS
		{"010 uimm[5] rd!=0 uimm[4:2|7:6] 10 C.LWSP", daTypeCSSa, emu32_C_LWSP, emu64_C_LWSP},                // CSS
		{"011 uimm[5] rd uimm[4:2|7:6] 10 C.FLWSP", daNone, emu32_C_FLWSP, emu64_C_FLWSP},                    // CSS
		{"100 0 rs1!=0 00000 10 C.JR", daTypeCJa, emu32_C_JR, emu64_C_JR},                                    // CJ
		{"100 0 rd!=0 rs2!=0 10 C.MV", daTypeCRa, emu32_C_MV, emu64_C_MV},                                    // CR
		{"100 1 00000 00000 10 C.EBREAK", daNone, emu32_C_EBREAK, emu64_C_EBREAK},                            // CI
		{"100 1 rs1!=0 00000 10 C.JALR", daNone, emu32_C_JALR, emu64_C_JALR},                                 // CJ
		{"100 1 rs1/rd!=0 rs2!=0 10 C.ADD", daTypeCRb, emu32_C_ADD, emu64_C_ADD},                             // CR
		{"101 uimm[5:3|8:6] rs2 10 C.FSDSP", daNone, emu32_C_FSDSP, emu64_C_FSDSP},                           // CSS
		{"110 uimm[5:2|7:6] rs2 10 C.SWSP", daTypeCSSb, emu32_C_SWSP, emu64_C_SWSP},                          // CSS
		{"111 uimm[5:2|7:6] rs2 10 C.FSWSP", daNone, emu32_C_FSWSP, emu64_C_FSWSP},                           // CSS
	},
}

//-----------------------------------------------------------------------------
// RV64 instructions (+ RV32)

// ISArv64i Integer
var ISArv64i = ISAModule{
	name: "rv64i",
	ilen: 32,
	defn: []insDefn{
		{"imm[11:0] rs1 110 rd 0000011 LWU", daTypeIa, emu32_Illegal, emu64_LWU},         // I
		{"imm[11:0] rs1 011 rd 0000011 LD", daTypeIa, emu32_Illegal, emu64_LD},           // I
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD", daTypeSa, emu32_Illegal, emu64_SD}, // S
		{"000000 shamt6 rs1 001 rd 0010011 SLLI", daNone, emu32_Illegal, emu64_SLLI},     // I
		{"000000 shamt6 rs1 101 rd 0010011 SRLI", daNone, emu32_Illegal, emu64_SRLI},     // I
		{"010000 shamt6 rs1 101 rd 0010011 SRAI", daNone, emu32_Illegal, emu64_SRAI},     // I
		{"imm[11:0] rs1 000 rd 0011011 ADDIW", daTypeIa, emu32_Illegal, emu64_ADDIW},     // I
		{"0000000 shamt5 rs1 001 rd 0011011 SLLIW", daNone, emu32_Illegal, emu64_SLLIW},  // I
		{"0000000 shamt5 rs1 101 rd 0011011 SRLIW", daNone, emu32_Illegal, emu64_SRLIW},  // I
		{"0100000 shamt5 rs1 101 rd 0011011 SRAIW", daNone, emu32_Illegal, emu64_SRAIW},  // I
		{"0000000 rs2 rs1 000 rd 0111011 ADDW", daNone, emu32_Illegal, emu64_ADDW},       // R
		{"0100000 rs2 rs1 000 rd 0111011 SUBW", daNone, emu32_Illegal, emu64_SUBW},       // R
		{"0000000 rs2 rs1 001 rd 0111011 SLLW", daNone, emu32_Illegal, emu64_SLLW},       // R
		{"0000000 rs2 rs1 101 rd 0111011 SRLW", daNone, emu32_Illegal, emu64_SRLW},       // R
		{"0100000 rs2 rs1 101 rd 0111011 SRAW", daNone, emu32_Illegal, emu64_SRAW},       // R
	},
}

// ISArv64m Integer Multiplication and Division
var ISArv64m = ISAModule{
	name: "rv64m",
	ilen: 32,
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0111011 MULW", daNone, emu32_Illegal, emu64_MULW},   // R
		{"0000001 rs2 rs1 100 rd 0111011 DIVW", daNone, emu32_Illegal, emu64_DIVW},   // R
		{"0000001 rs2 rs1 101 rd 0111011 DIVUW", daNone, emu32_Illegal, emu64_DIVUW}, // R
		{"0000001 rs2 rs1 110 rd 0111011 REMW", daNone, emu32_Illegal, emu64_REMW},   // R
		{"0000001 rs2 rs1 111 rd 0111011 REMUW", daNone, emu32_Illegal, emu64_REMUW}, // R
	},
}

// ISArv64a Atomics
var ISArv64a = ISAModule{
	name: "rv64a",
	ilen: 32,
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 011 rd 0101111 LR.D", daNone, emu32_Illegal, emu64_LR_D},         // R
		{"00011 aq rl rs2 rs1 011 rd 0101111 SC.D", daNone, emu32_Illegal, emu64_SC_D},           // R
		{"00001 aq rl rs2 rs1 011 rd 0101111 AMOSWAP.D", daNone, emu32_Illegal, emu64_AMOSWAP_D}, // R
		{"00000 aq rl rs2 rs1 011 rd 0101111 AMOADD.D", daNone, emu32_Illegal, emu64_AMOADD_D},   // R
		{"00100 aq rl rs2 rs1 011 rd 0101111 AMOXOR.D", daNone, emu32_Illegal, emu64_AMOXOR_D},   // R
		{"01100 aq rl rs2 rs1 011 rd 0101111 AMOAND.D", daNone, emu32_Illegal, emu64_AMOAND_D},   // R
		{"01000 aq rl rs2 rs1 011 rd 0101111 AMOOR.D", daNone, emu32_Illegal, emu64_AMOOR_D},     // R
		{"10000 aq rl rs2 rs1 011 rd 0101111 AMOMIN.D", daNone, emu32_Illegal, emu64_AMOMIN_D},   // R
		{"10100 aq rl rs2 rs1 011 rd 0101111 AMOMAX.D", daNone, emu32_Illegal, emu64_AMOMAX_D},   // R
		{"11000 aq rl rs2 rs1 011 rd 0101111 AMOMINU.D", daNone, emu32_Illegal, emu64_AMOMINU_D}, // R
		{"11100 aq rl rs2 rs1 011 rd 0101111 AMOMAXU.D", daNone, emu32_Illegal, emu64_AMOMAXU_D}, // R
	},
}

// ISArv64f Single-Precision Floating-Point
var ISArv64f = ISAModule{
	name: "rv64f",
	ilen: 32,
	defn: []insDefn{
		{"1100000 00010 rs1 rm rd 1010011 FCVT.L.S", daNone, emu32_Illegal, emu64_FCVT_L_S},   // R
		{"1100000 00011 rs1 rm rd 1010011 FCVT.LU.S", daNone, emu32_Illegal, emu64_FCVT_LU_S}, // R
		{"1101000 00010 rs1 rm rd 1010011 FCVT.S.L", daNone, emu32_Illegal, emu64_FCVT_S_L},   // R
		{"1101000 00011 rs1 rm rd 1010011 FCVT.S.LU", daNone, emu32_Illegal, emu64_FCVT_S_LU}, // R
	},
}

// ISArv64d Double-Precision Floating-Point
var ISArv64d = ISAModule{
	name: "rv64d",
	ilen: 32,
	defn: []insDefn{
		{"1100001 00010 rs1 rm rd 1010011 FCVT.L.D", daNone, emu32_Illegal, emu64_FCVT_L_D},   // R
		{"1100001 00011 rs1 rm rd 1010011 FCVT.LU.D", daNone, emu32_Illegal, emu64_FCVT_LU_D}, // R
		{"1110001 00000 rs1 000 rd 1010011 FMV.X.D", daNone, emu32_Illegal, emu64_FMV_X_D},    // R
		{"1101001 00010 rs1 rm rd 1010011 FCVT.D.L", daNone, emu32_Illegal, emu64_FCVT_D_L},   // R
		{"1101001 00011 rs1 rm rd 1010011 FCVT.D.LU", daNone, emu32_Illegal, emu64_FCVT_D_LU}, // R
		{"1111001 00000 rs1 000 rd 1010011 FMV.D.X", daNone, emu32_Illegal, emu64_FMV_D_X},    // R
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
