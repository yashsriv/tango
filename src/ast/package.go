package ast

import (
	"errors"
	"fmt"
	"tango/src/token"

	"github.com/awalterschulze/gographviz"
	uuid "github.com/satori/go.uuid"
)

// PackageName represents a package name
type PackageName struct {
	Value string
}

var _ ourAttrib = (*PackageName)(nil)

// PackageClause represents a package clause
type PackageClause struct {
	Name *PackageName
}

var _ ourAttrib = (*PackageClause)(nil)

// GenGraph is used to generate a graph
func (s *PackageName) GenGraph(g gographviz.Interface) string {
	u1 := fmt.Sprintf("%q", uuid.NewV4().String())
	g.AddNode("main", u1, map[string]string{"label": "PackageName"})
	u2 := fmt.Sprintf("%q", uuid.NewV4().String())
	g.AddNode("main", u2, map[string]string{"label": s.Value, "shape": "square"})
	g.AddEdge(u1, u2, true, nil)
	return u1
}

// GenGraph is used to generate a graph
func (s *PackageClause) GenGraph(g gographviz.Interface) string {
	u1 := fmt.Sprintf("%q", uuid.NewV4().String())
	g.AddNode("main", u1, map[string]string{"label": "PackageClause"})
	child := s.Name.GenGraph(g)
	g.AddEdge(u1, child, true, nil)
	return u1
}

// NewPackageClause creates an instance of a package clause
func NewPackageClause(packageName Attrib) (Attrib, error) {
	tempPackageName, ok := packageName.(*PackageName)
	if !ok {
		return nil, errors.New("Expected a package name")
	}
	return &PackageClause{
		Name: tempPackageName,
	}, nil
}

// NewPackageName creates an instance of a package name
func NewPackageName(id Attrib) (Attrib, error) {
	tokenVal, ok := id.(*token.Token)
	if !ok || tokenVal.Type != token.TokMap.Type("identifier") {
		return nil, errors.New("Expected an identifier token")
	}
	return &PackageName{
		Value: string(tokenVal.Lit),
	}, nil
}
