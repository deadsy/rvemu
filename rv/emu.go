//-----------------------------------------------------------------------------
/*

RISC-V CPU Emulation

Notes:

We use uint64 for integer registers.
The upper 32-bits is ignored for xlen == 32.

For RV32e (16 integer registers) an out of range register (>=16) will
generate an exception.

We use uint64 for float registers.
For CPUs that have only 32-bit float support the upper 32-bits are ignored.
For CPUs that have 32/64-bit float support the full 64-bits is used.
A 32-bit value written to the float register will set the upper 32-bits to
all ones for NaN boxing.

*/
//-----------------------------------------------------------------------------

package rv

import (
	"math"
	"sync"

	"github.com/deadsy/riscv/big"
	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/mem"
)

//-----------------------------------------------------------------------------
// rv32i

func emu_LUI(m *RV, ins uint) {
	imm, rd := decodeU(ins)
	m.wrX(rd, uint64(imm<<12))
	m.PC += 4
}

func emu_AUIPC(m *RV, ins uint) {
	imm, rd := decodeU(ins)
	m.wrX(rd, uint64(int(m.PC)+(imm<<12)))
	m.PC += 4
}

func emu_JAL(m *RV, ins uint) {
	imm, rd := decodeJ(ins)
	m.wrX(rd, m.PC+4)
	m.PC = uint64(int(m.PC) + int(imm))
}

func emu_JALR(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	t := m.PC + 4
	m.PC = uint64((int(m.rdX(rs1)) + imm) & ^1)
	m.wrX(rd, t)
}

func emu_BEQ(m *RV, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.rdX(rs1) == m.rdX(rs2) {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu_BNE(m *RV, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.rdX(rs1) != m.rdX(rs2) {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu_BLT(m *RV, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	var lt bool
	if m.xlen == 32 {
		lt = int32(m.rdX(rs1)) < int32(m.rdX(rs2))
	} else {
		lt = int64(m.rdX(rs1)) < int64(m.rdX(rs2))
	}
	if lt {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu_BGE(m *RV, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	var ge bool
	if m.xlen == 32 {
		ge = int32(m.rdX(rs1)) >= int32(m.rdX(rs2))
	} else {
		ge = int64(m.rdX(rs1)) >= int64(m.rdX(rs2))
	}
	if ge {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu_BLTU(m *RV, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.rdX(rs1) < m.rdX(rs2) {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu_BGEU(m *RV, ins uint) {
	imm, rs2, rs1 := decodeB(ins)
	if m.rdX(rs1) >= m.rdX(rs2) {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
}

func emu_LB(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd8(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int8(val)))
	m.PC += 4
}

func emu_LH(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd16(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int16(val)))
	m.PC += 4
}

func emu_LW(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int(val)))
	m.PC += 4
}

func emu_LBU(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd8(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(val))
	m.PC += 4
}

func emu_LHU(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd16(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(val))
	m.PC += 4
}

func emu_SB(m *RV, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr8(adr, uint8(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.PC += 4
}

func emu_SH(m *RV, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr16(adr, uint16(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.PC += 4
}

func emu_SW(m *RV, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.PC += 4
}

func emu_ADDI(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, uint64(int(m.rdX(rs1))+imm))
	m.PC += 4
}

func emu_SLTI(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	var lt bool
	if m.xlen == 32 {
		lt = int32(m.rdX(rs1)) < int32(imm)
	} else {
		lt = int64(m.rdX(rs1)) < int64(imm)
	}
	var result uint64
	if lt {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu_SLTIU(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	var lt bool
	if m.xlen == 32 {
		lt = uint32(m.rdX(rs1)) < uint32(imm)
	} else {
		lt = m.rdX(rs1) < uint64(imm)
	}
	var result uint64
	if lt {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu_XORI(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.rdX(rs1)^uint64(imm))
	m.PC += 4
}

func emu_ORI(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.rdX(rs1)|uint64(imm))
	m.PC += 4
}

func emu_ANDI(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.rdX(rs1)&uint64(imm))
	m.PC += 4
}

// rv32i/rv64i
func emu_SLLI(m *RV, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if m.xlen == 32 && shamt > 31 {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, m.rdX(rs1)<<shamt)
	m.PC += 4
}

// rv32i/rv64i
func emu_SRLI(m *RV, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if m.xlen == 32 && shamt > 31 {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, m.rdX(rs1)>>shamt)
	m.PC += 4
}

// rv32i/rv64i
func emu_SRAI(m *RV, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if m.xlen == 32 && shamt > 31 {
		m.illegalInstruction()
		return
	}
	if m.xlen == 32 {
		m.wrX(rd, uint64(int32(m.rdX(rs1))>>shamt))
	} else {
		m.wrX(rd, uint64(int64(m.rdX(rs1))>>shamt))
	}
	m.PC += 4
}

func emu_ADD(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)+m.rdX(rs2))
	m.PC += 4
}

func emu_SUB(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)-m.rdX(rs2))
	m.PC += 4
}

func emu_SLL(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var shamt uint64
	if m.xlen == 32 {
		shamt = m.rdX(rs2) & 31
	} else {
		shamt = m.rdX(rs2) & 63
	}
	m.wrX(rd, m.rdX(rs1)<<shamt)
	m.PC += 4
}

func emu_SLT(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var lt bool
	if m.xlen == 32 {
		lt = int32(m.rdX(rs1)) < int32(m.rdX(rs2))
	} else {
		lt = int64(m.rdX(rs1)) < int64(m.rdX(rs2))
	}
	var result uint64
	if lt {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu_SLTU(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var result uint64
	if m.rdX(rs1) < m.rdX(rs2) {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu_XOR(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)^m.rdX(rs2))
	m.PC += 4
}

func emu_SRL(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.rdX(rs2) & 63
	m.wrX(rd, m.rdX(rs1)>>shamt)
	m.PC += 4
}

func emu_SRA(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		shamt := m.rdX(rs2) & 31
		x = uint64(int32(m.rdX(rs1)) >> shamt)
	} else {
		shamt := m.rdX(rs2) & 63
		x = uint64(int64(m.rdX(rs1)) >> shamt)
	}
	m.wrX(rd, x)
	m.PC += 4
}

func emu_OR(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)|m.rdX(rs2))
	m.PC += 4
}

func emu_AND(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)&m.rdX(rs2))
	m.PC += 4
}

func emu_FENCE(m *RV, ins uint) {
	// no-op for a sw emulator
	m.PC += 4
}

func emu_FENCE_I(m *RV, ins uint) {
	// no-op for a sw emulator
	m.PC += 4
}

func emu_ECALL(m *RV, ins uint) {
	if m.ecall == nil {
		m.err.N = ErrEcall
		return
	}
	m.ecall.Call(m)
	m.PC += 4
}

func emu_EBREAK(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_CSRRW(m *RV, ins uint) {
	csr, rs1, rd := decodeIb(ins)
	var t uint64
	if rd != 0 {
		t = m.rdCSR(csr)
	}
	m.wrCSR(csr, m.rdX(rs1))
	m.wrX(rd, t)
	m.PC += 4
}

func emu_CSRRS(m *RV, ins uint) {
	csr, rs1, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	if rs1 != 0 {
		m.wrCSR(csr, t|m.rdX(rs1))
	}
	m.wrX(rd, t)
	m.PC += 4
}

func emu_CSRRC(m *RV, ins uint) {
	csr, rs1, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	if rs1 != 0 {
		m.wrCSR(csr, t & ^m.rdX(rs1))
	}
	m.wrX(rd, t)
	m.PC += 4
}

func emu_CSRRWI(m *RV, ins uint) {
	csr, zimm, rd := decodeIb(ins)
	if rd != 0 {
		m.wrX(rd, m.rdCSR(csr))
	}
	m.wrCSR(csr, uint64(zimm))
	m.PC += 4
}

func emu_CSRRSI(m *RV, ins uint) {
	csr, zimm, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	m.wrCSR(csr, t|uint64(zimm))
	m.wrX(rd, t)
	m.PC += 4
}

func emu_CSRRCI(m *RV, ins uint) {
	csr, zimm, rd := decodeIb(ins)
	t := m.rdCSR(csr)
	m.wrCSR(csr, t & ^uint64(zimm))
	m.wrX(rd, t)
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32i privileged

func emu_URET(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_SRET(m *RV, ins uint) {
	pc, err := m.CSR.SRET()
	m.setError(ErrCSR, err)
	m.PC = uint64(pc)
}

func emu_MRET(m *RV, ins uint) {
	pc, err := m.CSR.MRET()
	m.setError(ErrCSR, err)
	m.PC = uint64(pc)
}

func emu_WFI(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_SFENCE_VMA(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_HFENCE_BVMA(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_HFENCE_GVMA(m *RV, ins uint) {
	m.err.N = ErrTodo
}

//-----------------------------------------------------------------------------
// rv32m

func emu_MUL(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)*m.rdX(rs2))
	m.PC += 4
}

func emu_MULH(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		a := int64(int32(m.X[rs1]))
		b := int64(int32(m.X[rs2]))
		c := (a * b) >> 32
		x = uint64(c)
	} else {
		a := big.Int128FromInt(int64(m.rdX(rs1)))
		b := big.Int128FromInt(int64(m.rdX(rs2)))
		c := a.Mul(b)
		x = c.Hi
	}
	m.wrX(rd, x)
	m.PC += 4
}

func emu_MULHSU(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		a := int64(int32(m.X[rs1]))
		b := int64(m.X[rs2])
		c := (a * b) >> 32
		x = uint64(c)
	} else {
		a := big.Int128FromInt(int64(m.rdX(rs1)))
		b := big.Int128FromUint(m.rdX(rs2))
		c := a.Mul(b)
		x = c.Hi
	}
	m.wrX(rd, x)
	m.PC += 4
}

func emu_MULHU(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		a := uint64(m.X[rs1])
		b := uint64(m.X[rs2])
		c := (a * b) >> 32
		x = uint64(c)
	} else {
		a := big.Uint128FromUint(m.rdX(rs1))
		b := big.Uint128FromUint(m.rdX(rs2))
		c := a.Mul(b)
		x = c.Hi
	}
	m.wrX(rd, x)
	m.PC += 4
}

func emu_DIV(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		result := int32(-1)
		a := int32(m.X[rs1])
		b := int32(m.X[rs2])
		if b != 0 {
			result = a / b
		}
		x = uint64(result)
	} else {
		result := int64(-1)
		a := int64(m.rdX(rs1))
		b := int64(m.rdX(rs2))
		if b != 0 {
			result = a / b
		}
		x = uint64(result)
	}
	m.wrX(rd, x)
	m.PC += 4
}

func emu_DIVU(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	result := uint64((1 << 64) - 1)
	if m.rdX(rs2) != 0 {
		result = m.rdX(rs1) / m.rdX(rs2)
	}
	m.wrX(rd, result)
	m.PC += 4
}

func emu_REM(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		result := int32(m.X[rs1])
		b := int32(m.X[rs2])
		if b != 0 {
			result %= b
		}
		x = uint64(result)
	} else {
		result := int64(m.rdX(rs1))
		b := int64(m.rdX(rs2))
		if b != 0 {
			result %= b
		}
		x = uint64(result)
	}
	m.wrX(rd, x)
	m.PC += 4
}

func emu_REMU(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	result := m.rdX(rs1)
	if m.rdX(rs2) != 0 {
		result %= m.rdX(rs2)
	}
	m.wrX(rd, result)
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32a

func emu_LR_W(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_SC_W(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOSWAP_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

func emu_AMOADD_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, t+uint32(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

func emu_AMOXOR_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, t^uint32(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

func emu_AMOAND_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, t&uint32(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

func emu_AMOOR_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, t|uint32(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

func emu_AMOMIN_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, uint32(minInt32(int32(t), int32(m.rdX(rs2)))))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

func emu_AMOMAX_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, uint32(maxInt32(int32(t), int32(m.rdX(rs2)))))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

func emu_AMOMINU_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, uint32(minUint32(t, uint32(m.rdX(rs2)))))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

func emu_AMOMAXU_W(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	err = m.Mem.Wr32(adr, uint32(maxUint32(t, uint32(m.rdX(rs2)))))
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32f

func emu_FLW(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	x, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FSW(m *RV, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr32(adr, uint32(m.F[rs2]))
	m.setError(ErrMemory, err)
	m.PC += 4
}

func emu_FMADD_S(m *RV, ins uint) {
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_s(uint32(m.F[rs1]), uint32(m.F[rs2]), uint32(m.F[rs3]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FMSUB_S(m *RV, ins uint) {
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_s(uint32(m.F[rs1]), uint32(m.F[rs2]), neg32(uint32(m.F[rs3])), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FNMSUB_S(m *RV, ins uint) {
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_s(neg32(uint32(m.F[rs1])), uint32(m.F[rs2]), uint32(m.F[rs3]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FNMADD_S(m *RV, ins uint) {
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_s(neg32(uint32(m.F[rs1])), uint32(m.F[rs2]), neg32(uint32(m.F[rs3])), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FADD_S(m *RV, ins uint) {
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fadd_s(uint32(m.F[rs1]), uint32(m.F[rs2]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FSUB_S(m *RV, ins uint) {
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fsub_s(uint32(m.F[rs1]), uint32(m.F[rs2]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FMUL_S(m *RV, ins uint) {
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fmul_s(uint32(m.F[rs1]), uint32(m.F[rs2]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FDIV_S(m *RV, ins uint) {
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fdiv_s(uint32(m.F[rs1]), uint32(m.F[rs2]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FSQRT_S(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fsqrt_s(uint32(m.F[rs1]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FSGNJ_S(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := m.F[rs2] & f32SignMask
	m.F[rd] = sign | (m.F[rs1] & mask30to0) | upper32
	m.PC += 4
}

func emu_FSGNJN_S(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := ^m.F[rs2] & f32SignMask
	m.F[rd] = sign | (m.F[rs1] & mask30to0) | upper32
	m.PC += 4
}

func emu_FSGNJX_S(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := (m.F[rs1] ^ m.F[rs2]) & f32SignMask
	m.F[rd] = sign | (m.F[rs1] & mask30to0) | upper32
	m.PC += 4
}

func emu_FMIN_S(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.F[rd] = uint64(fmin_s(uint32(m.F[rs1]), uint32(m.F[rs2]), m.CSR)) | upper32
	m.PC += 4
}

func emu_FMAX_S(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.F[rd] = uint64(fmax_s(uint32(m.F[rs1]), uint32(m.F[rs2]), m.CSR)) | upper32
	m.PC += 4
}

func emu_FCVT_W_S(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_w_s(uint32(m.F[rs1]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, uint64(x))
	m.PC += 4
}

func emu_FCVT_WU_S(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_wu_s(uint32(m.F[rs1]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, uint64(int32(x)))
	m.PC += 4
}

func emu_FMV_X_W(m *RV, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.F[rs1])))
	m.PC += 4
}

func emu_FEQ_S(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(feq_s(uint32(m.F[rs1]), uint32(m.F[rs2]), m.CSR)))
	m.PC += 4
}

func emu_FLT_S(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(flt_s(uint32(m.F[rs1]), uint32(m.F[rs2]), m.CSR)))
	m.PC += 4
}

func emu_FLE_S(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(fle_s(uint32(m.F[rs1]), uint32(m.F[rs2]), m.CSR)))
	m.PC += 4
}

func emu_FCLASS_S(m *RV, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(fclass_s(uint32(m.F[rs1]))))
	m.PC += 4
}

func emu_FCVT_S_W(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_s_w(int32(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FCVT_S_WU(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_s_wu(uint32(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FMV_W_X(m *RV, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	m.F[rd] = uint64(uint32(m.rdX(rs1))) | upper32
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32d

func emu_FLD(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	x, err := m.Mem.Rd64(adr)
	m.setError(ErrMemory, err)
	m.F[rd] = x
	m.PC += 4
}

func emu_FSD(m *RV, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr64(adr, m.F[rs2])
	m.setError(ErrMemory, err)
	m.PC += 4
}

func emu_FMADD_D(m *RV, ins uint) {
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_d(m.F[rs1], m.F[rs2], m.F[rs3], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FMSUB_D(m *RV, ins uint) {
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_d(m.F[rs1], m.F[rs2], neg64(m.F[rs3]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FNMSUB_D(m *RV, ins uint) {
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_d(neg64(m.F[rs1]), m.F[rs2], m.F[rs3], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FNMADD_D(m *RV, ins uint) {
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_d(neg64(m.F[rs1]), m.F[rs2], neg64(m.F[rs3]), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FADD_D(m *RV, ins uint) {
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fadd_d(m.F[rs1], m.F[rs2], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FSUB_D(m *RV, ins uint) {
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fsub_d(m.F[rs1], m.F[rs2], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FMUL_D(m *RV, ins uint) {
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fmul_d(m.F[rs1], m.F[rs2], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FDIV_D(m *RV, ins uint) {
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fdiv_d(m.F[rs1], m.F[rs2], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FSQRT_D(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fsqrt_d(m.F[rs1], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FSGNJ_D(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := m.F[rs2] & f64SignMask
	m.F[rd] = sign | (m.F[rs1] & mask62to0)
	m.PC += 4
}

func emu_FSGNJN_D(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := ^m.F[rs2] & f64SignMask
	m.F[rd] = sign | (m.F[rs1] & mask62to0)
	m.PC += 4
}

func emu_FSGNJX_D(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	sign := (m.F[rs1] ^ m.F[rs2]) & f64SignMask
	m.F[rd] = sign | (m.F[rs1] & mask62to0)
	m.PC += 4
}

func emu_FMIN_D(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.F[rd] = fmin_d(m.F[rs1], m.F[rs2], m.CSR)
	m.PC += 4
}

func emu_FMAX_D(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.F[rd] = fmax_d(m.F[rs1], m.F[rs2], m.CSR)
	m.PC += 4
}

func emu_FCVT_S_D(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_s_d(m.F[rs1], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = uint64(x) | upper32
	m.PC += 4
}

func emu_FCVT_D_S(m *RV, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	m.F[rd] = fcvt_d_s(uint32(m.F[rs1]), m.CSR)
	m.PC += 4
}

func emu_FEQ_D(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(feq_d(m.F[rs1], m.F[rs2], m.CSR)))
	m.PC += 4
}

func emu_FLT_D(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(flt_d(m.F[rs1], m.F[rs2], m.CSR)))
	m.PC += 4
}

func emu_FLE_D(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(fle_d(m.F[rs1], m.F[rs2], m.CSR)))
	m.PC += 4
}

func emu_FCLASS_D(m *RV, ins uint) {
	_, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(fclass_d(m.F[rs1])))
	m.PC += 4
}

func emu_FCVT_W_D(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_w_d(m.F[rs1], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, uint64(x))
	m.PC += 4
}

func emu_FCVT_WU_D(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_wu_d(m.F[rs1], rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, uint64(int32(x)))
	m.PC += 4
}

func emu_FCVT_D_W(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_d_w(int32(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

func emu_FCVT_D_WU(m *RV, ins uint) {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_d_wu(uint32(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		m.illegalInstruction()
		return
	}
	m.F[rd] = x
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv32c

func emu_C_ILLEGAL(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_ADDI4SPN(m *RV, ins uint) {
	uimm, rd := decodeCIW(ins)
	m.wrX(rd, m.rdX(RegSp)+uint64(uimm))
	m.PC += 2
}

func emu_C_LW(m *RV, ins uint) {
	uimm, rs1, rd := decodeCS(ins)
	adr := uint(m.rdX(rs1)) + uimm
	val, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int(val)))
	m.PC += 2
}

func emu_C_SW(m *RV, ins uint) {
	uimm, rs1, rs2 := decodeCS(ins)
	adr := uint(m.rdX(rs1)) + uimm
	err := m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.PC += 2
}

func emu_C_NOP(m *RV, ins uint) {
	m.PC += 2
}

func emu_C_ADDI(m *RV, ins uint) {
	imm, rd := decodeCIa(ins)
	if rd != 0 {
		m.wrX(rd, uint64(int(m.rdX(rd))+imm))
	}
	m.PC += 2
}

func emu_C_LI(m *RV, ins uint) {
	imm, rd := decodeCIa(ins)
	m.wrX(rd, uint64(imm))
	m.PC += 2
}

func emu_C_ADDI16SP(m *RV, ins uint) {
	imm := decodeCIb(ins)
	m.wrX(RegSp, uint64(int(m.rdX(RegSp))+imm))
	m.PC += 2
}

func emu_C_LUI(m *RV, ins uint) {
	imm, rd := decodeCIf(ins)
	if imm == 0 {
		m.illegalInstruction()
		return
	}
	if rd != 0 && rd != 2 {
		m.wrX(rd, uint64(imm<<12))
	}
	m.PC += 2
}

func emu_C_SRLI(m *RV, ins uint) {
	shamt, rd := decodeCIc(ins)
	if m.xlen == 32 && shamt > 31 {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, m.rdX(rd)>>shamt)
	m.PC += 2
}

func emu_C_SRAI(m *RV, ins uint) {
	shamt, rd := decodeCIc(ins)
	var x uint64
	if m.xlen == 32 {
		if shamt > 31 {
			m.illegalInstruction()
			return
		}
		x = uint64(int32(m.X[rd]) >> shamt)
	} else {
		x = uint64(int64(m.rdX(rd)) >> shamt)
	}
	m.wrX(rd, x)
	m.PC += 2
}

func emu_C_ANDI(m *RV, ins uint) {
	imm, rd := decodeCIe(ins)
	m.wrX(rd, uint64(int(m.rdX(rd))&imm))
	m.PC += 2
}

func emu_C_SUB(m *RV, ins uint) {
	rd, rs := decodeCRa(ins)
	m.wrX(rd, m.rdX(rd)-m.rdX(rs))
	m.PC += 2
}

func emu_C_XOR(m *RV, ins uint) {
	rd, rs := decodeCRa(ins)
	m.wrX(rd, m.rdX(rd)^m.rdX(rs))
	m.PC += 2
}

func emu_C_OR(m *RV, ins uint) {
	rd, rs := decodeCRa(ins)
	m.wrX(rd, m.rdX(rd)|m.rdX(rs))
	m.PC += 2
}

func emu_C_AND(m *RV, ins uint) {
	rd, rs := decodeCRa(ins)
	m.wrX(rd, m.rdX(rd)&m.rdX(rs))
	m.PC += 2
}

func emu_C_J(m *RV, ins uint) {
	imm := decodeCJ(ins)
	m.PC = uint64(int(m.PC) + imm)
}

func emu_C_BEQZ(m *RV, ins uint) {
	imm, rs := decodeCB(ins)
	if m.rdX(rs) == 0 {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 2
	}
}

func emu_C_BNEZ(m *RV, ins uint) {
	imm, rs := decodeCB(ins)
	if m.rdX(rs) != 0 {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 2
	}
}

func emu_C_SLLI(m *RV, ins uint) {
	shamt, rd := decodeCId(ins)
	if rd != 0 && shamt != 0 {
		m.wrX(rd, m.rdX(rd)<<shamt)
	}
	m.PC += 2
}

func emu_C_SLLI64(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_LWSP(m *RV, ins uint) {
	uimm, rd := decodeCSSa(ins)
	if rd == 0 {
		m.illegalInstruction()
		return
	}
	adr := uint(m.rdX(RegSp)) + uimm
	val, err := m.Mem.Rd32(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, uint64(int(val)))
	m.PC += 2
}

func emu_C_JR(m *RV, ins uint) {
	rs1, _ := decodeCR(ins)
	if rs1 == 0 {
		m.illegalInstruction()
		return
	}
	m.PC = m.rdX(rs1)
}

func emu_C_MV(m *RV, ins uint) {
	rd, rs := decodeCR(ins)
	if rs != 0 {
		m.wrX(rd, m.rdX(rs))
	}
	m.PC += 2
}

func emu_C_EBREAK(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_JALR(m *RV, ins uint) {
	rs1, _ := decodeCR(ins)
	if rs1 == 0 {
		m.illegalInstruction()
		return
	}
	t := m.PC + 2
	m.PC = m.rdX(rs1)
	m.wrX(RegRa, t)
}

func emu_C_ADD(m *RV, ins uint) {
	rd, rs := decodeCR(ins)
	m.wrX(rd, m.rdX(rd)+m.rdX(rs))
	m.PC += 2
}

func emu_C_SWSP(m *RV, ins uint) {
	uimm, rs2 := decodeCSSb(ins)
	adr := uint(m.rdX(RegSp)) + uimm
	err := m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	m.setError(ErrMemory, err)
	m.PC += 2
}

//-----------------------------------------------------------------------------
// rv32c-only

func emu_C_JAL(m *RV, ins uint) {
	imm := decodeCJ(ins)
	m.wrX(RegRa, m.PC+2)
	m.PC = uint64(int(m.PC) + imm)
}

//-----------------------------------------------------------------------------
// rv32fc

func emu_C_FLW(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_FLWSP(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_FSW(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_FSWSP(m *RV, ins uint) {
	m.err.N = ErrTodo
}

//-----------------------------------------------------------------------------
// rv32dc

func emu_C_FLD(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_FLDSP(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_FSD(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_FSDSP(m *RV, ins uint) {
	m.err.N = ErrTodo
}

//-----------------------------------------------------------------------------
// rv64i

func emu_LWU(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_LD(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd64(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, val)
	m.PC += 4
}

func emu_SD(m *RV, ins uint) {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr64(adr, m.rdX(rs2))
	m.setError(ErrMemory, err)
	m.PC += 4
}

func emu_ADDIW(m *RV, ins uint) {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, uint64(int32(int(m.rdX(rs1))+imm)))
	m.PC += 4
}

func emu_SLLIW(m *RV, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, uint64(int32(m.rdX(rs1))<<shamt))
	m.PC += 4
}

func emu_SRLIW(m *RV, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, uint64(int32(uint32(m.rdX(rs1))>>shamt)))
	m.PC += 4
}

func emu_SRAIW(m *RV, ins uint) {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		m.illegalInstruction()
		return
	}
	m.wrX(rd, uint64(int32(m.rdX(rs1))>>shamt))
	m.PC += 4
}

func emu_ADDW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.rdX(rs1)+m.rdX(rs2))))
	m.PC += 4
}

func emu_SUBW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.rdX(rs1)-m.rdX(rs2))))
	m.PC += 4
}

func emu_SLLW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.rdX(rs2) & 31
	m.wrX(rd, uint64(int32(m.rdX(rs1)<<shamt)))
	m.PC += 4
}

func emu_SRLW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.rdX(rs2) & 31
	m.wrX(rd, uint64(int32(uint32(m.rdX(rs1))>>shamt)))
	m.PC += 4
}

func emu_SRAW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.rdX(rs2) & 31
	m.wrX(rd, uint64(int32(m.rdX(rs1))>>shamt))
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv64m

func emu_MULW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	result := int32(m.rdX(rs1) * m.rdX(rs2))
	m.wrX(rd, uint64(result))
	m.PC += 4
}

func emu_DIVW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	result := int32(m.rdX(rs1))
	divisor := int32(m.rdX(rs2))
	if divisor == -1 && result == math.MinInt32 {
		// overflow
	} else if divisor == 0 {
		result = -1
	} else {
		result /= divisor
	}
	m.wrX(rd, uint64(result))
	m.PC += 4
}

func emu_DIVUW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	dividend := uint32(m.rdX(rs1))
	divisor := uint32(m.rdX(rs2))
	result := int32(-1)
	if divisor != 0 {
		result = int32(dividend / divisor)
	}
	m.wrX(rd, uint64(result))
	m.PC += 4
}

func emu_REMW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	result := int32(m.rdX(rs1))
	divisor := int32(m.rdX(rs2))
	if divisor == -1 && result == math.MinInt32 {
		// overflow
		result = 0
	} else if divisor == 0 {
		// nop
	} else {
		result %= divisor
	}
	m.wrX(rd, uint64(result))
	m.PC += 4
}

func emu_REMUW(m *RV, ins uint) {
	rs2, rs1, _, rd := decodeR(ins)
	dividend := uint32(m.rdX(rs1))
	divisor := uint32(m.rdX(rs2))
	result := int32(dividend)
	if divisor != 0 {
		result = int32(dividend % divisor)
	}
	m.wrX(rd, uint64(result))
	m.PC += 4
}

//-----------------------------------------------------------------------------
// rv64a

func emu_LR_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_SC_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOSWAP_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOADD_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOXOR_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOAND_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOOR_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOMIN_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOMAX_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOMINU_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_AMOMAXU_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

//-----------------------------------------------------------------------------
// rv64f

func emu_FCVT_L_S(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_FCVT_LU_S(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_FCVT_S_L(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_FCVT_S_LU(m *RV, ins uint) {
	m.err.N = ErrTodo
}

//-----------------------------------------------------------------------------
// rv64d

func emu_FCVT_L_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_FCVT_LU_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_FMV_X_D(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_FCVT_D_L(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_FCVT_D_LU(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_FMV_D_X(m *RV, ins uint) {
	m.err.N = ErrTodo
}

//-----------------------------------------------------------------------------
// rv64c

func emu_C_ADDIW(m *RV, ins uint) {
	imm, rd := decodeCIa(ins)
	if rd != 0 {
		m.wrX(rd, uint64(int32(int(m.rdX(rd))+imm)))
	} else {
		m.illegalInstruction()
	}
	m.PC += 2
}

func emu_C_LDSP(m *RV, ins uint) {
	uimm, rd := decodeCIg(ins)
	adr := uint(m.rdX(RegSp)) + uimm
	val, err := m.Mem.Rd64(adr)
	m.setError(ErrMemory, err)
	if rd != 0 {
		m.wrX(rd, val)
	} else {
		m.illegalInstruction()
	}
	m.PC += 2
}

func emu_C_SDSP(m *RV, ins uint) {
	uimm, rs2 := decodeCSSc(ins)
	adr := uint(m.rdX(RegSp)) + uimm
	err := m.Mem.Wr64(adr, m.rdX(rs2))
	m.setError(ErrMemory, err)
	m.PC += 2
}

func emu_C_LD(m *RV, ins uint) {
	uimm, rs1, rd := decodeCSa(ins)
	adr := uint(m.rdX(rs1)) + uimm
	val, err := m.Mem.Rd64(adr)
	m.setError(ErrMemory, err)
	m.wrX(rd, val)
	m.PC += 2
}

func emu_C_SD(m *RV, ins uint) {
	uimm, rs1, rs2 := decodeCSa(ins)
	adr := uint(m.rdX(rs1)) + uimm
	err := m.Mem.Wr64(adr, m.rdX(rs2))
	m.setError(ErrMemory, err)
	m.PC += 2
}

func emu_C_SUBW(m *RV, ins uint) {
	m.err.N = ErrTodo
}

func emu_C_ADDW(m *RV, ins uint) {
	m.err.N = ErrTodo
}

//-----------------------------------------------------------------------------
// Integer Register Access

// wrX writes an integer register
func (m *RV) wrX(i uint, val uint64) {
	if i == 0 {
		// no writes to zero
		return
	}
	if m.nreg == 16 && i >= 16 {
		m.err.N = ErrBadReg
		return
	}
	if m.xlen == 32 {
		val = uint64(uint32(val))
	}
	m.X[i] = val
}

// rdX reads an integer register
func (m *RV) rdX(i uint) uint64 {
	if m.nreg == 16 && i >= 16 {
		m.err.N = ErrBadReg
		return 0
	}
	if m.xlen == 32 {
		return uint64(uint32(m.X[i]))
	}
	return m.X[i]
}

//-----------------------------------------------------------------------------

// illegalInstruction
func (m *RV) illegalInstruction() {
	m.err.N = ErrIllegal
	m.CSR.SetException(csr.ExInsIllegal)
}

//-----------------------------------------------------------------------------
// CSR access

// rdCSR reads a CSR.
func (m *RV) rdCSR(reg uint) uint64 {
	val, err := m.CSR.Rd(reg)
	m.setError(ErrCSR, err)
	return uint64(val)
}

// wrCSR writes a CSR.
func (m *RV) wrCSR(reg uint, val uint64) {
	err := m.CSR.Wr(reg, uint(val))
	m.setError(ErrCSR, err)
}

//-----------------------------------------------------------------------------

func maxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func maxUint32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func minUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

// RV is a RISC-V CPU.
type RV struct {
	X        [32]uint64  // integer registers
	F        [32]uint64  // float registers
	PC       uint64      // program counter
	Mem      *mem.Memory // memory of the target system
	CSR      *csr.State  // CSR state
	amo      sync.Mutex  // lock for atomic operations
	insCount uint        // number of instructions run
	lastPC   uint64      // stuck PC detection
	xlen     uint        // bit length of integer registers
	nreg     uint        // number of integer registers
	err      Error       // emulation errors
	isa      *ISA        // ISA implemented for the CPU
	ecall    Ecall       // ecall interface
}

// NewRV64 returns a 64-bit RISC-V CPU.
func NewRV64(isa *ISA, mem *mem.Memory, ecall Ecall) *RV {
	m := RV{
		xlen:  64,
		nreg:  32,
		Mem:   mem,
		isa:   isa,
		ecall: ecall,
	}
	m.Reset()
	return &m
}

// NewRV32 returns a 32-bit RISC-V CPU.
func NewRV32(isa *ISA, mem *mem.Memory, ecall Ecall) *RV {
	m := RV{
		xlen:  32,
		nreg:  32,
		Mem:   mem,
		isa:   isa,
		ecall: ecall,
	}
	m.Reset()
	return &m
}

// NewRV32E returns a 32-bit embedded RISC-V CPU (16 integer registers).
func NewRV32E(isa *ISA, mem *mem.Memory, ecall Ecall) *RV {
	m := RV{
		xlen:  32,
		nreg:  16,
		Mem:   mem,
		isa:   isa,
		ecall: ecall,
	}
	m.Reset()
	return &m
}

// Run the CPU for a single instruction.
func (m *RV) Run() error {

	// set the pc for the exception (if there is one)
	m.err.pc = uint(m.PC)

	// read the next instruction
	ins, err := m.Mem.RdIns(uint(m.PC))
	if err != nil {
		m.setError(ErrMemory, err)
		return &m.err
	}

	// lookup and emulate the instruction
	im := m.isa.lookup(ins)
	if im != nil {
		im.defn.emu(m, ins)
		m.insCount++
	} else {
		m.illegalInstruction()
	}

	// check exception flags
	if m.err.N != 0 {
		return &m.err
	}

	// stuck PC detection
	if m.PC == m.lastPC {
		m.err.N = ErrStuck
		return &m.err
	}
	m.lastPC = m.PC

	return nil
}

//-----------------------------------------------------------------------------

// IntRegs returns a display string for the integer registers.
func (m *RV) IntRegs() string {
	reg := make([]uint, m.nreg)
	for i := range reg {
		reg[i] = uint(m.rdX(uint(i)))
	}
	return intRegString(reg, uint(m.PC), m.xlen)
}

// FloatRegs returns a display string for the float registers.
func (m *RV) FloatRegs() string {
	return floatRegString(m.F[:])
}

// Disassemble the instruction at the address.
func (m *RV) Disassemble(adr uint) *Disassembly {
	return m.isa.Disassemble(m.Mem, adr)
}

//-----------------------------------------------------------------------------

// Reset the CPU.
func (m *RV) Reset() {
	m.PC = m.Mem.Entry
	m.wrX(RegSp, uint64(uint(1<<32)-16))
	m.CSR = csr.NewState(m.xlen)
	m.insCount = 0
	m.lastPC = 0
	m.err = Error{alen: int(m.xlen)}
}

// Exit sets a status code and exits the emulation
func (m *RV) Exit(status uint64) {
	m.err.N = ErrExit
}

// SetBreak sets the break flag.
func (m *RV) SetBreak() {
	m.err.N = ErrBreak
}

//-----------------------------------------------------------------------------
