package container

import (
	"errors"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("cmd.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					Cmd: []interface{}{
						"$()",
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				os.TempDir(),
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to interpret $() to string; error was unable to interpret '' as reference; '' not in scope")))
		})
	})

	Context("dirs.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			identifier := "identifier"

			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{
					identifier: &model.Value{
						Socket: new(string),
					},
				},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					Dirs: map[string]interface{}{
						"/something": fmt.Sprintf("$(%s)", identifier),
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				os.TempDir(),
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to bind /something to $(identifier); error was unable to interpret $(identifier) to dir; error was unable to coerce socket to dir; incompatible types")))
		})
	})

	Context("envVars.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					EnvVars: "$()",
				},
				"dummyContainerID",
				"dummyOpPath",
				os.TempDir(),
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to interpret '$()' as envVars; error was unable to interpret $() to object; error was unable to interpret '' as reference; '' not in scope")))
		})
	})

	Context("files.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{
					"not": &model.Value{
						Socket: new(string),
					},
				},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					Files: map[string]interface{}{
						"/something": "$(not)",
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				os.TempDir(),
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to bind /something to $(not); error was unable to coerce '{\"socket\":\"\"}' to file")))
		})
	})

	Context("image.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "$()",
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				os.TempDir(),
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to interpret $() to string; error was unable to interpret '' as reference; '' not in scope")))
		})
	})

	It("should return expected result", func() {
		/* arrange */
		providedContainerID := "providedContainerID"
		providedOpPath := "providedOpPath"

		expectedRef := "docker.io/library/ref"

		expectedResult := model.ContainerCall{
			BaseCall: model.BaseCall{
				OpPath: providedOpPath,
			},
			ContainerID: providedContainerID,
			Cmd:         []string{},
			Dirs:        map[string]string{},
			Files:       map[string]string{},
			Image: &model.ContainerCallImage{
				Ref: &expectedRef,
			},
			Sockets: map[string]string{},
			WorkDir: "",
		}

		/* act */
		actualResult, actualErr := Interpret(
			map[string]*model.Value{},
			&model.ContainerCallSpec{
				Image: &model.ContainerCallImageSpec{
					Ref: "ref",
				},
			},
			providedContainerID,
			providedOpPath,
			os.TempDir(),
		)

		/* assert */
		Expect(actualErr).To(BeNil())
		Expect(*actualResult).To(Equal(expectedResult))
	})
})
