package client

import (
	"errors"
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/url"
	"strings"
	"time"
)

var _ = Context("GetEventStream", func() {

	It("should call wsDialer.Dial() w/ expected args", func() {

		/* arrange */
		providedSince := time.Now().UTC()
		providedReq := &model.GetEventStreamReq{
			Filter: model.EventFilter{
				Since: &providedSince,
				Roots: []string{
					"dummyRoot",
				},
			},
		}

		// construct expected URL
		expectedReqUrl := url.URL{}
		expectedReqUrl.Scheme = "ws"
		expectedReqUrl.Path = api.URLEvents_Stream

		queryValues := expectedReqUrl.Query()
		if nil != providedReq.Filter.Since {
			queryValues.Add("since", providedReq.Filter.Since.Format(time.RFC3339))
		}
		if nil != providedReq.Filter.Roots {
			queryValues.Add("roots", strings.Join(providedReq.Filter.Roots, ","))
		}
		expectedReqUrl.RawQuery = queryValues.Encode()

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
