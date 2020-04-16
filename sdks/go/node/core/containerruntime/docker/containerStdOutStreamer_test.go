package docker

import (
	"bytes"
	"errors"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/internal/fakes"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
)

var _ = Context("containerStdOutStreamer", func() {
	Context("Stream", func() {
		It("should call dockerClient.ContainerLogs w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedContainerName := "dummyContainerName"

			fakeDockerClient := new(FakeCommonAPIClient)
			// err to trigger immediate return
			fakeDockerClient.ContainerLogsReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _containerStdOutStreamer{
				dockerClient: fakeDockerClient,
			}

			expectedOptions := types.ContainerLogsOptions{
				Follow:     true,
				ShowStdout: true,
			}

			/* act */
			objectUnderTest.Stream(
				providedCtx,
				providedContainerName,
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			actualContext,
				actualContainerName,
				actualOptions := fakeDockerClient.ContainerLogsArgsForCall(0)

			Expect(actualContext).To(Equal(providedCtx))
			Expect(actualContainerName).To(Equal(providedContainerName))
			Expect(actualOptions).To(Equal(expectedOptions))
		})
		Context("dockerClient.ContainerLogs errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeDockerClient := new(FakeCommonAPIClient)
				expectedErr := errors.New("dummyErr")
				fakeDockerClient.ContainerLogsReturns(nil, expectedErr)

				objectUnderTest := _containerStdErrStreamer{
					dockerClient: fakeDockerClient,
				}

				/* act */
				actualErr := objectUnderTest.Stream(
					context.Background(),
					"dummyContainerName",
					nopWriteCloser{ioutil.Discard},
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("dockerClient.ContainerLogs doesn't err", func() {
			It("should write expected logs to writeCloser", func() {
				/* arrange */
				providedWriter := bytes.NewBufferString("")

				fakeDockerClient := new(FakeCommonAPIClient)
				expectedErr := errors.New("dummyErr")
				fakeDockerClient.ContainerLogsReturns(nil, expectedErr)

				expectedLogs := "dummyLogs"
				fakeDockerClient.ContainerLogsStub = func(
					ctx context.Context,
					container string,
					options types.ContainerLogsOptions,
				) (io.ReadCloser, error) {
					return ioutil.NopCloser(bytes.NewBufferString(expectedLogs)), nil
				}

				objectUnderTest := _containerStdErrStreamer{
					dockerClient: fakeDockerClient,
				}

				/* act */
				objectUnderTest.Stream(
					context.Background(),
					"dummyContainerName",
					providedWriter,
				)

				/* assert */
				Expect(providedWriter.String()).To(Equal(expectedLogs))
			})
		})
	})
})
