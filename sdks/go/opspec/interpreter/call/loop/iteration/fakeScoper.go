// Code generated by counterfeiter. DO NOT EDIT.
package iteration

import (
	"sync"

	"github.com/opctl/opctl/sdks/go/types"
)

type FakeScoper struct {
	ScopeStub        func(int, map[string]*types.Value, interface{}, *types.SCGLoopVars, types.DataHandle) (map[string]*types.Value, error)
	scopeMutex       sync.RWMutex
	scopeArgsForCall []struct {
		arg1 int
		arg2 map[string]*types.Value
		arg3 interface{}
		arg4 *types.SCGLoopVars
		arg5 types.DataHandle
	}
	scopeReturns struct {
		result1 map[string]*types.Value
		result2 error
	}
	scopeReturnsOnCall map[int]struct {
		result1 map[string]*types.Value
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeScoper) Scope(arg1 int, arg2 map[string]*types.Value, arg3 interface{}, arg4 *types.SCGLoopVars, arg5 types.DataHandle) (map[string]*types.Value, error) {
	fake.scopeMutex.Lock()
	ret, specificReturn := fake.scopeReturnsOnCall[len(fake.scopeArgsForCall)]
	fake.scopeArgsForCall = append(fake.scopeArgsForCall, struct {
		arg1 int
		arg2 map[string]*types.Value
		arg3 interface{}
		arg4 *types.SCGLoopVars
		arg5 types.DataHandle
	}{arg1, arg2, arg3, arg4, arg5})
	fake.recordInvocation("Scope", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.scopeMutex.Unlock()
	if fake.ScopeStub != nil {
		return fake.ScopeStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.scopeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeScoper) ScopeCallCount() int {
	fake.scopeMutex.RLock()
	defer fake.scopeMutex.RUnlock()
	return len(fake.scopeArgsForCall)
}

func (fake *FakeScoper) ScopeCalls(stub func(int, map[string]*types.Value, interface{}, *types.SCGLoopVars, types.DataHandle) (map[string]*types.Value, error)) {
	fake.scopeMutex.Lock()
	defer fake.scopeMutex.Unlock()
	fake.ScopeStub = stub
}

func (fake *FakeScoper) ScopeArgsForCall(i int) (int, map[string]*types.Value, interface{}, *types.SCGLoopVars, types.DataHandle) {
	fake.scopeMutex.RLock()
	defer fake.scopeMutex.RUnlock()
	argsForCall := fake.scopeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeScoper) ScopeReturns(result1 map[string]*types.Value, result2 error) {
	fake.scopeMutex.Lock()
	defer fake.scopeMutex.Unlock()
	fake.ScopeStub = nil
	fake.scopeReturns = struct {
		result1 map[string]*types.Value
		result2 error
	}{result1, result2}
}

func (fake *FakeScoper) ScopeReturnsOnCall(i int, result1 map[string]*types.Value, result2 error) {
	fake.scopeMutex.Lock()
	defer fake.scopeMutex.Unlock()
	fake.ScopeStub = nil
	if fake.scopeReturnsOnCall == nil {
		fake.scopeReturnsOnCall = make(map[int]struct {
			result1 map[string]*types.Value
			result2 error
		})
	}
	fake.scopeReturnsOnCall[i] = struct {
		result1 map[string]*types.Value
		result2 error
	}{result1, result2}
}

func (fake *FakeScoper) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.scopeMutex.RLock()
	defer fake.scopeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeScoper) recordInvocation(key string, args []interface{}) {
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

var _ Scoper = new(FakeScoper)
