package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("caller", func() {
	Context("newCaller", func() {
		It("should return caller", func() {
			/* arrange/act/assert */
			Expect(newCaller(new(fakeContainerCaller))).To(Not(BeNil()))
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
					"dummyCallId",
					map[string]*model.Value{},
					nil,
					new(data.FakeHandle),
					"dummyRootOpId",
				)
			})
		})
		Context("Container SCG", func() {
			It("should call containerCaller.Call w/ expected args", func() {
				/* arrange */
				fakeContainerCaller := new(fakeContainerCaller)

				providedCallId := "dummyCallId"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Container: &model.SCGContainerCall{},
				}
				providedOpDirHandle := new(data.FakeHandle)
				providedRootOpId := "dummyRootOpId"

				objectUnderTest := newCaller(fakeContainerCaller)

				/* act */
				objectUnderTest.Call(
					providedCallId,
					providedArgs,
					providedSCG,
					providedOpDirHandle,
					providedRootOpId,
				)

				/* assert */
				actualArgs,
					actualCallId,
					actualSCG,
					actualOpDirHandle,
					actualRootOpId := fakeContainerCaller.CallArgsForCall(0)
				Expect(actualCallId).To(Equal(providedCallId))
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualSCG).To(Equal(providedSCG.Container))
				Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			})
		})
		Context("Op SCG", func() {
			It("should call opCaller.Call w/ expected args", func() {
				/* arrange */
				fakeOpCaller := new(fakeOpCaller)

				providedCallId := "dummyCallId"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Op: &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: "dummyPkgRef",
						},
					},
				}
				providedOpDirHandle := new(data.FakeHandle)
				providedRootOpId := "dummyRootOpId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setOpCaller(fakeOpCaller)

				/* act */
				objectUnderTest.Call(
					providedCallId,
					providedArgs,
					providedSCG,
					providedOpDirHandle,
					providedRootOpId,
				)

				/* assert */
				actualArgs,
					actualCallId,
					actualPkgRef,
					actualRootOpId,
					actualSCG := fakeOpCaller.CallArgsForCall(0)
				Expect(actualCallId).To(Equal(providedCallId))
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualPkgRef).To(Equal(providedOpDirHandle))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
				Expect(actualSCG).To(Equal(providedSCG.Op))
			})
		})
		Context("Parallel SCG", func() {
			It("should call parallelCaller.Call w/ expected args", func() {
				/* arrange */
				fakeParallelCaller := new(fakeParallelCaller)

				providedCallId := "dummyCallId"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Parallel: []*model.SCG{
						{Container: &model.SCGContainerCall{}},
					},
				}
				providedOpDirHandle := new(data.FakeHandle)
				providedRootOpId := "dummyRootOpId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setParallelCaller(fakeParallelCaller)

				/* act */
				objectUnderTest.Call(
					providedCallId,
					providedArgs,
					providedSCG,
					providedOpDirHandle,
					providedRootOpId,
				)

				/* assert */
				providedCallId,
					actualArgs,
					actualRootOpId,
					actualOpDirHandle,
					actualSCG := fakeParallelCaller.CallArgsForCall(0)
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
				Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
				Expect(actualSCG).To(Equal(providedSCG.Parallel))
			})
		})
		Context("Serial SCG", func() {

			It("should call serialCaller.Call w/ expected args", func() {
				/* arrange */
				fakeSerialCaller := new(fakeSerialCaller)

				providedCallId := "dummyCallId"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Serial: []*model.SCG{
						{Container: &model.SCGContainerCall{}},
					},
				}
				providedOpDirHandle := new(data.FakeHandle)
				providedRootOpId := "dummyRootOpId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setSerialCaller(fakeSerialCaller)

				/* act */
				objectUnderTest.Call(
					providedCallId,
					providedArgs,
					providedSCG,
					providedOpDirHandle,
					providedRootOpId,
				)

				/* assert */
				actualCallId,
					actualArgs,
					actualRootOpId,
					actualOpDirHandle,
					actualSCG := fakeSerialCaller.CallArgsForCall(0)
				Expect(actualCallId).To(Equal(providedCallId))
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
				Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
				Expect(actualSCG).To(Equal(providedSCG.Serial))
			})
		})
		Context("No SCG", func() {
			It("should error", func() {
				/* arrange */
				fakeSerialCaller := new(fakeSerialCaller)

				providedCallId := "dummyCallId"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{}
				providedOpDirHandle := new(data.FakeHandle)
				providedRootOpId := "dummyRootOpId"
				expectedError := fmt.Errorf("Invalid call graph %+v\n", providedSCG)

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setSerialCaller(fakeSerialCaller)

				/* act */
				actualError := objectUnderTest.Call(
					providedCallId,
					providedArgs,
					providedSCG,
					providedOpDirHandle,
					providedRootOpId,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
