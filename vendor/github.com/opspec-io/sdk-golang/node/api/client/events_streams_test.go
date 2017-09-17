package client

import (
	"encoding/json"
	"errors"
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/url"
)

var _ = Context("GetEventStream", func() {

	It("should call wsDialer.Dial() w/ expected args", func() {

		/* arrange */
		providedReq := &model.GetEventStreamReq{
			Filter: &model.EventFilter{
				RootOpIds: []string{
					"dummyROID1",
				},
			},
		}

		// construct expected URL
		expectedReqUrl := url.URL{}
		expectedReqUrl.Scheme = "ws"
		expectedReqUrl.Path = api.URLEvents_Stream

		if nil != providedReq.Filter {
			// add non-nil filter
			var filterBytes []byte
			filterBytes, err := json.Marshal(providedReq.Filter)
			if nil != err {
				panic(err)
			}
			queryValues := expectedReqUrl.Query()
			queryValues.Add("filter", string(filterBytes))

			expectedReqUrl.RawQuery = queryValues.Encode()
		}

		fakeWSDialer := new(iwebsocket.FakeDialer)
		//error to trigger immediate retur
		fakeWSDialer.DialReturns(nil, nil, errors.New("dummyError"))

		objectUnderTest := client{
			wsDialer: fakeWSDialer,
		}

		/* act */
		objectUnderTest.GetEventStream(providedReq)

		/* assert */
		actualReqUrl, _ := fakeWSDialer.DialArgsForCall(0)
		Expect(actualReqUrl).To(Equal(expectedReqUrl.String()))

	})
})
