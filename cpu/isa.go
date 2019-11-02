//-----------------------------------------------------------------------------
/*

RISC-V ISA Definition

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

type ISAModule struct {
	name string   // name of module
	defn []string // instruction definition strings (from the standard)
}

//-----------------------------------------------------------------------------
// RV32 instructions

var ISArv32i = ISAModule{
	name: "rv32i",
	defn: []string{
		"imm[31:12] rd 0110111 LUI",
		"imm[31:12] rd 0010111 AUIPC",
		"imm[20|10:1|11|19:12] rd 1101111 JAL",
		"imm[11:0] rs1 000 rd 1100111 JALR",
		"imm[12|10:5] rs2 rs1 000 imm[4:1|11] 1100011 BEQ",
		"imm[12|10:5] rs2 rs1 001 imm[4:1|11] 1100011 BNE",
		"imm[12|10:5] rs2 rs1 100 imm[4:1|11] 1100011 BLT",
		"imm[12|10:5] rs2 rs1 101 imm[4:1|11] 1100011 BGE",
		"imm[12|10:5] rs2 rs1 110 imm[4:1|11] 1100011 BLTU",
		"imm[12|10:5] rs2 rs1 111 imm[4:1|11] 1100011 BGEU",
		"imm[11:0] rs1 000 rd 0000011 LB",
		"imm[11:0] rs1 001 rd 0000011 LH",
		"imm[11:0] rs1 010 rd 0000011 LW",
		"imm[11:0] rs1 100 rd 0000011 LBU",
		"imm[11:0] rs1 101 rd 0000011 LHU",
		"imm[11:5] rs2 rs1 000 imm[4:0] 0100011 SB",
		"imm[11:5] rs2 rs1 001 imm[4:0] 0100011 SH",
		"imm[11:5] rs2 rs1 010 imm[4:0] 0100011 SW",
		"imm[11:0] rs1 000 rd 0010011 ADDI",
		"imm[11:0] rs1 010 rd 0010011 SLTI",
		"imm[11:0] rs1 011 rd 0010011 SLTIU",
		"imm[11:0] rs1 100 rd 0010011 XORI",
		"imm[11:0] rs1 110 rd 0010011 ORI",
		"imm[11:0] rs1 111 rd 0010011 ANDI",
		"0000000 shamt5 rs1 001 rd 0010011 SLLI",
		"0000000 shamt5 rs1 101 rd 0010011 SRLI",
		"0100000 shamt5 rs1 101 rd 0010011 SRAI",
		"0000000 rs2 rs1 000 rd 0110011 ADD",
		"0100000 rs2 rs1 000 rd 0110011 SUB",
		"0000000 rs2 rs1 001 rd 0110011 SLL",
		"0000000 rs2 rs1 010 rd 0110011 SLT",
		"0000000 rs2 rs1 011 rd 0110011 SLTU",
		"0000000 rs2 rs1 100 rd 0110011 XOR",
		"0000000 rs2 rs1 101 rd 0110011 SRL",
		"0100000 rs2 rs1 101 rd 0110011 SRA",
		"0000000 rs2 rs1 110 rd 0110011 OR",
		"0000000 rs2 rs1 111 rd 0110011 AND",
		"0000 pred succ 00000 000 00000 0001111 FENCE",
		"0000 0000 0000 00000 001 00000 0001111 FENCE.I",
		"000000000000 00000 000 00000 1110011 ECALL",
		"000000000001 00000 000 00000 1110011 EBREAK",
		"csr rs1 001 rd 1110011 CSRRW",
		"csr rs1 010 rd 1110011 CSRRS",
		"csr rs1 011 rd 1110011 CSRRC",
		"csr zimm 101 rd 1110011 CSRRWI",
		"csr zimm 110 rd 1110011 CSRRSI",
		"csr zimm 111 rd 1110011 CSRRCI",
	},
}

var ISArv32m = ISAModule{
	name: "rv32m",
	defn: []string{
		"0000001 rs2 rs1 000 rd 0110011 MUL",
		"0000001 rs2 rs1 001 rd 0110011 MULH",
		"0000001 rs2 rs1 010 rd 0110011 MULHSU",
		"0000001 rs2 rs1 011 rd 0110011 MULHU",
		"0000001 rs2 rs1 100 rd 0110011 DIV",
		"0000001 rs2 rs1 101 rd 0110011 DIVU",
		"0000001 rs2 rs1 110 rd 0110011 REM",
		"0000001 rs2 rs1 111 rd 0110011 REMU",
	},
}

var ISArv32a = ISAModule{
	name: "rv32a",
	defn: []string{
		"00010 aq rl 00000 rs1 010 rd 0101111 LR.W",
		"00011 aq rl rs2 rs1 010 rd 0101111 SC.W",
		"00001 aq rl rs2 rs1 010 rd 0101111 AMOSWAP.W",
		"00000 aq rl rs2 rs1 010 rd 0101111 AMOADD.W",
		"00100 aq rl rs2 rs1 010 rd 0101111 AMOXOR.W",
		"01100 aq rl rs2 rs1 010 rd 0101111 AMOAND.W",
		"01000 aq rl rs2 rs1 010 rd 0101111 AMOOR.W",
		"10000 aq rl rs2 rs1 010 rd 0101111 AMOMIN.W",
		"10100 aq rl rs2 rs1 010 rd 0101111 AMOMAX.W",
		"11000 aq rl rs2 rs1 010 rd 0101111 AMOMINU.W",
		"11100 aq rl rs2 rs1 010 rd 0101111 AMOMAXU.W",
	},
}

var ISArv32f = ISAModule{
	name: "rv32f",
	defn: []string{
		"imm[11:0] rs1 010 rd 0000111 FLW",
		"imm[11:5] rs2 rs1 010 imm[4:0] 0100111 FSW",
		"rs3 00 rs2 rs1 rm rd 1000011 FMADD.S",
		"rs3 00 rs2 rs1 rm rd 1000111 FMSUB.S",
		"rs3 00 rs2 rs1 rm rd 1001011 FNMSUB.S",
		"rs3 00 rs2 rs1 rm rd 1001111 FNMADD.S",
		"0000000 rs2 rs1 rm rd 1010011 FADD.S",
		"0000100 rs2 rs1 rm rd 1010011 FSUB.S",
		"0001000 rs2 rs1 rm rd 1010011 FMUL.S",
		"0001100 rs2 rs1 rm rd 1010011 FDIV.S",
		"0101100 00000 rs1 rm rd 1010011 FSQRT.S",
		"0010000 rs2 rs1 000 rd 1010011 FSGNJ.S",
		"0010000 rs2 rs1 001 rd 1010011 FSGNJN.S",
		"0010000 rs2 rs1 010 rd 1010011 FSGNJX.S",
		"0010100 rs2 rs1 000 rd 1010011 FMIN.S",
		"0010100 rs2 rs1 001 rd 1010011 FMAX.S",
		"1100000 00000 rs1 rm rd 1010011 FCVT.W.S",
		"1100000 00001 rs1 rm rd 1010011 FCVT.WU.S",
		"1110000 00000 rs1 000 rd 1010011 FMV.X.W",
		"1010000 rs2 rs1 010 rd 1010011 FEQ.S",
		"1010000 rs2 rs1 001 rd 1010011 FLT.S",
		"1010000 rs2 rs1 000 rd 1010011 FLE.S",
		"1110000 00000 rs1 001 rd 1010011 FCLASS.S",
		"1101000 00000 rs1 rm rd 1010011 FCVT.S.W",
		"1101000 00001 rs1 rm rd 1010011 FCVT.S.WU",
		"1111000 00000 rs1 000 rd 1010011 FMV.W.X",
	},
}

var ISArv32d = ISAModule{
	name: "rv32d",
	defn: []string{
		"imm[11:0] rs1 011 rd 0000111 FLD",
		"imm[11:5] rs2 rs1 011 imm[4:0] 0100111 FSD",
		"rs3 01 rs2 rs1 rm rd 1000011 FMADD.D",
		"rs3 01 rs2 rs1 rm rd 1000111 FMSUB.D",
		"rs3 01 rs2 rs1 rm rd 1001011 FNMSUB.D",
		"rs3 01 rs2 rs1 rm rd 1001111 FNMADD.D",
		"0000001 rs2 rs1 rm rd 1010011 FADD.D",
		"0000101 rs2 rs1 rm rd 1010011 FSUB.D",
		"0001001 rs2 rs1 rm rd 1010011 FMUL.D",
		"0001101 rs2 rs1 rm rd 1010011 FDIV.D",
		"0101101 00000 rs1 rm rd 1010011 FSQRT.D",
		"0010001 rs2 rs1 000 rd 1010011 FSGNJ.D",
		"0010001 rs2 rs1 001 rd 1010011 FSGNJN.D",
		"0010001 rs2 rs1 010 rd 1010011 FSGNJX.D",
		"0010101 rs2 rs1 000 rd 1010011 FMIN.D",
		"0010101 rs2 rs1 001 rd 1010011 FMAX.D",
		"0100000 00001 rs1 rm rd 1010011 FCVT.S.D",
		"0100001 00000 rs1 rm rd 1010011 FCVT.D.S",
		"1010001 rs2 rs1 010 rd 1010011 FEQ.D",
		"1010001 rs2 rs1 001 rd 1010011 FLT.D",
		"1010001 rs2 rs1 000 rd 1010011 FLE.D",
		"1110001 00000 rs1 001 rd 1010011 FCLASS.D",
		"1100001 00000 rs1 rm rd 1010011 FCVT.W.D",
		"1100001 00001 rs1 rm rd 1010011 FCVT.WU.D",
		"1101001 00000 rs1 rm rd 1010011 FCVT.D.W",
		"1101001 00001 rs1 rm rd 1010011 FCVT.D.WU",
	},
}

//-----------------------------------------------------------------------------
// RV64 instructions (+ RV32)

var ISArv64i = ISAModule{
	name: "rv64i",
	defn: []string{
		"imm[11:0] rs1 110 rd 0000011 LWU",
		"imm[11:0] rs1 011 rd 0000011 LD",
		"imm[11:5] rs2 rs1 011 imm[4:0] 0100011 SD",
		"000000 shamt6 rs1 001 rd 0010011 SLLI",
		"000000 shamt6 rs1 101 rd 0010011 SRLI",
		"010000 shamt6 rs1 101 rd 0010011 SRAI",
		"imm[11:0] rs1 000 rd 0011011 ADDIW",
		"0000000 shamt5 rs1 001 rd 0011011 SLLIW",
		"0000000 shamt5 rs1 101 rd 0011011 SRLIW",
		"0100000 shamt5 rs1 101 rd 0011011 SRAIW",
		"0000000 rs2 rs1 000 rd 0111011 ADDW",
		"0100000 rs2 rs1 000 rd 0111011 SUBW",
		"0000000 rs2 rs1 001 rd 0111011 SLLW",
		"0000000 rs2 rs1 101 rd 0111011 SRLW",
		"0100000 rs2 rs1 101 rd 0111011 SRAW",
	},
}

var ISArv64m = ISAModule{
	name: "rv64m",
	defn: []string{
		"0000001 rs2 rs1 000 rd 0111011 MULW",
		"0000001 rs2 rs1 100 rd 0111011 DIVW",
		"0000001 rs2 rs1 101 rd 0111011 DIVUW",
		"0000001 rs2 rs1 110 rd 0111011 REMW",
		"0000001 rs2 rs1 111 rd 0111011 REMUW",
	},
}

var ISArv64a = ISAModule{
	name: "rv64a",
	defn: []string{
		"00010 aq rl 00000 rs1 011 rd 0101111 LR.D",
		"00011 aq rl rs2 rs1 011 rd 0101111 SC.D",
		"00001 aq rl rs2 rs1 011 rd 0101111 AMOSWAP.D",
		"00000 aq rl rs2 rs1 011 rd 0101111 AMOADD.D",
		"00100 aq rl rs2 rs1 011 rd 0101111 AMOXOR.D",
		"01100 aq rl rs2 rs1 011 rd 0101111 AMOAND.D",
		"01000 aq rl rs2 rs1 011 rd 0101111 AMOOR.D",
		"10000 aq rl rs2 rs1 011 rd 0101111 AMOMIN.D",
		"10100 aq rl rs2 rs1 011 rd 0101111 AMOMAX.D",
		"11000 aq rl rs2 rs1 011 rd 0101111 AMOMINU.D",
		"11100 aq rl rs2 rs1 011 rd 0101111 AMOMAXU.D",
	},
}

var ISArv64f = ISAModule{
	name: "rv64f",
	defn: []string{
		"1100000 00010 rs1 rm rd 1010011 FCVT.L.S",
		"1100000 00011 rs1 rm rd 1010011 FCVT.LU.S",
		"1101000 00010 rs1 rm rd 1010011 FCVT.S.L",
		"1101000 00011 rs1 rm rd 1010011 FCVT.S.LU",
	},
}

var ISArv64d = ISAModule{
	name: "rv64d",
	defn: []string{
		"1100001 00010 rs1 rm rd 1010011 FCVT.L.D",
		"1100001 00011 rs1 rm rd 1010011 FCVT.LU.D",
		"1110001 00000 rs1 000 rd 1010011 FMV.X.D",
		"1101001 00010 rs1 rm rd 1010011 FCVT.D.L",
		"1101001 00011 rs1 rm rd 1010011 FCVT.D.LU",
		"1111001 00000 rs1 000 rd 1010011 FMV.D.X",
	},
}

//-----------------------------------------------------------------------------

type insInfo struct {
	module    string     // ISA module to which the instruction belongs
	mneumonic string     // instruction mneumonic
	val       uint32     // value of the fixed bits in the instruction
	mask      uint32     // mask of the fixed bits in the instruction
	decode    decodeType // instruction decode type
}

// ISA is an instruction set
type ISA struct {
	name string
	ins  []*insInfo
}

// NewISA creates an empty instruction set.
func NewISA(name string) *ISA {
	return &ISA{
		name: name,
		ins:  make([]*insInfo, 0),
	}
}

// addInstruction adds an instruction to the ISA.
func (isa *ISA) addInstruction(ins string, module string) error {
	d, err := parseDefn(ins, module)
	if err != nil {
		return err
	}
	isa.ins = append(isa.ins, d)
	return nil
}

func (isa *ISA) Add(module ...ISAModule) error {
	for i := range module {
		for _, defn := range module[i].defn {
			err := isa.addInstruction(defn, module[i].name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//-----------------------------------------------------------------------------
