package stream

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	iwebsocket "github.com/golang-interfaces/github.com-gorilla-websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/urltemplates"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

var _ = Context("Handler", func() {
	Context("NewHandler", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewHandler(new(nodeFakes.FakeNode))).Should(Not(BeNil()))
		})
	})
	Context("Handle", func() {
		Context("upgrader.Upgrade errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				fakeUpgrader := new(iwebsocket.FakeUpgrader)
				fakeUpgrader.UpgradeReturns(nil, errors.New("dummyError"))

				objectUnderTest := _handler{
					node:     new(nodeFakes.FakeNode),
					upgrader: fakeUpgrader,
				}

				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest(http.MethodGet, urltemplates.Events_Stream, bytes.NewReader([]byte{}))
				if err != nil {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.Handle(providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusBadRequest))

			})
		})
		Context("nonempty since", func() {
			Context("time.Parse errors", func() {
				It("should return StatusCode of 400", func() {

					/* arrange */
					objectUnderTest := _handler{
						node: new(nodeFakes.FakeNode),
					}

					invalidSince := "notValidTime"

					providedHTTPResp := httptest.NewRecorder()

					providedHTTPReq, err := http.NewRequest(
						http.MethodGet,
						fmt.Sprintf("%v?since=%v", urltemplates.Events_Stream, invalidSince),
						bytes.NewReader([]byte{}),
					)
					if err != nil {
						panic(err.Error())
					}

					/* act */
					defer func() {
						// conn.Close() will panic so recover (no way to fake it)
						recover()
					}()
					objectUnderTest.Handle(providedHTTPResp, providedHTTPReq)

					/* assert */
					Expect(providedHTTPResp.Code).To(Equal(http.StatusBadRequest))

				})
			})
			Context("time.Parse doesn't error", func() {
				It("should call core.GetEventStream w/ expected args", func() {

					/* arrange */
					fakeCore := new(nodeFakes.FakeNode)
					eventChannel := make(chan model.Event)
					// close eventChannel to trigger immediate return
					close(eventChannel)
					fakeCore.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _handler{
						node: fakeCore,
					}

					expectedSince := time.Now().UTC()
					expectedReq := &model.GetEventStreamReq{
						Filter: model.EventFilter{
							Since: &expectedSince,
						},
					}

					providedHTTPReq, err := http.NewRequest(
						http.MethodGet,
						fmt.Sprintf("%v?since=%v", urltemplates.Events_Stream, expectedSince.Format(time.RFC3339)),
						bytes.NewReader([]byte{}),
					)
					if err != nil {
						panic(err.Error())
					}

					/* act */
					defer func() {
						// conn.Close() will panic so recover (no way to fake it)
						recover()
					}()
					objectUnderTest.Handle(httptest.NewRecorder(), providedHTTPReq)

					/* assert */
					_, actualReq := fakeCore.GetEventStreamArgsForCall(0)
					Expect(*actualReq).To(Equal(*expectedReq))

				})
			})
		})
		Context("nonempty roots", func() {
			It("should call core.GetEventStream w/ expected args", func() {

				/* arrange */
				fakeCore := new(nodeFakes.FakeNode)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakeCore.GetEventStreamReturns(eventChannel, nil)

				objectUnderTest := _handler{
					node: fakeCore,
				}

				root1 := "dummyRoot1"
				root2 := "dummyRoot2"

				expectedRoots := []string{root1, root2}
				expectedReq := &model.GetEventStreamReq{
					Filter: model.EventFilter{
						Roots: expectedRoots,
					},
				}

				providedHTTPReq, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%v?roots=%v,%v", urltemplates.Events_Stream, expectedRoots[0], expectedRoots[1]),
					bytes.NewReader([]byte{}),
				)
				if err != nil {
					panic(err.Error())
				}

				/* act */
				defer func() {
					// conn.Close() will panic so recover (no way to fake it)
					recover()
				}()
				objectUnderTest.Handle(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualReq := fakeCore.GetEventStreamArgsForCall(0)
				Expect(*actualReq).To(Equal(*expectedReq))

			})
		})
	})
})
