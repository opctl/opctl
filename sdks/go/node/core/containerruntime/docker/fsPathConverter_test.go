package docker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	iruntimeFakes "github.com/opctl/opctl/sdks/go/internal/iruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/hostruntime"
	uuid "github.com/satori/go.uuid"
)

var _ = Context("fsPathConverter", func() {
	Context("when runtime.GOOS == windows", func() {

		fakeRuntime := new(iruntimeFakes.FakeIRuntime)
		fakeRuntime.GOOSReturns("windows")
		objectUnderTest := _fsPathConverter{
			runtime:     fakeRuntime,
			hostRuntime: hostruntime.RuntimeInfo{},
		}

		Context("when path contains drive letter", func() {
			It("should prepend a slash", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				actual := objectUnderTest.LocalToEngine("c:/DummyPath")

				/* assert */
				Expect(actual).To(Equal(expected))
			})
			It("should convert the drive letter to lowercase", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				actual := objectUnderTest.LocalToEngine("C:/DummyPath")

				/* assert */
				Expect(actual).To(Equal(expected))
			})
			It("should strip colon from drive", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				result := objectUnderTest.LocalToEngine("c:/DummyPath")

				/* assert */
				Expect(result).To(Equal(expected))
			})
			It("should replace backslashes with forward slashes", func() {
				/* arrange */
				pathWithMultipleBackslashes := `c\\dummy\path`
				expected := `c//dummy/path`

				/* act */
				result := objectUnderTest.LocalToEngine(pathWithMultipleBackslashes)

				/* assert */
				Expect(result).To(Equal(expected))
			})
		})
		Context("when path doesn't contain a drive letter", func() {
			It("should replace backslashes with forward slashes", func() {
				/* arrange */
				pathWithMultipleBackslashes := `\\dummy\path`
				expected := `//dummy/path`

				/* act */
				actual := objectUnderTest.LocalToEngine(pathWithMultipleBackslashes)

				/* assert */
				Expect(actual).To(Equal(expected))
			})
		})
	})

	Context("when runtime.GOOS == linux", func() {
		Context("when running outside a container", func() {
			fakeRuntime := new(iruntimeFakes.FakeIRuntime)
			fakeRuntime.GOOSReturns("windows")
			objectUnderTest := _fsPathConverter{
				runtime:     fakeRuntime,
				hostRuntime: hostruntime.RuntimeInfo{},
			}

			It("should not modify path", func() {
				/* arrange */
				path := "/app"

				/* act */
				actual := objectUnderTest.LocalToEngine(path)

				/* assert */
				Expect(actual).To(Equal(path))
			})
		})

		Context("when running inside a container", func() {
			fakeRuntime := new(iruntimeFakes.FakeIRuntime)
			fakeRuntime.GOOSReturns("linux")
			dockerID, _ := uuid.NewV4()
			objectUnderTest := _fsPathConverter{
				runtime: fakeRuntime,
				hostRuntime: hostruntime.RuntimeInfo{
					InAContainer: true,
					DockerID:     dockerID.String(),
					HostPathMap: map[string]string{
						"/home/user/foo": "/foo",
					},
				},
			}

			Context("when path is mounted from the host", func() {
				It("should modify path", func() {
					/* arrange */
					given := "/foo"
					expected := "/home/user/foo"

					/* act */
					actual := objectUnderTest.LocalToEngine(given)

					/* assert */
					Expect(actual).To(Equal(expected))
				})
			})

			Context("when path is not mounted from the host", func() {
				It("should not modify path", func() {
					/* arrange */
					path := "/app"

					/* act */
					actual := objectUnderTest.LocalToEngine(path)

					/* assert */
					Expect(actual).To(Equal(path))
				})
			})
		})
	})
})
