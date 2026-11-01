package sdk

import (
	"encoding/json"
	"fmt"
)

func literalValue(v any) inputValue {
	data, err := json.Marshal(v)
	if err != nil {
		return inputValue{err: fmt.Errorf("literal handle: %w", err)}
	}
	return inputValue{literal: true, value: data}
}

// StringOf returns a String holding a literal, readable during plan.
func StringOf(v string) String { return String{literalValue(v)} }

// NumberOf returns a Number holding a finite literal, readable during plan.
func NumberOf(v float64) Number { return Number{literalValue(v)} }

// BooleanOf returns a Boolean holding a literal, readable during plan.
func BooleanOf(v bool) Boolean { return Boolean{literalValue(v)} }

// ArrayOf snapshots a slice as JSON in a literal handle, readable during plan.
func ArrayOf[T any](v []T) Array[T] { return Array[T]{literalValue(v)} }

// ObjectOf snapshots v as JSON, respecting struct json tags. Value reads it during plan.
func ObjectOf[T any](v T) Object[T] { return Object[T]{literalValue(v)} }
