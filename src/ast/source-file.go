package ast

import (
	"errors"
	"fmt"

	"github.com/awalterschulze/gographviz"
	uuid "github.com/satori/go.uuid"
)

// SourceFile represents a source file
type SourceFile struct {
	PackageClause *PackageClause
}

// var _ ourAttrib = (*SourceFile)(nil)

// GenGraph is used to generate a graph
func (s *SourceFile) GenGraph(g gographviz.Interface) string {
	u1 := fmt.Sprintf("%q", uuid.NewV4().String())
	g.AddNode("main", u1, map[string]string{"label": "SourceFile"})
	child := s.PackageClause.GenGraph(g)
	g.AddEdge(u1, child, true, nil)
	return u1
}

// NewSourceFile creates a new source file
func NewSourceFile(a Attrib) (Attrib, error) {
	packageClause, ok := a.(*PackageClause)
	if !ok {
		return nil, errors.New("Expected a package clause")
	}
	return &SourceFile{
		PackageClause: packageClause,
	}, nil
}
