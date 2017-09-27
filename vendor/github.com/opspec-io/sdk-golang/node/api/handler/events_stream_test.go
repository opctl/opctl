package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Context("GET /events/stream", func() {
	Context("upgrader.Upgrade errors", func() {
		It("should return StatusCode of 400", func() {

			/* arrange */
			fakeUpgrader := new(iwebsocket.FakeUpgrader)
			fakeUpgrader.UpgradeReturns(nil, errors.New("dummyError"))

			objectUnderTest := New(new(core.Fake))

			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(http.MethodGet, api.URLEvents_Stream, bytes.NewReader([]byte{}))
			if nil != err {
				panic(err.Error())
			}

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))

		})
	})
	Context("nonempty since", func() {
		Context("time.Parse errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				objectUnderTest := New(new(core.Fake))

				invalidSince := "notValidTime"

				recorder := httptest.NewRecorder()

				httpReq, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%v?since=%v", api.URLEvents_Stream, invalidSince),
					bytes.NewReader([]byte{}),
				)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				defer func() {
					// conn.Close() will panic so recover (no way to fake it)
					recover()
				}()
				objectUnderTest.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

			})
		})
		Context("time.Parse doesn't error", func() {
			It("should call core.GetEventStream w/ expected args", func() {

				/* arrange */
				fakeCore := new(core.Fake)
				fakeCore.GetEventStreamStub = func(req *model.GetEventStreamReq, eventChannel chan *model.Event) (err error) {
					// close eventChannel to trigger immediate return
					close(eventChannel)
					return
				}

				objectUnderTest := New(fakeCore)

				expectedSince := time.Now().UTC()
				expectedReq := &model.GetEventStreamReq{
					Filter: &model.EventFilter{
						Since: &expectedSince,
					},
				}

				httpReq, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%v?since=%v", api.URLEvents_Stream, expectedSince.Format(time.RFC3339)),
					bytes.NewReader([]byte{}),
				)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				defer func() {
					// conn.Close() will panic so recover (no way to fake it)
					recover()
				}()
				objectUnderTest.ServeHTTP(httptest.NewRecorder(), httpReq)

				/* assert */
				actualReq, _ := fakeCore.GetEventStreamArgsForCall(0)
				Expect(*actualReq).To(Equal(*expectedReq))

			})
		})
	})
	Context("nonempty roots", func() {
		It("should call core.GetEventStream w/ expected args", func() {

			/* arrange */
			fakeCore := new(core.Fake)
			fakeCore.GetEventStreamStub = func(req *model.GetEventStreamReq, eventChannel chan *model.Event) (err error) {
				// close eventChannel to trigger immediate return
				close(eventChannel)
				return
			}

			objectUnderTest := New(fakeCore)

			root1 := "dummyRoot1"
			root2 := "dummyRoot2"

			expectedRoots := []string{root1, root2}
			expectedReq := &model.GetEventStreamReq{
				Filter: &model.EventFilter{
					Roots: expectedRoots,
				},
			}

			httpReq, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%v?roots=%v,%v", api.URLEvents_Stream, expectedRoots[0], expectedRoots[1]),
				bytes.NewReader([]byte{}),
			)
			if nil != err {
				panic(err.Error())
			}

			/* act */
			defer func() {
				// conn.Close() will panic so recover (no way to fake it)
				recover()
			}()
			objectUnderTest.ServeHTTP(httptest.NewRecorder(), httpReq)

			/* assert */
			actualReq, _ := fakeCore.GetEventStreamArgsForCall(0)
			Expect(*actualReq).To(Equal(*expectedReq))

		})
	})
})
