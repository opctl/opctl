package git

import (
	"context"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Pull", func() {
	Context("parseRef errs", func() {
		It("should return error", func() {
			/* arrange */
			/* act */
			actualError := Pull(
				context.Background(),
				"dummyPath",
				"\\///%%&",
				nil,
			)

			fmt.Print(actualError.Error())

			/* assert */
			Expect(actualError).To(MatchError(`invalid git ref: parse "\\///%%&": invalid URL escape "%%&"`))
		})
	})
	Context("parseRef doesn't err", func() {
		Context("git.PlainClone errors", func() {
			Context("err.Error() returns git.ErrRepositoryAlreadyExists", func() {
				It("shouldn't error", func() {

					/* arrange */
					providedPath, err := ioutil.TempDir("", "")
					if err != nil {
						panic(err)
					}
					// some small public repo
					providedRef := "github.com/opspec-pkgs/_.op.create#3.2.0"

					/* act */
					firstErr := Pull(
						context.Background(),
						providedPath,
						providedRef,
						nil,
					)
					if nil != firstErr {
						panic(firstErr)
					}

					actualError := Pull(
						context.Background(),
						providedPath,
						providedRef,
						nil,
					)

					/* assert */
					Expect(actualError).To(BeNil())
				})
			})
			Context("err.Error() returns transport.ErrAuthenticationRequired error", func() {
				It("should return expected error", func() {

					/* arrange */
					providedPath, err := ioutil.TempDir("", "")
					if err != nil {
						panic(err)
					}

					// some small private repo
					providedRef := "github.com/Remitly/infra-ops#9.1.6"

					expectedError := model.ErrDataProviderAuthentication{}

					/* act */
					actualError := Pull(
						context.Background(),
						providedPath,
						providedRef,
						nil,
					)

					/* assert */
					Expect(actualError).To(MatchError(expectedError))
				})
			})
			Context("err.Error() returns transport.ErrAuthorizationFailed error", func() {
				It("should return expected error", func() {

					/* arrange */
					providedPath, err := ioutil.TempDir("", "")
					if err != nil {
						panic(err)
					}

					// gitlab cuz github returns 404 not 403
					providedRef := "gitlab.com/joetesterperson1/private#0.0.0"

					expectedError := model.ErrDataProviderAuthorization{}

					/* act */
					actualError := Pull(
						context.Background(),
						providedPath,
						providedRef,
						&model.Creds{
							Username: "joetesterperson",
							Password: "MWgQpun9TWUx2iFQctyJ",
						},
					)

					/* assert */
					Expect(actualError).To(MatchError(expectedError))
				})
			})
			Context("err.Error() returns other error", func() {
				It("should return error", func() {
					/* arrange */
					providedPath, err := ioutil.TempDir("", "")
					if err != nil {
						panic(err)
					}

					// non existent
					providedRef := "dummyDataRef#0.0.0"

					expectedMsg := `Get "https://dummyDataRef/info/refs?service=git-upload-pack": dial tcp: lookup dummyDataRef on 127.0.0.11:53: no such host`

					/* act */
					actualError := Pull(
						context.Background(),
						providedPath,
						providedRef,
						nil,
					)

					/* assert */
					Expect(actualError.Error()).To(Equal(expectedMsg))
				})
			})
		})
	})
})
