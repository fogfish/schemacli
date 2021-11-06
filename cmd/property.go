package cmd

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/spf13/cobra"
)

/*

Command processes rdf:Property and build type system

	{
		"@id": "schema:url",
		"@type": "rdf:Property",
		"rdfs:comment": "URL of the item.",
		"rdfs:label": "url",
		"schema:domainIncludes": {
			"@id": "schema:Thing"
		},
		"schema:rangeIncludes": {
			"@id": "schema:URL"
		}
	}

*/

func init() {
	rootCmd.AddCommand(coreTypeCmd)
}

var coreTypeCmd = &cobra.Command{
	Use:   "property",
	Short: "property parses ontology and builds core types",
	Example: `
schemaorg property -f schemaorg.json
`,
	SilenceUsage: true,
	RunE:         coretype,
}

func coretype(cmd *cobra.Command, args []string) error {
	schema, err := parseSchemaOrg()
	if err != nil {
		return err
	}
	keyval := indexSchemaOrg(schema)

	coreTypeSpec := &ast.File{
		Name:  &ast.Ident{Name: "schemaorg"},
		Decls: []ast.Decl{},
	}

	for _, spec := range schema.Graph {
		if isA(&spec, "rdf:Property") {
			if v := isRangeOf(keyval, &spec, "schema:DataType"); v != nil {
				//
				// Note: only schema:Text is supported
				if v.ID == "schema:Text" {
					coreTypeSpec.Decls = append(coreTypeSpec.Decls,
						declareTypeForProperty(&spec, "string"),
					)
				}
			}
		}
	}

	return stdout(coreTypeSpec)
}

//
//
func declareTypeForProperty(sc *Schema, goCoreType string) ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: fmt.Sprintf("\n/*\n\n%s is https://schema.org/%s\n\n%s\n*/",
						strings.Title(sc.Label[0]), sc.Label[0], sc.Comment[0]),
				},
			},
		},
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: strings.Title(sc.Label[0])},
				Type: &ast.Ident{Name: goCoreType},
			},
		},
	}
}
