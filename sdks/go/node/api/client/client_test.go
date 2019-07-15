package client

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Context("nodeApiClient", func() {

	Context("New()", func() {
		Context("opts nil", func() {
			It("should not return nil", func() {
				/* arrange/act/assert */
				Expect(New(url.URL{}, nil)).Should(Not(BeNil()))
			})
		})
		Context("opts not nil", func() {
			It("should not return nil", func() {
				/* arrange/act/assert */
				Expect(New(url.URL{}, &Opts{})).Should(Not(BeNil()))
			})
		})
	})

})
