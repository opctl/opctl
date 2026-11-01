package node

import (
	"context"
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("legacy YAML event bytes", func() {
	It("keeps start and completion event JSON unchanged without container results", func() {
		spec := &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}}}
		call := &model.Call{ID: "id", RootID: "root", Container: &model.ContainerCall{ContainerID: "id", Image: &model.ContainerCallImage{}, EnvVars: map[string]string{"OPCTL_OUTPUTS": "authored"}, Dirs: map[string]string{"/opctl/outputs": "/authored"}}}
		ctx := withPrivateCompletions(context.Background())
		Expect(privateResultsActive(ctx)).To(BeFalse())
		snapshot := callForEvent(ctx, call, spec)
		expectedCall := `{"container":{"opPath":"","containerId":"id","cmd":null,"dirs":{"/opctl/outputs":"/authored"},"envVars":{"OPCTL_OUTPUTS":"authored"},"files":null,"image":{"ref":null},"sockets":null,"workDir":""},"id":"id","isKilled":false,"rootId":"root"}`
		start, err := json.Marshal(model.Event{CallStarted: &model.CallStarted{Call: snapshot, Ref: "op"}, Timestamp: time.Time{}})
		Expect(err).NotTo(HaveOccurred())
		Expect(string(start)).To(Equal(`{"callStarted":{"call":` + expectedCall + `,"ref":"op"},"timestamp":"0001-01-01T00:00:00Z"}`))
		end, err := json.Marshal(model.Event{CallEnded: &model.CallEnded{Call: snapshot, Ref: "op", Outputs: map[string]*model.Value{}, Outcome: model.OpOutcomeSucceeded}})
		Expect(err).NotTo(HaveOccurred())
		Expect(string(end)).To(Equal(`{"callEnded":{"call":` + expectedCall + `,"ref":"op","outputs":{},"outcome":"SUCCEEDED"},"timestamp":"0001-01-01T00:00:00Z"}`))
	})
})
