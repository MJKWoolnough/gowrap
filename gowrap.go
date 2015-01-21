package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type TypeDef struct {
	t *ast.TypeSpec
	m map[string]*ast.FuncDecl
}

type Package struct {
	Funcs map[string]*ast.FuncDecl
	Types map[string]*TypeDef
}

func (p *Package) AddFunc(f *ast.FuncDecl) {
	p.Funcs[f.Name.Name] = f
}

func (p *Package) AddType(t *ast.TypeSpec) {
	typeName := t.Name.Name
	if _, ok := p.Types[typeName]; !ok {
		p.Types[typeName] = &TypeDef{m: make(map[string]*ast.FuncDecl)}
	}
	p.Types[typeName].t = t
}

func (p *Package) AddMethod(f *ast.FuncDecl) {
	rType := f.Recv.List[0].Type
	if r, ok := rType.(*ast.StarExpr); ok {
		rType = r.X
	}
	typeName := rType.(*ast.Ident).Name
	if _, ok := p.Types[typeName]; !ok {
		p.Types[typeName] = &TypeDef{m: make(map[string]*ast.FuncDecl)}
	}
	p.Types[typeName].m[f.Name.Name] = f
}

func main() {
	var packagePath string
	flag.StringVar(&packagePath, "pkg", "", "path to package")
	flag.Parse()
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	packages := make(map[string]Package)
	for name, pkg := range pkgs {
		p := Package{
			Funcs: make(map[string]*ast.FuncDecl),
			Types: make(map[string]*TypeDef),
		}
		for _, file := range pkg.Files {
			for _, d := range file.Decls {
				switch d := d.(type) {
				case *ast.GenDecl:
					switch d.Tok {
					case token.IMPORT:
						//ImportSpec
					case token.CONST:
						//ValueSpec
					case token.TYPE:
						for _, spec := range d.Specs {
							p.AddType(spec.(*ast.TypeSpec))
						}
					case token.VAR:
						//ValueSpec
					}
				case *ast.FuncDecl:
					if d.Recv == nil {
						p.AddFunc(d)
					} else {
						p.AddMethod(d)
					}
				}
			}
		}
		packages[name] = p
	}
	for name, pkg := range packages {
		fmt.Println("Package:", name)
		fmt.Println("\nFuncs: -\n")
		for fName := range pkg.Funcs {
			fmt.Println(fName)
		}
		fmt.Println("\nTypes: -\n")
		for tName, typ := range pkg.Types {
			fmt.Println(tName)
			for mName := range typ.m {
				fmt.Println("-", mName)
			}
			fmt.Println()
		}
	}
}
