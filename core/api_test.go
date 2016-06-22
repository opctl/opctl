package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opctl/engine/core/models"
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
  Context(".RunOp() method", func() {
    It("should invoke compositionRoot.runOpUseCase.Execute() with expected args & return result", func() {

      /* arrange */
      providedRunOpReq := models.NewRunOpReq("", map[string]string{})

      // wire up fakes
      fakeRunOpUseCase := new(fakeRunOpUseCase)

      fakeCompositionRoot := new(fakeCompositionRoot)
      fakeCompositionRoot.RunOpUseCaseReturns(fakeRunOpUseCase)

      objectUnderTest := &_api{
        compositionRoot:fakeCompositionRoot,
      }

      /* act */
      objectUnderTest.RunOp(*providedRunOpReq)

      /* assert */
      Expect(fakeRunOpUseCase.ExecuteArgsForCall(0)).To(Equal(*providedRunOpReq))
      Expect(fakeRunOpUseCase.ExecuteCallCount()).To(Equal(1))

    })
  })

})
