package codegen

// IRType is the type of an IR Instruction
type IRType int

// Enum IRType
const (
	INV IRType = iota
	LBL
	BOP
	UOP
	CBR
	JMP
	ASN
	KEY
)

// IROp is the type of an operation
type IROp string

// Binary Ops
const (
	ADD  IROp = "+"
	SUB       = "-"
	MUL       = "*"
	DIV       = "/"
	REM       = "%"
	BSL       = "<<"
	BSR       = ">>"
	AND       = "&&"
	OR        = "||"
	BAND      = "&"
	BOR       = "|"
	LT        = "<"
	LTE       = "<="
	GT        = ">"
	GTE       = ">="
	EQ        = "=="
	NEQ       = "!="
	XOR       = "^"
	TAKE      = "take"
	PUT       = "put"
)

// Unary Ops
const (
	NEG   IROp = "neg"
	PARAM      = "param"
	NOT        = "!"
	INC        = "inc"
	DEC        = "dec"
	BNOT       = "not"
	VAL        = "val"
	ADDR       = "addr"
)

// Assignment Operation
const (
	ASNO IROp = "="
)

// Branch Operations
const (
	JMPO  IROp = "jmp"
	BREQ       = "breq"
	BRNEQ      = "brneq"
	BRLT       = "brlt"
	BRLTE      = "brlte"
	BRGT       = "brgt"
	BRGTE      = "brgte"
)

// Key Operations
const (
	CALL      IROp = "call"
	RET            = "ret"
	HALT           = "halt"
	PRINTINT       = "printi"
	PRINTCHAR      = "printc"
	PRINTSTR       = "prints"
	SCANINT        = "scani"
	SCANCHAR       = "scanc"
	SCANSTR        = "scans"
)

// GetType of an IROp
func GetType(op IROp) IRType {

	if op == JMPO {
		return JMP
	}

	if op == ASNO {
		return ASN
	}

	if op == BREQ || op == BRNEQ || op == BRLT ||
		op == BRLTE || op == BRGT || op == BRGTE {
		return CBR
	}

	if op == NEG || op == NOT || op == INC ||
		op == DEC || op == BNOT || op == VAL ||
		op == ADDR {
		return UOP
	}

	if op == PARAM || op == CALL || op == RET ||
		op == HALT || op == PRINTINT || op == PRINTCHAR ||
		op == PRINTSTR || op == SCANINT || op == SCANCHAR ||
		op == SCANSTR {
		return KEY
	}

	if op == ADD || op == SUB || op == MUL ||
		op == DIV || op == REM || op == BSL ||
		op == BSR || op == AND || op == OR ||
		op == BAND || op == BOR || op == LT ||
		op == LTE || op == GT || op == GTE ||
		op == EQ || op == NEQ || op == XOR ||
		op == TAKE || op == PUT {
		return BOP
	}

	return INV
}
