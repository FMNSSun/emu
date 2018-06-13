package lib

/**

+----+------------+---------+---------+
| 00 | OPCODE [6] | SRC [4] | DST [4] |
+----+------------+---------+---------+

+----+------------+-------+-------+---------+----------+
| 01 | OPCODE [6] | A [4] | B [4] | DST [4] | IMM [12] |
+----+------------+-------+-------+---------+----------+

+----+------------+-----+---------+------------+
| 11 | OPCODE [6] | [4] | REG [4] | IMM16 [16] |
+----+------------+-----+---------+------------+
*/

const IF_REGS = uint8(0x00)
const IF_IMM12 = uint8(0x01)
const IF_IMM16 = uint8(0x03)

// 00
const OPC_HLT = uint8(0x00 | 0x00) // HALT
const OPC_NOP = uint8(0x00 | 0x01) // No Operation
const OPC_XOR = uint8(0x00 | 0x02) // XOR
const OPC_ADD = uint8(0x00 | 0x03) // ADD

// 01
const OPC_BEQ = uint8(0x40 | 0x00) // BEQ

// 11
const OPC_LDL = uint8(0xc0 | 0x00) // LoaD Low (to register)
const OPC_LDH = uint8(0xc0 | 0x01) // LoaD High (to register)
const OPC_WML = uint8(0xc0 | 0x02) // Write Memory Low (to memory)
const OPC_WMH = uint8(0xc0 | 0x03) // Write Memory High (to memory)
const OPC_WMB = uint8(0xc0 | 0x04) // Write Memory Byte (to memory)
const OPC_LDC = uint8(0xc0 | 0x05) // LoaD Constant (to register, signed, relative to PC)

const MASK_L16 = 0x0000FFFF
const MASK_H16 = 0xFFFF0000

const REG_RA = uint8(0x00)
const REG_RB = uint8(0x01)
const REG_RC = uint8(0x02)
const REG_RD = uint8(0x03)
const REG_RE = uint8(0x04)
const REG_RF = uint8(0x05)
const REG_RG = uint8(0x06)
const REG_RH = uint8(0x07)
const REG_RI = uint8(0x08)
const REG_RJ = uint8(0x09)
const REG_RK = uint8(0x0A)
const REG_RL = uint8(0x0B)
const REG_RM = uint8(0x0C)
const REG_RN = uint8(0x0D)
const REG_RO = uint8(0x0E)
const REG_RP = uint8(0x0F)

var OPC_Table_Opc2Name map[uint8]string = map[uint8]string{
	OPC_HLT: "hlt",
	OPC_NOP: "nop",
	OPC_LDL: "ldl",
	OPC_LDH: "ldh",
	OPC_XOR: "xor",
	OPC_ADD: "add",
	OPC_BEQ: "beq",
}

var OPC_Table_Name2Opc map[string]uint8 = map[string]uint8{
	"hlt": OPC_HLT,
	"nop": OPC_NOP,
	"ldl": OPC_LDL,
	"ldh": OPC_LDH,
	"xor": OPC_XOR,
	"add": OPC_ADD,
	"beq": OPC_BEQ,
}

var REG_Table_Reg2Name map[uint8]string = map[uint8]string{
	REG_RA: "ra",
	REG_RB: "rb",
	REG_RC: "rc",
	REG_RD: "rd",
	REG_RE: "re",
	REG_RF: "rf",
	REG_RG: "rg",
	REG_RH: "rh",
	REG_RI: "ri",
	REG_RJ: "rj",
	REG_RK: "rk",
	REG_RL: "rl",
	REG_RM: "rm",
	REG_RN: "rn",
	REG_RO: "ro",
	REG_RP: "rp",
}

var REG_Table_Name2Reg map[string]uint8 = map[string]uint8{
	"ra": REG_RA,
	"rb": REG_RB,
	"rc": REG_RC,
	"rd": REG_RD,
	"re": REG_RE,
	"rf": REG_RF,
	"rg": REG_RG,
	"rh": REG_RH,
	"ri": REG_RI,
	"rj": REG_RJ,
	"rk": REG_RK,
	"rl": REG_RL,
	"rm": REG_RM,
	"rn": REG_RN,
	"ro": REG_RO,
	"rp": REG_RP,
}
