package apiclient

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ApiClient

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/http"
	"github.com/sethgrid/pester"
)

type ApiClient interface {
	GetEventStream(
		req *model.GetEventStreamReq,
	) (
		stream chan model.Event,
		err error,
	)

	KillOp(
		req model.KillOpReq,
	) (
		err error,
	)

	StartOp(
		req model.StartOpReq,
	) (
		opId string,
		err error,
	)
}

func New() ApiClient {

	httpClient := pester.New()
	httpClient.Backoff = pester.ExponentialBackoff

	return &_apiClient{
		httpClient: httpClient,
		jsonFormat: format.NewJsonFormat(),
	}
}

type _apiClient struct {
	httpClient http.Client
	jsonFormat format.Format
}
