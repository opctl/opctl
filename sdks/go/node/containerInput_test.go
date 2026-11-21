package node

import (
	"context"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container input", func() {
	It("can release an unused claim without ending stdin", func() {
		registry := &containerInputs{}
		objectUnderTest := core{containerInputs: registry}
		_, unregister := registry.register("container-id")
		defer unregister()
		ctx := context.Background()

		input, release, err := objectUnderTest.ClaimContainerInput(ctx, "container-id")
		Expect(err).NotTo(HaveOccurred())
		_, err = objectUnderTest.AttachContainer(ctx, "container-id")
		Expect(err).To(MatchError(ErrContainerInputAttached))
		release()
		Expect(input.Done()).NotTo(BeClosed())

		next, err := objectUnderTest.AttachContainer(ctx, "container-id")
		Expect(err).NotTo(HaveOccurred())
		release()
		_, err = objectUnderTest.AttachContainer(ctx, "container-id")
		Expect(err).To(MatchError(ErrContainerInputAttached))
		Expect(next.Close()).To(Succeed())
	})

	It("allows one attached client and ends stdin when it disconnects", func() {
		registry := &containerInputs{}
		objectUnderTest := core{containerInputs: registry}
		ctx := context.Background()
		_, err := objectUnderTest.AttachContainer(ctx, "container-id")
		Expect(err).To(MatchError(ErrContainerInputNotFound))

		session, unregister := registry.register("container-id")
		defer unregister()
		input, err := objectUnderTest.AttachContainer(ctx, "container-id")
		Expect(err).NotTo(HaveOccurred())
		_, err = objectUnderTest.AttachContainer(ctx, "container-id")
		Expect(err).To(MatchError(ErrContainerInputAttached))

		result := make(chan string, 1)
		go func() {
			data, err := io.ReadAll(session.reader)
			if err != nil {
				result <- err.Error()
				return
			}
			result <- string(data)
		}()
		_, err = input.Write([]byte("hello"))
		Expect(err).NotTo(HaveOccurred())
		Expect(input.Close()).To(Succeed())
		Eventually(result).Should(Receive(Equal("hello")))

		unregister()
		_, err = objectUnderTest.AttachContainer(ctx, "container-id")
		Expect(err).To(MatchError(ErrContainerInputNotFound))
	})
})
