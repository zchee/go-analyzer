// Copyright 2022 The go-analyzer Authors
// SPDX-License-Identifier: BSD-3-Clause

package pkgerrors

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/zchee/go-analyzer/packagefact"
)

const Doc = `This analyzer analyzes and rewrites the github.com/pkg/errors (that has been deprecated) to the fmt.Errorf with %w verb provided after the go1.13.`

var Analyzer = &analysis.Analyzer{
	Name: "pkgerrors",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		packagefact.Analyzer,
		inspect.Analyzer,
	},
}

const (
	errosPath    = "errors"
	pkgerrosPath = "github.com/pkg/errors"
)

func run(pass *analysis.Pass) (interface{}, error) {
	for _, fact := range pass.AllPackageFacts() {
		pass.ImportPackageFact(fact.Package, fact.Fact)
	}
	for _, fact := range pass.AllObjectFacts() {
		pass.ImportObjectFact(fact.Object, fact.Fact)
	}

	inspected := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspected.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		// checks call expression is error
		if !isError(pass.TypesInfo.TypeOf(call)) {
			return
		}

		// filtered only of pkg/errors call expressions
		fnName, ok := isPkgErrorsCall(pass.TypesInfo, call)
		if !ok {
			return
		}

		// make a copy of the function declaration to avoid mutating the AST.
		callCopy := &ast.CallExpr{}
		*callCopy = *call
		callCopy.Args = call.Args

		var category string
		switch fnName {
		case "Cause": // errors.Cause
			callCopy.Fun.(*ast.SelectorExpr).Sel.Name = "Unwrap"
			category = "Cause"

		case "New": // errors.New
			// TODO(zchee): supprot replace stdlib errors.New
			return

		case "Errorf": // errors.Errorf
			callCopy.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name = "fmt"
			callCopy.Args = fixArgs(call.Args)
			category = "Errorf"

		case "WithStack": // errors.WithStack
			// not supported
			return

		case "WithMessage", "WithMessagef", "Wrap", "Wrapf": // errors.WithMessage{f}, errors.Wrap{f}
			callCopy.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name = "fmt"
			callCopy.Fun.(*ast.SelectorExpr).Sel.Name = "Errorf"
			callCopy.Args = reorderArgs(call.Args)
			category = fnName
		}

		var buf bytes.Buffer
		if err := format.Node(&buf, token.NewFileSet(), callCopy); err != nil {
			return
		}

		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			End:      call.End(),
			Category: category,
			Message:  fmt.Sprintf("found use location of the deprecated %s", pkgerrosPath),
			SuggestedFixes: []analysis.SuggestedFix{{
				Message: "Use fmt.Errorf with %%w verb instead",
				TextEdits: []analysis.TextEdit{{
					Pos:     callCopy.Pos(),
					End:     callCopy.End(),
					NewText: buf.Bytes(),
				}},
			}},
		})
	})

	return nil, nil
}

func rewriteImport(pass *analysis.Pass, ispec *ast.ImportSpec, imps map[string]bool) {
	if imps[unquote(ispec.Path.Value)] {
		return
	}
	defer func() {
		imps[unquote(ispec.Path.Value)] = true
	}()

	// make a copy of the function declaration to avoid mutating the AST.
	impCopy := &ast.ImportSpec{}
	*impCopy = *ispec

	switch {
	case ispec.Path.Value == strconv.Quote(pkgerrosPath):
		impCopy.Path = &ast.BasicLit{
			ValuePos: ispec.Path.ValuePos,
			Kind:     token.STRING,
			Value:    strconv.Quote(errosPath),
		}
		// overrides impCopy.Path.Pos
		impCopy.Path.Value = strconv.Quote(errosPath)

		var buf bytes.Buffer
		if err := format.Node(&buf, token.NewFileSet(), impCopy); err != nil {
			return
		}

		textEdits := []analysis.TextEdit{
			{
				Pos:     ispec.Pos(),
				End:     ispec.End(),
				NewText: buf.Bytes(),
			},
			{
				Pos:     ispec.Pos(),
				End:     ispec.End(),
				NewText: nil,
			},
		}
		if imps[errosPath] {
			// delete "github.com/pkg/errors" import if imported "errors" package
			for i := range textEdits {
				textEdits[i].NewText = nil
			}
		}

		pass.Report(analysis.Diagnostic{
			Pos:     ispec.Pos(),
			End:     ispec.End(),
			Message: fmt.Sprintf("found use %q package", pkgerrosPath),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message:   "Use fmt packace instead",
					TextEdits: textEdits,
				},
			},
		})
	}
}

// isError reports whether the typ is an error type.
func isError(typ types.Type) bool {
	if typ == nil {
		return false
	}

	return typ.String() == "error" || typ.Underlying().String() == "error"
}

func isPkgErrorsCall(info *types.Info, call *ast.CallExpr) (string, bool) {
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		obj := info.ObjectOf(fn.Sel)
		return obj.Name(), isPkgErrorsFunc(obj)

	case *ast.Ident:
		if declExpr, ok := findExpr(fn).(*ast.SelectorExpr); ok {
			obj := info.ObjectOf(declExpr.Sel)
			return obj.Name(), isPkgErrorsFunc(obj)
		}
	}

	return "", false
}

func isPkgErrorsFunc(obj types.Object) bool {
	if vendorlessPath(obj.Pkg().Path()) != pkgerrosPath {
		return false
	}

	switch obj.Name() {
	case
		"Cause",        // errors.Cause
		"New",          // errors.New
		"Errorf",       // errors.Errorf
		"WithMessage",  // errors.WithMessage
		"WithMessagef", // errors.WithMessagef
		"WithStack",    // errors.WithStack
		"Wrap",         // errors.Wrap
		"Wrapf":        // errors.Wrapf
		return true
	}

	return false
}

func findExpr(arg *ast.Ident) ast.Expr {
	if arg.Obj == nil {
		return nil
	}

	switch as := arg.Obj.Decl.(type) {
	case *ast.AssignStmt:
		if len(as.Lhs) != len(as.Rhs) {
			return nil
		}

		for i, lhs := range as.Lhs {
			lid, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}
			if lid.Obj == arg.Obj {
				return as.Rhs[i]
			}
		}

	case *ast.ValueSpec:
		if len(as.Names) != len(as.Values) {
			return nil
		}

		for i, name := range as.Names {
			if name.Obj == arg.Obj {
				return as.Values[i]
			}
		}
	}

	return nil
}

// vendorlessPath returns the devendorized version of the import path ipath.
// For example, VendorlessPath("foo/bar/vendor/a/b") returns "a/b".
//
// This function copid from https://github.com/golang/tools/blob/v0.1.10/internal/imports/fix.go#L1423-L1432
func vendorlessPath(ipath string) string {
	// Devendorize for use in import statement.
	if i := strings.LastIndex(ipath, "/vendor/"); i >= 0 {
		return ipath[i+len("/vendor/"):]
	}

	if strings.HasPrefix(ipath, "vendor/") {
		return ipath[len("vendor/"):]
	}

	return ipath
}

// fixArgs fixes pkg/errors args to fmt.Errorf format.
func fixArgs(exprs []ast.Expr) []ast.Expr {
	msg := exprs[0]

	lit, ok := msg.(*ast.BasicLit)
	if !ok {
		return exprs
	}

	s := lit.Value
	s = verb(unquote(s))

	msg.(*ast.BasicLit).Value = strconv.Quote(s) // re-quoted
	exprs[0] = msg

	return exprs
}

// reorderArgs re-orders pkg/errors args to fmt.Errorf format.
func reorderArgs(exprs []ast.Expr) []ast.Expr {
	errStmt := exprs[0]
	msg := exprs[1]
	args := exprs[2:]

	// adds %w verb to the end of msg
	s := msg.(*ast.BasicLit).Value
	s = verb(unquote(s) + ": %w")
	msg.(*ast.BasicLit).Value = strconv.Quote(s) // re-quoted

	return append(append([]ast.Expr{msg}, args...), errStmt)
}

// verb assumes unquoted msg.
func verb(msg string) string {
	if strings.Contains(msg, `%w`) {
		println("ignore:", msg)
		return msg
	}

	if strings.Contains(msg, `%v`) {
		println("replace:", msg)
		return strings.ReplaceAll(msg, `%v`, `%w`)
	}

	return msg
}

// unquote assumes quoted s.
func unquote(s string) string {
	if s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1] // skip first and last char
	}

	return s
}
