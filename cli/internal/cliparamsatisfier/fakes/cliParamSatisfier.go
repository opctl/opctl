// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
	"github.com/opctl/opctl/sdks/go/model"
)

type FakeCLIParamSatisfier struct {
	NewCliPromptInputSrcStub        func(map[string]*model.ParamSpec) inputsrc.InputSrc
	newCliPromptInputSrcMutex       sync.RWMutex
	newCliPromptInputSrcArgsForCall []struct {
		arg1 map[string]*model.ParamSpec
	}
	newCliPromptInputSrcReturns struct {
		result1 inputsrc.InputSrc
	}
	newCliPromptInputSrcReturnsOnCall map[int]struct {
		result1 inputsrc.InputSrc
	}
	NewEnvVarInputSrcStub        func() inputsrc.InputSrc
	newEnvVarInputSrcMutex       sync.RWMutex
	newEnvVarInputSrcArgsForCall []struct {
	}
	newEnvVarInputSrcReturns struct {
		result1 inputsrc.InputSrc
	}
	newEnvVarInputSrcReturnsOnCall map[int]struct {
		result1 inputsrc.InputSrc
	}
	NewParamDefaultInputSrcStub        func(map[string]*model.ParamSpec) inputsrc.InputSrc
	newParamDefaultInputSrcMutex       sync.RWMutex
	newParamDefaultInputSrcArgsForCall []struct {
		arg1 map[string]*model.ParamSpec
	}
	newParamDefaultInputSrcReturns struct {
		result1 inputsrc.InputSrc
	}
	newParamDefaultInputSrcReturnsOnCall map[int]struct {
		result1 inputsrc.InputSrc
	}
	NewSliceInputSrcStub        func([]string, string) inputsrc.InputSrc
	newSliceInputSrcMutex       sync.RWMutex
	newSliceInputSrcArgsForCall []struct {
		arg1 []string
		arg2 string
	}
	newSliceInputSrcReturns struct {
		result1 inputsrc.InputSrc
	}
	newSliceInputSrcReturnsOnCall map[int]struct {
		result1 inputsrc.InputSrc
	}
	NewYMLFileInputSrcStub        func(string) (inputsrc.InputSrc, error)
	newYMLFileInputSrcMutex       sync.RWMutex
	newYMLFileInputSrcArgsForCall []struct {
		arg1 string
	}
	newYMLFileInputSrcReturns struct {
		result1 inputsrc.InputSrc
		result2 error
	}
	newYMLFileInputSrcReturnsOnCall map[int]struct {
		result1 inputsrc.InputSrc
		result2 error
	}
	SatisfyStub        func(cliparamsatisfier.InputSourcer, map[string]*model.ParamSpec) (map[string]*ipld.Node, error)
	satisfyMutex       sync.RWMutex
	satisfyArgsForCall []struct {
		arg1 cliparamsatisfier.InputSourcer
		arg2 map[string]*model.ParamSpec
	}
	satisfyReturns struct {
		result1 map[string]*ipld.Node
		result2 error
	}
	satisfyReturnsOnCall map[int]struct {
		result1 map[string]*ipld.Node
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCLIParamSatisfier) NewCliPromptInputSrc(arg1 map[string]*model.ParamSpec) inputsrc.InputSrc {
	fake.newCliPromptInputSrcMutex.Lock()
	ret, specificReturn := fake.newCliPromptInputSrcReturnsOnCall[len(fake.newCliPromptInputSrcArgsForCall)]
	fake.newCliPromptInputSrcArgsForCall = append(fake.newCliPromptInputSrcArgsForCall, struct {
		arg1 map[string]*model.ParamSpec
	}{arg1})
	fake.recordInvocation("NewCliPromptInputSrc", []interface{}{arg1})
	fake.newCliPromptInputSrcMutex.Unlock()
	if fake.NewCliPromptInputSrcStub != nil {
		return fake.NewCliPromptInputSrcStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newCliPromptInputSrcReturns
	return fakeReturns.result1
}

func (fake *FakeCLIParamSatisfier) NewCliPromptInputSrcCallCount() int {
	fake.newCliPromptInputSrcMutex.RLock()
	defer fake.newCliPromptInputSrcMutex.RUnlock()
	return len(fake.newCliPromptInputSrcArgsForCall)
}

func (fake *FakeCLIParamSatisfier) NewCliPromptInputSrcCalls(stub func(map[string]*model.ParamSpec) inputsrc.InputSrc) {
	fake.newCliPromptInputSrcMutex.Lock()
	defer fake.newCliPromptInputSrcMutex.Unlock()
	fake.NewCliPromptInputSrcStub = stub
}

func (fake *FakeCLIParamSatisfier) NewCliPromptInputSrcArgsForCall(i int) map[string]*model.ParamSpec {
	fake.newCliPromptInputSrcMutex.RLock()
	defer fake.newCliPromptInputSrcMutex.RUnlock()
	argsForCall := fake.newCliPromptInputSrcArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCLIParamSatisfier) NewCliPromptInputSrcReturns(result1 inputsrc.InputSrc) {
	fake.newCliPromptInputSrcMutex.Lock()
	defer fake.newCliPromptInputSrcMutex.Unlock()
	fake.NewCliPromptInputSrcStub = nil
	fake.newCliPromptInputSrcReturns = struct {
		result1 inputsrc.InputSrc
	}{result1}
}

func (fake *FakeCLIParamSatisfier) NewCliPromptInputSrcReturnsOnCall(i int, result1 inputsrc.InputSrc) {
	fake.newCliPromptInputSrcMutex.Lock()
	defer fake.newCliPromptInputSrcMutex.Unlock()
	fake.NewCliPromptInputSrcStub = nil
	if fake.newCliPromptInputSrcReturnsOnCall == nil {
		fake.newCliPromptInputSrcReturnsOnCall = make(map[int]struct {
			result1 inputsrc.InputSrc
		})
	}
	fake.newCliPromptInputSrcReturnsOnCall[i] = struct {
		result1 inputsrc.InputSrc
	}{result1}
}

func (fake *FakeCLIParamSatisfier) NewEnvVarInputSrc() inputsrc.InputSrc {
	fake.newEnvVarInputSrcMutex.Lock()
	ret, specificReturn := fake.newEnvVarInputSrcReturnsOnCall[len(fake.newEnvVarInputSrcArgsForCall)]
	fake.newEnvVarInputSrcArgsForCall = append(fake.newEnvVarInputSrcArgsForCall, struct {
	}{})
	fake.recordInvocation("NewEnvVarInputSrc", []interface{}{})
	fake.newEnvVarInputSrcMutex.Unlock()
	if fake.NewEnvVarInputSrcStub != nil {
		return fake.NewEnvVarInputSrcStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newEnvVarInputSrcReturns
	return fakeReturns.result1
}

func (fake *FakeCLIParamSatisfier) NewEnvVarInputSrcCallCount() int {
	fake.newEnvVarInputSrcMutex.RLock()
	defer fake.newEnvVarInputSrcMutex.RUnlock()
	return len(fake.newEnvVarInputSrcArgsForCall)
}

func (fake *FakeCLIParamSatisfier) NewEnvVarInputSrcCalls(stub func() inputsrc.InputSrc) {
	fake.newEnvVarInputSrcMutex.Lock()
	defer fake.newEnvVarInputSrcMutex.Unlock()
	fake.NewEnvVarInputSrcStub = stub
}

func (fake *FakeCLIParamSatisfier) NewEnvVarInputSrcReturns(result1 inputsrc.InputSrc) {
	fake.newEnvVarInputSrcMutex.Lock()
	defer fake.newEnvVarInputSrcMutex.Unlock()
	fake.NewEnvVarInputSrcStub = nil
	fake.newEnvVarInputSrcReturns = struct {
		result1 inputsrc.InputSrc
	}{result1}
}

func (fake *FakeCLIParamSatisfier) NewEnvVarInputSrcReturnsOnCall(i int, result1 inputsrc.InputSrc) {
	fake.newEnvVarInputSrcMutex.Lock()
	defer fake.newEnvVarInputSrcMutex.Unlock()
	fake.NewEnvVarInputSrcStub = nil
	if fake.newEnvVarInputSrcReturnsOnCall == nil {
		fake.newEnvVarInputSrcReturnsOnCall = make(map[int]struct {
			result1 inputsrc.InputSrc
		})
	}
	fake.newEnvVarInputSrcReturnsOnCall[i] = struct {
		result1 inputsrc.InputSrc
	}{result1}
}

func (fake *FakeCLIParamSatisfier) NewParamDefaultInputSrc(arg1 map[string]*model.ParamSpec) inputsrc.InputSrc {
	fake.newParamDefaultInputSrcMutex.Lock()
	ret, specificReturn := fake.newParamDefaultInputSrcReturnsOnCall[len(fake.newParamDefaultInputSrcArgsForCall)]
	fake.newParamDefaultInputSrcArgsForCall = append(fake.newParamDefaultInputSrcArgsForCall, struct {
		arg1 map[string]*model.ParamSpec
	}{arg1})
	fake.recordInvocation("NewParamDefaultInputSrc", []interface{}{arg1})
	fake.newParamDefaultInputSrcMutex.Unlock()
	if fake.NewParamDefaultInputSrcStub != nil {
		return fake.NewParamDefaultInputSrcStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newParamDefaultInputSrcReturns
	return fakeReturns.result1
}

func (fake *FakeCLIParamSatisfier) NewParamDefaultInputSrcCallCount() int {
	fake.newParamDefaultInputSrcMutex.RLock()
	defer fake.newParamDefaultInputSrcMutex.RUnlock()
	return len(fake.newParamDefaultInputSrcArgsForCall)
}

func (fake *FakeCLIParamSatisfier) NewParamDefaultInputSrcCalls(stub func(map[string]*model.ParamSpec) inputsrc.InputSrc) {
	fake.newParamDefaultInputSrcMutex.Lock()
	defer fake.newParamDefaultInputSrcMutex.Unlock()
	fake.NewParamDefaultInputSrcStub = stub
}

func (fake *FakeCLIParamSatisfier) NewParamDefaultInputSrcArgsForCall(i int) map[string]*model.ParamSpec {
	fake.newParamDefaultInputSrcMutex.RLock()
	defer fake.newParamDefaultInputSrcMutex.RUnlock()
	argsForCall := fake.newParamDefaultInputSrcArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCLIParamSatisfier) NewParamDefaultInputSrcReturns(result1 inputsrc.InputSrc) {
	fake.newParamDefaultInputSrcMutex.Lock()
	defer fake.newParamDefaultInputSrcMutex.Unlock()
	fake.NewParamDefaultInputSrcStub = nil
	fake.newParamDefaultInputSrcReturns = struct {
		result1 inputsrc.InputSrc
	}{result1}
}

func (fake *FakeCLIParamSatisfier) NewParamDefaultInputSrcReturnsOnCall(i int, result1 inputsrc.InputSrc) {
	fake.newParamDefaultInputSrcMutex.Lock()
	defer fake.newParamDefaultInputSrcMutex.Unlock()
	fake.NewParamDefaultInputSrcStub = nil
	if fake.newParamDefaultInputSrcReturnsOnCall == nil {
		fake.newParamDefaultInputSrcReturnsOnCall = make(map[int]struct {
			result1 inputsrc.InputSrc
		})
	}
	fake.newParamDefaultInputSrcReturnsOnCall[i] = struct {
		result1 inputsrc.InputSrc
	}{result1}
}

func (fake *FakeCLIParamSatisfier) NewSliceInputSrc(arg1 []string, arg2 string) inputsrc.InputSrc {
	var arg1Copy []string
	if arg1 != nil {
		arg1Copy = make([]string, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.newSliceInputSrcMutex.Lock()
	ret, specificReturn := fake.newSliceInputSrcReturnsOnCall[len(fake.newSliceInputSrcArgsForCall)]
	fake.newSliceInputSrcArgsForCall = append(fake.newSliceInputSrcArgsForCall, struct {
		arg1 []string
		arg2 string
	}{arg1Copy, arg2})
	fake.recordInvocation("NewSliceInputSrc", []interface{}{arg1Copy, arg2})
	fake.newSliceInputSrcMutex.Unlock()
	if fake.NewSliceInputSrcStub != nil {
		return fake.NewSliceInputSrcStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newSliceInputSrcReturns
	return fakeReturns.result1
}

func (fake *FakeCLIParamSatisfier) NewSliceInputSrcCallCount() int {
	fake.newSliceInputSrcMutex.RLock()
	defer fake.newSliceInputSrcMutex.RUnlock()
	return len(fake.newSliceInputSrcArgsForCall)
}

func (fake *FakeCLIParamSatisfier) NewSliceInputSrcCalls(stub func([]string, string) inputsrc.InputSrc) {
	fake.newSliceInputSrcMutex.Lock()
	defer fake.newSliceInputSrcMutex.Unlock()
	fake.NewSliceInputSrcStub = stub
}

func (fake *FakeCLIParamSatisfier) NewSliceInputSrcArgsForCall(i int) ([]string, string) {
	fake.newSliceInputSrcMutex.RLock()
	defer fake.newSliceInputSrcMutex.RUnlock()
	argsForCall := fake.newSliceInputSrcArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCLIParamSatisfier) NewSliceInputSrcReturns(result1 inputsrc.InputSrc) {
	fake.newSliceInputSrcMutex.Lock()
	defer fake.newSliceInputSrcMutex.Unlock()
	fake.NewSliceInputSrcStub = nil
	fake.newSliceInputSrcReturns = struct {
		result1 inputsrc.InputSrc
	}{result1}
}

func (fake *FakeCLIParamSatisfier) NewSliceInputSrcReturnsOnCall(i int, result1 inputsrc.InputSrc) {
	fake.newSliceInputSrcMutex.Lock()
	defer fake.newSliceInputSrcMutex.Unlock()
	fake.NewSliceInputSrcStub = nil
	if fake.newSliceInputSrcReturnsOnCall == nil {
		fake.newSliceInputSrcReturnsOnCall = make(map[int]struct {
			result1 inputsrc.InputSrc
		})
	}
	fake.newSliceInputSrcReturnsOnCall[i] = struct {
		result1 inputsrc.InputSrc
	}{result1}
}

func (fake *FakeCLIParamSatisfier) NewYMLFileInputSrc(arg1 string) (inputsrc.InputSrc, error) {
	fake.newYMLFileInputSrcMutex.Lock()
	ret, specificReturn := fake.newYMLFileInputSrcReturnsOnCall[len(fake.newYMLFileInputSrcArgsForCall)]
	fake.newYMLFileInputSrcArgsForCall = append(fake.newYMLFileInputSrcArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("NewYMLFileInputSrc", []interface{}{arg1})
	fake.newYMLFileInputSrcMutex.Unlock()
	if fake.NewYMLFileInputSrcStub != nil {
		return fake.NewYMLFileInputSrcStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newYMLFileInputSrcReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCLIParamSatisfier) NewYMLFileInputSrcCallCount() int {
	fake.newYMLFileInputSrcMutex.RLock()
	defer fake.newYMLFileInputSrcMutex.RUnlock()
	return len(fake.newYMLFileInputSrcArgsForCall)
}

func (fake *FakeCLIParamSatisfier) NewYMLFileInputSrcCalls(stub func(string) (inputsrc.InputSrc, error)) {
	fake.newYMLFileInputSrcMutex.Lock()
	defer fake.newYMLFileInputSrcMutex.Unlock()
	fake.NewYMLFileInputSrcStub = stub
}

func (fake *FakeCLIParamSatisfier) NewYMLFileInputSrcArgsForCall(i int) string {
	fake.newYMLFileInputSrcMutex.RLock()
	defer fake.newYMLFileInputSrcMutex.RUnlock()
	argsForCall := fake.newYMLFileInputSrcArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCLIParamSatisfier) NewYMLFileInputSrcReturns(result1 inputsrc.InputSrc, result2 error) {
	fake.newYMLFileInputSrcMutex.Lock()
	defer fake.newYMLFileInputSrcMutex.Unlock()
	fake.NewYMLFileInputSrcStub = nil
	fake.newYMLFileInputSrcReturns = struct {
		result1 inputsrc.InputSrc
		result2 error
	}{result1, result2}
}

func (fake *FakeCLIParamSatisfier) NewYMLFileInputSrcReturnsOnCall(i int, result1 inputsrc.InputSrc, result2 error) {
	fake.newYMLFileInputSrcMutex.Lock()
	defer fake.newYMLFileInputSrcMutex.Unlock()
	fake.NewYMLFileInputSrcStub = nil
	if fake.newYMLFileInputSrcReturnsOnCall == nil {
		fake.newYMLFileInputSrcReturnsOnCall = make(map[int]struct {
			result1 inputsrc.InputSrc
			result2 error
		})
	}
	fake.newYMLFileInputSrcReturnsOnCall[i] = struct {
		result1 inputsrc.InputSrc
		result2 error
	}{result1, result2}
}

func (fake *FakeCLIParamSatisfier) Satisfy(arg1 cliparamsatisfier.InputSourcer, arg2 map[string]*model.ParamSpec) (map[string]*ipld.Node, error) {
	fake.satisfyMutex.Lock()
	ret, specificReturn := fake.satisfyReturnsOnCall[len(fake.satisfyArgsForCall)]
	fake.satisfyArgsForCall = append(fake.satisfyArgsForCall, struct {
		arg1 cliparamsatisfier.InputSourcer
		arg2 map[string]*model.ParamSpec
	}{arg1, arg2})
	fake.recordInvocation("Satisfy", []interface{}{arg1, arg2})
	fake.satisfyMutex.Unlock()
	if fake.SatisfyStub != nil {
		return fake.SatisfyStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.satisfyReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCLIParamSatisfier) SatisfyCallCount() int {
	fake.satisfyMutex.RLock()
	defer fake.satisfyMutex.RUnlock()
	return len(fake.satisfyArgsForCall)
}

func (fake *FakeCLIParamSatisfier) SatisfyCalls(stub func(cliparamsatisfier.InputSourcer, map[string]*model.ParamSpec) (map[string]*ipld.Node, error)) {
	fake.satisfyMutex.Lock()
	defer fake.satisfyMutex.Unlock()
	fake.SatisfyStub = stub
}

func (fake *FakeCLIParamSatisfier) SatisfyArgsForCall(i int) (cliparamsatisfier.InputSourcer, map[string]*model.ParamSpec) {
	fake.satisfyMutex.RLock()
	defer fake.satisfyMutex.RUnlock()
	argsForCall := fake.satisfyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCLIParamSatisfier) SatisfyReturns(result1 map[string]*ipld.Node, result2 error) {
	fake.satisfyMutex.Lock()
	defer fake.satisfyMutex.Unlock()
	fake.SatisfyStub = nil
	fake.satisfyReturns = struct {
		result1 map[string]*ipld.Node
		result2 error
	}{result1, result2}
}

func (fake *FakeCLIParamSatisfier) SatisfyReturnsOnCall(i int, result1 map[string]*ipld.Node, result2 error) {
	fake.satisfyMutex.Lock()
	defer fake.satisfyMutex.Unlock()
	fake.SatisfyStub = nil
	if fake.satisfyReturnsOnCall == nil {
		fake.satisfyReturnsOnCall = make(map[int]struct {
			result1 map[string]*ipld.Node
			result2 error
		})
	}
	fake.satisfyReturnsOnCall[i] = struct {
		result1 map[string]*ipld.Node
		result2 error
	}{result1, result2}
}

func (fake *FakeCLIParamSatisfier) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newCliPromptInputSrcMutex.RLock()
	defer fake.newCliPromptInputSrcMutex.RUnlock()
	fake.newEnvVarInputSrcMutex.RLock()
	defer fake.newEnvVarInputSrcMutex.RUnlock()
	fake.newParamDefaultInputSrcMutex.RLock()
	defer fake.newParamDefaultInputSrcMutex.RUnlock()
	fake.newSliceInputSrcMutex.RLock()
	defer fake.newSliceInputSrcMutex.RUnlock()
	fake.newYMLFileInputSrcMutex.RLock()
	defer fake.newYMLFileInputSrcMutex.RUnlock()
	fake.satisfyMutex.RLock()
	defer fake.satisfyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCLIParamSatisfier) recordInvocation(key string, args []interface{}) {
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

var _ cliparamsatisfier.CLIParamSatisfier = new(FakeCLIParamSatisfier)
