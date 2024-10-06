/**
 * Enum type for mow.cli
 *
 * Gist: https://gist.github.com/jawher/8b3720e5a166e5219397ce2b8c8f7ddb
 * Credit to Jawher Moussa (@jawher)
 */

package mowclitypes

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type EnumValue struct {
	Value string
	Desc  string
}

func NewEnum(defaultValue string, values ...EnumValue) *Enum {
	return &Enum{
		Values: values,
		Value:  defaultValue,
	}
}

type Enum struct {
	Values []EnumValue
	Value  string
}

func (e *Enum) Set(v string) error {
	for _, s := range e.Values {
		if v == s.Value {
			e.Value = v
			return nil
		}
	}

	possible := make([]string, len(e.Values))
	for i, v := range e.Values {
		possible[i] = fmt.Sprintf("%q", v.Value)
	}

	return errors.Errorf("invalid value %q. Must be one of: %s", v, strings.Join(possible, ", "))
}

func (e *Enum) Desc(desc string) string {
	var res strings.Builder
	res.WriteString(desc)
	res.WriteString("\n")

	for i, v := range e.Values {
		if i > 0 {
			res.WriteString("\n")
		}
		res.WriteString(" - ")
		res.WriteString(v.Value)
		res.WriteString(": ")
		res.WriteString(v.Desc)
	}

	return res.String()
}

func (e *Enum) String() string {
	return e.Value
}

func (e *Enum) IsDefault() bool {
	return e.Value == ""
}
