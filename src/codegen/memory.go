package codegen

import "fmt"

// MemoryLocation refers to a memory location where a
// variable can be stored.
type MemoryLocation interface {
	memoryLocationDummy()
	String() string
}

// GlobalMemory refers to a location in the global memory
// for a variable.
type GlobalMemory struct {
	Location string
}

func (g GlobalMemory) memoryLocationDummy() {}

func (g GlobalMemory) String() string {
	return fmt.Sprintf("%s(,1)", g.Location)
}

// StackMemory refers to a location for a variable in the stack
type StackMemory struct {
	BaseOffset int
}

func (s StackMemory) memoryLocationDummy() {}

func (s StackMemory) String() string {
	return fmt.Sprintf("%d(%%ebp)", s.BaseOffset)
}

// NoMemory refers to an invalid memory
// This is for types which shouldn't actually require a
// memory assignment (temporaries/register variables)
type NoMemory struct {
}

func (n NoMemory) memoryLocationDummy() {}

func (n NoMemory) String() string {
	panic("nomemory variable being assigned")
}
