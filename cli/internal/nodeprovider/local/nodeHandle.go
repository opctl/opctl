package local

import (
	"fmt"
	"net/url"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

func newNodeHandle(
	listenAddress string,
) (nodeprovider.NodeHandle, error) {
	apiBaseURL, err := url.Parse(fmt.Sprintf("http://%s/api", listenAddress))
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

func (nh _nodeHandle) APIClient() client.APIClient {
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
