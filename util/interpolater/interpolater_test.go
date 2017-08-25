package interpolater

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Context("Interpolate", func() {
	Context("no refs", func() {
		Context("single (", func() {
			Context("stand alone", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "("
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at start", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "(suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at end", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix("
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("within", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix(suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
		})
		Context("multiple (s", func() {
			Context("stand alone", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "(("
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at start", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "((suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at end", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix(("
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("within", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix((suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
		})
		Context("single )", func() {
			Context("stand alone", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := ")"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at start", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := ")suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at end", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix)"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("within", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix)suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
		})
		Context("multiple )s", func() {
			Context("stand alone", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "))"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at start", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "))suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at end", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix))"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("within", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix))suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
		})
		Context("single $(", func() {
			Context("stand alone", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "$("
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at start", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "$(dummyExpression"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at end", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "dummyExpression$("
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("within", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "dummy$(Expression"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
		})
		Context("multiple $(s", func() {
			Context("stand alone", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "$($("
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at start", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "$($(suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("at end", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix$($("
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
			Context("within", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "prefix$(sep$(suffix"
					objectUnderTest := _Interpolater{}

					/* act */
					actualResult, _ := objectUnderTest.Interpolate(
						providedExpression,
						nil,
					)

					/* assert */
					Expect(actualResult).To(Equal(providedExpression))
				})
			})
		})
	})
	Context("ref", func() {
		Context("stand alone", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedExpression := fmt.Sprintf("$(%v)", providedRef)
				objectUnderTest := _Interpolater{}

				refdValue := "dummyRefdValue"

				fakeDeReferencer := new(FakeDeReferencer)
				fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
					switch ref {
					case providedRef:
						return refdValue, true, nil
					default:
						return "", false, nil
					}
				}

				expectedValue := strings.Replace(
					providedExpression,
					fmt.Sprintf("$(%v)", providedRef),
					refdValue,
					-1,
				)

				/* act */
				actualResult, _ := objectUnderTest.Interpolate(
					providedExpression,
					fakeDeReferencer,
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedValue))
			})
		})
		Context("at start", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedExpression := fmt.Sprintf("$(%v)suffix", providedRef)
				objectUnderTest := _Interpolater{}

				refdValue := "dummyRefdValue"

				fakeDeReferencer := new(FakeDeReferencer)
				fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
					switch ref {
					case providedRef:
						return refdValue, true, nil
					default:
						return "", false, nil
					}
				}

				expectedValue := strings.Replace(
					providedExpression,
					fmt.Sprintf("$(%v)", providedRef),
					refdValue,
					-1,
				)

				/* act */
				actualResult, _ := objectUnderTest.Interpolate(
					providedExpression,
					fakeDeReferencer,
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedValue))
			})
		})
		Context("at end", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedExpression := fmt.Sprintf("prefix$(%v)", providedRef)
				objectUnderTest := _Interpolater{}

				refdValue := "dummyRefdValue"

				fakeDeReferencer := new(FakeDeReferencer)
				fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
					switch ref {
					case providedRef:
						return refdValue, true, nil
					default:
						return "", false, nil
					}
				}

				expectedValue := strings.Replace(
					providedExpression,
					fmt.Sprintf("$(%v)", providedRef),
					refdValue,
					-1,
				)

				/* act */
				actualResult, _ := objectUnderTest.Interpolate(
					providedExpression,
					fakeDeReferencer,
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedValue))
			})
		})
		Context("within", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedExpression := fmt.Sprintf("dummy$(%v)Expression", providedRef)
				objectUnderTest := _Interpolater{}

				refdValue := "dummyRefdValue"

				fakeDeReferencer := new(FakeDeReferencer)
				fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
					switch ref {
					case providedRef:
						return refdValue, true, nil
					default:
						return "", false, nil
					}
				}

				expectedValue := strings.Replace(
					providedExpression,
					fmt.Sprintf("$(%v)", providedRef),
					refdValue,
					-1,
				)

				/* act */
				actualResult, _ := objectUnderTest.Interpolate(
					providedExpression,
					fakeDeReferencer,
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedValue))
			})
		})
	})
	Context("refs", func() {
		Context("stand alone", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef1 := "dummyRef1"
				providedRef2 := "dummyRef2"
				providedExpression := fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2)
				objectUnderTest := _Interpolater{}

				refdValue1 := "dummyRefdValue1"
				refdValue2 := "dummyRefdValue2"

				fakeDeReferencer := new(FakeDeReferencer)
				fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
					switch ref {
					case providedRef1:
						return refdValue1, true, nil
					case providedRef2:
						return refdValue2, true, nil
					default:
						return "", false, nil
					}
				}

				expectedValue := strings.Replace(
					providedExpression,
					fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2),
					refdValue1+refdValue2,
					-1,
				)

				/* act */
				actualResult, _ := objectUnderTest.Interpolate(
					providedExpression,
					fakeDeReferencer,
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedValue))
			})
		})
		Context("at start", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef1 := "dummyRef1"
				providedRef2 := "dummyRef2"
				providedExpression := fmt.Sprintf("$(%v)$(%v)suffix", providedRef1, providedRef2)
				objectUnderTest := _Interpolater{}

				refdValue1 := "dummyRefdValue1"
				refdValue2 := "dummyRefdValue2"

				fakeDeReferencer := new(FakeDeReferencer)
				fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
					switch ref {
					case providedRef1:
						return refdValue1, true, nil
					case providedRef2:
						return refdValue2, true, nil
					default:
						return "", false, nil
					}
				}

				expectedValue := strings.Replace(
					providedExpression,
					fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2),
					refdValue1+refdValue2,
					-1,
				)

				/* act */
				actualResult, _ := objectUnderTest.Interpolate(
					providedExpression,
					fakeDeReferencer,
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedValue))
			})
		})
		Context("at end", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef1 := "dummyRef1"
				providedRef2 := "dummyRef2"
				providedExpression := fmt.Sprintf("prefix$(%v)$(%v)", providedRef1, providedRef2)
				objectUnderTest := _Interpolater{}

				refdValue1 := "dummyRefdValue1"
				refdValue2 := "dummyRefdValue2"

				fakeDeReferencer := new(FakeDeReferencer)
				fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
					switch ref {
					case providedRef1:
						return refdValue1, true, nil
					case providedRef2:
						return refdValue2, true, nil
					default:
						return "", false, nil
					}
				}

				expectedValue := strings.Replace(
					providedExpression,
					fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2),
					refdValue1+refdValue2,
					-1,
				)

				/* act */
				actualResult, _ := objectUnderTest.Interpolate(
					providedExpression,
					fakeDeReferencer,
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedValue))
			})
		})
		Context("within", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef1 := "dummyRef1"
				providedRef2 := "dummyRef2"
				providedExpression := fmt.Sprintf("prefix$(%v)$(%v)suffix", providedRef1, providedRef2)
				objectUnderTest := _Interpolater{}

				refdValue1 := "dummyRefdValue1"
				refdValue2 := "dummyRefdValue2"

				fakeDeReferencer := new(FakeDeReferencer)
				fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
					switch ref {
					case providedRef1:
						return refdValue1, true, nil
					case providedRef2:
						return refdValue2, true, nil
					default:
						return "", false, nil
					}
				}

				expectedValue := strings.Replace(
					providedExpression,
					fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2),
					refdValue1+refdValue2,
					-1,
				)

				/* act */
				actualResult, _ := objectUnderTest.Interpolate(
					providedExpression,
					fakeDeReferencer,
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedValue))
			})
		})
	})
	Context("nested ref", func() {
		Context("outer ref stand alone", func() {
      Context("inner ref stand alone", func() {
        It("should return expected result", func() {
          /* arrange */
          providedRef := "dummyRef"
          providedExpression := fmt.Sprintf("$($(%v))", providedRef)
          objectUnderTest := _Interpolater{}

          innerRefdValue := "innerRefdValue"
          outerRefdValue := "outerRefdValue"

          fakeDeReferencer := new(FakeDeReferencer)
          fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
            switch ref {
            case providedRef:
              return innerRefdValue, true, nil
            case innerRefdValue:
              return outerRefdValue, true, nil
            default:
              return "", false, nil
            }
          }

          expectedValue := strings.Replace(
            providedExpression,
            fmt.Sprintf("$($(%v))", providedRef),
            outerRefdValue,
            -1,
          )

          /* act */
          actualResult, _ := objectUnderTest.Interpolate(
            providedExpression,
            fakeDeReferencer,
          )

          /* assert */
          Expect(actualResult).To(Equal(expectedValue))
        })
        Context("outer ref isn't in scope", func() {
          It("should return expected result", func() {
            /* arrange */
            providedRef := "dummyRef"
            providedExpression := fmt.Sprintf("$($(%v))", providedRef)
            objectUnderTest := _Interpolater{}

            innerRefdValue := "innerRefdValue"

            fakeDeReferencer := new(FakeDeReferencer)
            fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
              switch ref {
              case providedRef:
                return innerRefdValue, true, nil
              default:
                return "", false, nil
              }
            }

            expectedValue := strings.Replace(
              providedExpression,
              fmt.Sprintf("$(%v)", providedRef),
              innerRefdValue,
              -1,
            )

            /* act */
            actualResult, _ := objectUnderTest.Interpolate(
              providedExpression,
              fakeDeReferencer,
            )

            /* assert */
            Expect(actualResult).To(Equal(expectedValue))
          })
        })
      })
		})
		Context("outer ref at start", func() {
      Context("inner ref stand alone", func() {
        It("should return expected result", func() {
          /* arrange */
          providedRef := "dummyRef"
          providedExpression := fmt.Sprintf("$($(%v))suffix", providedRef)
          objectUnderTest := _Interpolater{}

          innerRefdValue := "innerRefdValue"
          outerRefdValue := "outerRefdValue"

          fakeDeReferencer := new(FakeDeReferencer)
          fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
            switch ref {
            case providedRef:
              return innerRefdValue, true, nil
            case innerRefdValue:
              return outerRefdValue, true, nil
            default:
              return "", false, nil
            }
          }

          expectedValue := strings.Replace(
            providedExpression,
            fmt.Sprintf("$($(%v))", providedRef),
            outerRefdValue,
            -1,
          )

          /* act */
          actualResult, _ := objectUnderTest.Interpolate(
            providedExpression,
            fakeDeReferencer,
          )

          /* assert */
          Expect(actualResult).To(Equal(expectedValue))
        })
        Context("outer ref isn't in scope", func() {
          It("should return expected result", func() {
            /* arrange */
            providedRef := "dummyRef"
            providedExpression := fmt.Sprintf("$($(%v))suffix", providedRef)
            objectUnderTest := _Interpolater{}

            innerRefdValue := "innerRefdValue"

            fakeDeReferencer := new(FakeDeReferencer)
            fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
              switch ref {
              case providedRef:
                return innerRefdValue, true, nil
              default:
                return "", false, nil
              }
            }

            expectedValue := strings.Replace(
              providedExpression,
              fmt.Sprintf("$(%v)", providedRef),
              innerRefdValue,
              -1,
            )

            /* act */
            actualResult, _ := objectUnderTest.Interpolate(
              providedExpression,
              fakeDeReferencer,
            )

            /* assert */
            Expect(actualResult).To(Equal(expectedValue))
          })
        })
      })
		})
		Context("outer ref at end", func() {
      Context("inner ref stand alone", func() {
        It("should return expected result", func() {
          /* arrange */
          providedRef := "dummyRef"
          providedExpression := fmt.Sprintf("prefix$($(%v))", providedRef)
          objectUnderTest := _Interpolater{}

          innerRefdValue := "innerRefdValue"
          outerRefdValue := "outerRefdValue"

          fakeDeReferencer := new(FakeDeReferencer)
          fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
            switch ref {
            case providedRef:
              return innerRefdValue, true, nil
            case innerRefdValue:
              return outerRefdValue, true, nil
            default:
              return "", false, nil
            }
          }

          expectedValue := strings.Replace(
            providedExpression,
            fmt.Sprintf("$($(%v))", providedRef),
            outerRefdValue,
            -1,
          )

          /* act */
          actualResult, _ := objectUnderTest.Interpolate(
            providedExpression,
            fakeDeReferencer,
          )

          /* assert */
          Expect(actualResult).To(Equal(expectedValue))
        })
        Context("outer ref isn't in scope", func() {
          It("should return expected result", func() {
            /* arrange */
            providedRef := "dummyRef"
            providedExpression := fmt.Sprintf("prefix$($(%v))", providedRef)
            objectUnderTest := _Interpolater{}

            innerRefdValue := "innerRefdValue"

            fakeDeReferencer := new(FakeDeReferencer)
            fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
              switch ref {
              case providedRef:
                return innerRefdValue, true, nil
              default:
                return "", false, nil
              }
            }

            expectedValue := strings.Replace(
              providedExpression,
              fmt.Sprintf("$(%v)", providedRef),
              innerRefdValue,
              -1,
            )

            /* act */
            actualResult, _ := objectUnderTest.Interpolate(
              providedExpression,
              fakeDeReferencer,
            )

            /* assert */
            Expect(actualResult).To(Equal(expectedValue))
          })
        })
      })
		})
		Context("outer ref within", func() {
      Context("inner ref stand alone", func(){
        It("should return expected result", func() {
          /* arrange */
          providedRef := "dummyRef"
          providedExpression := fmt.Sprintf("prefix$($(%v))suffix", providedRef)
          objectUnderTest := _Interpolater{}

          innerRefdValue := "innerRefdValue"
          outerRefdValue := "outerRefdValue"

          fakeDeReferencer := new(FakeDeReferencer)
          fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
            switch ref {
            case providedRef:
              return innerRefdValue, true, nil
            case innerRefdValue:
              return outerRefdValue, true, nil
            default:
              return "", false, nil
            }
          }

          expectedValue := strings.Replace(
            providedExpression,
            fmt.Sprintf("$($(%v))", providedRef),
            outerRefdValue,
            -1,
          )

          /* act */
          actualResult, _ := objectUnderTest.Interpolate(
            providedExpression,
            fakeDeReferencer,
          )

          /* assert */
          Expect(actualResult).To(Equal(expectedValue))
        })
        Context("outer ref isn't in scope", func() {
          It("should return expected result", func() {
            /* arrange */
            providedRef := "dummyRef"
            providedExpression := fmt.Sprintf("prefix$($(%v))suffix", providedRef)
            objectUnderTest := _Interpolater{}

            innerRefdValue := "innerRefdValue"

            fakeDeReferencer := new(FakeDeReferencer)
            fakeDeReferencer.DeReferenceStub = func(ref string) (string, bool, error) {
              switch ref {
              case providedRef:
                return innerRefdValue, true, nil
              default:
                return "", false, nil
              }
            }

            expectedValue := strings.Replace(
              providedExpression,
              fmt.Sprintf("$(%v)", providedRef),
              innerRefdValue,
              -1,
            )

            /* act */
            actualResult, _ := objectUnderTest.Interpolate(
              providedExpression,
              fakeDeReferencer,
            )

            /* assert */
            Expect(actualResult).To(Equal(expectedValue))
          })
        })
      })
		})
	})
})
