package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/vfs"
	"github.com/opspec-io/opctl/util/vruntime"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

// dummy FileInfo
type fileInfo struct {
	name    string      // base name of the file
	size    int64       // length in bytes for regular files; system-dependent for others
	mode    os.FileMode // file mode bits
	modTime time.Time   // modification time
	isDir   bool        // abbreviation for Mode().IsDir()
	sys     interface{} // underlying data source (can return nil)
}

func (this fileInfo) Name() string {
	return this.name
}

func (this fileInfo) Size() int64 {
	return this.size
}
func (this fileInfo) Mode() os.FileMode {
	return this.mode
}
func (this fileInfo) ModTime() time.Time {
	return this.modTime
}
func (this fileInfo) IsDir() bool {
	return this.isDir
}
func (this fileInfo) Sys() interface{} {
	return &this.sys
}

var _ = Context("InspectContainerIfExists", func() {
	It("should call dockerClient.ContainerInspect w/ expected args", func() {
		/* arrange */
		fakeDockerClient := new(fakeDockerClient)

		providedContainerId := "dummyContainerId"
		expectedContainerId := providedContainerId
		fakeDockerClient.ContainerInspectReturns(
			types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{},
				Config:            &container.Config{},
				NetworkSettings: &types.NetworkSettings{
					DefaultNetworkSettings: types.DefaultNetworkSettings{},
				},
			},
			nil,
		)

		objectUnderTest := _containerProvider{
			dockerClient: fakeDockerClient,
		}

		/* act */
		objectUnderTest.InspectContainerIfExists(providedContainerId)

		/* assert */
		_, actualContainerId := fakeDockerClient.ContainerInspectArgsForCall(0)
		Expect(actualContainerId).To(Equal(expectedContainerId))
	})
	Context("dockerClient.ContainerInspect errors", func() {
		Context("is NotFoundError", func() {
			It("shouldn't return error", func() {
				/* arrange */
				fakeDockerClient := new(fakeDockerClient)
				fakeDockerClient.ContainerInspectReturns(
					types.ContainerJSON{
						ContainerJSONBase: &types.ContainerJSONBase{},
						Config:            &container.Config{},
						NetworkSettings: &types.NetworkSettings{
							DefaultNetworkSettings: types.DefaultNetworkSettings{},
						},
					},
					dockerNotFoundError{},
				)

				objectUnderTest := _containerProvider{
					dockerClient: fakeDockerClient,
				}

				/* act */
				_, actualError := objectUnderTest.InspectContainerIfExists("dummyContainerId")

				/* assert */
				Expect(actualError).To(BeNil())
			})
		})
		Context("isn't NotFoundError", func() {
			It("should return error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")
				fakeDockerClient := new(fakeDockerClient)
				fakeDockerClient.ContainerInspectReturns(types.ContainerJSON{}, expectedError)

				objectUnderTest := _containerProvider{
					dockerClient: fakeDockerClient,
				}

				/* act */
				_, actualError := objectUnderTest.InspectContainerIfExists("dummyContainerId")

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
	Context("dockerClient.ContainerInspect doesn't error", func() {
		It("should return expected container", func() {
			/* arrange */
			fakeDockerClient := new(fakeDockerClient)
			fakeVfs := new(vfs.Fake)

			providedContainerId := "dummyContainerId"
			dockerContainer := types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					Image: "dummyImage",
				},
				Config: &container.Config{
					Env: []string{
						"dummyEnvVar1Name=dummyEnvVar1Value",
						"dummyEnvVar2Name=dummyEnvVar2Value",
					},
					Entrypoint: []string{"dummyEntrypoint"},
					WorkingDir: "dummyWorkDir",
				},
				Mounts: []types.MountPoint{
					{
						Source:      "dummyFile1Src",
						Destination: "dummyFile1Dst",
					},
					{
						Source:      "dummyDir1Src",
						Destination: "dummyDir1Dst",
					},
					{
						Source:      "dummySocket1Src",
						Destination: "dummySocket1Dst",
					},
					{
						Source:      "dummyNamedPipe1Src",
						Destination: "dummyNamedPipe1Dst",
					},
				},
				NetworkSettings: &types.NetworkSettings{
					DefaultNetworkSettings: types.DefaultNetworkSettings{
						IPAddress: "dummyIpAddress",
					},
				},
			}
			fakeDockerClient.ContainerInspectReturns(
				dockerContainer,
				nil,
			)

			fakeVfs.StatStub = func(name string) (os.FileInfo, error) {
				switch {
				case strings.Contains(name, "Dir"):
					return fileInfo{isDir: true}, nil
				case strings.Contains(name, "File"):
					return fileInfo{}, nil
				case strings.Contains(name, "Socket"):
					return fileInfo{mode: os.ModeSocket}, nil
				case strings.Contains(name, "NamedPipe"):
					return fileInfo{mode: os.ModeNamedPipe}, nil
				default:
					panic("invalid test data")
				}
			}

			expectedContainer := &model.DcgContainerCall{
				Cmd: dockerContainer.Config.Entrypoint,
				Dirs: map[string]string{
					"dummyDir1Dst": "dummyDir1Src",
				},
				EnvVars: map[string]string{
					"dummyEnvVar1Name": "dummyEnvVar1Value",
					"dummyEnvVar2Name": "dummyEnvVar2Value",
				},
				Files: map[string]string{
					"dummyFile1Dst": "dummyFile1Src",
				},
				Image: dockerContainer.Image,
				Sockets: map[string]string{
					"dummySocket1Dst":    "dummySocket1Src",
					"dummyNamedPipe1Dst": "dummyNamedPipe1Src",
				},
				WorkDir:   dockerContainer.Config.WorkingDir,
				IpAddress: dockerContainer.NetworkSettings.IPAddress,
			}

			objectUnderTest := _containerProvider{
				dockerClient: fakeDockerClient,
				fs:           fakeVfs,
				runtime:      new(vruntime.Fake),
			}

			/* act */
			actualContainer, _ := objectUnderTest.InspectContainerIfExists(providedContainerId)

			/* assert */
			Expect(actualContainer).To(Equal(expectedContainer))
		})
	})
})
