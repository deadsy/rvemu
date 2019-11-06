//-----------------------------------------------------------------------------
/*

RISC-V ISA Definition

*/
//-----------------------------------------------------------------------------

package rv

//-----------------------------------------------------------------------------

// daFunc is the disassembler function for a specific instruction.
type daFunc func(name string, adr, ins uint32) (string, string)

type insDefn struct {
	defn string     // instruction definition string (from the standard)
	dt   decodeType // decode type
	da   daFunc
}

type ISAModule struct {
	name string    // name of module
	defn []insDefn // instruction definitions
}

//-----------------------------------------------------------------------------
// RV32 instructions

var ISArv32i = ISAModule{
	name: "rv32i",
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

var ISArv32m = ISAModule{
	name: "rv32m",
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

var ISArv32a = ISAModule{
	name: "rv32a",
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 010 rd 0101111 LR.W", decodeTypeR, daNone},
		{"00011 aq rl rs2 rs1 010 rd 0101111 SC.W", decodeTypeR, daNone},
		{"00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W", decodeTypeR, daNone},
		{"00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W", decodeTypeR, daNone},
		{"00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W", decodeTypeR, daNone},
		{"01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W", decodeTypeR, daNone},
		{"01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W", decodeTypeR, daNone},
		{"10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W", decodeTypeR, daNone},
		{"10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W", decodeTypeR, daNone},
		{"11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W", decodeTypeR, daNone},
		{"11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W", decodeTypeR, daNone},
	},
}

var ISArv32f = ISAModule{
	name: "rv32f",
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

var ISArv32d = ISAModule{
	name: "rv32d",
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

//-----------------------------------------------------------------------------
// RV64 instructions (+ RV32)

var ISArv64i = ISAModule{
	name: "rv64i",
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

var ISArv64m = ISAModule{
	name: "rv64m",
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0111011 MULW", decodeTypeR, daNone},
		{"0000001 rs2 rs1 100 rd 0111011 DIVW", decodeTypeR, daNone},
		{"0000001 rs2 rs1 101 rd 0111011 DIVUW", decodeTypeR, daNone},
		{"0000001 rs2 rs1 110 rd 0111011 REMW", decodeTypeR, daNone},
		{"0000001 rs2 rs1 111 rd 0111011 REMUW", decodeTypeR, daNone},
	},
}

var ISArv64a = ISAModule{
	name: "rv64a",
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

var ISArv64f = ISAModule{
	name: "rv64f",
	defn: []insDefn{
		{"1100000 00010 rs1 rm rd 1010011 FCVT.L.S", decodeTypeR, daNone},
		{"1100000 00011 rs1 rm rd 1010011 FCVT.LU.S", decodeTypeR, daNone},
		{"1101000 00010 rs1 rm rd 1010011 FCVT.S.L", decodeTypeR, daNone},
		{"1101000 00011 rs1 rm rd 1010011 FCVT.S.LU", decodeTypeR, daNone},
	},
}

var ISArv64d = ISAModule{
	name: "rv64d",
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
	val, mask uint32 // value and mask of fixed bits in the instruction
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

// addInstruction adds an instruction to the ISA.
func (isa *ISA) addInstruction(id *insDefn) error {
	ii, err := parseDefn(id)
	if err != nil {
		return err
	}
	isa.instruction = append(isa.instruction, ii)
	return nil
}

// Add a ISA sub-module to the ISA.
func (isa *ISA) Add(module ...ISAModule) error {
	for i := range module {
		for _, id := range module[i].defn {
			err := isa.addInstruction(&id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// lookup returns the instruction information for a given instruction.
func (isa *ISA) lookup(ins uint32) *insInfo {
	for _, ii := range isa.instruction {
		if ins&ii.mask == ii.val {
			return ii
		}
	}
	return nil
}

//-----------------------------------------------------------------------------
