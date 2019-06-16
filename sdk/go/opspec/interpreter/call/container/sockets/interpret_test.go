package sockets

import (
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/model"
	"path/filepath"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("should return expected dcg.Sockets", func() {
			/* arrange */
			providedCurrentScopeRef1 := "dummyScopeRef1"
			providedCurrentScopeRef1String := "dummyScopeRef1String"
			providedCurrentScopeRef2 := "/dummy-unix-socket"
			providedCurrentScopeRef2String := "dummyScopeRef2String"

			providedCurrentScope := map[string]*model.Value{
				providedCurrentScopeRef1: {Socket: &providedCurrentScopeRef1String},
			}

			providedSCGContainerCallSockets := map[string]string{
				// explicitly bound to scope
				providedCurrentScopeRef1: providedCurrentScopeRef1,
				// bound as unix socket
				providedCurrentScopeRef2: providedCurrentScopeRef2String,
			}

			providedScratchDirPath := "dummyScratchDirPath"

			expectedSockets := map[string]string{
				providedCurrentScopeRef1: providedCurrentScopeRef1String,
				providedCurrentScopeRef2: filepath.Join(providedScratchDirPath, providedCurrentScopeRef2),
			}

			objectUnderTest := _interpreter{
				os: new(ios.Fake),
			}

			/* act */
			actualDCGContainerCallSockets, err := objectUnderTest.Interpret(
				providedCurrentScope,
				providedSCGContainerCallSockets,
				providedScratchDirPath,
			)
			if nil != err {
				panic(err)
			}

			/* assert */
			Expect(actualDCGContainerCallSockets).To(Equal(expectedSockets))
		})
	})
})
