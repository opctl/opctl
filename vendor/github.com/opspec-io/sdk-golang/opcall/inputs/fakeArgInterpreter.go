// Code generated by counterfeiter. DO NOT EDIT.
package inputs

import (
	"sync"

	"github.com/opspec-io/sdk-golang/model"
)

type fakeArgInterpreter struct {
	InterpretStub        func(name string, value interface{}, param *model.Param, parentPkgHandle model.PkgHandle, scope map[string]*model.Value, opScratchDir string) (*model.Value, error)
	interpretMutex       sync.RWMutex
	interpretArgsForCall []struct {
		name            string
		value           interface{}
		param           *model.Param
		parentPkgHandle model.PkgHandle
		scope           map[string]*model.Value
		opScratchDir    string
	}
	interpretReturns struct {
		result1 *model.Value
		result2 error
	}
	interpretReturnsOnCall map[int]struct {
		result1 *model.Value
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *fakeArgInterpreter) Interpret(name string, value interface{}, param *model.Param, parentPkgHandle model.PkgHandle, scope map[string]*model.Value, opScratchDir string) (*model.Value, error) {
	fake.interpretMutex.Lock()
	ret, specificReturn := fake.interpretReturnsOnCall[len(fake.interpretArgsForCall)]
	fake.interpretArgsForCall = append(fake.interpretArgsForCall, struct {
		name            string
		value           interface{}
		param           *model.Param
		parentPkgHandle model.PkgHandle
		scope           map[string]*model.Value
		opScratchDir    string
	}{name, value, param, parentPkgHandle, scope, opScratchDir})
	fake.recordInvocation("Interpret", []interface{}{name, value, param, parentPkgHandle, scope, opScratchDir})
	fake.interpretMutex.Unlock()
	if fake.InterpretStub != nil {
		return fake.InterpretStub(name, value, param, parentPkgHandle, scope, opScratchDir)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.interpretReturns.result1, fake.interpretReturns.result2
}

func (fake *fakeArgInterpreter) InterpretCallCount() int {
	fake.interpretMutex.RLock()
	defer fake.interpretMutex.RUnlock()
	return len(fake.interpretArgsForCall)
}

func (fake *fakeArgInterpreter) InterpretArgsForCall(i int) (string, interface{}, *model.Param, model.PkgHandle, map[string]*model.Value, string) {
	fake.interpretMutex.RLock()
	defer fake.interpretMutex.RUnlock()
	return fake.interpretArgsForCall[i].name, fake.interpretArgsForCall[i].value, fake.interpretArgsForCall[i].param, fake.interpretArgsForCall[i].parentPkgHandle, fake.interpretArgsForCall[i].scope, fake.interpretArgsForCall[i].opScratchDir
}

func (fake *fakeArgInterpreter) InterpretReturns(result1 *model.Value, result2 error) {
	fake.InterpretStub = nil
	fake.interpretReturns = struct {
		result1 *model.Value
		result2 error
	}{result1, result2}
}

func (fake *fakeArgInterpreter) InterpretReturnsOnCall(i int, result1 *model.Value, result2 error) {
	fake.InterpretStub = nil
	if fake.interpretReturnsOnCall == nil {
		fake.interpretReturnsOnCall = make(map[int]struct {
			result1 *model.Value
			result2 error
		})
	}
	fake.interpretReturnsOnCall[i] = struct {
		result1 *model.Value
		result2 error
	}{result1, result2}
}

func (fake *fakeArgInterpreter) Invocations() map[string][][]interface{} {
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

func (fake *fakeArgInterpreter) recordInvocation(key string, args []interface{}) {
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
