package consumenodeapi

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/http"
	"io/ioutil"
	netHttp "net/http"
)

var _ = Describe("StartOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedStartOpReq := model.StartOpReq{
			Args:     map[string]*model.Data{},
			OpPkgRef: "dummyOpPkgRef",
		}

		expectedReqBytes, _ := format.NewJsonFormat().From(providedStartOpReq)
		expectedResult := "dummyOpId"

		expectedHttpReq, _ := netHttp.NewRequest(
			"POST",
			fmt.Sprintf("http://%v/ops/starts", "localhost:42224"),
			bytes.NewBuffer(expectedReqBytes),
		)

		fakeHttpClient := new(http.Fake)
		fakeHttpClient.DoReturns(&netHttp.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte(expectedResult)))}, nil)

		objectUnderTest := consumeNodeApi{
			httpClient: fakeHttpClient,
			jsonFormat: format.NewJsonFormat(),
		}

		/* act */
		actualResult, _ := objectUnderTest.StartOp(providedStartOpReq)

		/* assert */
		Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))
		Expect(expectedResult).To(Equal(actualResult))

	})
})
