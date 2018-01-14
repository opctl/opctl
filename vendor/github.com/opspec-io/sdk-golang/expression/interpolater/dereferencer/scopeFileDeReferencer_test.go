package dereferencer

import (
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
	"path/filepath"
)

var _ = Context("scopeFileDeReferencer", func() {
	Context("ref is scope file ref", func() {
		It("should call ioutil.ReadFile w/ expected args", func() {
			/* arrange */
			scopeRef := "dummyScopeRef"
			scopeValue := "/dummyScopeValue"
			pathSegment1 := "pathSegment1"
			pathSegment2 := "pathSegment2"

			providedRef := filepath.Join(scopeRef, pathSegment1, pathSegment2)

			expectedPath := filepath.Join(scopeValue, pathSegment1, pathSegment2)

			fakeIoUtil := new(iioutil.Fake)
			// err to trigger immediate return
			fakeIoUtil.ReadFileReturns([]byte{}, errors.New("dummyError"))

			objectUnderTest := _scopeFileDeReferencer{
				ioutil: fakeIoUtil,
			}

			/* act */

			objectUnderTest.DeReferenceScopeFile(
				providedRef,
				map[string]*model.Value{
					scopeRef: {Dir: &scopeValue},
				},
			)

			/* assert */
			actualPath := fakeIoUtil.ReadFileArgsForCall(0)

			Expect(actualPath).To(Equal(expectedPath))
		})
	})
})
