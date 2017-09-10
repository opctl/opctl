// Code generated by counterfeiter. DO NOT EDIT.
package data

import (
	"sync"

	"github.com/opspec-io/sdk-golang/model"
)

type fakeCoercer struct {
	CoerceToNumberStub        func(value *model.Value) (float64, error)
	coerceToNumberMutex       sync.RWMutex
	coerceToNumberArgsForCall []struct {
		value *model.Value
	}
	coerceToNumberReturns struct {
		result1 float64
		result2 error
	}
	coerceToNumberReturnsOnCall map[int]struct {
		result1 float64
		result2 error
	}
	CoerceToObjectStub        func(value *model.Value) (map[string]interface{}, error)
	coerceToObjectMutex       sync.RWMutex
	coerceToObjectArgsForCall []struct {
		value *model.Value
	}
	coerceToObjectReturns struct {
		result1 map[string]interface{}
		result2 error
	}
	coerceToObjectReturnsOnCall map[int]struct {
		result1 map[string]interface{}
		result2 error
	}
	CoerceToStringStub        func(value *model.Value) (string, error)
	coerceToStringMutex       sync.RWMutex
	coerceToStringArgsForCall []struct {
		value *model.Value
	}
	coerceToStringReturns struct {
		result1 string
		result2 error
	}
	coerceToStringReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *fakeCoercer) CoerceToNumber(value *model.Value) (float64, error) {
	fake.coerceToNumberMutex.Lock()
	ret, specificReturn := fake.coerceToNumberReturnsOnCall[len(fake.coerceToNumberArgsForCall)]
	fake.coerceToNumberArgsForCall = append(fake.coerceToNumberArgsForCall, struct {
		value *model.Value
	}{value})
	fake.recordInvocation("CoerceToNumber", []interface{}{value})
	fake.coerceToNumberMutex.Unlock()
	if fake.CoerceToNumberStub != nil {
		return fake.CoerceToNumberStub(value)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.coerceToNumberReturns.result1, fake.coerceToNumberReturns.result2
}

func (fake *fakeCoercer) CoerceToNumberCallCount() int {
	fake.coerceToNumberMutex.RLock()
	defer fake.coerceToNumberMutex.RUnlock()
	return len(fake.coerceToNumberArgsForCall)
}

func (fake *fakeCoercer) CoerceToNumberArgsForCall(i int) *model.Value {
	fake.coerceToNumberMutex.RLock()
	defer fake.coerceToNumberMutex.RUnlock()
	return fake.coerceToNumberArgsForCall[i].value
}

func (fake *fakeCoercer) CoerceToNumberReturns(result1 float64, result2 error) {
	fake.CoerceToNumberStub = nil
	fake.coerceToNumberReturns = struct {
		result1 float64
		result2 error
	}{result1, result2}
}

func (fake *fakeCoercer) CoerceToNumberReturnsOnCall(i int, result1 float64, result2 error) {
	fake.CoerceToNumberStub = nil
	if fake.coerceToNumberReturnsOnCall == nil {
		fake.coerceToNumberReturnsOnCall = make(map[int]struct {
			result1 float64
			result2 error
		})
	}
	fake.coerceToNumberReturnsOnCall[i] = struct {
		result1 float64
		result2 error
	}{result1, result2}
}

func (fake *fakeCoercer) CoerceToObject(value *model.Value) (map[string]interface{}, error) {
	fake.coerceToObjectMutex.Lock()
	ret, specificReturn := fake.coerceToObjectReturnsOnCall[len(fake.coerceToObjectArgsForCall)]
	fake.coerceToObjectArgsForCall = append(fake.coerceToObjectArgsForCall, struct {
		value *model.Value
	}{value})
	fake.recordInvocation("CoerceToObject", []interface{}{value})
	fake.coerceToObjectMutex.Unlock()
	if fake.CoerceToObjectStub != nil {
		return fake.CoerceToObjectStub(value)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.coerceToObjectReturns.result1, fake.coerceToObjectReturns.result2
}

func (fake *fakeCoercer) CoerceToObjectCallCount() int {
	fake.coerceToObjectMutex.RLock()
	defer fake.coerceToObjectMutex.RUnlock()
	return len(fake.coerceToObjectArgsForCall)
}

func (fake *fakeCoercer) CoerceToObjectArgsForCall(i int) *model.Value {
	fake.coerceToObjectMutex.RLock()
	defer fake.coerceToObjectMutex.RUnlock()
	return fake.coerceToObjectArgsForCall[i].value
}

func (fake *fakeCoercer) CoerceToObjectReturns(result1 map[string]interface{}, result2 error) {
	fake.CoerceToObjectStub = nil
	fake.coerceToObjectReturns = struct {
		result1 map[string]interface{}
		result2 error
	}{result1, result2}
}

func (fake *fakeCoercer) CoerceToObjectReturnsOnCall(i int, result1 map[string]interface{}, result2 error) {
	fake.CoerceToObjectStub = nil
	if fake.coerceToObjectReturnsOnCall == nil {
		fake.coerceToObjectReturnsOnCall = make(map[int]struct {
			result1 map[string]interface{}
			result2 error
		})
	}
	fake.coerceToObjectReturnsOnCall[i] = struct {
		result1 map[string]interface{}
		result2 error
	}{result1, result2}
}

func (fake *fakeCoercer) CoerceToString(value *model.Value) (string, error) {
	fake.coerceToStringMutex.Lock()
	ret, specificReturn := fake.coerceToStringReturnsOnCall[len(fake.coerceToStringArgsForCall)]
	fake.coerceToStringArgsForCall = append(fake.coerceToStringArgsForCall, struct {
		value *model.Value
	}{value})
	fake.recordInvocation("CoerceToString", []interface{}{value})
	fake.coerceToStringMutex.Unlock()
	if fake.CoerceToStringStub != nil {
		return fake.CoerceToStringStub(value)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.coerceToStringReturns.result1, fake.coerceToStringReturns.result2
}

func (fake *fakeCoercer) CoerceToStringCallCount() int {
	fake.coerceToStringMutex.RLock()
	defer fake.coerceToStringMutex.RUnlock()
	return len(fake.coerceToStringArgsForCall)
}

func (fake *fakeCoercer) CoerceToStringArgsForCall(i int) *model.Value {
	fake.coerceToStringMutex.RLock()
	defer fake.coerceToStringMutex.RUnlock()
	return fake.coerceToStringArgsForCall[i].value
}

func (fake *fakeCoercer) CoerceToStringReturns(result1 string, result2 error) {
	fake.CoerceToStringStub = nil
	fake.coerceToStringReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *fakeCoercer) CoerceToStringReturnsOnCall(i int, result1 string, result2 error) {
	fake.CoerceToStringStub = nil
	if fake.coerceToStringReturnsOnCall == nil {
		fake.coerceToStringReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.coerceToStringReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *fakeCoercer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.coerceToNumberMutex.RLock()
	defer fake.coerceToNumberMutex.RUnlock()
	fake.coerceToObjectMutex.RLock()
	defer fake.coerceToObjectMutex.RUnlock()
	fake.coerceToStringMutex.RLock()
	defer fake.coerceToStringMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *fakeCoercer) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
