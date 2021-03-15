package coerce

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("ToFile", func() {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	Context("Value is nil", func() {
		Context("ioutil.WriteFile doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedScratchDir := tmpDir

				/* act */
				actualValue, actualErr := ToFile(
					nil,
					providedScratchDir,
				)

				/* assert */
				Expect(*actualValue.File).To(HaveLen(32 + 1 + len(tmpDir)))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Array isn't nil", func() {
		Context("json.Marshal doesn't err", func() {
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := tmpDir

					/* act */
					actualValue, actualErr := ToFile(
						&model.Value{
							Array: new([]interface{}),
						},
						providedScratchDir,
					)

					/* assert */
					Expect(*actualValue.File).To(HaveLen(32 + 1 + len(tmpDir)))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
	Context("Value.Boolean isn't nil", func() {
		Context("ioutil.WriteFile doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedScratchDir := tmpDir

				/* act */
				actualValue, actualErr := ToFile(
					&model.Value{
						Boolean: new(bool),
					},
					providedScratchDir,
				)

				/* assert */
				Expect(*actualValue.File).To(HaveLen(32 + 1 + len(tmpDir)))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Dir isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScratchDir := tmpDir

			providedDir := "dummyValue"
			providedValue := &model.Value{
				Dir: &providedDir,
			}

			/* act */
			actualValue, actualErr := ToFile(providedValue, providedScratchDir)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce dir 'dummyValue' to file: incompatible types"))
		})
	})
	Context("Value.File isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScratchDir := tmpDir

			providedFile := "dummyFile"
			providedValue := &model.Value{
				File: &providedFile,
			}

			/* act */
			actualValue, actualErr := ToFile(providedValue, providedScratchDir)

			/* assert */
			Expect(actualValue).To(Equal(providedValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Number isn't nil", func() {
		Context("ioutil.WriteFile doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedScratchDir := tmpDir

				/* act */
				actualValue, actualErr := ToFile(
					&model.Value{
						Number: new(float64),
					},
					providedScratchDir,
				)

				/* assert */
				Expect(*actualValue.File).To(HaveLen(32 + 1 + len(tmpDir)))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Object isn't nil", func() {
		Context("json.Marshal doesn't err", func() {
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := tmpDir

					/* act */
					actualValue, actualErr := ToFile(
						&model.Value{
							Object: new(map[string]interface{}),
						},
						providedScratchDir,
					)

					/* assert */
					Expect(*actualValue.File).To(HaveLen(32 + 1 + len(tmpDir)))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
	Context("Value.String isn't nil", func() {
		Context("ioutil.WriteFile doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedScratchDir := tmpDir

				/* act */
				actualValue, actualErr := ToFile(
					&model.Value{
						String: new(string),
					},
					providedScratchDir,
				)

				/* assert */
				Expect(*actualValue.File).To(HaveLen(32 + 1 + len(tmpDir)))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Dir,File,Number,Object,String nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScratchDir := tmpDir

			providedValue := &model.Value{}

			/* act */
			actualValue, actualErr := ToFile(providedValue, providedScratchDir)

			/* assert */
			Expect(actualErr.Error()).To(Equal("unable to coerce '{}' to file"))
			Expect(actualValue).To(BeNil())
		})
	})
	Context("scratchDir doesn't exist", func() {
		Context("os.MkdirAll doesn't err", func() {
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := tmpDir

					/* act */
					actualValue, actualErr := ToFile(
						&model.Value{
							String: new(string),
						},
						providedScratchDir,
					)

					/* assert */
					Expect(*actualValue.File).To(HaveLen(32 + 1 + len(tmpDir)))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
