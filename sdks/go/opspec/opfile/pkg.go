package opfile

const (
	FileName = "op.yml"
	// SchemaVersion versions the OpSpec JSON schema, independently of opspec semver.
	// Bump it whenever the schema changes.
	SchemaVersion = 4
	// MinSchemaVersion is the oldest schema this node continues to read.
	MinSchemaVersion = 1
)
