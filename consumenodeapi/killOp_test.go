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

	It("should call httpClient.Do() with expected args", func() {

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
		Expect(expectedHTTPReq).To(Equal(fakeHTTPClient.DoArgsForCall(0)))

	})
})
