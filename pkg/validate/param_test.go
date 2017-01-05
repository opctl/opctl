package validate

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Describe("Param", func() {

	objectUnderTest := New()
	Context("invoked w/ nil param", func() {
		It("should panic", func() {
			/* arrange/act/assert */
			Expect(
				func() {
					objectUnderTest.Param(&model.Arg{}, nil)
				},
			).To(Panic())
		})
	})
	Context("invoked w/ non-nil param.String", func() {
		Context("& non-empty arg.String", func() {
			Context("equal-length to non-zero param.Constraints.Length.Max", func() {

				It("returns no errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						String: "dummyValue",
					}
					providedParam := &model.Param{
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								Length: &model.StringLengthConstraint{
									Max: len(providedArg.String),
								},
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("longer than non-zero param.Constraints.Length.Max", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						String: "dummyValue",
					}
					providedParam := &model.Param{
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								Length: &model.StringLengthConstraint{
									Max: len(providedArg.String) - 1,
								},
							},
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
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("shorter than non-zero param.Constraints.Length.Max", func() {

				It("returns no errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						String: "dummyValue",
					}
					providedParam := &model.Param{
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								Length: &model.StringLengthConstraint{
									Max: len(providedArg.String) + 1,
								},
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("equal-length to non-zero param.Constraints.Length.Min", func() {

				It("should return no errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						String: "dummyValue",
					}
					providedParam := &model.Param{
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								Length: &model.StringLengthConstraint{
									Min: len(providedArg.String),
								},
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("shorter than non-zero param.Constraints.Length.Min", func() {

				It("should return expected errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						String: "dummyValue",
					}
					providedParam := &model.Param{
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								Length: &model.StringLengthConstraint{
									Min: len(providedArg.String) + 1,
								},
							},
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
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("longer than non-zero param.Constraints.Length.Min", func() {

				It("should return no errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						String: "dummyValue",
					}
					providedParam := &model.Param{
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								Length: &model.StringLengthConstraint{
									Min: len(providedArg.String) - 1,
								},
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("not matching non-empty param.Constraints.Patterns", func() {

				It("should return expected errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						String: "dummyValue",
					}
					providedParam := &model.Param{
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								Patterns: []*model.StringPatternConstraint{
									{
										Regex: "^$",
									},
								},
							},
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
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("matching non-empty param.Constraints.Patterns", func() {

				It("should return no errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						String: "dummyValue",
					}
					providedParam := &model.Param{
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								Patterns: []*model.StringPatternConstraint{
									{
										Regex: ".$",
									},
								},
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("& empty arg.String", func() {
			Context("and non empty Default", func() {
				Context("equal-length to non-zero param.Constraints.Length.Max", func() {

					It("returns no errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							String: "",
						}
						providedDefault := "dummyDefault"
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Length: &model.StringLengthConstraint{
										Max: len(providedDefault),
									},
								},
								Default: providedDefault,
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("longer than non-zero param.Constraints.Length.Max", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							String: "",
						}
						providedDefault := "dummyDefault"
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Length: &model.StringLengthConstraint{
										Max: len(providedDefault) - 1,
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
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("shorter than non-zero param.Constraints.Length.Max", func() {

					It("returns no errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							String: "",
						}
						providedDefault := "dummyDefault"
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Length: &model.StringLengthConstraint{
										Max: len(providedDefault) + 1,
									},
								},
								Default: providedDefault,
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("equal-length to non-zero param.Constraints.Length.Min", func() {

					It("should return no errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							String: "",
						}
						providedDefault := "dummyDefault"
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Length: &model.StringLengthConstraint{
										Min: len(providedDefault),
									},
								},
								Default: providedDefault,
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("shorter than non-zero param.Constraints.Length.Min", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							String: "",
						}
						providedDefault := "dummyDefault"
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Length: &model.StringLengthConstraint{
										Min: len(providedDefault) + 1,
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
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("longer than non-zero param.Constraints.Length.Min", func() {

					It("should return no errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							String: "",
						}
						providedDefault := "dummyDefault"
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Length: &model.StringLengthConstraint{
										Min: len(providedDefault) - 1,
									},
								},
								Default: providedDefault,
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("not matching non-empty param.Constraints.Patterns", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							String: "",
						}
						providedDefault := "dummyDefault"
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Patterns: []*model.StringPatternConstraint{
										{
											Regex: "^$",
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
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("matching non-empty param.Constraints.Patterns", func() {

					It("should return no errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							String: "",
						}
						providedDefault := "dummyDefault"
						providedParam := &model.Param{
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Patterns: []*model.StringPatternConstraint{
										{
											Regex: ".$",
										},
									},
								},
								Default: providedDefault,
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
		Context("& nil arg", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedParam := &model.Param{
					String: &model.StringParam{},
				}

				expectedErrors := []error{
					fmt.Errorf("%v required", providedParam.String.Name),
				}

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})
	Context("invoked w/ non-nil param.NetSocket", func() {
		Context("& non-nil arg.NetSocket", func() {
			Context("w/ non-empty Port", func() {
				Context("equal to 0", func() {
					It("should return expected errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							NetSocket: &model.NetSocketArg{
								Host: "dummyHost",
								Port: 0,
							},
						}
						providedParam := &model.Param{
							NetSocket: &model.NetSocketParam{},
						}

						expectedErrors := []error{
							fmt.Errorf("%v.port must be > 0", providedParam.NetSocket.Name),
						}

						/* act */
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("between 0 and 65536", func() {
					It("should return no errors", func() {
						var i uint = 1
						for i < 65536 {

							/* arrange */
							providedArg := &model.Arg{
								NetSocket: &model.NetSocketArg{
									Host: "dummyHost",
									Port: i,
								},
							}
							providedParam := &model.Param{
								NetSocket: &model.NetSocketParam{},
							}

							expectedErrors := []error{}

							/* act */
							actualErrors := objectUnderTest.Param(providedArg, providedParam)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

							i += i
						}

					})
				})
				Context("equal to 65536", func() {
					It("should return expected errors", func() {

						/* arrange */
						providedArg := &model.Arg{
							NetSocket: &model.NetSocketArg{
								Host: "dummyHost",
								Port: 65536,
							},
						}
						providedParam := &model.Param{
							NetSocket: &model.NetSocketParam{},
						}

						expectedErrors := []error{
							fmt.Errorf("%v.port must be <= 65535", providedParam.NetSocket.Name),
						}

						/* act */
						actualErrors := objectUnderTest.Param(providedArg, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("w/ empty Port", func() {
				It("should return expected errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						NetSocket: &model.NetSocketArg{
							Host: "dummyHost",
						},
					}
					providedParam := &model.Param{
						NetSocket: &model.NetSocketParam{},
					}

					expectedErrors := []error{
						fmt.Errorf("%v.port must be > 0", providedParam.NetSocket.Name),
					}

					/* act */
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("w/ empty Host", func() {
				It("should return expected errors", func() {

					/* arrange */
					providedArg := &model.Arg{
						NetSocket: &model.NetSocketArg{
							Port: 80,
						},
					}
					providedParam := &model.Param{
						NetSocket: &model.NetSocketParam{},
					}

					expectedErrors := []error{
						fmt.Errorf("%v.host required", providedParam.NetSocket.Name),
					}

					/* act */
					actualErrors := objectUnderTest.Param(providedArg, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("& nil arg.NetSocket", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Arg{}
				providedParam := &model.Param{
					NetSocket: &model.NetSocketParam{},
				}

				expectedErrors := []error{
					fmt.Errorf("%v required", providedParam.NetSocket.Name),
				}

				/* act */
				actualErrors := objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("& nil arg", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedParam := &model.Param{
					NetSocket: &model.NetSocketParam{},
				}

				expectedErrors := []error{
					fmt.Errorf("%v required", providedParam.NetSocket.Name),
				}

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})

})
