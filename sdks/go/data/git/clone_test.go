package git

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Clone", func() {
	Context("git.PlainClone errors", func() {
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

				u, err := url.Parse(testServer.URL)
				if err != nil {
					panic(err)
				}

				providedRef := &ref{
					Name: u.Host,
				}

				providedPath, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}

				expectedError := model.ErrDataProviderAuthentication{}

				/* act */
				actualError := Clone(
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

				u, err := url.Parse(testServer.URL)
				if err != nil {
					panic(err)
				}

				providedRef := &ref{
					Name: u.Host,
				}

				providedPath, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}

				expectedError := model.ErrDataProviderAuthorization{}

				/* act */
				actualError := Clone(
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

				u, err := url.Parse(testServer.URL)
				if err != nil {
					panic(err)
				}

				providedRef := &ref{
					Name: u.Host,
				}

				providedPath, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}

				/* act */
				actualError := Clone(
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
