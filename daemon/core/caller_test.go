package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
	"path/filepath"
)

var _ = Describe("caller", func() {
	Context("newCaller", func() {
		It("should return caller", func() {
			/* arrange/act/assert */
			Expect(newCaller(new(fakeContainerCaller))).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		Context("Container SCG", func() {
			It("should call containerCaller.Call w/ expected args", func() {
				/* arrange */
				fakeContainerCaller := new(fakeContainerCaller)

				providedNodeId := "dummyNodeId"
				providedArgs := map[string]*model.Data{}
				providedScg := &model.Scg{
					Container: &model.ScgContainerCall{},
				}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				objectUnderTest := newCaller(fakeContainerCaller)

				/* act */
				objectUnderTest.Call(
					providedNodeId,
					providedArgs,
					providedScg,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				actualArgs,
					actualNodeId,
					actualScg,
					actualOpRef,
					actualOpGraphId := fakeContainerCaller.CallArgsForCall(0)
				Expect(actualNodeId).To(Equal(providedNodeId))
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualScg).To(Equal(providedScg.Container))
				Expect(actualOpRef).To(Equal(providedOpRef))
				Expect(actualOpGraphId).To(Equal(providedOpGraphId))
			})
			Context("containerCaller.Call errors", func() {

			})
			Context("containerCaller.Call doesn't error", func() {

			})
		})
		Context("Op SCG", func() {
			It("should call opCaller.Call w/ expected args", func() {
				/* arrange */
				fakeOpCaller := new(fakeOpCaller)

				providedNodeId := "dummyNodeId"
				providedArgs := map[string]*model.Data{}
				providedScg := &model.Scg{
					Op: &model.ScgOpCall{
						Ref: "dummyOpRef",
					},
				}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setOpCaller(fakeOpCaller)

				/* act */
				objectUnderTest.Call(
					providedNodeId,
					providedArgs,
					providedScg,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				actualArgs,
					actualNodeId,
					actualOpRef,
					actualOpGraphId := fakeOpCaller.CallArgsForCall(0)
				Expect(actualNodeId).To(Equal(providedNodeId))
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualOpRef).To(Equal(path.Join(filepath.Dir(providedOpRef), providedScg.Op.Ref)))
				Expect(actualOpGraphId).To(Equal(providedOpGraphId))
			})
			Context("opCaller.Call errors", func() {

			})
			Context("opCaller.Call doesn't error", func() {

			})
		})
		Context("Parallel SCG", func() {
			It("should call parallelCaller.Call w/ expected args", func() {
				/* arrange */
				fakeParallelCaller := new(fakeParallelCaller)

				providedNodeId := "dummyNodeId"
				providedArgs := map[string]*model.Data{}
				providedScg := &model.Scg{
					Parallel: []*model.Scg{
						{Container: &model.ScgContainerCall{}},
					},
				}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setParallelCaller(fakeParallelCaller)

				/* act */
				objectUnderTest.Call(
					providedNodeId,
					providedArgs,
					providedScg,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				actualArgs,
					actualOpGraphId,
					actualOpRef,
					actualScg := fakeParallelCaller.CallArgsForCall(0)
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualOpGraphId).To(Equal(providedOpGraphId))
				Expect(actualOpRef).To(Equal(providedOpRef))
				Expect(actualScg).To(Equal(providedScg.Parallel))
			})
			Context("parallelCaller.Call errors", func() {

			})
			Context("parallelCaller.Call doesn't error", func() {

			})
		})
		Context("Serial SCG", func() {

			It("should call serialCaller.Call w/ expected args", func() {
				/* arrange */
				fakeSerialCaller := new(fakeSerialCaller)

				providedNodeId := "dummyNodeId"
				providedArgs := map[string]*model.Data{}
				providedScg := &model.Scg{
					Serial: []*model.Scg{
						{Container: &model.ScgContainerCall{}},
					},
				}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setSerialCaller(fakeSerialCaller)

				/* act */
				objectUnderTest.Call(
					providedNodeId,
					providedArgs,
					providedScg,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				actualArgs,
					actualOpGraphId,
					actualOpRef,
					actualScg := fakeSerialCaller.CallArgsForCall(0)
				Expect(actualArgs).To(Equal(providedArgs))
				Expect(actualOpGraphId).To(Equal(providedOpGraphId))
				Expect(actualOpRef).To(Equal(providedOpRef))
				Expect(actualScg).To(Equal(providedScg.Serial))
			})
			Context("serialCaller.Call errors", func() {

			})
			Context("serialCaller.Call doesn't error", func() {

			})
		})
		Context("No SCG", func() {
			It("should error", func() {
				/* arrange */
				fakeSerialCaller := new(fakeSerialCaller)

				providedNodeId := "dummyNodeId"
				providedArgs := map[string]*model.Data{}
				providedScg := &model.Scg{}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"
				expectedError := fmt.Errorf("Invalid call graph %+v\n", providedScg)

				objectUnderTest := newCaller(new(fakeContainerCaller))
				objectUnderTest.setSerialCaller(fakeSerialCaller)

				/* act */
				_, actualError := objectUnderTest.Call(
					providedNodeId,
					providedArgs,
					providedScg,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
