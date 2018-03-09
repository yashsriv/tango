package html

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"runtime"

	"tango/src/ast"
	"tango/src/token"
)

func nodeIsNil(args ...interface{}) bool {
	var n *ast.Node
	if len(args) == 1 {
		n, _ = args[0].(*ast.Node)
	}
	return n == nil
}

func processStack(args ...interface{}) string {
	var stack ast.Stack
	if len(args) == 1 {
		stack, _ = args[0].(ast.Stack)
	}
	var str string
	for _, attrib := range stack {
		switch v := attrib.(type) {
		case *token.Token:
			str += fmt.Sprintf("<span class=\"token\">%s</span> ", v)
			if token.TokMap.Type("stmt_end") == v.Type {
				str += "\n"
			}
		case *ast.Node:
			str += fmt.Sprintf("<span class=\"node\">%s</span> ", v)
		default:
			panic(fmt.Sprintf("Unknown Type: %T", v))
		}
	}
	return str
}

func isToken(args ...interface{}) bool {
	ok := false
	if len(args) == 1 {
		_, ok = args[0].(*token.Token)
	}
	return ok
}

func isNode(args ...interface{}) bool {
	ok := false
	if len(args) == 1 {
		_, ok = args[0].(*ast.Node)
	}
	return ok
}

// Output writes the html
func Output(entries []Entry, wr io.Writer) {
	// TODO: Load a template html file and populate it
	// See https://astaxie.gitbooks.io/build-web-application-with-golang/en/07.4.html
	_, filename, _, _ := runtime.Caller(0)
	stackTemplName := filepath.Join(filepath.Dir(filename), "stack_templ.html")
	mainTemplName := filepath.Join(filepath.Dir(filename), "main_templ.html")

	t := template.New("test template")
	t = t.Funcs(template.FuncMap{
		"nodeNil":      nodeIsNil,
		"processStack": processStack,
		"isToken":      isToken,
		"isNode":       isNode,
	})
	t, err := t.ParseFiles(stackTemplName, mainTemplName)
	if err != nil {
		panic(err)
	}
	fmt.Println(t.DefinedTemplates())

	err = t.ExecuteTemplate(wr, "main_templ.html", entries)
	if err != nil {
		panic(err)
	}

	for _, val := range entries {
		fmt.Printf("%s _%s_ %s\n", val.Prefix, val.Node, val.Suffix)
	}
}
