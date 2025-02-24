package jsonschema

import "fmt"

// Function to extract error messages from the schema
func extractErrorMessages(
	schema map[string]interface{},
	path string,
	messages map[string]string,
) {
	if errMsg, ok := schema["errorMessage"]; ok {
		if msg, valid := errMsg.(string); valid {
			messages[path] = msg
		}
	}

	for key, value := range schema {
		switch v := value.(type) {
		case map[string]interface{}:
			newPath := path
			if key == "properties" || key == "items" {
				newPath = path + "/" + key
			} else if path != "" {
				newPath = path + "/" + key
			} else {
				newPath = key
			}
			extractErrorMessages(v, newPath, messages)
		case []interface{}:
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					newPath := fmt.Sprintf("%s/%s[%d]", path, key, i)
					extractErrorMessages(itemMap, newPath, messages)
				}
			}
		}
	}
}
