// @TODO remove this package if future version of golang embed supports either symlinks or parent dir references
package webapp

import "embed"

//go:embed build
var Build embed.FS
