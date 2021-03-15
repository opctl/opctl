package container

import (
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("cmd.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

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
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret $() to string: unable to interpret '' as reference: '' not in scope"))
		})
	})

	Context("dirs.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			identifier := "identifier"
			dataDir, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

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
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to bind /something to $(identifier): unable to interpret $(identifier) to dir: unable to coerce socket to dir: incompatible types"))
		})
	})

	Context("envVars.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

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
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret '$()' as envVars: unable to interpret $() to object: unable to interpret '' as reference: '' not in scope"))
		})
	})

	Context("files.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

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
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to bind /something to $(not): unable to coerce '{\"socket\":\"\"}' to file"))
		})
	})

	Context("image.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

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
			Dirs:        map[string]string{},
			Files:       map[string]string{},
			Image: &model.ContainerCallImage{
				Ref: &expectedRef,
			},
			Sockets: map[string]string{},
			WorkDir: "",
		}

		dataDir, err := ioutil.TempDir("", "")
		if err != nil {
			panic(err)
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
			dataDir,
		)

		/* assert */
		Expect(actualErr).To(BeNil())
		Expect(*actualResult).To(Equal(expectedResult))
	})
})
