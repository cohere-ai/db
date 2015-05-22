package sqlgen

import (
	"fmt"
	"strings"
)

// ColumnValue represents a bundle between a column and a corresponding value.
type ColumnValue struct {
	Column
	Operator string
	Value    cc
	hash     string
}

type columnValueT struct {
	Column   string
	Operator string
	Value    string
}

// Hash returns a unique identifier.
func (c *ColumnValue) Hash() string {
	if c.hash == "" {
		c.hash = fmt.Sprintf(`ColumnValue{Name:%q, Operator:%q, Value:%q}`, c.Column.Hash(), c.Operator, c.Value.Hash())
	}
	return c.hash
}

// Compile transforms the ColumnValue into its equivalent SQL representation.
func (c *ColumnValue) Compile(layout *Template) (compiled string) {

	if z, ok := layout.Read(c); ok {
		return z
	}

	data := columnValueT{
		c.Column.Compile(layout),
		c.Operator,
		c.Value.Compile(layout),
	}

	compiled = mustParse(layout.ColumnValue, data)

	layout.Write(c, compiled)

	return
}

// ColumnValues represents an array of ColumnValue
type ColumnValues struct {
	ColumnValues []ColumnValue
	hash         string
}

// NewColumnValues returns an array of ColumnValue
func NewColumnValues(values ...ColumnValue) *ColumnValues {
	return &ColumnValues{ColumnValues: values}
}

// Hash returns a unique identifier.
func (c *ColumnValues) Hash() string {
	if c.hash == "" {
		s := make([]string, len(c.ColumnValues))
		for i := range c.ColumnValues {
			s[i] = c.ColumnValues[i].Hash()
		}
		c.hash = fmt.Sprintf("ColumnValues{ColumnValues:{%s}}", strings.Join(s, ", "))
	}
	return c.hash
}

// Compile transforms the ColumnValues into its SQL representation.
func (c *ColumnValues) Compile(layout *Template) (compiled string) {

	if z, ok := layout.Read(c); ok {
		return z
	}

	l := len(c.ColumnValues)

	out := make([]string, l)

	for i := range c.ColumnValues {
		out[i] = c.ColumnValues[i].Compile(layout)
	}

	compiled = strings.Join(out, layout.IdentifierSeparator)

	layout.Write(c, compiled)

	return
}
