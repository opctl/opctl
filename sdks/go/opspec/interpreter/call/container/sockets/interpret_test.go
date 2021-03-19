package sockets

import (
	"path/filepath"

	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	It("should return expected dcg.Sockets", func() {
		/* arrange */
		providedCurrentScopeRef1 := "scopeRef1"
		providedCurrentScopeRef1String := "scopeRef1String"
		providedCurrentScopeRef2 := "/unix-socket"
		providedCurrentScopeRef2String := "scopeRef2String"

		providedCurrentScope := map[string]*model.Value{
			providedCurrentScopeRef1: {Socket: &providedCurrentScopeRef1String},
		}

		providedContainerCallSpecSockets := map[string]string{
			// explicitly bound to scope
			providedCurrentScopeRef1: providedCurrentScopeRef1,
			// bound as unix socket
			providedCurrentScopeRef2: providedCurrentScopeRef2String,
		}

		providedScratchDirPath, err := ioutil.TempDir("", "")
		if err != nil {
			panic(err)
		}

		expectedSockets := map[string]string{
			providedCurrentScopeRef1: providedCurrentScopeRef1String,
			providedCurrentScopeRef2: filepath.Join(providedScratchDirPath, providedCurrentScopeRef2),
		}

		/* act */
		actualContainerCallSockets, err := Interpret(
			providedCurrentScope,
			providedContainerCallSpecSockets,
			providedScratchDirPath,
		)
		if err != nil {
			panic(err)
		}

		/* assert */
		Expect(actualContainerCallSockets).To(Equal(expectedSockets))
	})
})
