// Copyright 2022 The go-analyzer Authors
// SPDX-License-Identifier: BSD-3-Clause

package packagefact

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:       "packagefact",
	Doc:        "gather name/value pairs from constant declarations",
	Run:        run,
	FactTypes:  []analysis.Fact{new(pairsFact)},
	ResultType: reflect.TypeOf(map[string]string{}),
}

func run(pass *analysis.Pass) (interface{}, error) {
	result := make(map[string]string)

	// At each import, print the fact from the imported
	// package and accumulate its information into the result.
	// (Warning: accumulation leads to quadratic growth of work.)
	doImport := func(spec *ast.ImportSpec) {
		pkg := imported(pass.TypesInfo, spec)
		var fact pairsFact
		if pass.ImportPackageFact(pkg, &fact) {
			for _, pair := range fact {
				eq := strings.IndexByte(pair, '=')
				result[pair[:eq]] = pair[1+eq:]
			}
			pass.ReportRangef(spec, "%s", strings.Join(fact, " "))
		}
	}

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range decl.Specs {
					switch decl.Tok {
					case token.IMPORT:
						doImport(spec.(*ast.ImportSpec))
					}
				}
			}
		}
	}

	// Sort/deduplicate the result and save it as a package fact.
	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var fact pairsFact
	for _, key := range keys {
		fact = append(fact, fmt.Sprintf("%s=%s", key, result[key]))
	}
	if len(fact) > 0 {
		pass.ExportPackageFact(&fact)
	}

	return result, nil
}

func imported(info *types.Info, spec *ast.ImportSpec) *types.Package {
	obj, ok := info.Implicits[spec]
	if !ok {
		obj = info.Defs[spec.Name] // renaming import
	}

	return obj.(*types.PkgName).Imported()
}

// pairsFact a pairsFact is a package-level fact that records
// an set of key=value strings accumulated from constant
// declarations in this package and its dependencies.
// Elements are ordered by keys, which are unique.
type pairsFact []string

func (f *pairsFact) AFact()         {}
func (f *pairsFact) String() string { return "pairs(" + strings.Join(*f, ", ") + ")" }
