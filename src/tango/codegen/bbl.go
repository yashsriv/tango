package codegen

import "fmt"

// BBLEntry is Single entry in BBL List
type BBLEntry struct {
	Block []IRIns
	Info  []map[*SymbolTableRegisterEntry]UseInfo
}

func (b *BBLEntry) String() string {
	repr := "\n<BBL Begin>\n"
	for _, ins := range b.Block {
		repr += fmt.Sprintf("%s\n", ins.String())
	}
	repr += "<BBL End>\n\n"
	return repr
}

// UseInfo stores life and next Use Information of a variable
type UseInfo struct {
	Live    bool
	NextUse int
}

// BBLList is the list of all the Basic Blocks
var BBLList []BBLEntry

var symbolInfo = make(map[*SymbolTableRegisterEntry]UseInfo)

// GenBBLList takes the IRCode (list of IRIns) as input & creates list of basic blocks
func GenBBLList(IRCode []IRIns) {
	if len(IRCode) == 0 {
		return
	}
	prevIndex := 0
	for index, ins := range IRCode {
		if arg1 := ins.Arg1; arg1 != nil {
			if arg1, isRegister := arg1.(*SymbolTableRegisterEntry); isRegister {
				symbolInfo[arg1] = UseInfo{true, -1}
			}
		}
		if arg2 := ins.Arg2; arg2 != nil {
			if arg2, isRegister := arg2.(*SymbolTableRegisterEntry); isRegister {
				symbolInfo[arg2] = UseInfo{true, -1}
			}
		}
		if dst := ins.Dst; dst != nil {
			if dst, isRegister := dst.(*SymbolTableRegisterEntry); isRegister {
				symbolInfo[dst] = UseInfo{true, -1}
			}
		}
		if ins.Typ == LBL && index != prevIndex {
			bbl := BBLEntry{Block: IRCode[prevIndex:index]}
			bbl = addUseInfo(bbl)
			BBLList = append(BBLList, bbl)
			prevIndex = index
		} else if ins.Typ == CBR || ins.Typ == JMP || ins.Typ == KEY || index == len(IRCode)-1 {
			bbl := BBLEntry{Block: IRCode[prevIndex : index+1]}
			bbl = addUseInfo(bbl)
			BBLList = append(BBLList, bbl)
			if index != len(IRCode)-1 {
				prevIndex = index + 1
			}
		}
	}
}

// Adds Operands' UseInfo in the BBL
func addUseInfo(bbl BBLEntry) BBLEntry {
	bbl.Info = make([]map[*SymbolTableRegisterEntry]UseInfo, len(bbl.Block))
	infomap := make(map[*SymbolTableRegisterEntry]UseInfo)
	for i := len(bbl.Block) - 1; i >= 0; i-- {
		if dst := bbl.Block[i].Dst; dst != nil {
			if dst, isRegister := dst.(*SymbolTableRegisterEntry); isRegister {
				infomap[dst] = symbolInfo[dst]
				symbolInfo[dst] = UseInfo{false, -1}
			}
		}
		if arg1 := bbl.Block[i].Arg1; arg1 != nil {
			if arg1, isRegister := arg1.(*SymbolTableRegisterEntry); isRegister {
				infomap[arg1] = symbolInfo[arg1]
				symbolInfo[arg1] = UseInfo{true, i}
			}
		}
		if arg2 := bbl.Block[i].Arg2; arg2 != nil {
			if arg2, isRegister := arg2.(*SymbolTableRegisterEntry); isRegister {
				infomap[arg2] = symbolInfo[arg2]
				symbolInfo[arg2] = UseInfo{true, i}
			}
		}
		ninfomap := make(map[*SymbolTableRegisterEntry]UseInfo)
		for k, v := range infomap {
			ninfomap[k] = v
		}
		bbl.Info[i] = ninfomap
	}
	return bbl
}
