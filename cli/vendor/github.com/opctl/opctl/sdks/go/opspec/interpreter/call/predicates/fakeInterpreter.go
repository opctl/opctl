// Code generated by counterfeiter. DO NOT EDIT.
package predicates

import (
	"sync"

	"github.com/opctl/opctl/sdks/go/types"
)

type FakeInterpreter struct {
	InterpretStub        func(types.DataHandle, []*types.SCGPredicate, map[string]*types.Value) (bool, error)
	interpretMutex       sync.RWMutex
	interpretArgsForCall []struct {
		arg1 types.DataHandle
		arg2 []*types.SCGPredicate
		arg3 map[string]*types.Value
	}
	interpretReturns struct {
		result1 bool
		result2 error
	}
	interpretReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeInterpreter) Interpret(arg1 types.DataHandle, arg2 []*types.SCGPredicate, arg3 map[string]*types.Value) (bool, error) {
	var arg2Copy []*types.SCGPredicate
	if arg2 != nil {
		arg2Copy = make([]*types.SCGPredicate, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.interpretMutex.Lock()
	ret, specificReturn := fake.interpretReturnsOnCall[len(fake.interpretArgsForCall)]
	fake.interpretArgsForCall = append(fake.interpretArgsForCall, struct {
		arg1 types.DataHandle
		arg2 []*types.SCGPredicate
		arg3 map[string]*types.Value
	}{arg1, arg2Copy, arg3})
	fake.recordInvocation("Interpret", []interface{}{arg1, arg2Copy, arg3})
	fake.interpretMutex.Unlock()
	if fake.InterpretStub != nil {
		return fake.InterpretStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.interpretReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeInterpreter) InterpretCallCount() int {
	fake.interpretMutex.RLock()
	defer fake.interpretMutex.RUnlock()
	return len(fake.interpretArgsForCall)
}

func (fake *FakeInterpreter) InterpretCalls(stub func(types.DataHandle, []*types.SCGPredicate, map[string]*types.Value) (bool, error)) {
	fake.interpretMutex.Lock()
	defer fake.interpretMutex.Unlock()
	fake.InterpretStub = stub
}

func (fake *FakeInterpreter) InterpretArgsForCall(i int) (types.DataHandle, []*types.SCGPredicate, map[string]*types.Value) {
	fake.interpretMutex.RLock()
	defer fake.interpretMutex.RUnlock()
	argsForCall := fake.interpretArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeInterpreter) InterpretReturns(result1 bool, result2 error) {
	fake.interpretMutex.Lock()
	defer fake.interpretMutex.Unlock()
	fake.InterpretStub = nil
	fake.interpretReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeInterpreter) InterpretReturnsOnCall(i int, result1 bool, result2 error) {
	fake.interpretMutex.Lock()
	defer fake.interpretMutex.Unlock()
	fake.InterpretStub = nil
	if fake.interpretReturnsOnCall == nil {
		fake.interpretReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.interpretReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeInterpreter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.interpretMutex.RLock()
	defer fake.interpretMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeInterpreter) recordInvocation(key string, args []interface{}) {
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

var _ Interpreter = new(FakeInterpreter)
