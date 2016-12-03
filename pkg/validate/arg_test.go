package validate

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "fmt"
)

var _ = Describe("Arg", func() {

  objectUnderTest := New()
  Context("invoked with non nil arg.String", func() {

    Context("equal-length to non-zero param maxLength", func() {

      It("returns no errors", func() {

        /* arrange */
        providedArg := &model.Arg{
          String:"dummyStringArg",
        }
        providedParam := &model.Param{
          String:&model.StringParam{
            MaxLength: len(providedArg.String),
          },
        }

        expectedErrors := []error{}

        /* act */
        actualErrors := objectUnderTest.Arg(providedArg, providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
    Context("longer than non-zero param maxLength", func() {

      It("returns expected errors", func() {

        /* arrange */
        providedArg := &model.Arg{
          String:"dummyStringArg",
        }
        providedParam := &model.Param{
          String:&model.StringParam{
            MaxLength: len(providedArg.String) - 1,
          },
        }

        expectedErrors := []error{
          fmt.Errorf(
            "%v must be <= %v characters",
            providedParam.String.Name,
            providedParam.String.MaxLength,
          ),
        }

        /* act */
        actualErrors := objectUnderTest.Arg(providedArg, providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
    Context("shorter than non-zero param maxLength", func() {

      It("returns no errors", func() {

        /* arrange */
        providedArg := &model.Arg{
          String:"dummyStringArg",
        }
        providedParam := &model.Param{
          String:&model.StringParam{
            MaxLength: len(providedArg.String) + 1,
          },
        }

        expectedErrors := []error{}

        /* act */
        actualErrors := objectUnderTest.Arg(providedArg, providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })

    Context("equal-length to non-zero param minLength", func() {

      It("should return no errors", func() {

        /* arrange */
        providedArg := &model.Arg{
          String:"dummyStringArg",
        }
        providedParam := &model.Param{
          String:&model.StringParam{
            MinLength: len(providedArg.String),
          },
        }

        expectedErrors := []error{}

        /* act */
        actualErrors := objectUnderTest.Arg(providedArg, providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
    Context("shorter than non-zero param minLength", func() {

      It("should return expected errors", func() {

        /* arrange */
        providedArg := &model.Arg{
          String:"dummyStringArg",
        }
        providedParam := &model.Param{
          String:&model.StringParam{
            MinLength: len(providedArg.String) + 1,
          },
        }

        expectedErrors := []error{
          fmt.Errorf(
            "%v must be >= %v characters",
            providedParam.String.Name,
            providedParam.String.MinLength,
          ),
        }

        /* act */
        actualErrors := objectUnderTest.Arg(providedArg, providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
    Context("longer than non-zero param minLength", func() {

      It("should return no errors", func() {

        /* arrange */
        providedArg := &model.Arg{
          String:"dummyStringArg",
        }
        providedParam := &model.Param{
          String:&model.StringParam{
            MinLength: len(providedArg.String) - 1,
          },
        }

        expectedErrors := []error{}

        /* act */
        actualErrors := objectUnderTest.Arg(providedArg, providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
    Context("not matching non-empty param pattern", func() {

      It("should return expected errors", func() {

        /* arrange */
        providedArg := &model.Arg{
          String:"dummyStringArg",
        }
        providedParam := &model.Param{
          String:&model.StringParam{
            Pattern: "^$",
          },
        }

        expectedErrors := []error{
          fmt.Errorf(
            "%v must match pattern %v",
            providedParam.String.Name,
            providedParam.String.Pattern,
          ),
        }

        /* act */
        actualErrors := objectUnderTest.Arg(providedArg, providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
    Context("matching non-empty param pattern", func() {

      It("should return no errors", func() {

        /* arrange */
        providedArg := &model.Arg{
          String:"dummyStringArg",
        }
        providedParam := &model.Param{
          String:&model.StringParam{
            Pattern: ".*",
          },
        }

        expectedErrors := []error{}

        /* act */
        actualErrors := objectUnderTest.Arg(providedArg, providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
  })
})
