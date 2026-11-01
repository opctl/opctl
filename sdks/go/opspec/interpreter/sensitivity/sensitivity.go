// Package sensitivity carries private runtime metadata without serializing values.
package sensitivity

import (
	"github.com/opctl/opctl/sdks/go/model"
	"reflect"
	"strings"
)

func Secret(param *model.ParamSpec) bool {
	return param != nil && ((param.String != nil && param.String.IsSecret) || (param.Number != nil && param.Number.IsSecret) || (param.Array != nil && param.Array.IsSecret) || (param.Object != nil && param.Object.IsSecret) || (param.Dir != nil && param.Dir.IsSecret) || (param.File != nil && param.File.IsSecret) || (param.Socket != nil && param.Socket.IsSecret))
}
func Mark(value *model.Value, sensitive bool) *model.Value {
	if value == nil || !sensitive || value.Sensitive {
		return value
	}
	copy := *value
	copy.Sensitive = true
	return &copy
}
func References(expr interface{}, scope map[string]*model.Value) bool {
	var visit func(reflect.Value) bool
	visit = func(value reflect.Value) bool {
		if !value.IsValid() {
			return false
		}
		switch value.Kind() {
		case reflect.Interface, reflect.Pointer:
			if !value.IsNil() {
				return visit(value.Elem())
			}
		case reflect.String:
			text := value.String()
			for name, value := range scope {
				if value == nil || !value.Sensitive {
					continue
				}
				for _, suffix := range []string{")", "/", ".", "["} {
					if strings.Contains(text, "$("+name+suffix) {
						return true
					}
				}
			}
		case reflect.Struct:
			for i := 0; i < value.NumField(); i++ {
				if visit(value.Field(i)) {
					return true
				}
			}
		case reflect.Map:
			for _, key := range value.MapKeys() {
				if visit(key) || visit(value.MapIndex(key)) {
					return true
				}
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < value.Len(); i++ {
				if visit(value.Index(i)) {
					return true
				}
			}
		}
		return false
	}
	return visit(reflect.ValueOf(expr))
}
func Inputs(args map[string]interface{}, params map[string]*model.ParamSpec, scope map[string]*model.Value) map[string]bool {
	protected := map[string]bool{}
	for name, param := range params {
		expr, ok := args[name]
		if ok && (expr == nil || expr == "") {
			expr = "$(" + name + ")"
		}
		protected[name] = Secret(param) || References(expr, scope)
	}
	return protected
}
