// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generation based on GopherCon UK 2019 talk by Paul Jolly:
// Write Less (Code), Generate More
// https://www.youtube.com/watch?v=xcpboZZy-64

package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/gertd/go-pluralize"
	"golang.org/x/tools/go/packages"
)

const (
	inFile          = "enums.go"
	outFile         = "generated_" + inFile
	outFileTemplate = `` +
		`// Copyright (c) Juniper Networks, Inc., 2024-{{.Year}}.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// contents of this file are auto-generated by ./generator/generator.go - DO NOT EDIT

package enum

import oenum "github.com/orsinium-labs/enum"
{{ range $key, $value := .NameToTypeInfo }}
var _ enum = (*{{ $key }})(nil)

func (o {{ $key }}) String() string {
	return o.Value
}

func (o *{{ $key }}) FromString(s string) error {
	if {{ $value.Plural }}.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}
{{ end }}
var ({{ range  $key, $value := .NameToTypeInfo }}
	_ enum = new({{ $key }})
	{{ $value.Plural }} = oenum.New({{ range $v := $value.Values }}
		{{ $v }},{{ end }}
	)
{{ end }})
`
)

type TypeInfo struct {
	Plural string   // used with oenum.New() - Things
	Values []string // each enum value       - Thing1, Thing2, ... ThingN
}

var (
	Pluralize      *pluralize.Client
	TypeNameToInfo map[string]TypeInfo
)

func main() {
	Pluralize = pluralize.NewClient()
	TypeNameToInfo = make(map[string]TypeInfo)

	cfg := packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}

	pkgs, err := packages.Load(&cfg)
	if err != nil {
		panic(fmt.Errorf("while loading packages - %w", err))
	}

	if len(pkgs) != 1 {
		panic(fmt.Errorf("expected 1 package, got %d packages", len(pkgs)))
	}
	pkg := pkgs[0]

	for _, file := range pkg.Syntax {
		absPath := pkg.Fset.Position(file.Package).Filename
		if path.Base(absPath) != inFile {
			continue
		}

		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			if gd.Tok != token.VAR {
				continue
			}

			err = handleVar(gd)
			if err != nil {
				panic(err)
			}
		}
	}

	err = render()
	if err != nil {
		panic(fmt.Errorf("while rendering template - %w", err))
	}
}

func render() error {
	var tmplData struct {
		Year           string
		NameToTypeInfo map[string]TypeInfo
	}

	tmplData.Year = time.Now().Format("2006")
	tmplData.NameToTypeInfo = TypeNameToInfo

	f, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("while creating file for generated code - %w", err)
	}

	tmpl, err := template.New("").Parse(outFileTemplate)
	if err != nil {
		return fmt.Errorf("while parsing template - %w", err)
	}

	err = tmpl.Execute(f, tmplData)
	if err != nil {
		return fmt.Errorf("while executing template - %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("while closing file %q - %w", outFile, err)
	}

	return nil
}

func handleVar(gd *ast.GenDecl) error {
	for _, spec := range gd.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			return fmt.Errorf("spec should have been a ValueSpec, got %v", spec)
		}

		if len(spec.Names) != 1 {
			return fmt.Errorf("expected 1 spec name, got %d spec names", len(spec.Names))
		}

		if len(spec.Values) != 1 {
			return fmt.Errorf("expected 1 spec value, got %d spec values", len(spec.Values))
		}

		name := spec.Names[0]
		value, ok := spec.Values[0].(*ast.CompositeLit)
		if !ok {
			continue
		}

		valueType, ok := value.Type.(*ast.Ident)
		if !ok {
			continue
		}

		tName := valueType.Name

		// Fill the TypeNameToInfo map which is used by render()
		var typeInfo TypeInfo
		if typeInfo, ok = TypeNameToInfo[tName]; !ok {
			// entry not found - create one with the plural version of the type name included
			typeInfo.Plural = Pluralize.Plural(tName)
			if tName == typeInfo.Plural {
				return fmt.Errorf("cannot pluralize - plural of %q is %q", tName, typeInfo.Plural)
			}
		}
		typeInfo.Values = append(typeInfo.Values, name.Name)
		TypeNameToInfo[tName] = typeInfo
	}

	return nil
}