package cmd

import (
	"encoding/json"
	"fmt"
)

// String is either single string or array of strings
type String []string

// UnmarshalJSON handler either Seq<Identity> or Identity parseing
func (seq *String) UnmarshalJSON(b []byte) error {
	var (
		one  string
		many []string
	)

	isSeq, isObj := isJSONArray(b)
	switch {
	case isSeq:
		if err := json.Unmarshal(b, &many); err != nil {
			return err
		}
		*seq = many
		return nil
	case isObj:
		var obj struct {
			Value string `json:"@value"`
		}
		if err := json.Unmarshal(b, &obj); err != nil {
			return err
		}
		*seq = append(many, obj.Value)
		return nil
	default:
		if err := json.Unmarshal(b, &one); err != nil {
			fmt.Println("====")
			fmt.Println(err)
			fmt.Println(string(b))
			return err
		}

		*seq = append(many, one)
		return nil
	}
}

/*

Identity of other graph node

  "rdfs:subClassOf": {
    "@id": "schema:Text"
  }
*/
type Identity struct {
	ID string `json:"@id,omitempty"`
}

// Identities is sequence
type Identities []Identity

// UnmarshalJSON handler either Seq<Identity> or Identity parseing
func (seq *Identities) UnmarshalJSON(b []byte) error {
	var (
		one  Identity
		many []Identity
	)

	isSeq, _ := isJSONArray(b)

	if isSeq {
		if err := json.Unmarshal(b, &many); err != nil {
			return err
		}
		*seq = many
		return nil
	}

	if err := json.Unmarshal(b, &one); err != nil {
		return err
	}

	*seq = append(many, one)
	return nil
}

// check if input is JSON array
func isJSONArray(b []byte) (isSeq, isObj bool) {
	for _, c := range b {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}
		isSeq = c == '['
		isObj = c == '{'
		break
	}

	return
}

/*

Schema specification of entity from schema.org

	{
		"@id": "schema:URL",
		"@type": "rdfs:Class",
		"rdfs:comment": "Data type: URL.",
		"rdfs:label": "URL",
		"rdfs:subClassOf": {
			"@id": "schema:Text"
		}
	}

*/
type Schema struct {
	ID            string     `json:"@id,omitempty"`
	Type          String     `json:"@type,omitempty"`
	Comment       String     `json:"rdfs:comment,omitempty"`
	Label         String     `json:"rdfs:label,omitempty"`
	SubPropertyOf Identities `json:"rdfs:subPropertyOf,omitempty"`
	SubClassOf    Identities `json:"rdfs:subClassOf,omitempty"`
	Domain        Identities `json:"schema:domainIncludes,omitempty"`
	Range         Identities `json:"schema:rangeIncludes,omitempty"`
}

// SchemaOrg is JSON-LD representation of schema.org
type SchemaOrg struct {
	Graph []Schema `json:"@graph,omitempty"`
}

//
func isA(sc *Schema, t string) bool {
	for _, tp := range sc.Type {
		if t == tp {
			return true
		}
	}

	return false
}

//
func isSubClassOf(sc *Schema, t string) bool {
	for _, tp := range sc.SubClassOf {
		if t == tp.ID {
			return true
		}
	}

	return false
}

func isTypeOf(schema map[string]*Schema, t string, rangeType string) *Schema {
	sc := schema[t]
	if sc == nil {
		return nil
	}

	if isA(sc, rangeType) {
		return sc
	}

	for _, subClass := range sc.SubClassOf {
		if v := isTypeOf(schema, subClass.ID, rangeType); v != nil {
			return v
		}
	}

	return nil
}

func isRangeOf(schema map[string]*Schema, sc *Schema, rangeType string) *Schema {
	for _, typeOf := range sc.Range {
		if v := isTypeOf(schema, typeOf.ID, rangeType); v != nil {
			return v
		}
	}
	return nil
}
