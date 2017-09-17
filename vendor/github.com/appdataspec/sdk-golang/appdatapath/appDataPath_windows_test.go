package appdatapath

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("appdata", func() {
	Context("Global", func() {
		Context("PROGRAMDATA env var exists", func() {
			It("should return expected path", func() {
				/* arrange */
				expectedGlobal := "dummyGlobal"

				fakeOS := new(ios.Fake)
				fakeOS.GetenvStub = func(key string) string {
					switch key {
					case `PROGRAMDATA`:
						return expectedGlobal
					default:
						return ""
					}
				}

				objectUnderTest := appDataPath{
					os: fakeOS,
				}

				/* act */
				result, _ := objectUnderTest.Global()

				/* assert */
				Expect(result).To(Equal(expectedGlobal))
			})
		})
		Context("PROGRAMDATA env var doesn't exist", func() {
			It("should panic w/ expected message", func() {
				/* arrange */
				expectedError := errors.New("unable to determine per user app data path. Error was: PROGRAMDATA env var required")

				objectUnderTest := appDataPath{
					os: new(ios.Fake),
				}

				/* act */
				_, actualError := objectUnderTest.Global()

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
	Context("PerUser", func() {
		Context("LOCALAPPDATA env var exists", func() {
			It("should return expected path", func() {
				/* arrange */
				expectedPerUser := "dummyHomeDirPath"

				fakeOS := new(ios.Fake)
				fakeOS.GetenvStub = func(key string) string {
					switch key {
					case `LOCALAPPDATA`:
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
		Context("LOCALAPPDATA env var doesn't exist", func() {
			It("should panic w/ expected message", func() {
				/* arrange */
				expectedError := errors.New("unable to determine per user app data path. Error was: LOCALAPPDATA env var required")

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
