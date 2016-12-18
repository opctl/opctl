package validate

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "fmt"
  "errors"
)

var _ = Describe("Param", func() {

  objectUnderTest := New()
  Context("invoked w/ non nil param.String", func() {
    Context("w/ non-empty Value", func() {
      Context("equal-length to non-zero param.Constraints.Length.Max", func() {

        It("returns no errors", func() {

          /* arrange */
          providedValue := "dummyValue"
          providedParam := &model.Param{
            String:&model.StringParam{
              Constraints: &model.StringConstraints{
                Length:&model.StringLengthConstraint{
                  Max:len(providedValue),
                },
              },
              Value: providedValue,
            },
          }

          expectedErrors := []error{}

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
      Context("longer than non-zero param.Constraints.Length.Max", func() {

        It("returns expected errors", func() {

          /* arrange */
          providedValue := "dummyValue"
          providedParam := &model.Param{
            String:&model.StringParam{
              Constraints: &model.StringConstraints{
                Length:&model.StringLengthConstraint{
                  Max:len(providedValue) - 1,
                },
              },
              Value: providedValue,
            },
          }

          expectedErrors := []error{
            fmt.Errorf(
              "%v must be <= %v characters",
              providedParam.String.Name,
              providedParam.String.Constraints.Length.Max,
            ),
          }

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
      Context("shorter than non-zero param.Constraints.Length.Max", func() {

        It("returns no errors", func() {

          /* arrange */
          providedValue := "dummyValue"
          providedParam := &model.Param{
            String:&model.StringParam{
              Constraints: &model.StringConstraints{
                Length:&model.StringLengthConstraint{
                  Max:len(providedValue) + 1,
                },
              },
              Value: providedValue,
            },
          }

          expectedErrors := []error{}

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
      Context("equal-length to non-zero param.Constraints.Length.Min", func() {

        It("should return no errors", func() {

          /* arrange */
          providedValue := "dummyValue"
          providedParam := &model.Param{
            String:&model.StringParam{
              Constraints: &model.StringConstraints{
                Length:&model.StringLengthConstraint{
                  Min:len(providedValue),
                },
              },
              Value: providedValue,
            },
          }

          expectedErrors := []error{}

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
      Context("shorter than non-zero param.Constraints.Length.Min", func() {

        It("should return expected errors", func() {

          /* arrange */
          providedValue := "dummyValue"
          providedParam := &model.Param{
            String:&model.StringParam{
              Constraints: &model.StringConstraints{
                Length:&model.StringLengthConstraint{
                  Min:len(providedValue) + 1,
                },
              },
              Value: providedValue,
            },
          }

          expectedErrors := []error{
            fmt.Errorf(
              "%v must be >= %v characters",
              providedParam.String.Name,
              providedParam.String.Constraints.Length.Min,
            ),
          }

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
      Context("longer than non-zero param.Constraints.Length.Min", func() {

        It("should return no errors", func() {

          /* arrange */
          providedValue := "dummyValue"
          providedParam := &model.Param{
            String:&model.StringParam{
              Constraints: &model.StringConstraints{
                Length:&model.StringLengthConstraint{
                  Min:len(providedValue) - 1,
                },
              },
              Value: providedValue,
            },
          }

          expectedErrors := []error{}

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
      Context("not matching non-empty param.Constraints.Patterns", func() {

        It("should return expected errors", func() {

          /* arrange */
          providedValue := "dummyValue"
          providedParam := &model.Param{
            String:&model.StringParam{
              Constraints: &model.StringConstraints{
                Patterns:[]*model.StringPatternConstraint{
                  {
                    Regex:"^$",
                  },
                },
              },
              Value: providedValue,
            },
          }

          expectedErrors := []error{
            fmt.Errorf(
              "%v must match pattern %v",
              providedParam.String.Name,
              providedParam.String.Constraints.Patterns[0].Regex,
            ),
          }

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
      Context("matching non-empty param.Constraints.Patterns", func() {

        It("should return no errors", func() {

          /* arrange */
          providedValue := "dummyValue"
          providedParam := &model.Param{
            String:&model.StringParam{
              Constraints: &model.StringConstraints{
                Patterns:[]*model.StringPatternConstraint{
                  {
                    Regex:".$",
                  },
                },
              },
              Value: providedValue,
            },
          }

          expectedErrors := []error{}

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
    })
    Context("w/ empty Value", func() {
      Context("and non empty Default", func() {
        Context("equal-length to non-zero param.Constraints.Length.Max", func() {

          It("returns no errors", func() {

            /* arrange */
            providedDefault := "dummyDefault"
            providedParam := &model.Param{
              String:&model.StringParam{
                Constraints: &model.StringConstraints{
                  Length:&model.StringLengthConstraint{
                    Max:len(providedDefault),
                  },
                },
                Default: providedDefault,
              },
            }

            expectedErrors := []error{}

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

          })
        })
        Context("longer than non-zero param.Constraints.Length.Max", func() {

          It("returns expected errors", func() {

            /* arrange */
            providedDefault := "dummyDefault"
            providedParam := &model.Param{
              String:&model.StringParam{
                Constraints: &model.StringConstraints{
                  Length:&model.StringLengthConstraint{
                    Max:len(providedDefault) - 1,
                  },
                },
                Default: providedDefault,
              },
            }

            expectedErrors := []error{
              fmt.Errorf(
                "%v must be <= %v characters",
                providedParam.String.Name,
                providedParam.String.Constraints.Length.Max,
              ),
            }

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

          })
        })
        Context("shorter than non-zero param.Constraints.Length.Max", func() {

          It("returns no errors", func() {

            /* arrange */
            providedDefault := "dummyDefault"
            providedParam := &model.Param{
              String:&model.StringParam{
                Constraints: &model.StringConstraints{
                  Length:&model.StringLengthConstraint{
                    Max:len(providedDefault) + 1,
                  },
                },
                Default: providedDefault,
              },
            }

            expectedErrors := []error{}

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

          })
        })
        Context("equal-length to non-zero param.Constraints.Length.Min", func() {

          It("should return no errors", func() {

            /* arrange */
            providedDefault := "dummyDefault"
            providedParam := &model.Param{
              String:&model.StringParam{
                Constraints: &model.StringConstraints{
                  Length:&model.StringLengthConstraint{
                    Min:len(providedDefault),
                  },
                },
                Default: providedDefault,
              },
            }

            expectedErrors := []error{}

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

          })
        })
        Context("shorter than non-zero param.Constraints.Length.Min", func() {

          It("should return expected errors", func() {

            /* arrange */
            providedDefault := "dummyDefault"
            providedParam := &model.Param{
              String:&model.StringParam{
                Constraints: &model.StringConstraints{
                  Length:&model.StringLengthConstraint{
                    Min:len(providedDefault) + 1,
                  },
                },
                Default: providedDefault,
              },
            }

            expectedErrors := []error{
              fmt.Errorf(
                "%v must be >= %v characters",
                providedParam.String.Name,
                providedParam.String.Constraints.Length.Min,
              ),
            }

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

          })
        })
        Context("longer than non-zero param.Constraints.Length.Min", func() {

          It("should return no errors", func() {

            /* arrange */
            providedDefault := "dummyDefault"
            providedParam := &model.Param{
              String:&model.StringParam{
                Constraints: &model.StringConstraints{
                  Length:&model.StringLengthConstraint{
                    Min:len(providedDefault) - 1,
                  },
                },
                Default: providedDefault,
              },
            }

            expectedErrors := []error{}

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

          })
        })
        Context("not matching non-empty param.Constraints.Patterns", func() {

          It("should return expected errors", func() {

            /* arrange */
            providedDefault := "dummyDefault"
            providedParam := &model.Param{
              String:&model.StringParam{
                Constraints: &model.StringConstraints{
                  Patterns:[]*model.StringPatternConstraint{
                    {
                      Regex:"^$",
                    },
                  },
                },
                Default: providedDefault,
              },
            }

            expectedErrors := []error{
              fmt.Errorf(
                "%v must match pattern %v",
                providedParam.String.Name,
                providedParam.String.Constraints.Patterns[0].Regex,
              ),
            }

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

          })
        })
        Context("matching non-empty param.Constraints.Patterns", func() {

          It("should return no errors", func() {

            /* arrange */
            providedDefault := "dummyDefault"
            providedParam := &model.Param{
              String:&model.StringParam{
                Constraints: &model.StringConstraints{
                  Patterns:[]*model.StringPatternConstraint{
                    {
                      Regex:".$",
                    },
                  },
                },
                Default: providedDefault,
              },
            }

            expectedErrors := []error{}

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

          })
        })
      })
    })
  })
  Context("invoked w/ non nil param.NetSocket", func() {
    Context("w/ non-empty Value.Port", func() {
      Context("equal to 0", func() {
        It("should return expected errors", func() {

          /* arrange */
          providedValue := &model.NetSocketParamValue{
            Host:"dummyHost",
            Port: 0,
          }
          providedParam := &model.Param{
            NetSocket:&model.NetSocketParam{
              Value: providedValue,
            },
          }

          expectedErrors := []error{
            errors.New("Port must be > 0"),
          }

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
      Context("between 0 and 65536", func() {
        It("should return no errors", func() {
          var i uint = 1
          for i < 65536 {

            /* arrange */
            providedValue := &model.NetSocketParamValue{
              Host:"dummyHost",
              Port: i,
            }
            providedParam := &model.Param{
              NetSocket:&model.NetSocketParam{
                Value: providedValue,
              },
            }

            expectedErrors := []error{}

            /* act */
            actualErrors := objectUnderTest.Param(providedParam)

            /* assert */
            Expect(actualErrors).To(Equal(expectedErrors))

            i += i
          }

        })
      })
      Context("equal to 65536", func() {
        It("should return expected errors", func() {

          /* arrange */
          providedValue := &model.NetSocketParamValue{
            Host:"dummyHost",
            Port: 65536,
          }
          providedParam := &model.Param{
            NetSocket:&model.NetSocketParam{
              Value: providedValue,
            },
          }

          expectedErrors := []error{
            errors.New("Port must be <= 65535"),
          }

          /* act */
          actualErrors := objectUnderTest.Param(providedParam)

          /* assert */
          Expect(actualErrors).To(Equal(expectedErrors))

        })
      })
    })
    Context("w/ empty Value.Port", func() {
      It("should return expected errors", func() {

        /* arrange */
        providedValue := &model.NetSocketParamValue{
          Host:"dummyHost",
        }
        providedParam := &model.Param{
          NetSocket:&model.NetSocketParam{
            Value: providedValue,
          },
        }

        expectedErrors := []error{
          errors.New("Port must be > 0"),
        }

        /* act */
        actualErrors := objectUnderTest.Param(providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
    Context("w/ empty Value.Host", func() {
      It("should return expected errors", func() {

        /* arrange */
        providedValue := &model.NetSocketParamValue{
          Port: 80,
        }
        providedParam := &model.Param{
          NetSocket:&model.NetSocketParam{
            Value: providedValue,
          },
        }

        expectedErrors := []error{
          errors.New("Host required"),
        }

        /* act */
        actualErrors := objectUnderTest.Param(providedParam)

        /* assert */
        Expect(actualErrors).To(Equal(expectedErrors))

      })
    })
  })
})
