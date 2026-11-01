package sdk

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type paramSpec map[string]any
type inputField struct {
	name, kind string
	hasDefault bool
	secret     bool
	handle     inputHandle
	index      int
}

func inputSpecs[I any]() (map[string]paramSpec, []inputField, error) {
	return parameterSpecs[I]("input")
}

func parameterSpecs[P any](direction string) (map[string]paramSpec, []inputField, error) {
	typ := reflect.TypeOf((*P)(nil)).Elem()
	if typ.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("Go op %ss must be a struct", direction)
	}
	specs := map[string]paramSpec{}
	fields := []inputField{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			return nil, nil, fmt.Errorf("%s field %s must be exported", direction, field.Name)
		}
		handle, ok := reflect.Zero(field.Type).Interface().(inputHandle)
		if !ok || field.Type.Kind() == reflect.Pointer || reflect.TypeOf(handle.bind(inputValue{})) != field.Type {
			return nil, nil, fmt.Errorf("%s field %s must be an opctl handle value", direction, field.Name)
		}
		tag, err := parseTag(field.Tag.Get("opctl"))
		if err != nil {
			return nil, nil, fmt.Errorf("%s field %s: %w", direction, field.Name, err)
		}
		name := tag["name"]
		if strings.HasPrefix(name, "opctl_") {
			return nil, nil, fmt.Errorf("%s %q: prefix opctl_ is reserved for generated handles", direction, name)
		}
		delete(tag, "name")
		if _, ok := specs[name]; ok {
			return nil, nil, fmt.Errorf("duplicate %s %q", direction, name)
		}
		kind, valueType := handle.inputType()
		spec, err := param(kind, valueType, tag)
		if err != nil {
			return nil, nil, fmt.Errorf("%s %q: %w", direction, name, err)
		}
		specs[name] = spec
		_, hasDefault := tag["default"]
		fields = append(fields, inputField{hasDefault: hasDefault, name: name, kind: kind, secret: tag["secret"] == "true", handle: handle, index: i})
	}
	return specs, fields, nil
}

func bindInputs[I any](fields []inputField, values map[string]any) (I, error) {
	var in I
	target := reflect.ValueOf(&in).Elem()
	for name := range values {
		found := false
		for _, field := range fields {
			if field.name == name {
				found = true
				break
			}
		}
		if !found {
			return in, fmt.Errorf("unknown input %q", name)
		}
	}
	for _, field := range fields {
		value, exists := values[field.name]
		refOnly := field.secret || field.kind == "dir" || field.kind == "file" || field.kind == "socket"
		if refOnly && exists {
			return in, fmt.Errorf("input %q must be omitted from the plan request", field.name)
		}
		if !refOnly && !exists {
			return in, fmt.Errorf("missing value for input %q", field.name)
		}
		raw, err := json.Marshal(value)
		if err != nil {
			return in, fmt.Errorf("input %q cannot be encoded", field.name)
		}
		if !refOnly {
			_, typ := field.handle.inputType()
			if value == nil {
				return in, fmt.Errorf("input %q cannot be null", field.name)
			}
			if err := json.Unmarshal(raw, reflect.New(typ).Interface()); err != nil {
				return in, fmt.Errorf("input %q has an invalid value for its Go type", field.name)
			}
		}
		target.Field(field.index).Set(reflect.ValueOf(field.handle.bind(inputValue{name: field.name, value: raw, secret: field.secret})))
	}
	return in, nil
}

func parseTag(tag string) (map[string]string, error) {
	parts := []string{}
	var part strings.Builder
	quoted, escaped := false, false
	closed := false
	for _, r := range tag {
		if closed && r != ',' {
			return nil, fmt.Errorf("text after quoted tag value")
		}
		if escaped {
			if r != '\\' && r != '\'' {
				return nil, fmt.Errorf("invalid tag escape")
			}
			part.WriteRune(r)
			escaped = false
			continue
		}
		if quoted && r == '\\' {
			escaped = true
			continue
		}
		if r == '\'' {
			if !quoted && !strings.HasSuffix(part.String(), "=") {
				return nil, fmt.Errorf("quote must start a tag value")
			}
			closed = quoted
			quoted = !quoted
			part.WriteRune(r)
			continue
		}
		if r == ',' && !quoted {
			parts = append(parts, part.String())
			part.Reset()
			closed = false
		} else {
			part.WriteRune(r)
		}
	}
	if quoted || escaped {
		return nil, fmt.Errorf("unterminated quoted tag value")
	}
	parts = append(parts, part.String())
	if !regexp.MustCompile(`^[-_a-zA-Z0-9]+$`).MatchString(parts[0]) {
		return nil, fmt.Errorf("input needs a name containing only letters, numbers, - or _")
	}
	out := map[string]string{"name": parts[0]}
	for _, part := range parts[1:] {
		key, value, ok := strings.Cut(part, "=")
		if part == "secret" {
			key, value, ok = "secret", "true", true
		}
		if !ok || key == "name" || (key == "secret" && part != "secret") {
			return nil, fmt.Errorf("invalid opctl tag item %q", part)
		}
		if _, ok := out[key]; ok {
			return nil, fmt.Errorf("duplicate opctl tag key %s", key)
		}
		if strings.HasPrefix(value, "'") {
			if len(value) < 2 || !strings.HasSuffix(value, "'") {
				return nil, fmt.Errorf("invalid quoted tag value")
			}
			value = value[1 : len(value)-1]
		} else if strings.Contains(value, "'") {
			return nil, fmt.Errorf("invalid quoted tag value")
		}
		out[key] = value
	}
	return out, nil
}

func param(kind string, typ reflect.Type, tag map[string]string) (paramSpec, error) {
	body := map[string]any{}
	spec := paramSpec{kind: body}
	constraints := map[string]any{}
	if kind == "array" || kind == "object" {
		schema, err := typeSchema(typ, map[reflect.Type]bool{})
		if err != nil {
			return nil, err
		}
		if schema["type"] != kind {
			return nil, fmt.Errorf("%s handle needs a Go %s type", kind, kind)
		}
		delete(schema, "type")
		for _, key := range []string{"items", "additionalProperties"} {
			if child, ok := schema[key].(map[string]any); ok && len(child) == 0 {
				delete(schema, key)
			}
		}
		constraints = schema
	}
	flat := map[string]any{}
	for key, value := range tag {
		valid := true
		var parsed any = value
		var err error
		switch key {
		case "description":
			spec[key] = value
			continue
		case "secret":
			valid = kind != "boolean"
			body["isSecret"] = true
		case "default":
			valid = kind != "socket"
			if kind == "number" {
				parsed, err = strconv.ParseFloat(value, 64)
			}
			if kind == "boolean" {
				parsed, err = strconv.ParseBool(value)
			}
			if kind == "array" || kind == "object" {
				err = json.Unmarshal([]byte(value), &parsed)
				if err == nil {
					encoded, _ := json.Marshal(parsed)
					err = json.Unmarshal(encoded, reflect.New(typ).Interface())
					if parsed == nil {
						err = fmt.Errorf("null default")
					}
				}
			}
			body[key] = parsed
		case "enum":
			valid = kind == "string" || kind == "number"
			enum := []any{}
			for _, entry := range strings.Split(value, "|") {
				var v any = entry
				if kind == "number" {
					v, err = strconv.ParseFloat(entry, 64)
					if err != nil {
						break
					}
				}
				enum = append(enum, v)
			}
			flat[key] = enum
		case "format":
			valid = kind == "string" || kind == "number"
			if kind == "number" && value != "integer" {
				valid = false
			}
			if kind == "string" && !strings.Contains("|date-time|docker-image-ref|email|hostname|ipv4|ipv6|semver|uri|", "|"+value+"|") {
				valid = false
			}
			flat[key] = value
		case "pattern":
			valid = kind == "string"
			flat[key] = value
		case "minLength", "maxLength", "minItems", "maxItems", "minProperties", "maxProperties":
			valid = (strings.HasSuffix(key, "Length") && kind == "string") || (strings.HasSuffix(key, "Items") && kind == "array") || (strings.HasSuffix(key, "Properties") && kind == "object")
			var n int
			n, err = strconv.Atoi(value)
			if n < 0 {
				err = fmt.Errorf("negative limit")
			}
			flat[key] = n
		case "minimum", "maximum", "multipleOf":
			valid = kind == "number"
			flat[key], err = strconv.ParseFloat(value, 64)
		case "uniqueItems":
			valid = kind == "array"
			flat[key], err = strconv.ParseBool(value)
		case "required":
			valid = kind == "object"
			flat[key] = strings.Split(value, "|")
		case "constraints":
			valid = kind == "string" || kind == "number" || kind == "array" || kind == "object"
			var nested map[string]any
			err = json.Unmarshal([]byte(value), &nested)
			if nested == nil {
				err = fmt.Errorf("constraints must be a JSON object")
			}
			for k, v := range nested {
				if _, ok := constraints[k]; ok {
					return nil, fmt.Errorf("duplicate constraint %s", k)
				}
				constraints[k] = v
			}
		default:
			return nil, fmt.Errorf("unknown opctl tag key %s", key)
		}
		if !valid {
			return nil, fmt.Errorf("opctl tag key %s is invalid for %s", key, kind)
		}
		if err != nil {
			return nil, fmt.Errorf("invalid %s tag", key)
		}
	}
	for key, value := range flat {
		if _, ok := constraints[key]; ok {
			return nil, fmt.Errorf("duplicate constraint %s", key)
		}
		constraints[key] = value
	}
	if len(constraints) > 0 {
		body["constraints"] = constraints
	}
	return spec, nil
}

func typeSchema(typ reflect.Type, visiting map[reflect.Type]bool) (map[string]any, error) {
	if visiting[typ] {
		return nil, fmt.Errorf("recursive input types are not supported")
	}
	visiting[typ] = true
	defer delete(visiting, typ)
	schema := map[string]any{}
	switch typ.Kind() {
	case reflect.Interface:
		if typ.NumMethod() != 0 {
			return nil, fmt.Errorf("only any is supported as an input interface")
		}
	case reflect.String:
		schema["type"] = "string"
	case reflect.Bool:
		schema["type"] = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		schema["type"] = "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema["type"] = "integer"
		schema["minimum"] = 0
	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"
	case reflect.Slice:
		schema["type"] = "array"
		item, err := typeSchema(typ.Elem(), visiting)
		if err != nil {
			return nil, err
		}
		schema["items"] = item
	case reflect.Map:
		if typ.Key().Kind() != reflect.String {
			return nil, fmt.Errorf("object maps need string keys")
		}
		schema["type"] = "object"
		item, err := typeSchema(typ.Elem(), visiting)
		if err != nil {
			return nil, err
		}
		schema["additionalProperties"] = item
	case reflect.Struct:
		schema["type"] = "object"
		properties := map[string]any{}
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if !field.IsExported() {
				continue
			}
			tag := strings.Split(field.Tag.Get("json"), ",")
			name := tag[0]
			if name == "-" {
				continue
			}
			if field.Anonymous {
				return nil, fmt.Errorf("embedded object fields are not supported")
			}
			if name == "" {
				name = field.Name
			}
			if _, ok := properties[name]; ok {
				return nil, fmt.Errorf("duplicate object property %q", name)
			}
			property, err := typeSchema(field.Type, visiting)
			if err != nil {
				return nil, err
			}
			properties[name] = property
		}
		schema["properties"] = properties
	default:
		return nil, fmt.Errorf("unsupported Go input type %s", typ)
	}
	return schema, nil
}
