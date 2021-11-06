package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

/*

Command processes schema:DataType

	{
		"@id": "schema:Text",
		"@type": [
			"rdfs:Class",
			"schema:DataType"
		],
		"rdfs:comment": "Data type: Text.",
		"rdfs:label": "Text"
	}
*/

var (
	isaType string
)

func init() {
	rootCmd.AddCommand(isaCmd)
	isaCmd.Flags().StringVarP(&isaType, "type", "t", "schema:DataType", "type to query")
}

var isaCmd = &cobra.Command{
	Use:   "isa",
	Short: "isa query searches for nodes that belongs to the type",
	Example: `
schemaorg isa -f schemaorg.json -t schema:DataType
`,
	SilenceUsage: true,
	RunE:         isa,
}

func isa(cmd *cobra.Command, args []string) error {
	schema, err := parseSchemaOrg()
	if err != nil {
		return err
	}

	seq := []Schema{}

	for _, spec := range schema.Graph {
		if isA(&spec, isaType) {
			seq = append(seq, spec)
		}
	}

	b, err := json.MarshalIndent(seq, "", "  ")
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(b)
	return err
}
