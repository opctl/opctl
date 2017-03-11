package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
	"path/filepath"
)

var _ = Context("caller", func() {
	Context("newCaller", func() {
		It("should return caller", func() {
			/* arrange/act/assert */
			Expect(newCaller(new(fakeContainerCaller))).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		Context("Null SCG", func() {
			It("should not throw", func() {
				/* arrange */
				fakeContainerCaller := new(fakeContainerCaller)

				/* act */
				objectUnderTest := newCaller(fakeContainerCaller)

				/* assert */
				objectUnderTest.Call(
					"dummyNodeId",
					make(chan *variable, 150),
					make(chan *variable, 150),
					nil,
					"dummyPkgRef",
					"dummyRootOpId",
				)
			})
		})
		Context("Container SCG", func() {
			It("should call containerCaller.Call w/ expected args", func() {
				/* arrange */
				fakeContainerCaller := new(fakeContainerCaller)

				providedNodeId := "dummyNodeId"
				providedInputs := make(chan *variable, 150)
				providedOutputs := make(chan *variable, 150)
				providedScg := &model.Scg{
					Container: &model.ScgContainerCall{},
				}
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"

				objectUnderTest := newCaller(fakeContainerCaller)

				/* act */
				objectUnderTest.Call(
					providedNodeId,
					providedInputs,
					providedOutputs,
					providedScg,
					providedPkgRef,
					providedRootOpId,
				)

				/* assert */
				actualArgs,
					_,
					actualNodeId,
					actualScg,
					actualPkgRef,
					actualRootOpId := fakeContainerCaller.CallArgsForCall(0)
				Expect(actualNodeId).To(Equal(providedNodeId))
				Expect(actualArgs).To(Equal(providedInputs))
				Expect(actualScg).To(Equal(providedScg.Container))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			})
		})
		Context("Op SCG", func() {
			It("should call opCaller.Call w/ expected args", func() {
				/* arrange */
				fakeOpCaller := new(fakeOpCaller)

				providedNodeId := "dummyNodeId"
				providedInputs := make(chan *variable, 150)
				providedOutputs := make(chan *variable, 150)
				providedScg := &model.Scg{
					Op: &model.ScgOpCall{
						Ref: "dummyPkgRef",
					},
				}
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setOpCaller(fakeOpCaller)

				/* act */
				objectUnderTest.Call(
					providedNodeId,
					providedInputs,
					providedOutputs,
					providedScg,
					providedPkgRef,
					providedRootOpId,
				)

				/* assert */
				actualArgs,
					_,
					actualNodeId,
					actualPkgRef,
					actualRootOpId := fakeOpCaller.CallArgsForCall(0)
				Expect(actualNodeId).To(Equal(providedNodeId))
				Expect(actualArgs).To(Equal(providedInputs))
				Expect(actualPkgRef).To(Equal(path.Join(filepath.Dir(providedPkgRef), providedScg.Op.Ref)))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			})
		})
		Context("Parallel SCG", func() {
			It("should call parallelCaller.Call w/ expected args", func() {
				/* arrange */
				fakeParallelCaller := new(fakeParallelCaller)

				providedNodeId := "dummyNodeId"
				providedInputs := make(chan *variable, 150)
				providedOutputs := make(chan *variable, 150)
				providedScg := &model.Scg{
					Parallel: []*model.Scg{
						{Container: &model.ScgContainerCall{}},
					},
				}
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setParallelCaller(fakeParallelCaller)

				/* act */
				objectUnderTest.Call(
					providedNodeId,
					providedInputs,
					providedOutputs,
					providedScg,
					providedPkgRef,
					providedRootOpId,
				)

				/* assert */
				actualArgs,
					actualRootOpId,
					actualPkgRef,
					actualScg := fakeParallelCaller.CallArgsForCall(0)
				Expect(actualArgs).To(Equal(providedInputs))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualScg).To(Equal(providedScg.Parallel))
			})
		})
		Context("Serial SCG", func() {

			It("should call serialCaller.Call w/ expected args", func() {
				/* arrange */
				fakeSerialCaller := new(fakeSerialCaller)

				providedNodeId := "dummyNodeId"
				providedInputs := make(chan *variable, 150)
				providedOutputs := make(chan *variable, 150)
				providedScg := &model.Scg{
					Serial: []*model.Scg{
						{Container: &model.ScgContainerCall{}},
					},
				}
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setSerialCaller(fakeSerialCaller)

				/* act */
				objectUnderTest.Call(
					providedNodeId,
					providedInputs,
					providedOutputs,
					providedScg,
					providedPkgRef,
					providedRootOpId,
				)

				/* assert */
				actualArgs,
					_,
					actualRootOpId,
					actualPkgRef,
					actualScg := fakeSerialCaller.CallArgsForCall(0)
				Expect(actualArgs).To(Equal(providedInputs))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualScg).To(Equal(providedScg.Serial))
			})
		})
		Context("No SCG", func() {
			It("should error", func() {
				/* arrange */
				fakeSerialCaller := new(fakeSerialCaller)

				providedNodeId := "dummyNodeId"
				providedInputs := make(chan *variable, 150)
				providedOutputs := make(chan *variable, 150)
				providedScg := &model.Scg{}
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"
				expectedError := fmt.Errorf("Invalid call graph %+v\n", providedScg)

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setSerialCaller(fakeSerialCaller)

				/* act */
				actualError := objectUnderTest.Call(
					providedNodeId,
					providedInputs,
					providedOutputs,
					providedScg,
					providedPkgRef,
					providedRootOpId,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
