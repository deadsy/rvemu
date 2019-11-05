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
	defn string // instruction definition string (from the standard)
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
		{"imm[31:12] rd 0110111 LUI", daTypeUa},
		{"imm[31:12] rd 0010111 AUIPC", daTypeUa},
		{"imm[20|10:1|11|19:12] rd 1101111 JAL", daNone},
		{"imm[11:0] rs1 000 rd 1100111 JALR", daTypeIa},
		{"imm[12|10:5] rs2 rs1 000 imm[4:1|11] 1100011 BEQ", daNone},
		{"imm[12|10:5] rs2 rs1 001 imm[4:1|11] 1100011 BNE", daNone},
		{"imm[12|10:5] rs2 rs1 100 imm[4:1|11] 1100011 BLT", daNone},
		{"imm[12|10:5] rs2 rs1 101 imm[4:1|11] 1100011 BGE", daNone},
		{"imm[12|10:5] rs2 rs1 110 imm[4:1|11] 1100011 BLTU", daNone},
		{"imm[12|10:5] rs2 rs1 111 imm[4:1|11] 1100011 BGEU", daNone},
		{"imm[11:0] rs1 000 rd 0000011 LB", daTypeIc},
		{"imm[11:0] rs1 001 rd 0000011 LH", daTypeIc},
		{"imm[11:0] rs1 010 rd 0000011 LW", daTypeIc},
		{"imm[11:0] rs1 100 rd 0000011 LBU", daTypeIc},
		{"imm[11:0] rs1 101 rd 0000011 LHU", daTypeIc},
		{"imm[11:5] rs2 rs1 000 imm[4:0] 0100011 SB", daTypeSa},
		{"imm[11:5] rs2 rs1 001 imm[4:0] 0100011 SH", daTypeSa},
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100011 SW", daTypeSa},
		{"imm[11:0] rs1 000 rd 0010011 ADDI", daTypeIb},
		{"imm[11:0] rs1 010 rd 0010011 SLTI", daTypeIa},
		{"imm[11:0] rs1 011 rd 0010011 SLTIU", daTypeIa},
		{"imm[11:0] rs1 100 rd 0010011 XORI", daTypeIa},
		{"imm[11:0] rs1 110 rd 0010011 ORI", daTypeIa},
		{"imm[11:0] rs1 111 rd 0010011 ANDI", daTypeIa},
		{"0000000 shamt5 rs1 001 rd 0010011 SLLI", daNone},
		{"0000000 shamt5 rs1 101 rd 0010011 SRLI", daNone},
		{"0100000 shamt5 rs1 101 rd 0010011 SRAI", daNone},
		{"0000000 rs2 rs1 000 rd 0110011 ADD", daNone},
		{"0100000 rs2 rs1 000 rd 0110011 SUB", daNone},
		{"0000000 rs2 rs1 001 rd 0110011 SLL", daNone},
		{"0000000 rs2 rs1 010 rd 0110011 SLT", daNone},
		{"0000000 rs2 rs1 011 rd 0110011 SLTU", daNone},
		{"0000000 rs2 rs1 100 rd 0110011 XOR", daNone},
		{"0000000 rs2 rs1 101 rd 0110011 SRL", daNone},
		{"0100000 rs2 rs1 101 rd 0110011 SRA", daNone},
		{"0000000 rs2 rs1 110 rd 0110011 OR", daNone},
		{"0000000 rs2 rs1 111 rd 0110011 AND", daNone},
		{"0000 pred succ 00000 000 00000 0001111 FENCE", daNone},
		{"0000 0000 0000 00000 001 00000 0001111 FENCE.I", daNone},
		{"000000000000 00000 000 00000 1110011 ECALL", daNone},
		{"000000000001 00000 000 00000 1110011 EBREAK", daNone},
		{"csr rs1 001 rd 1110011 CSRRW", daNone},
		{"csr rs1 010 rd 1110011 CSRRS", daNone},
		{"csr rs1 011 rd 1110011 CSRRC", daNone},
		{"csr zimm 101 rd 1110011 CSRRWI", daNone},
		{"csr zimm 110 rd 1110011 CSRRSI", daNone},
		{"csr zimm 111 rd 1110011 CSRRCI", daNone},
	},
}

var ISArv32m = ISAModule{
	name: "rv32m",
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0110011 MUL", daNone},
		{"0000001 rs2 rs1 001 rd 0110011 MULH", daNone},
		{"0000001 rs2 rs1 010 rd 0110011 MULHSU", daNone},
		{"0000001 rs2 rs1 011 rd 0110011 MULHU", daNone},
		{"0000001 rs2 rs1 100 rd 0110011 DIV", daNone},
		{"0000001 rs2 rs1 101 rd 0110011 DIVU", daNone},
		{"0000001 rs2 rs1 110 rd 0110011 REM", daNone},
		{"0000001 rs2 rs1 111 rd 0110011 REMU", daNone},
	},
}

var ISArv32a = ISAModule{
	name: "rv32a",
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 010 rd 0101111 LR.W", daNone},
		{"00011 aq rl rs2 rs1 010 rd 0101111 SC.W", daNone},
		{"00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W", daNone},
		{"00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W", daNone},
		{"00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W", daNone},
		{"01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W", daNone},
		{"01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W", daNone},
		{"10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W", daNone},
		{"10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W", daNone},
		{"11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W", daNone},
		{"11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W", daNone},
	},
}

var ISArv32f = ISAModule{
	name: "rv32f",
	defn: []insDefn{
		{"imm[11:0] rs1 010 rd 0000111 FLW", daTypeIa},
		{"imm[11:5] rs2 rs1 010 imm[4:0] 0100111 FSW", daTypeSa},
		{"rs3 00 rs2 rs1 rm rd 1000011 FMADD.S", daNone},
		{"rs3 00 rs2 rs1 rm rd 1000111 FMSUB.S", daNone},
		{"rs3 00 rs2 rs1 rm rd 1001011 FNMSUB.S", daNone},
		{"rs3 00 rs2 rs1 rm rd 1001111 FNMADD.S", daNone},
		{"0000000 rs2 rs1 rm rd 1010011 FADD.S", daNone},
		{"0000100 rs2 rs1 rm rd 1010011 FSUB.S", daNone},
		{"0001000 rs2 rs1 rm rd 1010011 FMUL.S", daNone},
		{"0001100 rs2 rs1 rm rd 1010011 FDIV.S", daNone},
		{"0101100 00000 rs1 rm rd 1010011 FSQRT.S", daNone},
		{"0010000 rs2 rs1 000 rd 1010011 FSGNJ.S", daNone},
		{"0010000 rs2 rs1 001 rd 1010011 FSGNJN.S", daNone},
		{"0010000 rs2 rs1 010 rd 1010011 FSGNJX.S", daNone},
		{"0010100 rs2 rs1 000 rd 1010011 FMIN.S", daNone},
		{"0010100 rs2 rs1 001 rd 1010011 FMAX.S", daNone},
		{"1100000 00000 rs1 rm rd 1010011 FCVT.W.S", daNone},
		{"1100000 00001 rs1 rm rd 1010011 FCVT.WU.S", daNone},
		{"1110000 00000 rs1 000 rd 1010011 FMV.X.W", daNone},
		{"1010000 rs2 rs1 010 rd 1010011 FEQ.S", daNone},
		{"1010000 rs2 rs1 001 rd 1010011 FLT.S", daNone},
		{"1010000 rs2 rs1 000 rd 1010011 FLE.S", daNone},
		{"1110000 00000 rs1 001 rd 1010011 FCLASS.S", daNone},
		{"1101000 00000 rs1 rm rd 1010011 FCVT.S.W", daNone},
		{"1101000 00001 rs1 rm rd 1010011 FCVT.S.WU", daNone},
		{"1111000 00000 rs1 000 rd 1010011 FMV.W.X", daNone},
	},
}

var ISArv32d = ISAModule{
	name: "rv32d",
	defn: []insDefn{
		{"imm[11:0] rs1 011 rd 0000111 FLD", daTypeIa},
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100111 FSD", daTypeSa},
		{"rs3 01 rs2 rs1 rm rd 1000011 FMADD.D", daNone},
		{"rs3 01 rs2 rs1 rm rd 1000111 FMSUB.D", daNone},
		{"rs3 01 rs2 rs1 rm rd 1001011 FNMSUB.D", daNone},
		{"rs3 01 rs2 rs1 rm rd 1001111 FNMADD.D", daNone},
		{"0000001 rs2 rs1 rm rd 1010011 FADD.D", daNone},
		{"0000101 rs2 rs1 rm rd 1010011 FSUB.D", daNone},
		{"0001001 rs2 rs1 rm rd 1010011 FMUL.D", daNone},
		{"0001101 rs2 rs1 rm rd 1010011 FDIV.D", daNone},
		{"0101101 00000 rs1 rm rd 1010011 FSQRT.D", daNone},
		{"0010001 rs2 rs1 000 rd 1010011 FSGNJ.D", daNone},
		{"0010001 rs2 rs1 001 rd 1010011 FSGNJN.D", daNone},
		{"0010001 rs2 rs1 010 rd 1010011 FSGNJX.D", daNone},
		{"0010101 rs2 rs1 000 rd 1010011 FMIN.D", daNone},
		{"0010101 rs2 rs1 001 rd 1010011 FMAX.D", daNone},
		{"0100000 00001 rs1 rm rd 1010011 FCVT.S.D", daNone},
		{"0100001 00000 rs1 rm rd 1010011 FCVT.D.S", daNone},
		{"1010001 rs2 rs1 010 rd 1010011 FEQ.D", daNone},
		{"1010001 rs2 rs1 001 rd 1010011 FLT.D", daNone},
		{"1010001 rs2 rs1 000 rd 1010011 FLE.D", daNone},
		{"1110001 00000 rs1 001 rd 1010011 FCLASS.D", daNone},
		{"1100001 00000 rs1 rm rd 1010011 FCVT.W.D", daNone},
		{"1100001 00001 rs1 rm rd 1010011 FCVT.WU.D", daNone},
		{"1101001 00000 rs1 rm rd 1010011 FCVT.D.W", daNone},
		{"1101001 00001 rs1 rm rd 1010011 FCVT.D.WU", daNone},
	},
}

//-----------------------------------------------------------------------------
// RV64 instructions (+ RV32)

var ISArv64i = ISAModule{
	name: "rv64i",
	defn: []insDefn{
		{"imm[11:0] rs1 110 rd 0000011 LWU", daTypeIa},
		{"imm[11:0] rs1 011 rd 0000011 LD", daTypeIa},
		{"imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD", daTypeSa},
		{"000000 shamt6 rs1 001 rd 0010011 SLLI", daNone},
		{"000000 shamt6 rs1 101 rd 0010011 SRLI", daNone},
		{"010000 shamt6 rs1 101 rd 0010011 SRAI", daNone},
		{"imm[11:0] rs1 000 rd 0011011 ADDIW", daTypeIa},
		{"0000000 shamt5 rs1 001 rd 0011011 SLLIW", daNone},
		{"0000000 shamt5 rs1 101 rd 0011011 SRLIW", daNone},
		{"0100000 shamt5 rs1 101 rd 0011011 SRAIW", daNone},
		{"0000000 rs2 rs1 000 rd 0111011 ADDW", daNone},
		{"0100000 rs2 rs1 000 rd 0111011 SUBW", daNone},
		{"0000000 rs2 rs1 001 rd 0111011 SLLW", daNone},
		{"0000000 rs2 rs1 101 rd 0111011 SRLW", daNone},
		{"0100000 rs2 rs1 101 rd 0111011 SRAW", daNone},
	},
}

var ISArv64m = ISAModule{
	name: "rv64m",
	defn: []insDefn{
		{"0000001 rs2 rs1 000 rd 0111011 MULW", daNone},
		{"0000001 rs2 rs1 100 rd 0111011 DIVW", daNone},
		{"0000001 rs2 rs1 101 rd 0111011 DIVUW", daNone},
		{"0000001 rs2 rs1 110 rd 0111011 REMW", daNone},
		{"0000001 rs2 rs1 111 rd 0111011 REMUW", daNone},
	},
}

var ISArv64a = ISAModule{
	name: "rv64a",
	defn: []insDefn{
		{"00010 aq rl 00000 rs1 011 rd 0101111 LR.D", daNone},
		{"00011 aq rl rs2 rs1 011 rd 0101111 SC.D", daNone},
		{"00001 aq rl rs2 rs1 011 rd 0101111 AMOSWAP.D", daNone},
		{"00000 aq rl rs2 rs1 011 rd 0101111 AMOADD.D", daNone},
		{"00100 aq rl rs2 rs1 011 rd 0101111 AMOXOR.D", daNone},
		{"01100 aq rl rs2 rs1 011 rd 0101111 AMOAND.D", daNone},
		{"01000 aq rl rs2 rs1 011 rd 0101111 AMOOR.D", daNone},
		{"10000 aq rl rs2 rs1 011 rd 0101111 AMOMIN.D", daNone},
		{"10100 aq rl rs2 rs1 011 rd 0101111 AMOMAX.D", daNone},
		{"11000 aq rl rs2 rs1 011 rd 0101111 AMOMINU.D", daNone},
		{"11100 aq rl rs2 rs1 011 rd 0101111 AMOMAXU.D", daNone},
	},
}

var ISArv64f = ISAModule{
	name: "rv64f",
	defn: []insDefn{
		{"1100000 00010 rs1 rm rd 1010011 FCVT.L.S", daNone},
		{"1100000 00011 rs1 rm rd 1010011 FCVT.LU.S", daNone},
		{"1101000 00010 rs1 rm rd 1010011 FCVT.S.L", daNone},
		{"1101000 00011 rs1 rm rd 1010011 FCVT.S.LU", daNone},
	},
}

var ISArv64d = ISAModule{
	name: "rv64d",
	defn: []insDefn{
		{"1100001 00010 rs1 rm rd 1010011 FCVT.L.D", daNone},
		{"1100001 00011 rs1 rm rd 1010011 FCVT.LU.D", daNone},
		{"1110001 00000 rs1 000 rd 1010011 FMV.X.D", daNone},
		{"1101001 00010 rs1 rm rd 1010011 FCVT.D.L", daNone},
		{"1101001 00011 rs1 rm rd 1010011 FCVT.D.LU", daNone},
		{"1111001 00000 rs1 000 rd 1010011 FMV.D.X", daNone},
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
