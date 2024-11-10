// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"

	"github.com/opctl/opctl/sdks/go/model"
)

type FakeSerialLoopCaller struct {
	CallStub        func(context.Context, string, map[string]*ipld.Node, model.SerialLoopCallSpec, string, *string, string) (map[string]*ipld.Node, error)
	callMutex       sync.RWMutex
	callArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 map[string]*ipld.Node
		arg4 model.SerialLoopCallSpec
		arg5 string
		arg6 *string
		arg7 string
	}
	callReturns struct {
		result1 map[string]*ipld.Node
		result2 error
	}
	callReturnsOnCall map[int]struct {
		result1 map[string]*ipld.Node
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSerialLoopCaller) Call(arg1 context.Context, arg2 string, arg3 map[string]*ipld.Node, arg4 model.SerialLoopCallSpec, arg5 string, arg6 *string, arg7 string) (map[string]*ipld.Node, error) {
	fake.callMutex.Lock()
	ret, specificReturn := fake.callReturnsOnCall[len(fake.callArgsForCall)]
	fake.callArgsForCall = append(fake.callArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 map[string]*ipld.Node
		arg4 model.SerialLoopCallSpec
		arg5 string
		arg6 *string
		arg7 string
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.recordInvocation("Call", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.callMutex.Unlock()
	if fake.CallStub != nil {
		return fake.CallStub(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.callReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSerialLoopCaller) CallCallCount() int {
	fake.callMutex.RLock()
	defer fake.callMutex.RUnlock()
	return len(fake.callArgsForCall)
}

func (fake *FakeSerialLoopCaller) CallCalls(stub func(context.Context, string, map[string]*ipld.Node, model.SerialLoopCallSpec, string, *string, string) (map[string]*ipld.Node, error)) {
	fake.callMutex.Lock()
	defer fake.callMutex.Unlock()
	fake.CallStub = stub
}

func (fake *FakeSerialLoopCaller) CallArgsForCall(i int) (context.Context, string, map[string]*ipld.Node, model.SerialLoopCallSpec, string, *string, string) {
	fake.callMutex.RLock()
	defer fake.callMutex.RUnlock()
	argsForCall := fake.callArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7
}

func (fake *FakeSerialLoopCaller) CallReturns(result1 map[string]*ipld.Node, result2 error) {
	fake.callMutex.Lock()
	defer fake.callMutex.Unlock()
	fake.CallStub = nil
	fake.callReturns = struct {
		result1 map[string]*ipld.Node
		result2 error
	}{result1, result2}
}

func (fake *FakeSerialLoopCaller) CallReturnsOnCall(i int, result1 map[string]*ipld.Node, result2 error) {
	fake.callMutex.Lock()
	defer fake.callMutex.Unlock()
	fake.CallStub = nil
	if fake.callReturnsOnCall == nil {
		fake.callReturnsOnCall = make(map[int]struct {
			result1 map[string]*ipld.Node
			result2 error
		})
	}
	fake.callReturnsOnCall[i] = struct {
		result1 map[string]*ipld.Node
		result2 error
	}{result1, result2}
}

func (fake *FakeSerialLoopCaller) Invocations() map[string][][]interface{} {
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

func (fake *FakeSerialLoopCaller) recordInvocation(key string, args []interface{}) {
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
