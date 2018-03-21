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
					"dummyCallID",
					map[string]*model.Value{},
					nil,
					new(data.FakeHandle),
					"dummyRootOpID",
				)
			})
		})
		Context("Container SCG", func() {
			It("should call containerCaller.Call w/ expected args", func() {
				/* arrange */
				fakeContainerCaller := new(fakeContainerCaller)

				providedCallID := "dummyCallID"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Container: &model.SCGContainerCall{},
				}
				providedOpHandle := new(data.FakeHandle)
				providedRootOpID := "dummyRootOpID"

				objectUnderTest := newCaller(fakeContainerCaller)

				/* act */
				objectUnderTest.Call(
					providedCallID,
					providedArgs,
					providedSCG,
					providedOpHandle,
					providedRootOpID,
				)

				/* assert */
				actualArgs,
					actualCallID,
					actualSCG,
					actualOpHandle,
					actualRootOpID := fakeContainerCaller.CallArgsForCall(0)
				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualSCG).To(Equal(providedSCG.Container))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})
		Context("Op SCG", func() {
			It("should call opCaller.Call w/ expected args", func() {
				/* arrange */
				fakeOpCaller := new(fakeOpCaller)

				providedCallID := "dummyCallID"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Op: &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: "dummyOpRef",
						},
					},
				}
				providedOpHandle := new(data.FakeHandle)
				providedRootOpID := "dummyRootOpID"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setOpCaller(fakeOpCaller)

				/* act */
				objectUnderTest.Call(
					providedCallID,
					providedArgs,
					providedSCG,
					providedOpHandle,
					providedRootOpID,
				)

				/* assert */
				actualArgs,
					actualCallID,
					actualPkgRef,
					actualRootOpID,
					actualSCG := fakeOpCaller.CallArgsForCall(0)
				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualPkgRef).To(Equal(providedOpHandle))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualSCG).To(Equal(providedSCG.Op))
			})
		})
		Context("Parallel SCG", func() {
			It("should call parallelCaller.Call w/ expected args", func() {
				/* arrange */
				fakeParallelCaller := new(fakeParallelCaller)

				providedCallID := "dummyCallID"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Parallel: []*model.SCG{
						{Container: &model.SCGContainerCall{}},
					},
				}
				providedOpHandle := new(data.FakeHandle)
				providedRootOpID := "dummyRootOpID"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setParallelCaller(fakeParallelCaller)

				/* act */
				objectUnderTest.Call(
					providedCallID,
					providedArgs,
					providedSCG,
					providedOpHandle,
					providedRootOpID,
				)

				/* assert */
				providedCallID,
					actualArgs,
					actualRootOpID,
					actualOpHandle,
					actualSCG := fakeParallelCaller.CallArgsForCall(0)
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualSCG).To(Equal(providedSCG.Parallel))
			})
		})
		Context("Serial SCG", func() {

			It("should call serialCaller.Call w/ expected args", func() {
				/* arrange */
				fakeSerialCaller := new(fakeSerialCaller)

				providedCallID := "dummyCallID"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Serial: []*model.SCG{
						{Container: &model.SCGContainerCall{}},
					},
				}
				providedOpHandle := new(data.FakeHandle)
				providedRootOpID := "dummyRootOpID"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setSerialCaller(fakeSerialCaller)

				/* act */
				objectUnderTest.Call(
					providedCallID,
					providedArgs,
					providedSCG,
					providedOpHandle,
					providedRootOpID,
				)

				/* assert */
				actualCallID,
					actualArgs,
					actualRootOpID,
					actualOpHandle,
					actualSCG := fakeSerialCaller.CallArgsForCall(0)
				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualSCG).To(Equal(providedSCG.Serial))
			})
		})
		Context("No SCG", func() {
			It("should error", func() {
				/* arrange */
				fakeSerialCaller := new(fakeSerialCaller)

				providedCallID := "dummyCallID"
				providedArgs := map[string]*model.Value{}
				providedSCG := &model.SCG{}
				providedOpHandle := new(data.FakeHandle)
				providedRootOpID := "dummyRootOpID"
				expectedError := fmt.Errorf("Invalid call graph %+v\n", providedSCG)

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setSerialCaller(fakeSerialCaller)

				/* act */
				actualError := objectUnderTest.Call(
					providedCallID,
					providedArgs,
					providedSCG,
					providedOpHandle,
					providedRootOpID,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
