package docker

import (
	"bytes"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
)

var _ = Context("containerStdErrStreamer", func() {
	Context("Stream", func() {
		It("should call dockerClient.ContainerLogs w/ expected args", func() {
			/* arrange */
			providedContainerID := "dummyContainerID"

			fakeDockerClient := new(fakeDockerClient)
			// err to trigger immediate return
			fakeDockerClient.ContainerLogsReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _containerStdErrStreamer{
				dockerClient: fakeDockerClient,
			}

			expectedOptions := types.ContainerLogsOptions{
				Follow:     true,
				ShowStderr: true,
			}

			/* act */
			objectUnderTest.Stream(
				providedContainerID,
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			actualContext,
				actualContainerID,
				actualOptions := fakeDockerClient.ContainerLogsArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualContainerID).To(Equal(providedContainerID))
			Expect(actualOptions).To(Equal(expectedOptions))
		})
		Context("dockerClient.ContainerLogs errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeDockerClient := new(fakeDockerClient)
				expectedErr := errors.New("dummyErr")
				fakeDockerClient.ContainerLogsReturns(nil, expectedErr)

				objectUnderTest := _containerStdErrStreamer{
					dockerClient: fakeDockerClient,
				}

				/* act */
				actualErr := objectUnderTest.Stream(
					"dummyContainerID",
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

				fakeDockerClient := new(fakeDockerClient)
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
					"dummyContainerID",
					providedWriter,
				)

				/* assert */
				Expect(providedWriter.String()).To(Equal(expectedLogs))
			})
		})
	})
})
