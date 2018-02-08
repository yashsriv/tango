package codegen

type BBLEntry struct {
	Block []IRIns
	Info  []map[SymbolTableEntry]UseInfo
}

type UseInfo struct {
	Live    bool
	NextUse int
}

var BBLList []BBLEntry

var symbolMap = make(map[SymbolTableEntry]UseInfo)

func GenBBLList(IRCode []IRIns) {
	if len(IRCode) == 0 {
		return
	}
	prevIndex := 0
	for index, ins := range IRCode {
		if arg1 := ins.Arg1; arg1 != nil {
			if _, isRegister := arg1.(SymbolTableRegisterEntry); isRegister {
				symbolMap[arg1] = UseInfo{true, -1}
			}
		}
		if arg2 := ins.Arg2; arg2 != nil {
			if _, isRegister := arg2.(SymbolTableRegisterEntry); isRegister {
				symbolMap[arg2] = UseInfo{true, -1}
			}
		}
		if dst := ins.Dst; dst != nil {
			if _, isRegister := dst.(SymbolTableRegisterEntry); isRegister {
				symbolMap[dst] = UseInfo{true, -1}
			}
		}
		if ins.Label != "" && index != prevIndex {
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
	for i := len(bbl.Block) - 1; i >= 0; i-- {
		bbl.Info[i] = make(map[SymbolTableEntry]UseInfo)
		if dst := bbl.Block[i].Dst; dst != nil {
			if _, isRegister := dst.(SymbolTableRegisterEntry); isRegister {
				bbl.Info[i][dst] = symbolMap[dst]
				symbolMap[dst] = UseInfo{false, -1}
			}
		}
		if arg1 := bbl.Block[i].Arg1; arg1 != nil {
			if _, isRegister := dst.(SymbolTableRegisterEntry); isRegister {
				bbl.Info[i][arg1] = symbolMap[arg1]
				symbolMap[arg1] = UseInfo{true, i}
			}
		}
		if arg2 := bbl.Block[i].Arg2; arg2 != nil {
			if _, isRegister := dst.(SymbolTableRegisterEntry); isRegister {
				bbl.Info[i][arg2] = symbolMap[arg2]
				symbolMap[arg2] = UseInfo{true, i}
			}
		}
	}
	return bbl
}
