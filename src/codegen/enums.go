package codegen

// IRType is the type of an IR Instruction
type IRType int

// Enum IRType
const (
	INV IRType = iota
	LBL
	BOP
	LOP
	SOP
	DOP
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
	BAND      = "&"
	BOR       = "|"
	XOR       = "^"

	// Handled same as BAND and BOR
	AND = "&&"
	OR  = "||"

	TAKE = "take"
	PUT  = "put"
)

// Division Operations
const (
	DIV = "/"
	REM = "%"
)

// Shift Ops
const (
	BSL = "<<"
	BSR = ">>"
)

// Logical Ops
const (
	LT  = "<"
	LTE = "<="
	GT  = ">"
	GTE = ">="
	EQ  = "=="
	NEQ = "!="
)

// Unary Ops
const (
	NEG  IROp = "neg"
	NOT       = "!"
	BNOT      = "not"
	VAL       = "val"
	ADDR      = "addr"
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
	ARG            = "arg"
	INC            = "inc"
	DEC            = "dec"
	PARAM          = "param"
	RET            = "ret"
	RETI           = "reti"
	SETRET         = "setret"
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

	if op == NEG || op == NOT || op == BNOT ||
		op == VAL || op == ADDR {
		return UOP
	}

	if op == PARAM || op == CALL || op == RET || op == RETI ||
		op == SETRET || op == HALT || op == PRINTINT || op == PRINTCHAR ||
		op == PRINTSTR || op == SCANINT || op == SCANCHAR ||
		op == SCANSTR || op == INC || op == DEC {
		return KEY
	}

	if op == LT || op == LTE || op == GT ||
		op == GTE || op == EQ || op == NEQ {
		return LOP
	}

	if op == BSL || op == BSR {
		return SOP
	}

	if op == DIV || op == REM {
		return DOP
	}

	if op == ADD || op == SUB || op == MUL ||
		op == AND || op == OR || op == BAND ||
		op == BOR || op == XOR || op == TAKE ||
		op == PUT {
		return BOP
	}

	return INV
}
