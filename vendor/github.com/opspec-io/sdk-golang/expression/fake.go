// Code generated by counterfeiter. DO NOT EDIT.
package expression

import (
	"github.com/opspec-io/sdk-golang/model"
	"sync"
)

type Fake struct {
	EvalToFileStub        func(scope map[string]*model.Value, expression string, pkgHandle model.PkgHandle, scratchDir string) (*model.Value, error)
	evalToFileMutex       sync.RWMutex
	evalToFileArgsForCall []struct {
		scope      map[string]*model.Value
		expression string
		pkgHandle  model.PkgHandle
		scratchDir string
	}
	evalToFileReturns struct {
		result1 *model.Value
		result2 error
	}
	evalToFileReturnsOnCall map[int]struct {
		result1 *model.Value
		result2 error
	}
	EvalToNumberStub        func(scope map[string]*model.Value, expression string, pkgHandle model.PkgHandle) (float64, error)
	evalToNumberMutex       sync.RWMutex
	evalToNumberArgsForCall []struct {
		scope      map[string]*model.Value
		expression string
		pkgHandle  model.PkgHandle
	}
	evalToNumberReturns struct {
		result1 float64
		result2 error
	}
	evalToNumberReturnsOnCall map[int]struct {
		result1 float64
		result2 error
	}
	EvalToStringStub        func(scope map[string]*model.Value, expression string, pkgHandle model.PkgHandle) (string, error)
	evalToStringMutex       sync.RWMutex
	evalToStringArgsForCall []struct {
		scope      map[string]*model.Value
		expression string
		pkgHandle  model.PkgHandle
	}
	evalToStringReturns struct {
		result1 string
		result2 error
	}
	evalToStringReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Fake) EvalToFile(scope map[string]*model.Value, expression string, pkgHandle model.PkgHandle, scratchDir string) (*model.Value, error) {
	fake.evalToFileMutex.Lock()
	ret, specificReturn := fake.evalToFileReturnsOnCall[len(fake.evalToFileArgsForCall)]
	fake.evalToFileArgsForCall = append(fake.evalToFileArgsForCall, struct {
		scope      map[string]*model.Value
		expression string
		pkgHandle  model.PkgHandle
		scratchDir string
	}{scope, expression, pkgHandle, scratchDir})
	fake.recordInvocation("EvalToFile", []interface{}{scope, expression, pkgHandle, scratchDir})
	fake.evalToFileMutex.Unlock()
	if fake.EvalToFileStub != nil {
		return fake.EvalToFileStub(scope, expression, pkgHandle, scratchDir)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.evalToFileReturns.result1, fake.evalToFileReturns.result2
}

func (fake *Fake) EvalToFileCallCount() int {
	fake.evalToFileMutex.RLock()
	defer fake.evalToFileMutex.RUnlock()
	return len(fake.evalToFileArgsForCall)
}

func (fake *Fake) EvalToFileArgsForCall(i int) (map[string]*model.Value, string, model.PkgHandle, string) {
	fake.evalToFileMutex.RLock()
	defer fake.evalToFileMutex.RUnlock()
	return fake.evalToFileArgsForCall[i].scope, fake.evalToFileArgsForCall[i].expression, fake.evalToFileArgsForCall[i].pkgHandle, fake.evalToFileArgsForCall[i].scratchDir
}

func (fake *Fake) EvalToFileReturns(result1 *model.Value, result2 error) {
	fake.EvalToFileStub = nil
	fake.evalToFileReturns = struct {
		result1 *model.Value
		result2 error
	}{result1, result2}
}

func (fake *Fake) EvalToFileReturnsOnCall(i int, result1 *model.Value, result2 error) {
	fake.EvalToFileStub = nil
	if fake.evalToFileReturnsOnCall == nil {
		fake.evalToFileReturnsOnCall = make(map[int]struct {
			result1 *model.Value
			result2 error
		})
	}
	fake.evalToFileReturnsOnCall[i] = struct {
		result1 *model.Value
		result2 error
	}{result1, result2}
}

func (fake *Fake) EvalToNumber(scope map[string]*model.Value, expression string, pkgHandle model.PkgHandle) (float64, error) {
	fake.evalToNumberMutex.Lock()
	ret, specificReturn := fake.evalToNumberReturnsOnCall[len(fake.evalToNumberArgsForCall)]
	fake.evalToNumberArgsForCall = append(fake.evalToNumberArgsForCall, struct {
		scope      map[string]*model.Value
		expression string
		pkgHandle  model.PkgHandle
	}{scope, expression, pkgHandle})
	fake.recordInvocation("EvalToNumber", []interface{}{scope, expression, pkgHandle})
	fake.evalToNumberMutex.Unlock()
	if fake.EvalToNumberStub != nil {
		return fake.EvalToNumberStub(scope, expression, pkgHandle)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.evalToNumberReturns.result1, fake.evalToNumberReturns.result2
}

func (fake *Fake) EvalToNumberCallCount() int {
	fake.evalToNumberMutex.RLock()
	defer fake.evalToNumberMutex.RUnlock()
	return len(fake.evalToNumberArgsForCall)
}

func (fake *Fake) EvalToNumberArgsForCall(i int) (map[string]*model.Value, string, model.PkgHandle) {
	fake.evalToNumberMutex.RLock()
	defer fake.evalToNumberMutex.RUnlock()
	return fake.evalToNumberArgsForCall[i].scope, fake.evalToNumberArgsForCall[i].expression, fake.evalToNumberArgsForCall[i].pkgHandle
}

func (fake *Fake) EvalToNumberReturns(result1 float64, result2 error) {
	fake.EvalToNumberStub = nil
	fake.evalToNumberReturns = struct {
		result1 float64
		result2 error
	}{result1, result2}
}

func (fake *Fake) EvalToNumberReturnsOnCall(i int, result1 float64, result2 error) {
	fake.EvalToNumberStub = nil
	if fake.evalToNumberReturnsOnCall == nil {
		fake.evalToNumberReturnsOnCall = make(map[int]struct {
			result1 float64
			result2 error
		})
	}
	fake.evalToNumberReturnsOnCall[i] = struct {
		result1 float64
		result2 error
	}{result1, result2}
}

func (fake *Fake) EvalToString(scope map[string]*model.Value, expression string, pkgHandle model.PkgHandle) (string, error) {
	fake.evalToStringMutex.Lock()
	ret, specificReturn := fake.evalToStringReturnsOnCall[len(fake.evalToStringArgsForCall)]
	fake.evalToStringArgsForCall = append(fake.evalToStringArgsForCall, struct {
		scope      map[string]*model.Value
		expression string
		pkgHandle  model.PkgHandle
	}{scope, expression, pkgHandle})
	fake.recordInvocation("EvalToString", []interface{}{scope, expression, pkgHandle})
	fake.evalToStringMutex.Unlock()
	if fake.EvalToStringStub != nil {
		return fake.EvalToStringStub(scope, expression, pkgHandle)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.evalToStringReturns.result1, fake.evalToStringReturns.result2
}

func (fake *Fake) EvalToStringCallCount() int {
	fake.evalToStringMutex.RLock()
	defer fake.evalToStringMutex.RUnlock()
	return len(fake.evalToStringArgsForCall)
}

func (fake *Fake) EvalToStringArgsForCall(i int) (map[string]*model.Value, string, model.PkgHandle) {
	fake.evalToStringMutex.RLock()
	defer fake.evalToStringMutex.RUnlock()
	return fake.evalToStringArgsForCall[i].scope, fake.evalToStringArgsForCall[i].expression, fake.evalToStringArgsForCall[i].pkgHandle
}

func (fake *Fake) EvalToStringReturns(result1 string, result2 error) {
	fake.EvalToStringStub = nil
	fake.evalToStringReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *Fake) EvalToStringReturnsOnCall(i int, result1 string, result2 error) {
	fake.EvalToStringStub = nil
	if fake.evalToStringReturnsOnCall == nil {
		fake.evalToStringReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.evalToStringReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *Fake) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.evalToFileMutex.RLock()
	defer fake.evalToFileMutex.RUnlock()
	fake.evalToNumberMutex.RLock()
	defer fake.evalToNumberMutex.RUnlock()
	fake.evalToStringMutex.RLock()
	defer fake.evalToStringMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Fake) recordInvocation(key string, args []interface{}) {
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

var _ Expression = new(Fake)
