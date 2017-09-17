package appdatapath

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("appdata", func() {
	Context("Global", func() {
		It("should return expected path", func() {
			/* arrange */
			expected := "/var/lib"

			objectUnderTest := appDataPath{
				os: new(ios.Fake),
			}

			/* act */
			result, _ := objectUnderTest.Global()

			/* assert */
			Expect(result).To(Equal(expected))
		})
	})
	Context("PerUser", func() {
		Context("HOME env var exists", func() {
			It("should return expected path", func() {
				/* arrange */
				expectedPerUser := "dummyHomeDirPath"

				fakeOS := new(ios.Fake)
				fakeOS.GetenvStub = func(key string) string {
					switch key {
					case `HOME`:
						return expectedPerUser
					default:
						return ""
					}
				}

				objectUnderTest := appDataPath{
					os: fakeOS,
				}

				/* act */
				result, _ := objectUnderTest.PerUser()

				/* assert */
				Expect(result).To(Equal(expectedPerUser))
			})
		})
		Context("HOME env var doesn't exist", func() {
			It("should panic w/ expected message", func() {
				/* arrange */
				expectedError := errors.New("unable to determine per user app data path. Error was: HOME env var required")

				objectUnderTest := appDataPath{
					os: new(ios.Fake),
				}

				/* act */
				_, actualError := objectUnderTest.PerUser()

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
