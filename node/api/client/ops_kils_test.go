package client

import (
	"bytes"
	"encoding/json"
	"github.com/golang-interfaces/vhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/http"
	"net/url"
)

var _ = Describe("KillOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedReq := model.KillOpReq{
			OpId: "dummyRootOpId",
		}

		expectedReqUrl := url.URL{}
		expectedReqUrl.Path = api.Ops_KillsURLTpl

		expectedBytes, _ := json.Marshal(providedReq)

		expectedHttpReq, _ := http.NewRequest(
			"POST",
			expectedReqUrl.String(),
			bytes.NewBuffer(expectedBytes),
		)

		fakeHttpClient := new(vhttp.Fake)

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		objectUnderTest.KillOp(providedReq)

		/* assert */
		Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))

	})
})
