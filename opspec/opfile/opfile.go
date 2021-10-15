// @TODO remove this package if future version of golang embed supports either symlinks or parent dir references
package opfile

import _ "embed"

//go:embed jsonschema.json
var JsonSchemaBytes []byte
