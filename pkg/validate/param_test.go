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
					objectUnderTest.Param(&model.Data{}, nil)
				},
			).To(Panic())
		})
	})
	Context("invoked w/ non-nil param.String", func() {
		Context("& non-empty arg.String", func() {
			Context("equal-length to non-zero param.Constraints.Length.Max", func() {

				It("returns no errors", func() {

					/* arrange */
					providedArg := &model.Data{
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
					providedArg := &model.Data{
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
					providedArg := &model.Data{
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
					providedArg := &model.Data{
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
					providedArg := &model.Data{
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
					providedArg := &model.Data{
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
					providedArg := &model.Data{
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
					providedArg := &model.Data{
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
						providedArg := &model.Data{
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
						providedArg := &model.Data{
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
						providedArg := &model.Data{
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
						providedArg := &model.Data{
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
						providedArg := &model.Data{
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
						providedArg := &model.Data{
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
						providedArg := &model.Data{
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
						providedArg := &model.Data{
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
	Context("invoked w/ non-nil param.Socket", func() {
		Context("& non-empty arg.Socket", func() {
			It("should return no errors", func() {

				/* arrange */
				providedArg := &model.Data{
					Socket: "dummyValue",
				}
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{}

				/* act */
				actualErrors := objectUnderTest.Param(providedArg, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("& empty arg.Socket", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Data{}
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{
					fmt.Errorf("%v required", providedParam.Socket.Name),
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
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{
					fmt.Errorf("%v required", providedParam.Socket.Name),
				}

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})

})
