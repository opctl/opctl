package string

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
)

var _ = Context("Interpret", func() {
	Context("expression is ref", func() {
		Context("ref isn't in scope", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedExpression := fmt.Sprintf("$(%v)", providedRef)

				objectUnderTest := _interpreter{}

				/* act */
				actualString, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedExpression,
				)

				/* assert */
				Expect(actualString).To(BeEmpty())
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to interpret string; %v not in scope", providedRef)))
			})
		})
		Context("ref is in scope", func() {
			It("should call coercer.Coerce w/ expected args", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedExpression := fmt.Sprintf("$(%v)", providedRef)

				providedScopeValue := &model.Value{}
				providedScope := map[string]*model.Value{
					providedRef: providedScopeValue,
				}

				fakeCoercer := new(fakeCoercer)
				// err to trigger immediate return
				fakeCoercer.CoerceReturns("", errors.New("dummyError"))

				objectUnderTest := _interpreter{
					coercer: fakeCoercer,
				}

				/* act */
				objectUnderTest.Interpret(providedScope, providedExpression)

				/* assert */
				Expect(fakeCoercer.CoerceArgsForCall(0)).To(Equal(providedScopeValue))
			})
			Context("coercer.Coerce errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef := "dummyRef"
					providedExpression := fmt.Sprintf("$(%v)", providedRef)

					fakeCoercer := new(fakeCoercer)

					coerceError := errors.New("dummyError")
					fakeCoercer.CoerceReturns("", coerceError)

					objectUnderTest := _interpreter{
						coercer: fakeCoercer,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{
							providedRef: nil,
						},
						providedExpression,
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("Unable to interpret string; error was: %v", coerceError.Error())))
				})
			})
			Context("coercer.Coerce doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef := "dummyRef"
					providedExpression := fmt.Sprintf("$(%v)", providedRef)

					fakeCoercer := new(fakeCoercer)

					coercedString := "dummyString"
					fakeCoercer.CoerceReturns(coercedString, nil)

					objectUnderTest := _interpreter{
						coercer: fakeCoercer,
					}

					/* act */
					actualString, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{
							providedRef: nil,
						},
						providedExpression,
					)

					/* assert */
					Expect(actualString).To(Equal(coercedString))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
	Context("expression is refs", func() {
		Context("first ref not in scope", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef1 := "dummyRef1"
				providedExpression := fmt.Sprintf("$(%v)$(dummyRef2)", providedRef1)

				objectUnderTest := _interpreter{}

				/* act */
				actualString, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedExpression,
				)

				/* assert */
				Expect(actualString).To(BeEmpty())
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to interpret string; %v not in scope", providedRef1)))
			})
		})
		Context("second ref not in scope", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef1 := "dummyRef1"
				providedRef2 := "dummyRef2"
				providedExpression := fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2)

				fakeCoercer := new(fakeCoercer)

				objectUnderTest := _interpreter{
					coercer: fakeCoercer,
				}

				/* act */
				actualString, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{
						providedRef1: nil,
					},
					providedExpression,
				)

				/* assert */
				Expect(actualString).To(BeEmpty())
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to interpret string; %v not in scope", providedRef2)))
			})
		})
		Context("refs in scope", func() {
			It("should call coercer.Coerce w/ expected args", func() {
				/* arrange */
				providedRef1 := "dummyRef1"
				providedRef2 := "dummyRef2"
				providedExpression := fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2)

				providedScopeValue1 := &model.Value{}
				providedScopeValue2 := &model.Value{}
				providedScope := map[string]*model.Value{
					providedRef1: providedScopeValue1,
					providedRef2: providedScopeValue2,
				}

				fakeCoercer := new(fakeCoercer)

				objectUnderTest := _interpreter{
					coercer: fakeCoercer,
				}

				/* act */
				objectUnderTest.Interpret(providedScope, providedExpression)

				/* assert */
				Expect(fakeCoercer.CoerceArgsForCall(0)).To(Equal(providedScopeValue1))
				Expect(fakeCoercer.CoerceArgsForCall(1)).To(Equal(providedScopeValue2))
			})
			Context("first coercer.Coerce errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef1 := "dummyRef1"
					providedExpression := fmt.Sprintf("$(%v)$(dummyRef2)", providedRef1)

					fakeCoercer := new(fakeCoercer)

					coerceError := errors.New("dummyError")
					fakeCoercer.CoerceReturns("", coerceError)

					objectUnderTest := _interpreter{
						coercer: fakeCoercer,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{
							providedRef1: nil,
						},
						providedExpression,
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("Unable to interpret string; error was: %v", coerceError.Error())))
				})
			})
			Context("second coercer.Coerce errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef1 := "dummyRef1"
					providedRef2 := "dummyRef2"
					providedExpression := fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2)

					fakeCoercer := new(fakeCoercer)

					coerceError := errors.New("dummyError")
					fakeCoercer.CoerceReturnsOnCall(1, "", coerceError)

					objectUnderTest := _interpreter{
						coercer: fakeCoercer,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{
							providedRef1: nil,
							providedRef2: nil,
						},
						providedExpression,
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("Unable to interpret string; error was: %v", coerceError.Error())))
				})
			})
			Context("coercer.Coerce doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef1 := "dummyRef1"
					providedRef2 := "dummyRef2"
					providedExpression := fmt.Sprintf("$(%v)$(%v)", providedRef1, providedRef2)

					fakeCoercer := new(fakeCoercer)

					coercedString1 := "dummyString1"
					fakeCoercer.CoerceReturnsOnCall(0, coercedString1, nil)
					coercedString2 := "dummyString2"
					fakeCoercer.CoerceReturnsOnCall(1, coercedString2, nil)

					objectUnderTest := _interpreter{
						coercer: fakeCoercer,
					}

					/* act */
					actualString, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{
							providedRef1: nil,
							providedRef2: nil,
						},
						providedExpression,
					)

					/* assert */
					Expect(actualString).To(Equal(coercedString1 + coercedString2))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
	Context("expression doesn't contain ref", func() {
		It("should return input unmodified", func() {
			/* arrange */
			providedExpression := "dummyExpression"

			objectUnderTest := _interpreter{}

			/* act */
			actualString, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedExpression,
			)

			/* assert */
			Expect(actualString).To(Equal(providedExpression))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("expression contains refs", func() {
		Context("at beginning", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedSuffix := "dummySuffix"
				providedExpression := fmt.Sprintf("$(%v)%v", providedRef, providedSuffix)

				fakeCoercer := new(fakeCoercer)

				coercedString := "dummyString"
				fakeCoercer.CoerceReturns(coercedString, nil)

				objectUnderTest := _interpreter{
					coercer: fakeCoercer,
				}

				/* act */
				actualString, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{
						providedRef: nil,
					},
					providedExpression,
				)

				/* assert */
				Expect(actualString).To(Equal(coercedString + providedSuffix))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("in middle", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedPrefix := "dummyPrefix"
				providedSuffix := "dummySuffix"
				providedExpression := fmt.Sprintf("%v$(%v)%v", providedPrefix, providedRef, providedSuffix)

				fakeCoercer := new(fakeCoercer)

				coercedString := "dummyString"
				fakeCoercer.CoerceReturns(coercedString, nil)

				objectUnderTest := _interpreter{
					coercer: fakeCoercer,
				}

				/* act */
				actualString, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{
						providedRef: nil,
					},
					providedExpression,
				)

				/* assert */
				Expect(actualString).To(Equal(providedPrefix + coercedString + providedSuffix))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("at end", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedPrefix := "dummyPrefix"
				providedExpression := fmt.Sprintf("%v$(%v)", providedPrefix, providedRef)

				fakeCoercer := new(fakeCoercer)

				coercedString := "dummyString"
				fakeCoercer.CoerceReturns(coercedString, nil)

				objectUnderTest := _interpreter{
					coercer: fakeCoercer,
				}

				/* act */
				actualString, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{
						providedRef: nil,
					},
					providedExpression,
				)

				/* assert */
				Expect(actualString).To(Equal(providedPrefix + coercedString))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
