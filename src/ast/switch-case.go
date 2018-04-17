package ast

import (
	"fmt"
	"tango/src/codegen"
)

var switchCaseCount = 0

type caseDecl struct {
	exprList []*AddrCode
	body     *AddrCode
	isFall   bool
}

func EvalCaseDecl(a, b Attrib, isFall bool) (*caseDecl, error) {
	exprList, ok := a.([]*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[EvalCaseDecl] unable to type cast %v to []*AddrCode", a)
	}
	caseBody, err := MergeCodeList(b)
	if err != nil {
		return nil, err
	}

	return &caseDecl{exprList, caseBody, isFall}, nil
}

func EvalSwitch(a, b Attrib) (*AddrCode, error) {
	expr, ok := a.(*AddrCode)
	if !ok {
		return nil, fmt.Errorf("[EvalSwitch] unable to type cast %v to *AddrCode", a)
	}
	blocks, ok := b.([]*caseDecl)
	if !ok {
		return nil, fmt.Errorf("[EvalSwitch] unable to type cast %v to []*caseDecl", b)
	}
	code := make([]codegen.IRIns, 0)

	code = append(code, expr.Code...)

	endLbl := &codegen.TargetEntry{
		Target: fmt.Sprintf("#_switch_case_%d_end", switchCaseCount),
	}

	for i, block := range blocks {
		endBlockLbl := &codegen.TargetEntry{
			Target: fmt.Sprintf("#_switch_case_%d_%d_end", switchCaseCount, i),
		}
		if len(block.exprList) != 0 {
			// Not default case

			// Evaluate all expressions
			var orExpr *AddrCode
			for i, blockExpr := range block.exprList {
				// Eval expr
				tmp, err := RelOp(&AddrCode{Symbol: expr.Symbol}, "==", blockExpr)
				if err != nil {
					return nil, err
				}
				if i == 0 {
					orExpr = tmp
				} else {
					orExpr, err = OrOp(orExpr, tmp)
					if err != nil {
						return nil, err
					}
				}
			}
			code = append(code, orExpr.Code...)
			code = append(code, codegen.IRIns{
				Typ: codegen.CBR,
				Op:  codegen.BRNEQ,
				Dst: endBlockLbl,
				Arg1: &codegen.LiteralEntry{
					Value: 1,
				},
				Arg2: orExpr.Symbol,
			})
		}
		// Add body
		code = append(code, block.body.Code...)

		// Not fallthrough
		if !block.isFall {
			// Jump to end at the end of this switch case block
			code = append(code, codegen.IRIns{
				Typ:  codegen.JMP,
				Op:   codegen.JMPO,
				Arg1: endLbl,
			})
		}
		code = append(code, codegen.IRIns{
			Typ: codegen.LBL,
			Dst: endBlockLbl,
		})
	}
	code = append(code, codegen.IRIns{
		Typ: codegen.LBL,
		Dst: endLbl,
	})
	addrcode := &AddrCode{
		Code: code,
	}
	switchCaseCount++
	return addrcode, nil
}
