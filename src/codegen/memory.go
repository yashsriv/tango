package codegen

// MemoryLocation refers to a memory location where a
// variable can be stored.
type MemoryLocation interface {
	memoryLocationDummy()
}

// GlobalMemory refers to a location in the global memory
// for a variable.
type GlobalMemory struct {
	Location string
}

func (g GlobalMemory) memoryLocationDummy() {}

// StackMemory refers to a location for a variable in the stack
type StackMemory struct {
	BaseOffset int
}

func (s StackMemory) memoryLocationDummy() {}
