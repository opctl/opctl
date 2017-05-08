package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/url"
	"strings"
)

var _ = Describe("GetEventStream", func() {

	It("should call wsDialer.Dial() w/ expected args", func() {

		/* arrange */
		providedReq := &model.GetEventStreamReq{
			Filter: &model.EventFilter{
				RootOpIds: []string{
					"dummyROID1",
				},
			},
		}

		// construct query params
		queryParams := []string{}
		if filter := providedReq.Filter; nil != filter {
			var filterBytes []byte
			filterBytes, err := json.Marshal(filter)
			if nil != err {
				panic(err)
			}
			queryParams = append(
				queryParams,
				fmt.Sprintf("filter=%v", url.QueryEscape(string(filterBytes))),
			)
		}

		expectedReqUrl := url.URL{}
		expectedReqUrl.Scheme = "ws"
		expectedReqUrl.Path = api.Events_StreamsURLTpl
		expectedReqUrl.RawQuery = strings.Join(queryParams, "&")

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
