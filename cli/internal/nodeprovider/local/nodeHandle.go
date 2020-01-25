package local

import (
	"net/url"
	"os"

	"github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

func newNodeHandle() (model.NodeHandle, error) {
	apiBaseURLStr := os.Getenv("OPCTL_CLI_API_BASEURL")
	if "" == apiBaseURLStr {
		apiBaseURLStr = "http://localhost:42224/api"
	}
	apiBaseURL, err := url.Parse(apiBaseURLStr)
	if nil != err {
		return nil, err
	}

	return _nodeHandle{
		apiBaseURL: *apiBaseURL,
	}, nil
}

type _nodeHandle struct {
	apiBaseURL url.URL
}

func (nh _nodeHandle) APIClient() client.Client {
	apiClient := client.New(
		nh.apiBaseURL,
		&client.Opts{
			RetryLogHook: func(err error) {
				//cliOutput.Attention("request resulted in a recoverable error & will be retried; error was: %v", err)
			},
		},
	)

	return apiClient
}
