package core

import "github.com/opspec-io/sdk-golang/pkg/model"

func (this _core) GetEventStream(
	req *model.GetEventStreamReq,
	subscriberEventChannel chan model.Event,
) (err error) {

	this.pubSub.RegisterSubscriber(req.Filter, subscriberEventChannel)

	return
}
