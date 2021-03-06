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

	"github.com/deadsy/riscv/csr"
	"github.com/deadsy/riscv/mem"
)

//-----------------------------------------------------------------------------
// rv32i

func emu_LUI(m *RV, ins uint) error {
	imm, rd := decodeU(ins)
	m.wrX(rd, uint64(imm<<12))
	m.PC += 4
	return nil
}

func emu_AUIPC(m *RV, ins uint) error {
	imm, rd := decodeU(ins)
	m.wrX(rd, uint64(int(m.PC)+(imm<<12)))
	m.PC += 4
	return nil
}

func emu_JAL(m *RV, ins uint) error {
	imm, rd := decodeJ(ins)
	m.wrX(rd, m.PC+4)
	m.PC = uint64(int(m.PC) + int(imm))
	return nil
}

func emu_JALR(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	t := m.PC + 4
	m.PC = uint64((int(m.rdX(rs1)) + imm) & ^1)
	m.wrX(rd, t)
	return nil
}

func emu_BEQ(m *RV, ins uint) error {
	imm, rs2, rs1 := decodeB(ins)
	if m.rdX(rs1) == m.rdX(rs2) {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
	return nil
}

func emu_BNE(m *RV, ins uint) error {
	imm, rs2, rs1 := decodeB(ins)
	if m.rdX(rs1) != m.rdX(rs2) {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
	return nil
}

func emu_BLT(m *RV, ins uint) error {
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
	return nil
}

func emu_BGE(m *RV, ins uint) error {
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
	return nil
}

func emu_BLTU(m *RV, ins uint) error {
	imm, rs2, rs1 := decodeB(ins)
	if m.rdX(rs1) < m.rdX(rs2) {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
	return nil
}

func emu_BGEU(m *RV, ins uint) error {
	imm, rs2, rs1 := decodeB(ins)
	if m.rdX(rs1) >= m.rdX(rs2) {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 4
	}
	return nil
}

func emu_LB(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd8(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int8(val)))
	m.PC += 4
	return nil
}

func emu_LH(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd16(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int16(val)))
	m.PC += 4
	return nil
}

func emu_LW(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(val)))
	m.PC += 4
	return nil
}

func emu_LBU(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd8(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(val))
	m.PC += 4
	return nil
}

func emu_LHU(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd16(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(val))
	m.PC += 4
	return nil
}

func emu_SB(m *RV, ins uint) error {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr8(adr, uint8(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 4
	return nil
}

func emu_SH(m *RV, ins uint) error {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr16(adr, uint16(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 4
	return nil
}

func emu_SW(m *RV, ins uint) error {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 4
	return nil
}

func emu_ADDI(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, uint64(int(m.rdX(rs1))+imm))
	m.PC += 4
	return nil
}

func emu_SLTI(m *RV, ins uint) error {
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
	return nil
}

func emu_SLTIU(m *RV, ins uint) error {
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
	return nil
}

func emu_XORI(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.rdX(rs1)^uint64(imm))
	m.PC += 4
	return nil
}

func emu_ORI(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.rdX(rs1)|uint64(imm))
	m.PC += 4
	return nil
}

func emu_ANDI(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, m.rdX(rs1)&uint64(imm))
	m.PC += 4
	return nil
}

// rv32i/rv64i
func emu_SLLI(m *RV, ins uint) error {
	shamt, rs1, rd := decodeIc(ins)
	if m.xlen == 32 && shamt > 31 {
		return m.errIllegal(ins)
	}
	m.wrX(rd, m.rdX(rs1)<<shamt)
	m.PC += 4
	return nil
}

// rv32i/rv64i
func emu_SRLI(m *RV, ins uint) error {
	shamt, rs1, rd := decodeIc(ins)
	if m.xlen == 32 && shamt > 31 {
		return m.errIllegal(ins)
	}
	m.wrX(rd, m.rdX(rs1)>>shamt)
	m.PC += 4
	return nil
}

// rv32i/rv64i
func emu_SRAI(m *RV, ins uint) error {
	shamt, rs1, rd := decodeIc(ins)
	if m.xlen == 32 && shamt > 31 {
		return m.errIllegal(ins)
	}
	if m.xlen == 32 {
		m.wrX(rd, uint64(int32(m.rdX(rs1))>>shamt))
	} else {
		m.wrX(rd, uint64(int64(m.rdX(rs1))>>shamt))
	}
	m.PC += 4
	return nil
}

func emu_ADD(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)+m.rdX(rs2))
	m.PC += 4
	return nil
}

func emu_SUB(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)-m.rdX(rs2))
	m.PC += 4
	return nil
}

func emu_SLL(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	var shamt uint64
	if m.xlen == 32 {
		shamt = m.rdX(rs2) & 31
	} else {
		shamt = m.rdX(rs2) & 63
	}
	m.wrX(rd, m.rdX(rs1)<<shamt)
	m.PC += 4
	return nil
}

func emu_SLT(m *RV, ins uint) error {
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
	return nil
}

func emu_SLTU(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	var result uint64
	if m.rdX(rs1) < m.rdX(rs2) {
		result = 1
	}
	m.wrX(rd, result)
	m.PC += 4
	return nil
}

func emu_XOR(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)^m.rdX(rs2))
	m.PC += 4
	return nil
}

func emu_SRL(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.rdX(rs2) & 63
	m.wrX(rd, m.rdX(rs1)>>shamt)
	m.PC += 4
	return nil
}

func emu_SRA(m *RV, ins uint) error {
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
	return nil
}

func emu_OR(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)|m.rdX(rs2))
	m.PC += 4
	return nil
}

func emu_AND(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)&m.rdX(rs2))
	m.PC += 4
	return nil
}

func emu_FENCE(m *RV, ins uint) error {
	// no-op for a sw emulator
	m.PC += 4
	return nil
}

func emu_FENCE_I(m *RV, ins uint) error {
	// no-op for a sw emulator
	m.PC += 4
	return nil
}

func emu_ECALL(m *RV, ins uint) error {
	m.PC = m.CSR.ECALL(m.PC, 0)
	return nil
}

func emu_EBREAK(m *RV, ins uint) error {
	m.PC = m.CSR.Exception(m.PC, uint(csr.ExBreakpoint), uint(m.PC), false)
	return nil
}

func emu_CSRRW(m *RV, ins uint) error {
	csr, rs1, rd := decodeIb(ins)
	var t uint64
	var err error
	if rd != 0 {
		t, err = m.CSR.Rd(csr)
		if err != nil {
			return m.errCSR(err, ins)
		}
	}
	err = m.CSR.Wr(csr, m.rdX(rs1))
	if err != nil {
		return m.errCSR(err, ins)
	}
	m.wrX(rd, t)
	m.PC += 4
	return nil
}

func emu_CSRRS(m *RV, ins uint) error {
	csr, rs1, rd := decodeIb(ins)
	t, err := m.CSR.Rd(csr)
	if err != nil {
		return m.errCSR(err, ins)
	}
	if rs1 != 0 {
		err := m.CSR.Wr(csr, t|m.rdX(rs1))
		if err != nil {
			return m.errCSR(err, ins)
		}
	}
	m.wrX(rd, t)
	m.PC += 4
	return nil
}

func emu_CSRRC(m *RV, ins uint) error {
	csr, rs1, rd := decodeIb(ins)
	t, err := m.CSR.Rd(csr)
	if err != nil {
		return m.errCSR(err, ins)
	}
	if rs1 != 0 {
		err := m.CSR.Wr(csr, t & ^m.rdX(rs1))
		if err != nil {
			return m.errCSR(err, ins)
		}
	}
	m.wrX(rd, t)
	m.PC += 4
	return nil
}

func emu_CSRRWI(m *RV, ins uint) error {
	csr, zimm, rd := decodeIb(ins)
	if rd != 0 {
		t, err := m.CSR.Rd(csr)
		if err != nil {
			return m.errCSR(err, ins)
		}
		m.wrX(rd, t)
	}
	err := m.CSR.Wr(csr, uint64(zimm))
	if err != nil {
		return m.errCSR(err, ins)
	}
	m.PC += 4
	return nil
}

func emu_CSRRSI(m *RV, ins uint) error {
	csr, zimm, rd := decodeIb(ins)
	t, err := m.CSR.Rd(csr)
	if err != nil {
		return m.errCSR(err, ins)
	}
	err = m.CSR.Wr(csr, t|uint64(zimm))
	if err != nil {
		return m.errCSR(err, ins)
	}
	m.wrX(rd, t)
	m.PC += 4
	return nil
}

func emu_CSRRCI(m *RV, ins uint) error {
	csr, zimm, rd := decodeIb(ins)
	t, err := m.CSR.Rd(csr)
	if err != nil {
		return m.errCSR(err, ins)
	}
	err = m.CSR.Wr(csr, t & ^uint64(zimm))
	if err != nil {
		return m.errCSR(err, ins)
	}
	m.wrX(rd, t)
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv32i privileged

func emu_URET(m *RV, ins uint) error {
	m.PC = uint64(m.CSR.URET())
	return nil
}

func emu_SRET(m *RV, ins uint) error {
	m.PC = uint64(m.CSR.SRET())
	return nil
}

func emu_MRET(m *RV, ins uint) error {
	m.PC = uint64(m.CSR.MRET())
	return nil
}

func emu_WFI(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_SFENCE_VMA(m *RV, ins uint) error {
	m.PC += 4
	return nil
}

func emu_HFENCE_BVMA(m *RV, ins uint) error {
	m.PC += 4
	return nil
}

func emu_HFENCE_GVMA(m *RV, ins uint) error {
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv32m

func emu_MUL(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdX(rs1)*m.rdX(rs2))
	m.PC += 4
	return nil
}

func emu_MULH(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		a := int64(int32(m.rdX(rs1)))
		b := int64(int32(m.rdX(rs2)))
		c := (a * b) >> 32
		x = uint64(c)
	} else {
		x = uint64(mulhss(int64(m.rdX(rs1)), int64(m.rdX(rs2))))
	}
	m.wrX(rd, x)
	m.PC += 4
	return nil
}

func emu_MULHSU(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		a := int64(int32(m.rdX(rs1)))
		b := int64(m.rdX(rs2))
		c := (a * b) >> 32
		x = uint64(c)
	} else {
		x = uint64(mulhsu(int64(m.rdX(rs1)), m.rdX(rs2)))
	}
	m.wrX(rd, x)
	m.PC += 4
	return nil
}

func emu_MULHU(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		a := uint64(m.rdX(rs1))
		b := uint64(m.rdX(rs2))
		c := (a * b) >> 32
		x = uint64(c)
	} else {
		x = mulhuu(m.rdX(rs1), m.rdX(rs2))
	}
	m.wrX(rd, x)
	m.PC += 4
	return nil
}

func emu_DIV(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		result := int32(-1)
		a := int32(m.rdX(rs1))
		b := int32(m.rdX(rs2))
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
	return nil
}

func emu_DIVU(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	result := uint64((1 << 64) - 1)
	if m.rdX(rs2) != 0 {
		result = m.rdX(rs1) / m.rdX(rs2)
	}
	m.wrX(rd, result)
	m.PC += 4
	return nil
}

func emu_REM(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	var x uint64
	if m.xlen == 32 {
		result := int32(m.rdX(rs1))
		b := int32(m.rdX(rs2))
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
	return nil
}

func emu_REMU(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	result := m.rdX(rs1)
	if m.rdX(rs2) != 0 {
		result %= m.rdX(rs2)
	}
	m.wrX(rd, result)
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv32a

func emu_LR_W(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_SC_W(m *RV, ins uint) error {
	return m.errTodo()
	/*
	     rs2, rs1, _, rd := decodeR(ins)
	   	m.amo.Lock()
	   	adr := uint(m.rdX(rs1))
	   	err = m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	   	if err != nil {
	   		return m.errMemory(err)
	   	}
	   	m.amo.Unlock()
	   	m.PC += 4
	   	return nil
	*/
}

func emu_AMOSWAP_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOADD_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, t+uint32(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOXOR_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, t^uint32(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOAND_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, t&uint32(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOOR_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, t|uint32(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOMIN_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, uint32(minInt32(int32(t), int32(m.rdX(rs2)))))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOMAX_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, uint32(maxInt32(int32(t), int32(m.rdX(rs2)))))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOMINU_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, uint32(minUint32(t, uint32(m.rdX(rs2)))))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOMAXU_W(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr32(adr, uint32(maxUint32(t, uint32(m.rdX(rs2)))))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(t)))
	m.amo.Unlock()
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv32f

func emu_FLW(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	x, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FSW(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr32(adr, m.rdFS(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 4
	return nil
}

func emu_FMADD_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_s(m.rdFS(rs1), m.rdFS(rs2), m.rdFS(rs3), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FMSUB_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_s(m.rdFS(rs1), m.rdFS(rs2), neg32(m.rdFS(rs3)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FNMSUB_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_s(neg32(m.rdFS(rs1)), m.rdFS(rs2), m.rdFS(rs3), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FNMADD_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_s(neg32(m.rdFS(rs1)), m.rdFS(rs2), neg32(m.rdFS(rs3)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FADD_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fadd_s(m.rdFS(rs1), m.rdFS(rs2), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FSUB_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fsub_s(m.rdFS(rs1), m.rdFS(rs2), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FMUL_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fmul_s(m.rdFS(rs1), m.rdFS(rs2), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FDIV_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fdiv_s(m.rdFS(rs1), m.rdFS(rs2), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FSQRT_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fsqrt_s(m.rdFS(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FSGNJ_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	sign := m.rdFS(rs2) & f32SignMask
	m.wrFS(rd, sign|(m.rdFS(rs1)&mask30to0))
	m.PC += 4
	return nil
}

func emu_FSGNJN_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	sign := ^m.rdFS(rs2) & f32SignMask
	m.wrFS(rd, sign|(m.rdFS(rs1)&mask30to0))
	m.PC += 4
	return nil
}

func emu_FSGNJX_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	sign := (m.rdFS(rs1) ^ m.rdFS(rs2)) & f32SignMask
	m.wrFS(rd, sign|(m.rdFS(rs1)&mask30to0))
	m.PC += 4
	return nil
}

func emu_FMIN_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrFS(rd, fmin_s(m.rdFS(rs1), m.rdFS(rs2), m.CSR))
	m.PC += 4
	return nil
}

func emu_FMAX_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrFS(rd, fmax_s(m.rdFS(rs1), m.rdFS(rs2), m.CSR))
	m.PC += 4
	return nil
}

func emu_FCVT_W_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_w_s(m.rdFS(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(int64(x)))
	m.PC += 4
	return nil
}

func emu_FCVT_WU_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_wu_s(m.rdFS(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(int64(int32(x))))
	m.PC += 4
	return nil
}

func emu_FMV_X_W(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.rdFS(rs1))))
	m.PC += 4
	return nil
}

func emu_FEQ_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(feq_s(m.rdFS(rs1), m.rdFS(rs2), m.CSR)))
	m.PC += 4
	return nil
}

func emu_FLT_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(flt_s(m.rdFS(rs1), m.rdFS(rs2), m.CSR)))
	m.PC += 4
	return nil
}

func emu_FLE_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(fle_s(m.rdFS(rs1), m.rdFS(rs2), m.CSR)))
	m.PC += 4
	return nil
}

func emu_FCLASS_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(fclass_s(m.rdFS(rs1))))
	m.PC += 4
	return nil
}

func emu_FCVT_S_W(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_s_w(int32(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FCVT_S_WU(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_s_wu(uint32(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FMV_W_X(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, _, rd := decodeR(ins)
	m.wrFS(rd, uint32(m.rdX(rs1)))
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv32d

func emu_FLD(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	x, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FSD(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr64(adr, m.rdFD(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 4
	return nil
}

func emu_FMADD_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_d(m.rdFD(rs1), m.rdFD(rs2), m.rdFD(rs3), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FMSUB_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_d(m.rdFD(rs1), m.rdFD(rs2), neg64(m.rdFD(rs3)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FNMSUB_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_d(neg64(m.rdFD(rs1)), m.rdFD(rs2), m.rdFD(rs3), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FNMADD_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs3, rs2, rs1, rm, rd := decodeR4(ins)
	x, err := fmadd_d(neg64(m.rdFD(rs1)), m.rdFD(rs2), neg64(m.rdFD(rs3)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FADD_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fadd_d(m.rdFD(rs1), m.rdFD(rs2), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FSUB_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fsub_d(m.rdFD(rs1), m.rdFD(rs2), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FMUL_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fmul_d(m.rdFD(rs1), m.rdFD(rs2), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FDIV_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, rm, rd := decodeR(ins)
	x, err := fdiv_d(m.rdFD(rs1), m.rdFD(rs2), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FSQRT_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fsqrt_d(m.rdFD(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FSGNJ_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	sign := m.rdFD(rs2) & f64SignMask
	m.wrFD(rd, sign|(m.rdFD(rs1)&mask62to0))
	m.PC += 4
	return nil
}

func emu_FSGNJN_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	sign := ^m.rdFD(rs2) & f64SignMask
	m.wrFD(rd, sign|(m.rdFD(rs1)&mask62to0))
	m.PC += 4
	return nil
}

func emu_FSGNJX_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	sign := (m.rdFD(rs1) ^ m.rdFD(rs2)) & f64SignMask
	m.wrFD(rd, sign|(m.rdFD(rs1)&mask62to0))
	m.PC += 4
	return nil
}

func emu_FMIN_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrFD(rd, fmin_d(m.rdFD(rs1), m.rdFD(rs2), m.CSR))
	m.PC += 4
	return nil
}

func emu_FMAX_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrFD(rd, fmax_d(m.rdFD(rs1), m.rdFD(rs2), m.CSR))
	m.PC += 4
	return nil
}

func emu_FCVT_S_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_s_d(m.rdFD(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FCVT_D_S(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, _, rd := decodeR(ins)
	m.wrFD(rd, fcvt_d_s(m.rdFS(rs1), m.CSR))
	m.PC += 4
	return nil
}

func emu_FEQ_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(feq_d(m.rdFD(rs1), m.rdFD(rs2), m.CSR)))
	m.PC += 4
	return nil
}

func emu_FLT_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(flt_d(m.rdFD(rs1), m.rdFD(rs2), m.CSR)))
	m.PC += 4
	return nil
}

func emu_FLE_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(fle_d(m.rdFD(rs1), m.rdFD(rs2), m.CSR)))
	m.PC += 4
	return nil
}

func emu_FCLASS_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(fclass_d(m.rdFD(rs1))))
	m.PC += 4
	return nil
}

func emu_FCVT_W_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_w_d(m.rdFD(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(int64(x)))
	m.PC += 4
	return nil
}

func emu_FCVT_WU_D(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_wu_d(m.rdFD(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(int64(int32(x))))
	m.PC += 4
	return nil
}

func emu_FCVT_D_W(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_d_w(int32(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FCVT_D_WU(m *RV, ins uint) error {
	if m.CSR.IsFloatOff() {
		return m.errIllegal(ins)
	}
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_d_wu(uint32(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv32c

func emu_C_ILLEGAL(m *RV, ins uint) error {
	return m.errIllegal(ins)
}

func emu_C_ADDI4SPN(m *RV, ins uint) error {
	uimm, rd := decodeCIW(ins)
	m.wrX(rd, m.rdX(RegSp)+uint64(uimm))
	m.PC += 2
	return nil
}

func emu_C_LW(m *RV, ins uint) error {
	uimm, rs1, rd := decodeCS(ins)
	adr := uint(m.rdX(rs1)) + uimm
	val, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(val)))
	m.PC += 2
	return nil
}

func emu_C_SW(m *RV, ins uint) error {
	uimm, rs1, rs2 := decodeCS(ins)
	adr := uint(m.rdX(rs1)) + uimm
	err := m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 2
	return nil
}

func emu_C_NOP(m *RV, ins uint) error {
	m.PC += 2
	return nil
}

func emu_C_ADDI(m *RV, ins uint) error {
	imm, rd := decodeCIa(ins)
	if rd != 0 {
		m.wrX(rd, uint64(int(m.rdX(rd))+imm))
	}
	m.PC += 2
	return nil
}

func emu_C_LI(m *RV, ins uint) error {
	imm, rd := decodeCIa(ins)
	m.wrX(rd, uint64(imm))
	m.PC += 2
	return nil
}

func emu_C_ADDI16SP(m *RV, ins uint) error {
	imm := decodeCIb(ins)
	m.wrX(RegSp, uint64(int(m.rdX(RegSp))+imm))
	m.PC += 2
	return nil
}

func emu_C_LUI(m *RV, ins uint) error {
	imm, rd := decodeCIf(ins)
	if imm == 0 {
		return m.errIllegal(ins)
	}
	if rd != 0 && rd != 2 {
		m.wrX(rd, uint64(imm<<12))
	}
	m.PC += 2
	return nil
}

func emu_C_SRLI(m *RV, ins uint) error {
	shamt, rd := decodeCIc(ins)
	if m.xlen == 32 && shamt > 31 {
		return m.errIllegal(ins)
	}
	m.wrX(rd, m.rdX(rd)>>shamt)
	m.PC += 2
	return nil
}

func emu_C_SRAI(m *RV, ins uint) error {
	shamt, rd := decodeCIc(ins)
	var x uint64
	if m.xlen == 32 {
		if shamt > 31 {
			return m.errIllegal(ins)
		}
		x = uint64(int32(m.rdX(rd)) >> shamt)
	} else {
		x = uint64(int64(m.rdX(rd)) >> shamt)
	}
	m.wrX(rd, x)
	m.PC += 2
	return nil
}

func emu_C_ANDI(m *RV, ins uint) error {
	imm, rd := decodeCIe(ins)
	m.wrX(rd, uint64(int(m.rdX(rd))&imm))
	m.PC += 2
	return nil
}

func emu_C_SUB(m *RV, ins uint) error {
	rd, rs := decodeCRa(ins)
	m.wrX(rd, m.rdX(rd)-m.rdX(rs))
	m.PC += 2
	return nil
}

func emu_C_XOR(m *RV, ins uint) error {
	rd, rs := decodeCRa(ins)
	m.wrX(rd, m.rdX(rd)^m.rdX(rs))
	m.PC += 2
	return nil
}

func emu_C_OR(m *RV, ins uint) error {
	rd, rs := decodeCRa(ins)
	m.wrX(rd, m.rdX(rd)|m.rdX(rs))
	m.PC += 2
	return nil
}

func emu_C_AND(m *RV, ins uint) error {
	rd, rs := decodeCRa(ins)
	m.wrX(rd, m.rdX(rd)&m.rdX(rs))
	m.PC += 2
	return nil
}

func emu_C_J(m *RV, ins uint) error {
	imm := decodeCJ(ins)
	m.PC = uint64(int(m.PC) + imm)
	return nil
}

func emu_C_BEQZ(m *RV, ins uint) error {
	imm, rs := decodeCB(ins)
	if m.rdX(rs) == 0 {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 2
	}
	return nil
}

func emu_C_BNEZ(m *RV, ins uint) error {
	imm, rs := decodeCB(ins)
	if m.rdX(rs) != 0 {
		m.PC = uint64(int(m.PC) + imm)
	} else {
		m.PC += 2
	}
	return nil
}

func emu_C_SLLI(m *RV, ins uint) error {
	shamt, rd := decodeCId(ins)
	if rd != 0 && shamt != 0 {
		m.wrX(rd, m.rdX(rd)<<shamt)
	}
	m.PC += 2
	return nil
}

func emu_C_SLLI64(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_C_LWSP(m *RV, ins uint) error {
	uimm, rd := decodeCSSa(ins)
	if rd == 0 {
		return m.errIllegal(ins)
	}
	adr := uint(m.rdX(RegSp)) + uimm
	val, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(int32(val)))
	m.PC += 2
	return nil
}

func emu_C_JR(m *RV, ins uint) error {
	rs1, _ := decodeCR(ins)
	if rs1 == 0 {
		return m.errIllegal(ins)
	}
	m.PC = m.rdX(rs1)
	return nil
}

func emu_C_MV(m *RV, ins uint) error {
	rd, rs := decodeCR(ins)
	if rs != 0 {
		m.wrX(rd, m.rdX(rs))
	}
	m.PC += 2
	return nil
}

func emu_C_EBREAK(m *RV, ins uint) error {
	m.PC = m.CSR.Exception(m.PC, uint(csr.ExBreakpoint), uint(m.PC), false)
	return nil
}

func emu_C_JALR(m *RV, ins uint) error {
	rs1, _ := decodeCR(ins)
	if rs1 == 0 {
		return m.errIllegal(ins)
	}
	t := m.PC + 2
	m.PC = m.rdX(rs1)
	m.wrX(RegRa, t)
	return nil
}

func emu_C_ADD(m *RV, ins uint) error {
	rd, rs := decodeCR(ins)
	m.wrX(rd, m.rdX(rd)+m.rdX(rs))
	m.PC += 2
	return nil
}

func emu_C_SWSP(m *RV, ins uint) error {
	uimm, rs2 := decodeCSSb(ins)
	adr := uint(m.rdX(RegSp)) + uimm
	err := m.Mem.Wr32(adr, uint32(m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 2
	return nil
}

//-----------------------------------------------------------------------------
// rv32c-only

func emu_C_JAL(m *RV, ins uint) error {
	imm := decodeCJ(ins)
	m.wrX(RegRa, m.PC+2)
	m.PC = uint64(int(m.PC) + imm)
	return nil
}

//-----------------------------------------------------------------------------
// rv32fc

func emu_C_FLW(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_C_FLWSP(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_C_FSW(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_C_FSWSP(m *RV, ins uint) error {
	return m.errTodo()
}

//-----------------------------------------------------------------------------
// rv32dc

func emu_C_FLD(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_C_FLDSP(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_C_FSD(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_C_FSDSP(m *RV, ins uint) error {
	return m.errTodo()
}

//-----------------------------------------------------------------------------
// rv64i

func emu_LWU(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd32(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, uint64(val))
	m.PC += 4
	return nil
}

func emu_LD(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	val, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, val)
	m.PC += 4
	return nil
}

func emu_SD(m *RV, ins uint) error {
	imm, rs2, rs1 := decodeS(ins)
	adr := uint(int(m.rdX(rs1)) + imm)
	err := m.Mem.Wr64(adr, m.rdX(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 4
	return nil
}

func emu_ADDIW(m *RV, ins uint) error {
	imm, rs1, rd := decodeIa(ins)
	m.wrX(rd, uint64(int32(int(m.rdX(rs1))+imm)))
	m.PC += 4
	return nil
}

func emu_SLLIW(m *RV, ins uint) error {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(int32(m.rdX(rs1))<<shamt))
	m.PC += 4
	return nil
}

func emu_SRLIW(m *RV, ins uint) error {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(int32(uint32(m.rdX(rs1))>>shamt)))
	m.PC += 4
	return nil
}

func emu_SRAIW(m *RV, ins uint) error {
	shamt, rs1, rd := decodeIc(ins)
	if shamt&32 != 0 {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(int32(m.rdX(rs1))>>shamt))
	m.PC += 4
	return nil
}

func emu_ADDW(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.rdX(rs1)+m.rdX(rs2))))
	m.PC += 4
	return nil
}

func emu_SUBW(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.wrX(rd, uint64(int32(m.rdX(rs1)-m.rdX(rs2))))
	m.PC += 4
	return nil
}

func emu_SLLW(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.rdX(rs2) & 31
	m.wrX(rd, uint64(int32(m.rdX(rs1)<<shamt)))
	m.PC += 4
	return nil
}

func emu_SRLW(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.rdX(rs2) & 31
	m.wrX(rd, uint64(int32(uint32(m.rdX(rs1))>>shamt)))
	m.PC += 4
	return nil
}

func emu_SRAW(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	shamt := m.rdX(rs2) & 31
	m.wrX(rd, uint64(int32(m.rdX(rs1))>>shamt))
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv64m

func emu_MULW(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	result := int32(m.rdX(rs1) * m.rdX(rs2))
	m.wrX(rd, uint64(result))
	m.PC += 4
	return nil
}

func emu_DIVW(m *RV, ins uint) error {
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
	return nil
}

func emu_DIVUW(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	dividend := uint32(m.rdX(rs1))
	divisor := uint32(m.rdX(rs2))
	result := int32(-1)
	if divisor != 0 {
		result = int32(dividend / divisor)
	}
	m.wrX(rd, uint64(result))
	m.PC += 4
	return nil
}

func emu_REMW(m *RV, ins uint) error {
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
	return nil
}

func emu_REMUW(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	dividend := uint32(m.rdX(rs1))
	divisor := uint32(m.rdX(rs2))
	result := int32(dividend)
	if divisor != 0 {
		result = int32(dividend % divisor)
	}
	m.wrX(rd, uint64(result))
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv64a

func emu_LR_D(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_SC_D(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_AMOSWAP_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, m.rdX(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOADD_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, t+m.rdX(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOXOR_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, t^m.rdX(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOAND_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, t&m.rdX(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOOR_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, t|m.rdX(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOMIN_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, uint64(minInt64(int64(t), int64(m.rdX(rs2)))))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOMAX_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, uint64(maxInt64(int64(t), int64(m.rdX(rs2)))))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOMINU_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, minUint64(t, m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

func emu_AMOMAXU_D(m *RV, ins uint) error {
	rs2, rs1, _, rd := decodeR(ins)
	m.amo.Lock()
	adr := uint(m.rdX(rs1))
	t, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	err = m.Mem.Wr64(adr, maxUint64(t, m.rdX(rs2)))
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, t)
	m.amo.Unlock()
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv64f

func emu_FCVT_L_S(m *RV, ins uint) error {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_l_s(m.rdFS(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(x))
	m.PC += 4
	return nil
}

func emu_FCVT_LU_S(m *RV, ins uint) error {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_lu_s(m.rdFS(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrX(rd, x)
	m.PC += 4
	return nil
}

func emu_FCVT_S_L(m *RV, ins uint) error {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_s_l(int64(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

func emu_FCVT_S_LU(m *RV, ins uint) error {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_s_lu(m.rdX(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFS(rd, x)
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv64d

func emu_FCVT_L_D(m *RV, ins uint) error {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_l_d(m.rdFD(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrX(rd, uint64(x))
	m.PC += 4
	return nil
}

func emu_FCVT_LU_D(m *RV, ins uint) error {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_lu_d(m.rdFD(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrX(rd, x)
	m.PC += 4
	return nil
}

func emu_FMV_X_D(m *RV, ins uint) error {
	_, rs1, _, rd := decodeR(ins)
	m.wrX(rd, m.rdFD(rs1))
	m.PC += 4
	return nil
}

func emu_FCVT_D_L(m *RV, ins uint) error {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_d_l(int64(m.rdX(rs1)), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FCVT_D_LU(m *RV, ins uint) error {
	_, rs1, rm, rd := decodeR(ins)
	x, err := fcvt_d_lu(m.rdX(rs1), rm, m.CSR)
	if err != nil {
		return m.errIllegal(ins)
	}
	m.wrFD(rd, x)
	m.PC += 4
	return nil
}

func emu_FMV_D_X(m *RV, ins uint) error {
	_, rs1, _, rd := decodeR(ins)
	m.wrFD(rd, m.rdX(rs1))
	m.PC += 4
	return nil
}

//-----------------------------------------------------------------------------
// rv64c

func emu_C_ADDIW(m *RV, ins uint) error {
	imm, rd := decodeCIa(ins)
	if rd != 0 {
		m.wrX(rd, uint64(int32(int(m.rdX(rd))+imm)))
	} else {
		return m.errIllegal(ins)
	}
	m.PC += 2
	return nil
}

func emu_C_LDSP(m *RV, ins uint) error {
	uimm, rd := decodeCIg(ins)
	adr := uint(m.rdX(RegSp)) + uimm
	val, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	if rd != 0 {
		m.wrX(rd, val)
	} else {
		return m.errIllegal(ins)
	}
	m.PC += 2
	return nil
}

func emu_C_SDSP(m *RV, ins uint) error {
	uimm, rs2 := decodeCSSc(ins)
	adr := uint(m.rdX(RegSp)) + uimm
	err := m.Mem.Wr64(adr, m.rdX(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 2
	return nil
}

func emu_C_LD(m *RV, ins uint) error {
	uimm, rs1, rd := decodeCSa(ins)
	adr := uint(m.rdX(rs1)) + uimm
	val, err := m.Mem.Rd64(adr)
	if err != nil {
		return m.errMemory(err)
	}
	m.wrX(rd, val)
	m.PC += 2
	return nil
}

func emu_C_SD(m *RV, ins uint) error {
	uimm, rs1, rs2 := decodeCSa(ins)
	adr := uint(m.rdX(rs1)) + uimm
	err := m.Mem.Wr64(adr, m.rdX(rs2))
	if err != nil {
		return m.errMemory(err)
	}
	m.PC += 2
	return nil
}

func emu_C_SUBW(m *RV, ins uint) error {
	return m.errTodo()
}

func emu_C_ADDW(m *RV, ins uint) error {
	return m.errTodo()
}

//-----------------------------------------------------------------------------
// Integer Register Access

// wrX writes an integer register
func (m *RV) wrX(i uint, val uint64) {
	if i == 0 {
		// no writes to zero
		return
	}
	if m.xlen == 32 {
		val = uint64(uint32(val))
	}
	m.x[i] = val
}

// rdX reads an integer register
func (m *RV) rdX(i uint) uint64 {
	if i == 0 {
		// always reads as zero
		return 0
	}
	if m.xlen == 32 {
		return uint64(uint32(m.x[i]))
	}
	return m.x[i]
}

//-----------------------------------------------------------------------------
// Float Register Access

// wrFS writes a 32-bit float register.
func (m *RV) wrFS(i uint, val uint32) {
	m.f[i] = uint64(val) | upper32
}

// rdFS reads a 32-bit float register.
func (m *RV) rdFS(i uint) uint32 {
	return uint32(m.f[i])
}

// wrFD writes a 64-bit float register.
func (m *RV) wrFD(i uint, val uint64) {
	m.f[i] = val
}

// rdFD reads a 64-bit float register.
func (m *RV) rdFD(i uint) uint64 {
	return m.f[i]
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

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

// RV is a RISC-V CPU.
type RV struct {
	x      [32]uint64  // integer registers
	f      [32]uint64  // float registers
	PC     uint64      // program counter
	isa    *ISA        // ISA implemented for the CPU
	Mem    *mem.Memory // memory of the target system
	CSR    *csr.State  // CSR state
	amo    sync.Mutex  // lock for atomic operations
	lastPC uint64      // stuck PC detection
	xlen   uint        // bit length of integer registers
	err    *errBuffer  // buffer of handled/un-handled emulation errors
}

// Reset the CPU.
func (m *RV) Reset() {
	m.PC = m.Mem.Entry
	m.CSR.Reset()
	m.err.reset()
	m.lastPC = 0
}

// NewRV64 returns a 64-bit RISC-V CPU.
func NewRV64(isa *ISA, mem *mem.Memory, csr *csr.State) *RV {
	m := RV{
		xlen: 64,
		isa:  isa,
		Mem:  mem,
		CSR:  csr,
		err:  newErrBuffer(32),
	}
	m.Reset()
	return &m
}

// NewRV32 returns a 32-bit RISC-V CPU.
func NewRV32(isa *ISA, mem *mem.Memory, csr *csr.State) *RV {
	m := RV{
		xlen: 32,
		isa:  isa,
		Mem:  mem,
		CSR:  csr,
		err:  newErrBuffer(32),
	}
	m.Reset()
	return &m
}

//-----------------------------------------------------------------------------

func (m *RV) errHandler(err error) error {
	e := err.(*Error)

	// record the error
	m.err.write(e)

	// handle the error
	switch e.Type {
	case ErrIllegal, ErrCSR:
		m.PC = m.CSR.Exception(m.PC, uint(csr.ExInsIllegal), e.ins, false)
		return nil
	case ErrMemory:
		em := e.err.(*mem.Error)
		if em.Type&(mem.ErrBreak|mem.ErrEmpty) == 0 {
			m.PC = m.CSR.Exception(m.PC, uint(em.Ex), em.Addr, false)
			return nil
		}
	}

	return err
}

//-----------------------------------------------------------------------------

// Run the CPU for a single instruction.
func (m *RV) Run() error {

	// read the next instruction
	ins, err := m.Mem.RdIns(uint(m.PC))
	if err != nil {
		return m.errHandler(m.errMemory(err))
	}

	// check for break points
	err = m.Mem.GetBreak()
	if err != nil {
		return m.errMemory(err)
	}

	// lookup and emulate the instruction
	im := m.isa.lookup(ins)
	if im == nil {
		return m.errHandler(m.errIllegal(ins))
	}

	err = im.defn.emu(m, ins)
	if err != nil {
		return m.errHandler(err)
	}

	// Update the CSR registers
	m.CSR.IncInstructions()
	m.CSR.IncClockCycles(2)

	// check for breaks points
	err = m.Mem.GetBreak()
	if err != nil {
		return m.errMemory(err)
	}

	// stuck PC detection
	if m.PC == m.lastPC {
		return m.errStuckPC()
	}
	m.lastPC = m.PC

	return nil
}

//-----------------------------------------------------------------------------

// IntRegs returns a display string for the integer registers.
func (m *RV) IntRegs() string {
	reg := make([]uint, 32)
	for i := range reg {
		reg[i] = uint(m.rdX(uint(i)))
	}
	return intRegString(reg, uint(m.PC), m.xlen)
}

// FloatRegs returns a display string for the float registers.
func (m *RV) FloatRegs() string {
	return floatRegString(m.f[:])
}

// Disassemble the instruction at the address.
func (m *RV) Disassemble(addr uint) *Disassembly {
	return m.isa.Disassemble(m.Mem, addr)
}

//-----------------------------------------------------------------------------
