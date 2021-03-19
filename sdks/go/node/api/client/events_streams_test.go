package client

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	iwebsocket "github.com/golang-interfaces/github.com-gorilla-websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api"
)

var _ = Context("GetEventStream", func() {

	XIt("should call wsDialer.DialContext() w/ expected args", func() {

		/* arrange */
		providedCtx := context.Background()
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
		expectedReqURL := url.URL{}
		expectedReqURL.Scheme = "ws"
		expectedReqURL.Path = api.URLEvents_Stream

		queryValues := expectedReqURL.Query()
		if providedReq.Filter.Since != nil {
			queryValues.Add("since", providedReq.Filter.Since.Format(time.RFC3339))
		}
		if providedReq.Filter.Roots != nil {
			queryValues.Add("roots", strings.Join(providedReq.Filter.Roots, ","))
		}
		expectedReqURL.RawQuery = queryValues.Encode()

		fakeWSDialer := new(iwebsocket.FakeDialer)
		//error to trigger immediate retur
		fakeWSDialer.DialReturns(nil, nil, errors.New("dummyError"))

		objectUnderTest := apiClient{
			wsDialer: fakeWSDialer,
		}

		/* act */
		objectUnderTest.GetEventStream(
			providedCtx,
			providedReq,
		)

		/* assert */
		actualCtx,
			actualReqURL, _ := fakeWSDialer.DialContextArgsForCall(0)

		Expect(actualCtx).To(Equal(providedCtx))
		Expect(actualReqURL).To(Equal(expectedReqURL.String()))

	})
})
