package lib

import "testing"
import "math/rand"
import "bytes"

// Returns a *RunContext with 4096 of memory.
// Randomizes contents of registers and memory.
// PC_Init is set to zero.
func setupTest() (*RunContext, *bytes.Buffer) {
	rc := &RunContext{
		PC_Init:   0x0,
		Registers: make([]uint32, 16),
	}

	for i, _ := range rc.Registers {
		rc.Registers[i] = rand.Uint32()
	}

	return rc, &bytes.Buffer{}
}

func setMem(rc *RunContext, buf *bytes.Buffer) {
	// Write 4 bytes to ensure that there are always 4 bytes left.
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00})
	rc.Memory = buf.Bytes()
}

func setCode(t *testing.T, rc *RunContext, buf *bytes.Buffer, code []string) (int, error) {
	i, err := Put_Lines(buf, code)

	if err != nil {
		t.Errorf("ERR: %s", err.Error())
		return i, err
	}

	setMem(rc, buf)

	return -1, nil
}

func checkRegisters(t *testing.T, rc *RunContext, expectedValues []uint32) {
	for i, expectedValue := range expectedValues {
		got := rc.Registers[i]

		if got != expectedValue {
			t.Errorf("Expected REG[%x] to be %x but got %x", i, expectedValue, got)
		}
	}
}

func checkPCEnd(t *testing.T, rc *RunContext, expected_pc_end uint32) {
	if rc.PC_End != expected_pc_end {
		t.Errorf("Expected PC to be %x but got %x", expected_pc_end, rc.PC_End)
	}
}

func checkMemory(t *testing.T, rc *RunContext, expected_mem []uint8) {
	for i, _ := range expected_mem {
		if rc.Memory[i] != expected_mem[i] {
			t.Errorf("Expected MEM[%x] to be %x but got %x", i, expected_mem[i], rc.Memory[i])
		}
	}
}

func checkRegsMem(t *testing.T, rc *RunContext, expected_regs []uint32, expected_mem []uint8) {
	checkRegisters(t, rc, expected_regs)
	checkMemory(t, rc, expected_mem)
}

func copyRegisters(rc *RunContext) []uint32 {
	buf := make([]uint32, len(rc.Registers))
	copy(buf, rc.Registers)
	return buf
}

func copyMemory(rc *RunContext) []uint8 {
	buf := make([]uint8, len(rc.Memory))
	copy(buf, rc.Memory)
	return buf
}

func copyRegsMem(rc *RunContext) ([]uint32, []uint8) {
	return copyRegisters(rc), copyMemory(rc)
}

func TestHLT(t *testing.T) {
	rc, mem := setupTest()
	Put_Regs(mem, OPC_HLT, 0x00, 0x00)
	setMem(rc, mem)

	Run(rc)

	checkPCEnd(t, rc, 0x02)
}

func TestNOP(t *testing.T) {
	rc, mem := setupTest()
	Put_Regs(mem, OPC_NOP, 0x00, 0x00)
	Put_Regs(mem, OPC_HLT, 0x00, 0x00)
	setMem(rc, mem)

	expected_regs, expected_mem := copyRegsMem(rc)

	Run(rc)

	checkRegsMem(t, rc, expected_regs, expected_mem)
	checkPCEnd(t, rc, 0x04)
}

func TestLDL(t *testing.T) {
	rc, mem := setupTest()
	Put_Imm16(mem, OPC_LDL, 0x03, 0xFEED)
	Put_Regs(mem, OPC_HLT, 0x00, 0x00)

	setMem(rc, mem)

	rc.Registers[0x03] = 0xDEADBEEF

	e_regs, e_mem := copyRegsMem(rc)
	e_regs[0x03] = 0xDEADFEED

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
}

func TestLDH(t *testing.T) {
	rc, mem := setupTest()
	Put_Imm16(mem, OPC_LDH, 0x05, 0xFEED)
	Put_Regs(mem, OPC_HLT, 0x00, 0x00)
	setMem(rc, mem)

	rc.Registers[0x05] = 0xDEADBEEF

	e_regs, e_mem := copyRegsMem(rc)
	e_regs[0x05] = 0xFEEDBEEF

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
}

func TestXOR(t *testing.T) {
	rc, mem := setupTest()
	setCode(t, rc, mem, []string{
		"xor rc rd",
		"hlt",
	})

	e_regs, e_mem := copyRegsMem(rc)
	e_regs[REG_RD] ^= e_regs[REG_RC]

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
}

func TestADD(t *testing.T) {
	rc, mem := setupTest()
	setCode(t, rc, mem, []string{
		"add re rf",
		"hlt",
	})

	rc.Registers[REG_RE] = 0xBEEFDAD0
	rc.Registers[REG_RF] = 0xFEEDDEAF

	e_regs, e_mem := copyRegsMem(rc)
	e_regs[REG_RF] += e_regs[REG_RE]

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
}

func TestBEQ(t *testing.T) {
	// Beware that PC points to the next instruction
	// not the current one so when checking if the branch
	// jumped to hlt at 14 the PC will be offset by two (since
	// hlt is a two byte instruction).

	// Zero offset, Branch Taken
	rc, mem := setupTest()
	setCode(t, rc, mem, []string{
		"beq ra rb rc 0", // 0
		"nop",            // 4
		"nop",            // 6
		"hlt",            // 8
		"nop",            // 10
		"nop",            // 12
		"hlt",            // 14
	})

	rc.Registers[REG_RA] = 0xFEEDFEED
	rc.Registers[REG_RB] = 0xFEEDFEED
	rc.Registers[REG_RC] = 0x0000000E

	e_regs, e_mem := copyRegsMem(rc)

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
	checkPCEnd(t, rc, 0x0E+0x02)

	// Zero offset, Branch Not Taken
	rc, mem = setupTest()
	setCode(t, rc, mem, []string{
		"beq ra rb rc 0", // 0
		"nop",            // 4
		"nop",            // 6
		"hlt",            // 8
		"nop",            // 10
		"nop",            // 12
		"hlt",            // 14
	})

	rc.Registers[REG_RA] = 0xFEEDF1ED
	rc.Registers[REG_RB] = 0xFEEDFEED
	rc.Registers[REG_RC] = 0x0000000E

	e_regs, e_mem = copyRegsMem(rc)

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
	checkPCEnd(t, rc, 0x08+0x02)

	// Positive offset, Branch Taken
	rc, mem = setupTest()
	setCode(t, rc, mem, []string{
		"beq ra rb rc 2", // 0
		"nop",            // 4
		"nop",            // 6
		"hlt",            // 8
		"nop",            // 10
		"nop",            // 12
		"hlt",            // 14
	})

	rc.Registers[REG_RA] = 0xFEEDFEED
	rc.Registers[REG_RB] = 0xFEEDFEED
	rc.Registers[REG_RC] = 0x0000000C

	e_regs, e_mem = copyRegsMem(rc)

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
	checkPCEnd(t, rc, 0x0E+0x02)

	// Negative offset, Branch Taken
	rc, mem = setupTest()
	setCode(t, rc, mem, []string{
		"beq ra rb rc 802", // 0
		"nop",              // 4
		"nop",              // 6
		"hlt",              // 8
		"nop",              // 10
		"nop",              // 12
		"hlt",              // 14
	})

	rc.Registers[REG_RA] = 0xFEEDFEED
	rc.Registers[REG_RB] = 0xFEEDFEED
	rc.Registers[REG_RC] = 0x00000010

	e_regs, e_mem = copyRegsMem(rc)

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
	checkPCEnd(t, rc, 0x0E+0x02)

	// Negative offset, Branch Not Taken
	rc, mem = setupTest()
	setCode(t, rc, mem, []string{
		"beq ra rb rc 802", // 0
		"nop",              // 4
		"nop",              // 6
		"hlt",              // 8
		"nop",              // 10
		"nop",              // 12
		"hlt",              // 14
	})

	rc.Registers[REG_RA] = 0xFEEDF1ED
	rc.Registers[REG_RB] = 0xFEEDFEED
	rc.Registers[REG_RC] = 0x00000010

	e_regs, e_mem = copyRegsMem(rc)

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
	checkPCEnd(t, rc, 0x08+0x02)

	// Positive offset, Branch Not Taken
	rc, mem = setupTest()
	setCode(t, rc, mem, []string{
		"beq ra rb rc 2", // 0
		"nop",            // 4
		"nop",            // 6
		"hlt",            // 8
		"nop",            // 10
		"nop",            // 12
		"hlt",            // 14
	})

	rc.Registers[REG_RA] = 0xFEEDF1ED
	rc.Registers[REG_RB] = 0xFEEDFEED
	rc.Registers[REG_RC] = 0x0000000C

	e_regs, e_mem = copyRegsMem(rc)

	Run(rc)

	checkRegsMem(t, rc, e_regs, e_mem)
	checkPCEnd(t, rc, 0x08+0x02)
}
