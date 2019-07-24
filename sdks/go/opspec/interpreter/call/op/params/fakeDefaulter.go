// Code generated by counterfeiter. DO NOT EDIT.
package params

import (
	"sync"

	"github.com/opctl/opctl/sdks/go/types"
)

type FakeDefaulter struct {
	DefaultStub        func(args map[string]*types.Value, params map[string]*types.Param, opPath string) map[string]*types.Value
	defaultMutex       sync.RWMutex
	defaultArgsForCall []struct {
		args   map[string]*types.Value
		params map[string]*types.Param
		opPath string
	}
	defaultReturns struct {
		result1 map[string]*types.Value
	}
	defaultReturnsOnCall map[int]struct {
		result1 map[string]*types.Value
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDefaulter) Default(args map[string]*types.Value, params map[string]*types.Param, opPath string) map[string]*types.Value {
	fake.defaultMutex.Lock()
	ret, specificReturn := fake.defaultReturnsOnCall[len(fake.defaultArgsForCall)]
	fake.defaultArgsForCall = append(fake.defaultArgsForCall, struct {
		args   map[string]*types.Value
		params map[string]*types.Param
		opPath string
	}{args, params, opPath})
	fake.recordInvocation("Default", []interface{}{args, params, opPath})
	fake.defaultMutex.Unlock()
	if fake.DefaultStub != nil {
		return fake.DefaultStub(args, params, opPath)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.defaultReturns.result1
}

func (fake *FakeDefaulter) DefaultCallCount() int {
	fake.defaultMutex.RLock()
	defer fake.defaultMutex.RUnlock()
	return len(fake.defaultArgsForCall)
}

func (fake *FakeDefaulter) DefaultArgsForCall(i int) (map[string]*types.Value, map[string]*types.Param, string) {
	fake.defaultMutex.RLock()
	defer fake.defaultMutex.RUnlock()
	return fake.defaultArgsForCall[i].args, fake.defaultArgsForCall[i].params, fake.defaultArgsForCall[i].opPath
}

func (fake *FakeDefaulter) DefaultReturns(result1 map[string]*types.Value) {
	fake.DefaultStub = nil
	fake.defaultReturns = struct {
		result1 map[string]*types.Value
	}{result1}
}

func (fake *FakeDefaulter) DefaultReturnsOnCall(i int, result1 map[string]*types.Value) {
	fake.DefaultStub = nil
	if fake.defaultReturnsOnCall == nil {
		fake.defaultReturnsOnCall = make(map[int]struct {
			result1 map[string]*types.Value
		})
	}
	fake.defaultReturnsOnCall[i] = struct {
		result1 map[string]*types.Value
	}{result1}
}

func (fake *FakeDefaulter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.defaultMutex.RLock()
	defer fake.defaultMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDefaulter) recordInvocation(key string, args []interface{}) {
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

var _ Defaulter = new(FakeDefaulter)
