package bracketed

import (
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
)

// CoerceToArrayOrObject data to an object or an array value
func CoerceToArrayOrObject(
	data *ipld.Node,
) (*ipld.Node, error) {

	if dataAsArray, err := coerce.ToArray(data); err == nil {
		// array coercible type
		return dataAsArray, nil
	}

	return coerce.ToObject(data)
}
