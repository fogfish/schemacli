package cmd

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Execute is entry point for cobra cli application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}

var (
	fSchemaOrg string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&fSchemaOrg, "file", "f", "", "path to json-ld schema.org ontology")
}

var rootCmd = &cobra.Command{
	Use:     "schemaorg",
	Short:   "generates Go types for schema.org",
	Long:    `generates Go types for schema.org`,
	Run:     root,
	Version: "v0",
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}

//
// Helper functions
//

//
func parseSchemaOrg() (*SchemaOrg, error) {
	fd, err := os.Open(fSchemaOrg)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	bytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	var org SchemaOrg
	if err = json.Unmarshal(bytes, &org); err != nil {
		return nil, err
	}

	return &org, err
}

//
//
func indexSchemaOrg(schema *SchemaOrg) map[string]*Schema {
	idx := map[string]*Schema{}
	for _, sc := range schema.Graph {
		n := sc
		idx[sc.ID] = &n
	}
	return idx
}

//
//
func stdout(spec *ast.File) error {
	return printer.Fprint(
		os.Stdout,
		token.NewFileSet(),
		spec,
	)
}
