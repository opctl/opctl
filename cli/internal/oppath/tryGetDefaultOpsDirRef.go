package oppath

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

func TryGetDefaultOpsDirRef(
	ctx context.Context,
	currentPath string,
	dataResolver dataresolver.DataResolver,
	node node.Node,
) (*string, error) {
	originURL, err := tryGetOriginURL(
		currentPath,
	)
	if err != nil {
		return nil, err
	}

	if originURL == nil {
		return nil, nil
	}

	getConfigOpRef := filepath.Join(
		originURL.Host,
		path.Dir(originURL.Path),
		".opctl",
	) + "#/getConfig"

	startTime := time.Now()

	rootCallID, err := node.StartOp(
		ctx,
		model.StartOpReq{
			Op: model.StartOpReqOp{
				Ref: getConfigOpRef,
			},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "unauthenticated") ||
			strings.Contains(err.Error(), "unauthorized") {
			fmt.Printf("run \"opctl auth add %s -u <username> -p <password>\" to complete setup\n\n", originURL.Hostname())
			return nil, nil
		} else if strings.Contains(err.Error(), "repository not found") {
			return nil, nil
		}

		return nil, err
	}

	eventChan, err := node.GetEventStream(
		ctx,
		&model.GetEventStreamReq{
			Filter: model.EventFilter{
				Roots: []string{
					rootCallID,
				},
				Since: &startTime,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	for e := range eventChan {
		if e.CallEnded != nil && e.CallEnded.Call.ID == rootCallID {
			if e.CallEnded.Error != nil {
				fmt.Println(e.CallEnded.Error.Message)
			}
			if defaultOps, ok := e.CallEnded.Outputs["defaultOps"]; ok && defaultOps.Dir != nil {
				return defaultOps.Dir, nil
			}
			// error (we assume because not found)
			return nil, nil
		}
	}

	return nil, nil
}
