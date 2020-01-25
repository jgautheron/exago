package cov

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Function describes a function inner characteristics
type Function struct {
	// Name is the name of the function. If the function has a receiver, the
	// name will be of the form T.N, where T is the type and N is the name.
	Name string `json:"name"`
	// File is the full path to the file in which the function is defined.
	File string `json:"file"`
	// Start is the start offset of the function's signature.
	Start int `json:"start"`
	// End is the end offset of the function.
	End int `json:"end"`
	// Coverage
	Coverage float64 `json:"coverage"`
	// TLOC
	TLOC int64 `json:"tloc"`
	// Statements registered with this function, JSON output omit statements (for now)
	Statements []*Statement `json:"-"`
}

// FuncExtent describes a function's extent in the source by file and position.
type FuncExtent struct {
	extent
	name  string
	stmts []*StmtExtent
}

// Accumulate will accumulate the coverage information from the provided
// Function into this Function.
func (f *Function) Accumulate(f2 *Function) error {
	if f.Name != f2.Name {
		return fmt.Errorf("Names do not match: %q != %q", f.Name, f2.Name)
	}
	if f.File != f2.File {
		return fmt.Errorf("Files do not match: %q != %q", f.File, f2.File)
	}
	if f.Start != f2.Start || f.End != f2.End {
		return fmt.Errorf("Source ranges do not match: %d-%d != %d-%d", f.Start, f.End, f2.Start, f2.End)
	}
	if len(f.Statements) != len(f2.Statements) {
		return fmt.Errorf("Number of statements do not match: %d != %d", len(f.Statements), len(f2.Statements))
	}
	for i, s := range f.Statements {
		err := s.Accumulate(f2.Statements[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// FuncVisitor implements the visitor that builds the function position list for a file.
type FuncVisitor struct {
	fset  *token.FileSet
	funcs []*FuncExtent
}

// Visit implements the ast.Visitor interface.
func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	var body *ast.BlockStmt
	var name string
	switch n := node.(type) {
	case *ast.FuncLit:
		body = n.Body
	case *ast.FuncDecl:
		body = n.Body
		name = n.Name.Name
		// Function name is prepended with "T." if there is a receiver, where
		// T is the type of the receiver, dereferenced if it is a pointer.
		if n.Recv != nil {
			field := n.Recv.List[0]
			switch recv := field.Type.(type) {
			case *ast.StarExpr:
				name = recv.X.(*ast.Ident).Name + "." + name
			case *ast.Ident:
				name = recv.Name + "." + name
			}
		}
	}
	if body != nil {
		start := v.fset.Position(node.Pos())
		end := v.fset.Position(node.End())
		if name == "" {
			name = fmt.Sprintf("@%d:%d", start.Line, start.Column)
		}
		fe := &FuncExtent{
			name: name,
			extent: extent{
				startOffset: start.Offset,
				startLine:   start.Line,
				startCol:    start.Column,
				endOffset:   end.Offset,
				endLine:     end.Line,
				endCol:      end.Column,
			},
		}

		sv := StmtVisitor{fset: v.fset, function: fe}
		sv.VisitStmt(body)

		v.funcs = append(v.funcs, fe)
	}

	return v
}
