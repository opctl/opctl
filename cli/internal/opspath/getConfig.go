package opspath

import (
	"context"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

type Config struct {
	Ops model.Value
}

func getConfig(
	ctx context.Context,
	startPath string,
	node node.Node,
) (Config, error) {
	originURL, err := tryGetOriginURL(startPath)
	if err != nil {
		return Config{}, err
	}

	startTime := time.Now()

	rootCallID, err := node.StartOp(
		ctx,
		model.StartOpReq{
			Op: model.StartOpReqOp{
				Ref: filepath.Join(
					originURL.Host,
					path.Dir(originURL.Path),
					model.DotOpctlDirName,
				) + "#1.0.0/getConfig",
			},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve") {
			err = nil
		}
		return Config{}, err
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
		return Config{}, err
	}

	for e := range eventChan {
		if e.CallEnded != nil && e.CallEnded.Call.ID == rootCallID {
			if opsDir, ok := e.CallEnded.Outputs["ops"]; ok && opsDir.Dir != nil {
				return Config{
					Ops: *opsDir,
				}, nil
			}
			// error (we assume because not found)
			return Config{}, nil
		}
	}

	return Config{}, nil
}
