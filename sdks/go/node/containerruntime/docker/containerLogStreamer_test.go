package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
)

type nopWriteCloser struct {
	io.Writer
}

func (w nopWriteCloser) Close() error { return nil }

// orderedWriter records which stream each write came from, in order.
type orderedWriter struct {
	name string
	log  *[]string
}

func (w orderedWriter) Write(p []byte) (int, error) {
	*w.log = append(*w.log, w.name+":"+strings.TrimSpace(string(p)))
	return len(p), nil
}

var _ = Context("containerLogStreamer", func() {
	Context("Stream", func() {
		It("should follow stdout and stderr in one log stream", func() {
			/* arrange */
			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerLogsReturns(nil, errors.New("dummyErr"))

			/* act */
			_containerLogStreamer{fakeDockerClient}.Stream(context.Background(), "dummyContainerName", io.Discard, io.Discard)

			/* assert */
			_, actualContainerName, actualOptions := fakeDockerClient.ContainerLogsArgsForCall(0)
			Expect(actualContainerName).To(Equal("dummyContainerName"))
			Expect(actualOptions).To(Equal(container.LogsOptions{Follow: true, ShowStdout: true, ShowStderr: true}))
		})
		It("should return ContainerLogs errors", func() {
			/* arrange */
			fakeDockerClient := new(FakeCommonAPIClient)
			expectedErr := errors.New("dummyErr")
			fakeDockerClient.ContainerLogsReturns(nil, expectedErr)

			/* act */
			actualErr := _containerLogStreamer{fakeDockerClient}.Stream(context.Background(), "", io.Discard, io.Discard)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
		It("should demux frames to stdout and stderr in the order they were written", func() {
			/* arrange */
			framed := new(bytes.Buffer)
			outFrames := stdcopy.NewStdWriter(framed, stdcopy.Stdout)
			errFrames := stdcopy.NewStdWriter(framed, stdcopy.Stderr)
			outFrames.Write([]byte("out1\n"))
			errFrames.Write([]byte("err1\n"))
			outFrames.Write([]byte("out2\n"))

			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerLogsReturns(io.NopCloser(framed), nil)

			var log []string

			/* act */
			err := _containerLogStreamer{fakeDockerClient}.Stream(
				context.Background(),
				"",
				orderedWriter{"stdout", &log},
				orderedWriter{"stderr", &log},
			)

			/* assert */
			Expect(err).To(BeNil())
			Expect(log).To(Equal([]string{"stdout:out1", "stderr:err1", "stdout:out2"}))
		})
	})
})
