package node

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

var _ = Describe("Go output collection", func() {
	var dir string
	BeforeEach(func() {
		var err error
		dir, err = os.MkdirTemp("", "go-outputs-")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(filepath.Join(dir, "op.go"), []byte("package main"), 0600)).To(Succeed())
	})
	AfterEach(func() { os.RemoveAll(dir) })
	It("resolves bindings against merged inputs and child results, applies defaults, aliases, and root bindings", func() {
		for _, kind := range []string{"container", "op", "serial", "parallel"} {
			bootstrap := newGoBootstrap()
			child := new(fakes.FakeCaller)
			desc := &model.OpSpec{Inputs: map[string]*model.ParamSpec{"text": {String: &model.StringParamSpec{}}}, Outputs: map[string]*model.ParamSpec{"echoed": {String: &model.StringParamSpec{Constraints: model.Constraints{"minLength": 3}}}, "alias": {String: &model.StringParamSpec{}}, "artifact": {Dir: &model.DirParamSpec{}}, "fallback": {String: &model.StringParamSpec{Default: "default"}}}}
			bootstrap.DescribeReturns(desc, nil)
			run := &model.CallSpec{Container: &model.ContainerCallSpec{}}
			switch kind {
			case "op":
				run = &model.CallSpec{Op: &model.OpCallSpec{Ref: "./yaml"}}
			case "serial":
				children := []*model.CallSpec{{Container: &model.ContainerCallSpec{}}}
				run = &model.CallSpec{Serial: &children}
			case "parallel":
				children := []*model.CallSpec{{Container: &model.ContainerCallSpec{}}}
				run = &model.CallSpec{Parallel: &children}
			}
			bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: run}, OutputBindings: map[string]string{"echoed": "$(text)", "alias": "$(text)", "artifact": "$(opctl_dir_1)"}}, nil)
			artifact := filepath.Join(dir, "artifact")
			Expect(os.MkdirAll(artifact, 0700)).To(Succeed())
			child.CallReturns(map[string]*model.Value{"opctl_dir_1": {Dir: &artifact}}, nil)
			text := "from-input"
			spec := &model.OpCallSpec{Inputs: map[string]interface{}{"text": "$(incoming)"}}
			returned, err := (_opCaller{caller: child, bootstrapper: bootstrap, dataDirPath: dir}).Call(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}, OpID: "root", ChildCallID: "child"}, nil, "root", spec, map[string]*model.Value{"incoming": {String: &text}})
			Expect(err).NotTo(HaveOccurred(), kind)
			Expect(returned).To(HaveLen(4))
			Expect(*returned["echoed"].String).To(Equal(text))
			Expect(*returned["alias"].String).To(Equal(text))
			Expect(*returned["artifact"].Dir).To(Equal(artifact))
			Expect(*returned["fallback"].String).To(Equal("default"))
			Expect(spec.Outputs).To(BeNil())
		}
	})
	It("rejects incomplete bindings before running and missing references after running", func() {
		for _, bindings := range []map[string]string{{}, {"unknown": "$(text)"}, {"out": "$(missing)"}} {
			bootstrap := newGoBootstrap()
			child := new(fakes.FakeCaller)
			bootstrap.DescribeReturns(&model.OpSpec{Outputs: map[string]*model.ParamSpec{"out": {Dir: &model.DirParamSpec{}}}}, nil)
			bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{}}}, OutputBindings: bindings}, nil)
			_, err := (_opCaller{caller: child, bootstrapper: bootstrap, dataDirPath: dir}).Call(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}, OpID: "root"}, nil, "root", &model.OpCallSpec{}, nil)
			Expect(err).To(HaveOccurred())
			if bindings["out"] == "" {
				Expect(child.CallCallCount()).To(BeZero())
			}
		}
	})
	It("sanitizes only the failing secret default", func() {
		for _, constraint := range []bool{false, true} {
			bootstrap := newGoBootstrap()
			param := &model.ParamSpec{Number: &model.NumberParamSpec{IsSecret: true, Default: "private-invalid-default"}}
			if constraint {
				param = &model.ParamSpec{String: &model.StringParamSpec{IsSecret: true, Default: "private-invalid-default", Constraints: model.Constraints{"enum": []interface{}{"allowed"}}}}
			}
			bootstrap.DescribeReturns(&model.OpSpec{Outputs: map[string]*model.ParamSpec{"private": param}}, nil)
			bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{}}}}, nil)
			_, err := (_opCaller{caller: new(fakes.FakeCaller), bootstrapper: bootstrap, dataDirPath: dir}).Call(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}, OpID: "root"}, nil, "root", &model.OpCallSpec{}, nil)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "private-invalid-default")).To(BeFalse())
			Expect(err.Error()).To(ContainSubstring("private"))
		}
	})

	It("preserves public output errors when another output is secret", func() {
		for _, unknown := range []bool{false, true} {
			bootstrap := newGoBootstrap()
			child := new(fakes.FakeCaller)
			bootstrap.DescribeReturns(&model.OpSpec{Outputs: map[string]*model.ParamSpec{
				"private": {String: &model.StringParamSpec{IsSecret: true}},
				"count":   {Number: &model.NumberParamSpec{}},
			}}, nil)
			bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{}}}, OutputBindings: map[string]string{"private": "$(secret)", "count": "$(count)"}}, nil)
			secret, bad := "private-output", "not-a-number"
			number := 3.0
			count := &model.Value{String: &bad}
			requested := map[string]string{"count": "$(count)"}
			if unknown {
				count = &model.Value{Number: &number}
				requested["typo"] = "$(typo)"
			}
			child.CallReturns(map[string]*model.Value{"secret": {String: &secret}, "count": count}, nil)
			parent := "parent"
			_, err := (_opCaller{caller: child, bootstrapper: bootstrap, dataDirPath: dir}).Call(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}, OpID: "op"}, &parent, "root", &model.OpCallSpec{Outputs: requested}, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).NotTo(ContainSubstring("sensitive output"))
			if unknown {
				Expect(err.Error()).To(ContainSubstring("unknown output 'typo'"))
			} else {
				Expect(err.Error()).To(ContainSubstring("count"))
			}
		}
	})

	It("exports only requested nested outputs and rejects sensitive values bound to public Go inputs", func() {
		bootstrap := newGoBootstrap()
		child := new(fakes.FakeCaller)
		bootstrap.DescribeReturns(&model.OpSpec{Inputs: map[string]*model.ParamSpec{"text": {String: &model.StringParamSpec{}}}, Outputs: map[string]*model.ParamSpec{"out": {String: &model.StringParamSpec{}}}}, nil)
		bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{}}}, OutputBindings: map[string]string{"out": "$(text)"}}, nil)
		text := "private-runtime"
		parent := "parent"
		op := &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}, OpID: "op"}
		caller := _opCaller{caller: child, bootstrapper: bootstrap, dataDirPath: dir}
		args := &model.OpCallSpec{Inputs: map[string]interface{}{"text": "$(source)"}, Outputs: map[string]string{"out": "$(renamed)"}}
		returned, err := caller.Call(context.Background(), op, &parent, "root", args, map[string]*model.Value{"source": {String: &text}})
		Expect(err).NotTo(HaveOccurred())
		Expect(returned).To(HaveLen(1))
		Expect(*returned["renamed"].String).To(Equal(text))
		_, err = caller.Call(context.Background(), op, &parent, "root", args, map[string]*model.Value{"source": {String: &text, Sensitive: true}})
		Expect(err).To(MatchError(ContainSubstring("secret")))
		Expect(bootstrap.PlanCallCount()).To(Equal(1))
	})
})
