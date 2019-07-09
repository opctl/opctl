---
title: Golang
sidebar_label: Golang
---

## Examples

Run an op using the [Golang SDKs](https://github.com/opctl/sdk-golang) API client.

```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api/client"
)

func constructAPIClient(
) client.Client {
	// get client
	nodeURL, err := url.Parse("http://localhost:42224/api")
	if nil != err {
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
	if nil != err {
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
	if nil != err {
		panic(err)
	}

	for event := range eventChan {

		// print events
		eventAsJSON, err := json.MarshalIndent(event, "", "    ")
		if nil != err {
			panic(err)
		}
		fmt.Println(string(eventAsJSON))

		if nil != event.CallEnded &&
			event.CallEnded.CallID == rootID {

			// close event stream on root op ended
			close(eventChan)

		}
	}
}
```