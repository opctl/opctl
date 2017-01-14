package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("Data", func() {

	Context("when formatting to/from json", func() {
		json := format.NewJsonFormat()

		Context("with non-nil $.dir", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedData := Data{
					Dir: "dummyDirRef",
				}

				/* act */
				providedJson, err := json.From(expectedData)
				if nil != err {
					panic(err)
				}

				actualData := Data{}
				json.To(providedJson, &actualData)

				/* assert */
				Expect(actualData).To(Equal(expectedData))

			})

		})

		Context("with non-nil $.file", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedData := Data{
					File: "dummyFileRef",
				}

				/* act */
				providedJson, err := json.From(expectedData)
				if nil != err {
					panic(err)
				}

				actualData := Data{}
				json.To(providedJson, &actualData)

				/* assert */
				Expect(actualData).To(Equal(expectedData))

			})

		})

		Context("with non-nil $.netSocket", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedData := Data{
					NetSocket: &NetSocket{
						Host: "dummyName",
						Port: 1,
					},
				}

				/* act */
				providedJson, err := json.From(expectedData)
				if nil != err {
					panic(err)
				}

				actualData := Data{}
				json.To(providedJson, &actualData)

				/* assert */
				Expect(actualData).To(Equal(expectedData))

			})

		})

		Context("with non-nil $.string", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedData := Data{
					String: "dummyString",
				}

				/* act */
				providedJson, err := json.From(expectedData)
				if nil != err {
					panic(err)
				}

				actualData := Data{}
				json.To(providedJson, &actualData)

				/* assert */
				Expect(actualData).To(Equal(expectedData))

			})

		})

	})

})
