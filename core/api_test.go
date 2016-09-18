package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
)

var _ = Describe("_api", func() {
  Context(".GetEventStream() method", func() {
    It("should invoke compositionRoot.getEventStreamUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedGetEventStreamChannel := make(chan models.Event)

      // wire up fakes
      fakeGetEventStreamUseCase := new(fakeGetEventStreamUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.GetEventStreamUseCaseReturns(fakeGetEventStreamUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.GetEventStream(providedGetEventStreamChannel)

      /* assert */
      Expect(fakeGetEventStreamUseCase.ExecuteArgsForCall(0)).To(Equal(providedGetEventStreamChannel))
      Expect(fakeGetEventStreamUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".KillOpRun() method", func() {
    It("should invoke compositionRoot.killOpRunUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedKillOpRunReq := models.NewKillOpRunReq("dummyOpRunId")

      // wire up fakes
      fakeKillOpRunUseCase := new(fakeKillOpRunUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.KillOpRunUseCaseReturns(fakeKillOpRunUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.KillOpRun(*providedKillOpRunReq)

      /* assert */
      Expect(fakeKillOpRunUseCase.ExecuteArgsForCall(0)).To(Equal(*providedKillOpRunReq))
      Expect(fakeKillOpRunUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })
  Context(".StartOpRun() method", func() {
    It("should invoke compositionRoot.startOpRunUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedStartOpRunReq := models.NewStartOpRunReq("", map[string]string{})

      // wire up fakes
      fakeStartOpRunUseCase := new(fakeStartOpRunUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.StartOpRunUseCaseReturns(fakeStartOpRunUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.StartOpRun(*providedStartOpRunReq)

      /* assert */
      Expect(fakeStartOpRunUseCase.ExecuteArgsForCall(0)).To(Equal(*providedStartOpRunReq))
      Expect(fakeStartOpRunUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })

})
