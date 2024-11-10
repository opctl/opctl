package container

import (
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
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				map[string]*ipld.Node{},
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
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret [$()] to array: unable to interpret '$()' as array initializer item: unable to interpret '' as reference: '' not in scope"))
		})
	})

	Context("dirs.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			identifier := "identifier"
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				map[string]*ipld.Node{
					identifier: {
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
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to bind directory /something to $(identifier): unable to interpret $(identifier) to dir: unable to coerce socket to dir: incompatible types"))
		})
	})

	Context("envVars.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				map[string]*ipld.Node{},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					EnvVars: "$()",
				},
				"dummyContainerID",
				"dummyOpPath",
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret '$()' as envVars: unable to interpret $() to object: unable to interpret '' as reference: '' not in scope"))
		})
	})

	Context("files.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				map[string]*ipld.Node{
					"not": {
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
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to bind file /something to $(not): unable to coerce '{\"socket\":\"\"}' to file"))
		})
	})

	Context("image.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				map[string]*ipld.Node{},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "$()",
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret $() to string: unable to interpret '' as reference: '' not in scope"))
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
			Dirs:        model.NewStringMap(map[string]string{}),
			Files:       model.NewStringMap(map[string]string{}),
			Image: &model.ContainerCallImage{
				Ref: &expectedRef,
			},
			Sockets: model.NewStringMap(map[string]string{}),
			WorkDir: "",
		}

		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}

		/* act */
		actualResult, actualErr := Interpret(
			map[string]*ipld.Node{},
			&model.ContainerCallSpec{
				Image: &model.ContainerCallImageSpec{
					Ref: "ref",
				},
			},
			providedContainerID,
			providedOpPath,
			dataDir,
		)

		/* assert */
		Expect(actualErr).To(BeNil())
		Expect(*actualResult).To(Equal(expectedResult))
	})
})
