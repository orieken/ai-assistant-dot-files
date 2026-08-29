package tools

import "encoding/json"

func projectPathProperty() map[string]any {
	return map[string]any{
		"type":        "string",
		"description": "Absolute path to the project root",
	}
}

// objectSchema builds a raw JSON Schema for an object with the given required
// keys and properties.
func objectSchema(required []string, properties map[string]any) json.RawMessage {
	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return mustMarshalSchema(schema)
}

// mustMarshalSchema panics on failure: schemas are static literals, so a
// marshal error is a programmer error surfaced at registration, not a runtime
// condition to handle.
func mustMarshalSchema(schema map[string]any) json.RawMessage {
	body, err := json.Marshal(schema)
	if err != nil {
		panic("tools: static schema failed to marshal: " + err.Error())
	}
	return body
}

func projectPathOnlySchema() json.RawMessage {
	return objectSchema([]string{"projectPath"}, map[string]any{
		"projectPath": projectPathProperty(),
	})
}
