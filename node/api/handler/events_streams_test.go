package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Context("getEventStreamHandler", func() {
	Context("ServeHTTP() method", func() {
		Describe("upgrader.Upgrade errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				fakeUpgrader := new(iwebsocket.FakeUpgrader)
				fakeUpgrader.UpgradeReturns(nil, errors.New("dummyError"))

				objectUnderTest := getEventStreamHandler{
					upgrader: fakeUpgrader,
				}

				recorder := httptest.NewRecorder()
				m := mux.NewRouter()
				m.Handle(api.Events_StreamsURLTpl, objectUnderTest)

				httpReq, err := http.NewRequest(http.MethodGet, api.Events_StreamsURLTpl, bytes.NewReader([]byte{}))
				if nil != err {
					Fail(err.Error())
				}

				/* act */
				m.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

			})
		})
		Describe("nonempty filter", func() {
			Describe("json.Unmarshal errors", func() {
				It("should return StatusCode of 400", func() {

					/* arrange */
					objectUnderTest := getEventStreamHandler{
						upgrader: new(iwebsocket.FakeUpgrader),
					}

					invalidFilter := "notvalidjson"

					recorder := httptest.NewRecorder()
					m := mux.NewRouter()
					m.Handle(api.Events_StreamsURLTpl, objectUnderTest)

					httpReq, err := http.NewRequest(
						http.MethodGet,
						fmt.Sprintf("%v?filter=%v", api.Events_StreamsURLTpl, invalidFilter),
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
					m.ServeHTTP(recorder, httpReq)

					/* assert */
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				})
			})
			Describe("json.Unmarshal doesn't error", func() {
				It("should call core.GetEventStream w/ expected args", func() {

					/* arrange */
					fakeCore := new(core.Fake)
					fakeCore.GetEventStreamStub = func(req *model.GetEventStreamReq, eventChannel chan *model.Event) (err error) {
						// close eventChannel to trigger immediate return
						close(eventChannel)
						return
					}

					objectUnderTest := getEventStreamHandler{
						core:     fakeCore,
						upgrader: new(iwebsocket.FakeUpgrader),
					}

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

					m := mux.NewRouter()
					m.Handle(api.Events_StreamsURLTpl, objectUnderTest)

					httpReq, err := http.NewRequest(
						http.MethodGet,
						fmt.Sprintf("%v?filter=%v", api.Events_StreamsURLTpl, string(filterBytes)),
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
					m.ServeHTTP(httptest.NewRecorder(), httpReq)

					/* assert */
					actualReq, _ := fakeCore.GetEventStreamArgsForCall(0)
					Expect(*actualReq).To(Equal(*expectedReq))

				})
			})
		})
	})
})
