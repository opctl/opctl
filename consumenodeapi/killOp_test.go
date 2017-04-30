package consumenodeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/vhttp"
	netHttp "net/http"
)

var _ = Describe("KillOp", func() {

	It("should call httpClient.Do() w/ expected args", func() {

		/* arrange */
		providedKillOpReq := model.KillOpReq{
			OpId: "dummyRootOpId",
		}

		expectedBytes, _ := json.Marshal(providedKillOpReq)

		expectedHTTPReq, _ := netHttp.NewRequest(
			"POST",
			fmt.Sprintf("http://%v/ops/kills", "localhost:42224"),
			bytes.NewBuffer(expectedBytes),
		)

		fakeHTTPClient := new(vhttp.Fake)

		objectUnderTest := consumeNodeApi{
			httpClient: fakeHTTPClient,
		}

		/* act */
		objectUnderTest.KillOp(providedKillOpReq)

		/* assert */
		// can't simply assert on req due to non-public http.Request state
		actualHTTPReq := fakeHTTPClient.DoArgsForCall(0)
		Expect(expectedHTTPReq.Method).To(Equal(actualHTTPReq.Method))
		Expect(expectedHTTPReq.URL).To(Equal(actualHTTPReq.URL))
		Expect(expectedHTTPReq.Proto).To(Equal(actualHTTPReq.Proto))
		Expect(expectedHTTPReq.Body).To(Equal(actualHTTPReq.Body))

	})
})
