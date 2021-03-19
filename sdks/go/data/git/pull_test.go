package git

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

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
					if firstErr != nil {
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
					testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusUnauthorized)
					}))
					defer testServer.Close()

					// ignore unknown certificate signatory in mock tls server
					http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
					defer func() {
						http.DefaultTransport.(*http.Transport).TLSClientConfig = nil
					}()

					providedRef := fmt.Sprintf("%s#version", testServer.URL)

					providedPath, err := ioutil.TempDir("", "")
					if err != nil {
						panic(err)
					}

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
					testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusForbidden)
					}))
					defer testServer.Close()

					// ignore unknown certificate signatory in mock tls server
					http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
					defer func() {
						http.DefaultTransport.(*http.Transport).TLSClientConfig = nil
					}()

					providedRef := fmt.Sprintf("%s#version", testServer.URL)

					providedPath, err := ioutil.TempDir("", "")
					if err != nil {
						panic(err)
					}

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
					testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
					defer testServer.Close()

					// ignore unknown certificate signatory in mock tls server
					http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
					defer func() {
						http.DefaultTransport.(*http.Transport).TLSClientConfig = nil
					}()

					providedRef := fmt.Sprintf("%s#version", testServer.URL)

					providedPath, err := ioutil.TempDir("", "")
					if err != nil {
						panic(err)
					}

					/* act */
					actualError := Pull(
						context.Background(),
						providedPath,
						providedRef,
						nil,
					)

					fmt.Println(actualError.Error())

					/* assert */
					Expect(actualError).To(MatchError(fmt.Sprintf(`unexpected client error: unexpected requesting "%s/info/refs?service=git-upload-pack" status code: 500`, testServer.URL)))
				})
			})
		})
	})
})
