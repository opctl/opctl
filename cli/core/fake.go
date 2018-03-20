// Code generated by counterfeiter. DO NOT EDIT.
package core

import (
	"context"
	"sync"
)

type Fake struct {
	EventsStub            func()
	eventsMutex           sync.RWMutex
	eventsArgsForCall     []struct{}
	NodeCreateStub        func()
	nodeCreateMutex       sync.RWMutex
	nodeCreateArgsForCall []struct{}
	NodeKillStub          func()
	nodeKillMutex         sync.RWMutex
	nodeKillArgsForCall   []struct{}
	LsStub                func(ctx context.Context, path string)
	lsMutex               sync.RWMutex
	lsArgsForCall         []struct {
		ctx  context.Context
		path string
	}
	OpCreateStub        func(path string, description string, name string)
	opCreateMutex       sync.RWMutex
	opCreateArgsForCall []struct {
		path        string
		description string
		name        string
	}
	OpInstallStub        func(ctx context.Context, path, opRef, username, password string)
	opInstallMutex       sync.RWMutex
	opInstallArgsForCall []struct {
		ctx      context.Context
		path     string
		opRef    string
		username string
		password string
	}
	OpKillStub        func(ctx context.Context, opId string)
	opKillMutex       sync.RWMutex
	opKillArgsForCall []struct {
		ctx  context.Context
		opId string
	}
	OpValidateStub        func(ctx context.Context, opRef string)
	opValidateMutex       sync.RWMutex
	opValidateArgsForCall []struct {
		ctx   context.Context
		opRef string
	}
	RunStub        func(ctx context.Context, opRef string, opts *RunOpts)
	runMutex       sync.RWMutex
	runArgsForCall []struct {
		ctx   context.Context
		opRef string
		opts  *RunOpts
	}
	SelfUpdateStub        func(channel string)
	selfUpdateMutex       sync.RWMutex
	selfUpdateArgsForCall []struct {
		channel string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Fake) Events() {
	fake.eventsMutex.Lock()
	fake.eventsArgsForCall = append(fake.eventsArgsForCall, struct{}{})
	fake.recordInvocation("Events", []interface{}{})
	fake.eventsMutex.Unlock()
	if fake.EventsStub != nil {
		fake.EventsStub()
	}
}

func (fake *Fake) EventsCallCount() int {
	fake.eventsMutex.RLock()
	defer fake.eventsMutex.RUnlock()
	return len(fake.eventsArgsForCall)
}

func (fake *Fake) NodeCreate() {
	fake.nodeCreateMutex.Lock()
	fake.nodeCreateArgsForCall = append(fake.nodeCreateArgsForCall, struct{}{})
	fake.recordInvocation("NodeCreate", []interface{}{})
	fake.nodeCreateMutex.Unlock()
	if fake.NodeCreateStub != nil {
		fake.NodeCreateStub()
	}
}

func (fake *Fake) NodeCreateCallCount() int {
	fake.nodeCreateMutex.RLock()
	defer fake.nodeCreateMutex.RUnlock()
	return len(fake.nodeCreateArgsForCall)
}

func (fake *Fake) NodeKill() {
	fake.nodeKillMutex.Lock()
	fake.nodeKillArgsForCall = append(fake.nodeKillArgsForCall, struct{}{})
	fake.recordInvocation("NodeKill", []interface{}{})
	fake.nodeKillMutex.Unlock()
	if fake.NodeKillStub != nil {
		fake.NodeKillStub()
	}
}

func (fake *Fake) NodeKillCallCount() int {
	fake.nodeKillMutex.RLock()
	defer fake.nodeKillMutex.RUnlock()
	return len(fake.nodeKillArgsForCall)
}

func (fake *Fake) Ls(ctx context.Context, path string) {
	fake.lsMutex.Lock()
	fake.lsArgsForCall = append(fake.lsArgsForCall, struct {
		ctx  context.Context
		path string
	}{ctx, path})
	fake.recordInvocation("Ls", []interface{}{ctx, path})
	fake.lsMutex.Unlock()
	if fake.LsStub != nil {
		fake.LsStub(ctx, path)
	}
}

func (fake *Fake) LsCallCount() int {
	fake.lsMutex.RLock()
	defer fake.lsMutex.RUnlock()
	return len(fake.lsArgsForCall)
}

func (fake *Fake) LsArgsForCall(i int) (context.Context, string) {
	fake.lsMutex.RLock()
	defer fake.lsMutex.RUnlock()
	return fake.lsArgsForCall[i].ctx, fake.lsArgsForCall[i].path
}

func (fake *Fake) OpCreate(path string, description string, name string) {
	fake.opCreateMutex.Lock()
	fake.opCreateArgsForCall = append(fake.opCreateArgsForCall, struct {
		path        string
		description string
		name        string
	}{path, description, name})
	fake.recordInvocation("OpCreate", []interface{}{path, description, name})
	fake.opCreateMutex.Unlock()
	if fake.OpCreateStub != nil {
		fake.OpCreateStub(path, description, name)
	}
}

func (fake *Fake) OpCreateCallCount() int {
	fake.opCreateMutex.RLock()
	defer fake.opCreateMutex.RUnlock()
	return len(fake.opCreateArgsForCall)
}

func (fake *Fake) OpCreateArgsForCall(i int) (string, string, string) {
	fake.opCreateMutex.RLock()
	defer fake.opCreateMutex.RUnlock()
	return fake.opCreateArgsForCall[i].path, fake.opCreateArgsForCall[i].description, fake.opCreateArgsForCall[i].name
}

func (fake *Fake) OpInstall(ctx context.Context, path string, opRef string, username string, password string) {
	fake.opInstallMutex.Lock()
	fake.opInstallArgsForCall = append(fake.opInstallArgsForCall, struct {
		ctx      context.Context
		path     string
		opRef    string
		username string
		password string
	}{ctx, path, opRef, username, password})
	fake.recordInvocation("OpInstall", []interface{}{ctx, path, opRef, username, password})
	fake.opInstallMutex.Unlock()
	if fake.OpInstallStub != nil {
		fake.OpInstallStub(ctx, path, opRef, username, password)
	}
}

func (fake *Fake) OpInstallCallCount() int {
	fake.opInstallMutex.RLock()
	defer fake.opInstallMutex.RUnlock()
	return len(fake.opInstallArgsForCall)
}

func (fake *Fake) OpInstallArgsForCall(i int) (context.Context, string, string, string, string) {
	fake.opInstallMutex.RLock()
	defer fake.opInstallMutex.RUnlock()
	return fake.opInstallArgsForCall[i].ctx, fake.opInstallArgsForCall[i].path, fake.opInstallArgsForCall[i].opRef, fake.opInstallArgsForCall[i].username, fake.opInstallArgsForCall[i].password
}

func (fake *Fake) OpKill(ctx context.Context, opId string) {
	fake.opKillMutex.Lock()
	fake.opKillArgsForCall = append(fake.opKillArgsForCall, struct {
		ctx  context.Context
		opId string
	}{ctx, opId})
	fake.recordInvocation("OpKill", []interface{}{ctx, opId})
	fake.opKillMutex.Unlock()
	if fake.OpKillStub != nil {
		fake.OpKillStub(ctx, opId)
	}
}

func (fake *Fake) OpKillCallCount() int {
	fake.opKillMutex.RLock()
	defer fake.opKillMutex.RUnlock()
	return len(fake.opKillArgsForCall)
}

func (fake *Fake) OpKillArgsForCall(i int) (context.Context, string) {
	fake.opKillMutex.RLock()
	defer fake.opKillMutex.RUnlock()
	return fake.opKillArgsForCall[i].ctx, fake.opKillArgsForCall[i].opId
}

func (fake *Fake) OpValidate(ctx context.Context, opRef string) {
	fake.opValidateMutex.Lock()
	fake.opValidateArgsForCall = append(fake.opValidateArgsForCall, struct {
		ctx   context.Context
		opRef string
	}{ctx, opRef})
	fake.recordInvocation("OpValidate", []interface{}{ctx, opRef})
	fake.opValidateMutex.Unlock()
	if fake.OpValidateStub != nil {
		fake.OpValidateStub(ctx, opRef)
	}
}

func (fake *Fake) OpValidateCallCount() int {
	fake.opValidateMutex.RLock()
	defer fake.opValidateMutex.RUnlock()
	return len(fake.opValidateArgsForCall)
}

func (fake *Fake) OpValidateArgsForCall(i int) (context.Context, string) {
	fake.opValidateMutex.RLock()
	defer fake.opValidateMutex.RUnlock()
	return fake.opValidateArgsForCall[i].ctx, fake.opValidateArgsForCall[i].opRef
}

func (fake *Fake) Run(ctx context.Context, opRef string, opts *RunOpts) {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct {
		ctx   context.Context
		opRef string
		opts  *RunOpts
	}{ctx, opRef, opts})
	fake.recordInvocation("Run", []interface{}{ctx, opRef, opts})
	fake.runMutex.Unlock()
	if fake.RunStub != nil {
		fake.RunStub(ctx, opRef, opts)
	}
}

func (fake *Fake) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *Fake) RunArgsForCall(i int) (context.Context, string, *RunOpts) {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return fake.runArgsForCall[i].ctx, fake.runArgsForCall[i].opRef, fake.runArgsForCall[i].opts
}

func (fake *Fake) SelfUpdate(channel string) {
	fake.selfUpdateMutex.Lock()
	fake.selfUpdateArgsForCall = append(fake.selfUpdateArgsForCall, struct {
		channel string
	}{channel})
	fake.recordInvocation("SelfUpdate", []interface{}{channel})
	fake.selfUpdateMutex.Unlock()
	if fake.SelfUpdateStub != nil {
		fake.SelfUpdateStub(channel)
	}
}

func (fake *Fake) SelfUpdateCallCount() int {
	fake.selfUpdateMutex.RLock()
	defer fake.selfUpdateMutex.RUnlock()
	return len(fake.selfUpdateArgsForCall)
}

func (fake *Fake) SelfUpdateArgsForCall(i int) string {
	fake.selfUpdateMutex.RLock()
	defer fake.selfUpdateMutex.RUnlock()
	return fake.selfUpdateArgsForCall[i].channel
}

func (fake *Fake) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.eventsMutex.RLock()
	defer fake.eventsMutex.RUnlock()
	fake.nodeCreateMutex.RLock()
	defer fake.nodeCreateMutex.RUnlock()
	fake.nodeKillMutex.RLock()
	defer fake.nodeKillMutex.RUnlock()
	fake.lsMutex.RLock()
	defer fake.lsMutex.RUnlock()
	fake.opCreateMutex.RLock()
	defer fake.opCreateMutex.RUnlock()
	fake.opInstallMutex.RLock()
	defer fake.opInstallMutex.RUnlock()
	fake.opKillMutex.RLock()
	defer fake.opKillMutex.RUnlock()
	fake.opValidateMutex.RLock()
	defer fake.opValidateMutex.RUnlock()
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	fake.selfUpdateMutex.RLock()
	defer fake.selfUpdateMutex.RUnlock()
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

var _ Core = new(Fake)
