package client

import (
	"fmt"
	"os"
	"path/filepath"

	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("embedLocalFilesAndDirs", func() {
	initialMaxEmbedBytes := maxEmbedBytes
	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}

	BeforeEach(func() {
		maxEmbedBytes = initialMaxEmbedBytes
	})

	Describe("Args contains file > maxEmbedBytes", func() {
		It("should return expected result", func() {
			// arrange
			maxEmbedBytes = 0
			testDataFilePath := filepath.Join(wd, "testdata/embedLocalFilesAndDirs/rootfile1.txt")

			args := map[string]*model.Value{
				"dummyFile": &model.Value{
					File: &testDataFilePath,
				},
			}

			// act
			actualErr := embedLocalFilesAndDirs(
				args,
			)

			// assert
			Expect(actualErr).To(MatchError(
				fmt.Sprintf("%s is 1.33514404296875e-05Mb but cannot be bigger than 0Mb", testDataFilePath),
			))
		})
	})
	Describe("Args contains dir > maxEmbedBytes", func() {
		It("should return expected result", func() {
			// arrange
			maxEmbedBytes = 15
			testDataFilePath := filepath.Join(wd, "testdata/embedLocalFilesAndDirs")

			args := map[string]*model.Value{
				"dummyDir": &model.Value{
					Dir: &testDataFilePath,
				},
			}

			// act
			actualErr := embedLocalFilesAndDirs(
				args,
			)

			// assert
			Expect(actualErr).To(MatchError(
				fmt.Sprintf("embedding failed: %s cannot exceed 1.430511474609375e-05Mb", testDataFilePath),
			))
		})
	})
	Describe("Args contains files combined > maxEmbedBytes", func() {
		It("should return expected result", func() {
			// arrange
			maxEmbedBytes = 15
			testDataFile1Path := filepath.Join(wd, "testdata/embedLocalFilesAndDirs/rootfile1.txt")
			testDataFile2Path := filepath.Join(wd, "testdata/embedLocalFilesAndDirs/rootfile2.txt")

			args := map[string]*model.Value{
				"testDataFile1": &model.Value{
					File: &testDataFile1Path,
				},
				"testDataFile2": &model.Value{
					File: &testDataFile2Path,
				},
			}

			// act
			actualErr := embedLocalFilesAndDirs(
				args,
			)

			// assert
			Expect(actualErr).To(MatchError(
				errors.New("embedding failed: combined size of files/dirs cannot exceed 1.430511474609375e-05Mb"),
			))
		})
	})
	Describe("Args contains non existent dir ref", func() {
		It("should return expected result", func() {
			// arrange
			pathDoesntExist := "path/doesnt/exist"

			args := map[string]*model.Value{
				"pathDoesntExist": &model.Value{
					Dir: &pathDoesntExist,
				},
			}

			// act
			actualErr := embedLocalFilesAndDirs(
				args,
			)

			// assert
			Expect(actualErr).To(MatchError(fmt.Sprintf("stat %s: no such file or directory", pathDoesntExist)))
		})
	})
	Describe("Args contains non existent file ref", func() {
		It("should return expected result", func() {
			// arrange
			pathDoesntExist := "path/doesnt/exist"

			args := map[string]*model.Value{
				"pathDoesntExist": &model.Value{
					File: &pathDoesntExist,
				},
			}

			// act
			actualErr := embedLocalFilesAndDirs(
				args,
			)

			// assert
			Expect(actualErr).To(MatchError(fmt.Sprintf("stat %s: no such file or directory", pathDoesntExist)))
		})
	})
	Describe("Args contains non existent subdir ref", func() {
		It("should return expected result", func() {
			// arrange
			pathDoesntExist := filepath.Join(wd, "testdata/embedLocalFilesAndDirs", "path/doesnt/exist")

			args := map[string]*model.Value{
				"pathDoesntExist": &model.Value{
					File: &pathDoesntExist,
				},
			}

			// act
			actualErr := embedLocalFilesAndDirs(
				args,
			)

			// assert
			Expect(actualErr).To(MatchError(fmt.Sprintf("stat %s: no such file or directory", pathDoesntExist)))
		})
	})

	It("should return expected result", func() {
		// arrange
		testDataDirPath := filepath.Join(wd, "testdata/embedLocalFilesAndDirs")

		providedDirKey := "dir"
		providedDirValue := &model.Value{
			Dir: &testDataDirPath,
		}

		testDataFilePath := filepath.Join(wd, "testdata/embedLocalFilesAndDirs/rootfile1.txt")
		providedFileKey := "file"
		providedFileValue := &model.Value{
			File: &testDataFilePath,
		}

		providedStringKey := "string"
		providedStringValue := &model.Value{
			String: new(string),
		}

		args := map[string]*model.Value{
			providedDirKey:    providedDirValue,
			providedFileKey:   providedFileValue,
			providedStringKey: providedStringValue,
		}

		// act
		actualErr := embedLocalFilesAndDirs(
			args,
		)

		// assert
		Expect(actualErr).To(BeNil())
		Expect(args).To(
			BeEquivalentTo(
				map[string]*model.Value{
					providedDirKey: {
						Object: &map[string]interface{}{
							"/subdir1": map[string]interface{}{
								"/.gitkeep": map[string]interface{}{
									"data": "",
								},
							},
							"/subdir2": map[string]interface{}{
								"/.gitkeep": map[string]interface{}{
									"data": "",
								},
							},
							"/rootfile1.txt": map[string]interface{}{
								"data": "rootfile1 text",
							},
							"/rootfile2.txt": map[string]interface{}{
								"data": "rootfile2 text",
							},
						},
					},
					providedFileKey: {
						Object: &map[string]interface{}{
							"data": "rootfile1 text",
						},
					},
					providedStringKey: providedStringValue,
				},
			),
		)
	})
})
