package local

import (
	"fmt"
	"net/url"

	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

func newAPIClientNode(
	listenAddress string,
) (node.Node, error) {
	apiBaseURL, err := url.Parse(fmt.Sprintf("http://%s/api", listenAddress))
	if nil != err {
		return nil, err
	}

	return client.New(
		*apiBaseURL,
		&client.Opts{
			RetryLogHook: func(err error) {
				//fmt.Println("request resulted in a recoverable error & will be retried", err)
			},
		},
	), nil
}
