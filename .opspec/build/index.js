const jsonSchemaRefParser = require('json-schema-ref-parser')
const fs = require('fs')

async function bundleJsonSchemas() {
  const bundledJsonSchema = await jsonSchemaRefParser.bundle('/src/op-definition-format/jsonschema/root.json')

  fs.writeFileSync(
    '/src/op-definition-format/jsonschema.json',
    JSON.stringify(bundledJsonSchema, null, 2)
  )
}

bundleJsonSchemas().catch(err => {
  console.log(err)
  process.exit(1)
})
