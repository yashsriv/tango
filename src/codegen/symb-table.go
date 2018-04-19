package codegen

// symbolTable represents a table
type symbolTable struct {
	symbolMap map[string]SymbolTableEntry
	parent    *symbolTable
	offset    int
	curOffset int
}

// SymbolTable is current table
var SymbolTable *symbolTable

// rootTable refers to the global rootTable
var rootTable *symbolTable

// tableStack maintains a stack of symbol tables
var tableStack []*symbolTable

// Initialize data structures
func init() {
	rootTable = &symbolTable{
		symbolMap: make(map[string]SymbolTableEntry),
		parent:    nil,
	}
	SymbolTable = rootTable
}

// push to table stack
func pushToStack() {
	tableStack = append(tableStack, SymbolTable)
}

// pop from table stack
func popFromStack() (table *symbolTable, err error) {
	l := len(tableStack)
	if l == 0 {
		err = ErrEmptyTableStack
		return
	}
	table = tableStack[l-1]
	tableStack = tableStack[:l-1]
	return
}

func (s *symbolTable) InsertSymbol(key string, value SymbolTableEntry) error {
	// Check if already exists in current scope
	_, ok := s.symbolMap[key]
	if ok {
		return ErrAlreadyExists
	}

	// Insert otherwise
	s.symbolMap[key] = value

	return nil
}

// CreateAddrDescEntry creates an addrDesc entry for a variable
func CreateAddrDescEntry(value *VariableEntry) {
	addrDesc[value] = address{
		regLocation: "",
		memLocation: value.MemoryLocation,
	}
}

func (s *symbolTable) GetSymbol(key string) (SymbolTableEntry, error) {
	// Check if symbol exists in current scope
	x, ok := s.symbolMap[key]
	if !ok {
		if s.parent != nil {
			// If not, check in higher scopes
			return s.parent.GetSymbol(key)
		}
		return nil, ErrDoesntExist
	}
	// If highest scope, then result of this layer is the result
	return x, nil
}

func (s *symbolTable) IsRoot() bool {
	return s.parent == nil
}

func (s *symbolTable) Alloc(size int) int {
	s.curOffset -= size
	return s.curOffset
}

func (s *symbolTable) UnAlloc(size int) {
	s.curOffset += size
}

// NewScope creates a new scope
func NewScope() {
	pushToStack()
	SymbolTable = &symbolTable{
		symbolMap: make(map[string]SymbolTableEntry),
		parent:    SymbolTable,
		offset:    SymbolTable.curOffset,
		curOffset: SymbolTable.curOffset,
	}
}

// EndScope ends a scope
func EndScope() (err error) {
	SymbolTable, err = popFromStack()
	return
}
