package consumenodeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-interfaces/vhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
	"net/http"
)

var _ = Describe("StartOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedStartOpReq := model.StartOpReq{
			Args:   map[string]*model.Data{},
			PkgRef: "dummyPkgRef",
		}

		expectedReqBytes, _ := json.Marshal(providedStartOpReq)
		expectedResult := "dummyOpId"

		expectedHttpReq, _ := http.NewRequest(
			"POST",
			fmt.Sprintf("http://%v/ops/starts", "localhost:42224"),
			bytes.NewBuffer(expectedReqBytes),
		)

		fakeHttpClient := new(vhttp.Fake)
		fakeHttpClient.DoReturns(&http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte(expectedResult)))}, nil)

		objectUnderTest := consumeNodeApi{
			httpClient: fakeHttpClient,
		}

		/* act */
		actualResult, _ := objectUnderTest.StartOp(providedStartOpReq)

		/* assert */
		Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))
		Expect(expectedResult).To(Equal(actualResult))

	})
})
