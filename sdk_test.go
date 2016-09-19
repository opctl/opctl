package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
  "github.com/opspec-io/sdk-golang/adapters"
)

var _ = Describe("_sdk", func() {

  Context("newSdk()", func() {
    It("should not return nil", func() {

      /* arrange/act */
      objectUnderTest := newSdk(new(FakeFilesystem), new(adapters.FakeEngineHost))

      /* assert */
      Expect(objectUnderTest).To(Not(BeNil()))

    })
  })
  Context(".CreateCollection() method", func() {
    It("should invoke compositionRoot.createCollectionUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedCreateCollectionReq := models.NewCreateCollectionReq("", "", "")

      // wire up fakes
      fakeCreateCollectionUseCase := new(fakeCreateCollectionUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.CreateCollectionUseCaseReturns(fakeCreateCollectionUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.CreateCollection(*providedCreateCollectionReq)

      /* assert */
      Expect(fakeCreateCollectionUseCase.ExecuteArgsForCall(0)).To(Equal(*providedCreateCollectionReq))
      Expect(fakeCreateCollectionUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".CreateOp() method", func() {
    It("should invoke compositionRoot.createOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedCreateOpReq := models.NewCreateOpReq("", "", "")

      // wire up fakes
      fakeCreateOpUseCase := new(fakeCreateOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.CreateOpUseCaseReturns(fakeCreateOpUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.CreateOp(*providedCreateOpReq)

      /* assert */
      Expect(fakeCreateOpUseCase.ExecuteArgsForCall(0)).To(Equal(*providedCreateOpReq))
      Expect(fakeCreateOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".GetEventStream() method", func() {
    It("should invoke compositionRoot.getEventStreamUseCase.Execute() with expected args & return result", func() {

      /* arrange */

      // wire up fakes
      fakeGetEventStreamUseCase := new(fakeGetEventStreamUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.GetEventStreamUseCaseReturns(fakeGetEventStreamUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.GetEventStream()

      /* assert */
      Expect(fakeGetEventStreamUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".GetCollection() method", func() {
    It("should invoke compositionRoot.getCollectionUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedCollectionBundlePath := ""

      // wire up fakes
      fakeGetCollectionUseCase := new(fakeGetCollectionUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.GetCollectionUseCaseReturns(fakeGetCollectionUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.GetCollection(providedCollectionBundlePath)

      /* assert */
      Expect(fakeGetCollectionUseCase.ExecuteArgsForCall(0)).To(Equal(providedCollectionBundlePath))
      Expect(fakeGetCollectionUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".GetOp() method", func() {
    It("should invoke compositionRoot.getOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedCollectionBundlePath := ""

      // wire up fakes
      fakeGetOpUseCase := new(fakeGetOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.GetOpUseCaseReturns(fakeGetOpUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.GetOp(providedCollectionBundlePath)

      /* assert */
      Expect(fakeGetOpUseCase.ExecuteArgsForCall(0)).To(Equal(providedCollectionBundlePath))
      Expect(fakeGetOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".KillOpRun() method", func() {
    It("should invoke compositionRoot.killOpRunUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedKillOpRunReq := models.KillOpRunReq{}

      // wire up fakes
      fakeKillOpRunUseCase := new(fakeKillOpRunUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.KillOpRunUseCaseReturns(fakeKillOpRunUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.KillOpRun(providedKillOpRunReq)

      /* assert */
      Expect(fakeKillOpRunUseCase.ExecuteArgsForCall(0)).To(Equal(providedKillOpRunReq))
      Expect(fakeKillOpRunUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetCollectionDescription() method", func() {
    It("should invoke compositionRoot.setCollectionDescriptionUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetCollectionDescriptionReq := models.NewSetCollectionDescriptionReq("", "")

      // wire up fakes
      fakeSetCollectionDescriptionUseCase := new(fakeSetCollectionDescriptionUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SetCollectionDescriptionUseCaseReturns(fakeSetCollectionDescriptionUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SetCollectionDescription(*providedSetCollectionDescriptionReq)

      /* assert */
      Expect(fakeSetCollectionDescriptionUseCase.ExecuteArgsForCall(0)).To(Equal(*providedSetCollectionDescriptionReq))
      Expect(fakeSetCollectionDescriptionUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".SetOpDescription() method", func() {
    It("should invoke compositionRoot.setOpDescriptionUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedSetOpDescriptionReq := models.NewSetOpDescriptionReq("", "")

      // wire up fakes
      fakeSetOpDescriptionUseCase := new(fakeSetOpDescriptionUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.SetOpDescriptionUseCaseReturns(fakeSetOpDescriptionUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.SetOpDescription(*providedSetOpDescriptionReq)

      /* assert */
      Expect(fakeSetOpDescriptionUseCase.ExecuteArgsForCall(0)).To(Equal(*providedSetOpDescriptionReq))
      Expect(fakeSetOpDescriptionUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".StartOpRun() method", func() {
    It("should invoke compositionRoot.startOpRunUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedStartOpRunReq := models.StartOpRunReq{}

      // wire up fakes
      fakeStartOpRunUseCase := new(fakeStartOpRunUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.StartOpRunUseCaseReturns(fakeStartOpRunUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.StartOpRun(providedStartOpRunReq)

      /* assert */
      Expect(fakeStartOpRunUseCase.ExecuteArgsForCall(0)).To(Equal(providedStartOpRunReq))
      Expect(fakeStartOpRunUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".TryResolveDefaultCollection() method", func() {
    It("should invoke compositionRoot.tryResolveDefaultCollectionUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedTryResolveDefaultCollectionReq := models.NewTryResolveDefaultCollectionReq("")

      // wire up fakes
      fakeTryResolveDefaultCollectionUseCase := new(fakeTryResolveDefaultCollectionUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.TryResolveDefaultCollectionUseCaseReturns(fakeTryResolveDefaultCollectionUseCase)

      objectUnderTest := &_sdk{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.TryResolveDefaultCollection(*providedTryResolveDefaultCollectionReq)

      /* assert */
      Expect(fakeTryResolveDefaultCollectionUseCase.ExecuteArgsForCall(0)).To(Equal(*providedTryResolveDefaultCollectionReq))
      Expect(fakeTryResolveDefaultCollectionUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })

})
