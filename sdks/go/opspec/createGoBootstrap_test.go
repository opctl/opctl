package opspec_test

import (
	"os"

	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/internal/gooptest"
)

func bootstrapLocally(opDir string, inputs map[string]interface{}) (*model.OpSpec, *goop.Plan) {
	described, planned, err := gooptest.Bootstrap(opDir, inputs)
	Expect(err).NotTo(HaveOccurred())
	return described, planned
}

var _ = Context("CreateGo bootstrap", func() {
	It("scaffolds an op that describes and plans with the in-repo SDK", func() {
		sdk, err := filepath.Abs("sdk")
		Expect(err).NotTo(HaveOccurred())
		parent, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(parent)
		path := filepath.Join(parent, "hello")
		Expect(opspec.CreateGo(path, "hello", "Greets the caller.", opspec.GoSDK{})).To(Succeed())
		goMod, err := os.OpenFile(filepath.Join(path, "go.mod"), os.O_APPEND|os.O_WRONLY, 0)
		Expect(err).NotTo(HaveOccurred())
		_, err = goMod.WriteString("\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => " + sdk + "\n")
		Expect(err).NotTo(HaveOccurred())
		Expect(goMod.Close()).To(Succeed())

		described, planned := bootstrapLocally(path, map[string]interface{}{"who": "world"})

		Expect(described.Name).To(Equal("hello"))
		Expect(described.Description).To(Equal("Greets the caller."))
		Expect(described.Inputs).To(HaveKey("who"))
		Expect(described.Inputs["who"].String.Default).To(Equal("world"))
		Expect(planned.Run.Container.Image.Ref).To(Equal("alpine:3"))
		Expect(planned.Run.Container.EnvVars).To(Equal(map[string]interface{}{"WHO": "$(who)"}))
	})
})

var _ = Context("run-a-go-service example", func() {
	It("describes and plans its Go dev op, which calls the YAML mysql op", func() {
		described, planned := bootstrapLocally("../../../examples/run-a-go-service/.opspec/dev", map[string]interface{}{
			"goBuildCachePath": "Library/Caches/go-build", "MYSQL_USER": "testuser", "MYSQL_PASSWORD": "testpassword", "MYSQL_DATABASE": "testapp", "MYSQL_HOST": "run-a-go-service-mysql",
		})

		Expect(described.Name).To(Equal("dev"))
		Expect(described.Description).To(Equal("Runs the example API with a seeded MySQL database for development."))
		Expect(described.Inputs).To(HaveKey("HOME"))
		Expect(described.Inputs["HOME"].Dir).NotTo(BeNil())
		Expect(described.Inputs["MYSQL_HOST"].String.Default).To(Equal("run-a-go-service-mysql"))

		calls := *planned.Run.Parallel
		Expect(calls).To(HaveLen(2))
		Expect(calls[0].Op.Ref).To(Equal("../mysql"))
		Expect(calls[0].Op.Inputs).To(HaveKeyWithValue("MYSQL_HOST", "$(MYSQL_HOST)"))
		api := calls[1].Container
		Expect(*api.Name).To(Equal("run-a-go-service-api"))
		Expect(api.Dirs).To(Equal(map[string]interface{}{
			"/src": "$(../..)", "/go/pkg": "$(HOME/go/pkg)", "/root/.cache/go-build": "$(HOME/Library/Caches/go-build)",
		}))
	})
})
