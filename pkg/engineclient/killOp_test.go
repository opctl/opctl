package engineclient

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/http"
	netHttp "net/http"
)

var _ = Describe("KillOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedKillOpReq := model.KillOpReq{
			OpGraphId: "dummyOpGraphId",
		}

		expectedBytes, _ := format.NewJsonFormat().From(providedKillOpReq)

		expectedHttpReq, _ := netHttp.NewRequest(
			"POST",
			fmt.Sprintf("http://%v/instances/kills", "localhost:42224"),
			bytes.NewBuffer(expectedBytes),
		)

		fakeHttpClient := new(http.Fake)

		objectUnderTest := _engineClient{
			httpClient: fakeHttpClient,
			jsonFormat: format.NewJsonFormat(),
		}

		/* act */
		objectUnderTest.KillOp(providedKillOpReq)

		/* assert */
		Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))

	})
})
