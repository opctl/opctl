// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/opctl/opctl/cli/internal/clicolorer"
)

type FakeCliColorer struct {
	AttentionStub        func(string) string
	attentionMutex       sync.RWMutex
	attentionArgsForCall []struct {
		arg1 string
	}
	attentionReturns struct {
		result1 string
	}
	attentionReturnsOnCall map[int]struct {
		result1 string
	}
	DisableColorStub        func()
	disableColorMutex       sync.RWMutex
	disableColorArgsForCall []struct {
	}
	ErrorStub        func(string) string
	errorMutex       sync.RWMutex
	errorArgsForCall []struct {
		arg1 string
	}
	errorReturns struct {
		result1 string
	}
	errorReturnsOnCall map[int]struct {
		result1 string
	}
	InfoStub        func(string) string
	infoMutex       sync.RWMutex
	infoArgsForCall []struct {
		arg1 string
	}
	infoReturns struct {
		result1 string
	}
	infoReturnsOnCall map[int]struct {
		result1 string
	}
	MutedStub        func(string) string
	mutedMutex       sync.RWMutex
	mutedArgsForCall []struct {
		arg1 string
	}
	mutedReturns struct {
		result1 string
	}
	mutedReturnsOnCall map[int]struct {
		result1 string
	}
	SuccessStub        func(string) string
	successMutex       sync.RWMutex
	successArgsForCall []struct {
		arg1 string
	}
	successReturns struct {
		result1 string
	}
	successReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCliColorer) Attention(arg1 string) string {
	fake.attentionMutex.Lock()
	ret, specificReturn := fake.attentionReturnsOnCall[len(fake.attentionArgsForCall)]
	fake.attentionArgsForCall = append(fake.attentionArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Attention", []interface{}{arg1})
	fake.attentionMutex.Unlock()
	if fake.AttentionStub != nil {
		return fake.AttentionStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.attentionReturns
	return fakeReturns.result1
}

func (fake *FakeCliColorer) AttentionCallCount() int {
	fake.attentionMutex.RLock()
	defer fake.attentionMutex.RUnlock()
	return len(fake.attentionArgsForCall)
}

func (fake *FakeCliColorer) AttentionCalls(stub func(string) string) {
	fake.attentionMutex.Lock()
	defer fake.attentionMutex.Unlock()
	fake.AttentionStub = stub
}

func (fake *FakeCliColorer) AttentionArgsForCall(i int) string {
	fake.attentionMutex.RLock()
	defer fake.attentionMutex.RUnlock()
	argsForCall := fake.attentionArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCliColorer) AttentionReturns(result1 string) {
	fake.attentionMutex.Lock()
	defer fake.attentionMutex.Unlock()
	fake.AttentionStub = nil
	fake.attentionReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) AttentionReturnsOnCall(i int, result1 string) {
	fake.attentionMutex.Lock()
	defer fake.attentionMutex.Unlock()
	fake.AttentionStub = nil
	if fake.attentionReturnsOnCall == nil {
		fake.attentionReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.attentionReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) DisableColor() {
	fake.disableColorMutex.Lock()
	fake.disableColorArgsForCall = append(fake.disableColorArgsForCall, struct {
	}{})
	fake.recordInvocation("DisableColor", []interface{}{})
	fake.disableColorMutex.Unlock()
	if fake.DisableColorStub != nil {
		fake.DisableColorStub()
	}
}

func (fake *FakeCliColorer) DisableColorCallCount() int {
	fake.disableColorMutex.RLock()
	defer fake.disableColorMutex.RUnlock()
	return len(fake.disableColorArgsForCall)
}

func (fake *FakeCliColorer) DisableColorCalls(stub func()) {
	fake.disableColorMutex.Lock()
	defer fake.disableColorMutex.Unlock()
	fake.DisableColorStub = stub
}

func (fake *FakeCliColorer) Error(arg1 string) string {
	fake.errorMutex.Lock()
	ret, specificReturn := fake.errorReturnsOnCall[len(fake.errorArgsForCall)]
	fake.errorArgsForCall = append(fake.errorArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Error", []interface{}{arg1})
	fake.errorMutex.Unlock()
	if fake.ErrorStub != nil {
		return fake.ErrorStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.errorReturns
	return fakeReturns.result1
}

func (fake *FakeCliColorer) ErrorCallCount() int {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return len(fake.errorArgsForCall)
}

func (fake *FakeCliColorer) ErrorCalls(stub func(string) string) {
	fake.errorMutex.Lock()
	defer fake.errorMutex.Unlock()
	fake.ErrorStub = stub
}

func (fake *FakeCliColorer) ErrorArgsForCall(i int) string {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	argsForCall := fake.errorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCliColorer) ErrorReturns(result1 string) {
	fake.errorMutex.Lock()
	defer fake.errorMutex.Unlock()
	fake.ErrorStub = nil
	fake.errorReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) ErrorReturnsOnCall(i int, result1 string) {
	fake.errorMutex.Lock()
	defer fake.errorMutex.Unlock()
	fake.ErrorStub = nil
	if fake.errorReturnsOnCall == nil {
		fake.errorReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.errorReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) Info(arg1 string) string {
	fake.infoMutex.Lock()
	ret, specificReturn := fake.infoReturnsOnCall[len(fake.infoArgsForCall)]
	fake.infoArgsForCall = append(fake.infoArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Info", []interface{}{arg1})
	fake.infoMutex.Unlock()
	if fake.InfoStub != nil {
		return fake.InfoStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.infoReturns
	return fakeReturns.result1
}

func (fake *FakeCliColorer) InfoCallCount() int {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return len(fake.infoArgsForCall)
}

func (fake *FakeCliColorer) InfoCalls(stub func(string) string) {
	fake.infoMutex.Lock()
	defer fake.infoMutex.Unlock()
	fake.InfoStub = stub
}

func (fake *FakeCliColorer) InfoArgsForCall(i int) string {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	argsForCall := fake.infoArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCliColorer) InfoReturns(result1 string) {
	fake.infoMutex.Lock()
	defer fake.infoMutex.Unlock()
	fake.InfoStub = nil
	fake.infoReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) InfoReturnsOnCall(i int, result1 string) {
	fake.infoMutex.Lock()
	defer fake.infoMutex.Unlock()
	fake.InfoStub = nil
	if fake.infoReturnsOnCall == nil {
		fake.infoReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.infoReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) Muted(arg1 string) string {
	fake.mutedMutex.Lock()
	ret, specificReturn := fake.mutedReturnsOnCall[len(fake.mutedArgsForCall)]
	fake.mutedArgsForCall = append(fake.mutedArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Muted", []interface{}{arg1})
	fake.mutedMutex.Unlock()
	if fake.MutedStub != nil {
		return fake.MutedStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.mutedReturns
	return fakeReturns.result1
}

func (fake *FakeCliColorer) MutedCallCount() int {
	fake.mutedMutex.RLock()
	defer fake.mutedMutex.RUnlock()
	return len(fake.mutedArgsForCall)
}

func (fake *FakeCliColorer) MutedCalls(stub func(string) string) {
	fake.mutedMutex.Lock()
	defer fake.mutedMutex.Unlock()
	fake.MutedStub = stub
}

func (fake *FakeCliColorer) MutedArgsForCall(i int) string {
	fake.mutedMutex.RLock()
	defer fake.mutedMutex.RUnlock()
	argsForCall := fake.mutedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCliColorer) MutedReturns(result1 string) {
	fake.mutedMutex.Lock()
	defer fake.mutedMutex.Unlock()
	fake.MutedStub = nil
	fake.mutedReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) MutedReturnsOnCall(i int, result1 string) {
	fake.mutedMutex.Lock()
	defer fake.mutedMutex.Unlock()
	fake.MutedStub = nil
	if fake.mutedReturnsOnCall == nil {
		fake.mutedReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.mutedReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) Success(arg1 string) string {
	fake.successMutex.Lock()
	ret, specificReturn := fake.successReturnsOnCall[len(fake.successArgsForCall)]
	fake.successArgsForCall = append(fake.successArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Success", []interface{}{arg1})
	fake.successMutex.Unlock()
	if fake.SuccessStub != nil {
		return fake.SuccessStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.successReturns
	return fakeReturns.result1
}

func (fake *FakeCliColorer) SuccessCallCount() int {
	fake.successMutex.RLock()
	defer fake.successMutex.RUnlock()
	return len(fake.successArgsForCall)
}

func (fake *FakeCliColorer) SuccessCalls(stub func(string) string) {
	fake.successMutex.Lock()
	defer fake.successMutex.Unlock()
	fake.SuccessStub = stub
}

func (fake *FakeCliColorer) SuccessArgsForCall(i int) string {
	fake.successMutex.RLock()
	defer fake.successMutex.RUnlock()
	argsForCall := fake.successArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCliColorer) SuccessReturns(result1 string) {
	fake.successMutex.Lock()
	defer fake.successMutex.Unlock()
	fake.SuccessStub = nil
	fake.successReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) SuccessReturnsOnCall(i int, result1 string) {
	fake.successMutex.Lock()
	defer fake.successMutex.Unlock()
	fake.SuccessStub = nil
	if fake.successReturnsOnCall == nil {
		fake.successReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.successReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeCliColorer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.attentionMutex.RLock()
	defer fake.attentionMutex.RUnlock()
	fake.disableColorMutex.RLock()
	defer fake.disableColorMutex.RUnlock()
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	fake.mutedMutex.RLock()
	defer fake.mutedMutex.RUnlock()
	fake.successMutex.RLock()
	defer fake.successMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCliColorer) recordInvocation(key string, args []interface{}) {
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

var _ clicolorer.CliColorer = new(FakeCliColorer)
