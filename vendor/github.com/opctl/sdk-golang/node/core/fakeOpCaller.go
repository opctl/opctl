// Code generated by counterfeiter. DO NOT EDIT.
package core

import (
	"context"
	"sync"

	"github.com/opctl/sdk-golang/model"
)

type fakeOpCaller struct {
	CallStub        func(context.Context, *model.DCGOpCall, map[string]*model.Value, *string, *model.SCGOpCall) error
	callMutex       sync.RWMutex
	callArgsForCall []struct {
		arg1 context.Context
		arg2 *model.DCGOpCall
		arg3 map[string]*model.Value
		arg4 *string
		arg5 *model.SCGOpCall
	}
	callReturns struct {
		result1 error
	}
	callReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *fakeOpCaller) Call(arg1 context.Context, arg2 *model.DCGOpCall, arg3 map[string]*model.Value, arg4 *string, arg5 *model.SCGOpCall) error {
	fake.callMutex.Lock()
	ret, specificReturn := fake.callReturnsOnCall[len(fake.callArgsForCall)]
	fake.callArgsForCall = append(fake.callArgsForCall, struct {
		arg1 context.Context
		arg2 *model.DCGOpCall
		arg3 map[string]*model.Value
		arg4 *string
		arg5 *model.SCGOpCall
	}{arg1, arg2, arg3, arg4, arg5})
	fake.recordInvocation("Call", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.callMutex.Unlock()
	if fake.CallStub != nil {
		return fake.CallStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.callReturns
	return fakeReturns.result1
}

func (fake *fakeOpCaller) CallCallCount() int {
	fake.callMutex.RLock()
	defer fake.callMutex.RUnlock()
	return len(fake.callArgsForCall)
}

func (fake *fakeOpCaller) CallCalls(stub func(context.Context, *model.DCGOpCall, map[string]*model.Value, *string, *model.SCGOpCall) error) {
	fake.callMutex.Lock()
	defer fake.callMutex.Unlock()
	fake.CallStub = stub
}

func (fake *fakeOpCaller) CallArgsForCall(i int) (context.Context, *model.DCGOpCall, map[string]*model.Value, *string, *model.SCGOpCall) {
	fake.callMutex.RLock()
	defer fake.callMutex.RUnlock()
	argsForCall := fake.callArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *fakeOpCaller) CallReturns(result1 error) {
	fake.callMutex.Lock()
	defer fake.callMutex.Unlock()
	fake.CallStub = nil
	fake.callReturns = struct {
		result1 error
	}{result1}
}

func (fake *fakeOpCaller) CallReturnsOnCall(i int, result1 error) {
	fake.callMutex.Lock()
	defer fake.callMutex.Unlock()
	fake.CallStub = nil
	if fake.callReturnsOnCall == nil {
		fake.callReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.callReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *fakeOpCaller) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.callMutex.RLock()
	defer fake.callMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *fakeOpCaller) recordInvocation(key string, args []interface{}) {
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