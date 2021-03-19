
> *Be advised: this project is currently at Major version zero. Per the
> semantic versioning spec: "Major version zero (0.y.z) is for initial
> development. Anything may change at any time. The public API should
> not be considered stable."*

# Problem statement
Go SDK for [opctl](https://opctl.io).


# Documentation
Documentation for SDK packages are maintained in golang's native go doc format; which is web browsable via the [godoc webpage](https://pkg.go.dev/mod/github.com/opctl/opctl/sdks/go)


# Usage

## Run an op using an API client.
```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

func constructAPIClient(
) client.Client {
	// get client
	nodeURL, err := url.Parse("http://localhost:42224/api")
	if err != nil {
		panic(err)
	}

	return client.New(
		*nodeURL,
		&client.Opts{
			RetryLogHook: func(err error) {
				fmt.Printf("unable to connect to node; error was: %v\n", err)
			},
		},
	)
}

func main() {
    ctx := context.Background()
	startTime := time.Now()
	apiClient := constructAPIClient()

	// start op
	rootID, err := apiClient.StartOp(
		ctx,
		model.StartOpReq{
			Args: map[string]*model.Value{},
			Op: model.StartOpReqOp{
				Ref: "github.com/opspec-pkgs/uuid.v4.generate#1.1.0",
			},
		},
	)
	if err != nil {
		panic(err)
	}

	// get event stream filtered to events from our op
	eventChan, err := apiClient.GetEventStream(
		ctx,
		&model.GetEventStreamReq{
			Filter: model.EventFilter{
				Roots: []string{rootID},
				Since: &startTime,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	for event := range eventChan {

		// print events
		eventAsJSON, err := json.MarshalIndent(event, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(eventAsJSON))

		if event.CallEnded != nil &&
			event.CallEnded.CallID == rootID {

			// close event stream on root op ended
			close(eventChan)

		}
	}
}
```


# Contributing
see [CONTRIBUTING.md](CONTRIBUTING.md)