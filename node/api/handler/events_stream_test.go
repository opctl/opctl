package handler

import (
	"bytes"
	"encoding/json"
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
				Fail(err.Error())
			}

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))

		})
	})
	Context("nonempty filter", func() {
		Context("json.Unmarshal errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				objectUnderTest := New(new(core.Fake))

				invalidFilter := "notvalidjson"

				recorder := httptest.NewRecorder()

				httpReq, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%v?filter=%v", api.URLEvents_Stream, invalidFilter),
					bytes.NewReader([]byte{}),
				)
				if nil != err {
					Fail(err.Error())
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
		Context("json.Unmarshal doesn't error", func() {
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
						RootOpIds: []string{
							"dummyROID1",
						},
						Since: &expectedSince,
					},
				}

				filterBytes, err := json.Marshal(expectedReq.Filter)
				if nil != err {
					panic(err)
				}

				httpReq, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%v?filter=%v", api.URLEvents_Stream, string(filterBytes)),
					bytes.NewReader([]byte{}),
				)
				if nil != err {
					Fail(err.Error())
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
})
