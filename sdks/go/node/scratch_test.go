package node

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("scratch", func() {
	const (
		rootID  = "00000000000000000000000000000001"
		childID = "00000000000000000000000000000002"
	)

	// newTree starts a root with one child call, each with scratch holding a secret file, beside the
	// events db.
	newTree := func(ttl time.Duration, clean bool) (*scratch, string) {
		dataDir, err := os.MkdirTemp("", "")
		Expect(err).To(BeNil())

		s := newScratch(dataDir, ttl)
		for _, dir := range []string{"events", rootID, filepath.Join(childID, "fs")} {
			Expect(os.MkdirAll(filepath.Join(s.dcg, dir), 0o700)).To(Succeed())
		}
		Expect(os.WriteFile(filepath.Join(s.dcg, childID, "fs", "vars.json"), []byte("secret"), 0o600)).To(Succeed())

		s.begin(rootID, clean)
		s.track(rootID, rootID)
		s.track(rootID, childID)

		return s, dataDir
	}

	exists := func(s *scratch, name string) bool {
		_, err := os.Stat(filepath.Join(s.dcg, name))
		return err == nil
	}

	Context("a root that asked to be cleaned", func() {
		It("removes its whole tree when it ends, and nothing else", func() {
			s, dataDir := newTree(0, true)
			defer os.RemoveAll(dataDir)

			s.end(rootID)

			Expect(exists(s, rootID)).To(BeFalse())
			Expect(exists(s, childID)).To(BeFalse())
			Expect(exists(s, rootID+journalSuffix)).To(BeFalse())
			Expect(exists(s, "events")).To(BeTrue())
		})

		It("is removed on the next node start when the node died before it ended", func() {
			s, dataDir := newTree(0, true)
			defer os.RemoveAll(dataDir)

			newScratch(dataDir, 0).sweep(time.Now())

			Expect(exists(s, childID)).To(BeFalse())
			Expect(exists(s, "events")).To(BeTrue())
		})
	})

	Context("a root that didn't ask", func() {
		It("is kept after it ends, so API callers can still read its outputs", func() {
			s, dataDir := newTree(0, false)
			defer os.RemoveAll(dataDir)

			s.end(rootID)
			s.sweep(time.Now().Add(24 * time.Hour))

			Expect(exists(s, filepath.Join(childID, "fs", "vars.json"))).To(BeTrue())
			_, _, ended := s.read(rootID)
			Expect(ended).To(BeTrue())
		})

		It("is removed once the TTL has passed since it ended, and not before", func() {
			s, dataDir := newTree(time.Hour, false)
			defer os.RemoveAll(dataDir)

			s.end(rootID)

			s.removeExpired(time.Now().Add(30 * time.Minute))
			Expect(exists(s, childID)).To(BeTrue())

			s.removeExpired(time.Now().Add(2 * time.Hour))
			Expect(exists(s, childID)).To(BeFalse())
			Expect(exists(s, rootID+journalSuffix)).To(BeFalse())
		})

		It("is never removed by the TTL while it runs", func() {
			s, dataDir := newTree(time.Hour, false)
			defer os.RemoveAll(dataDir)

			s.removeExpired(time.Now().Add(48 * time.Hour))

			Expect(exists(s, childID)).To(BeTrue())
		})

		It("counts the TTL from the next node start when the node died before it ended", func() {
			s, dataDir := newTree(time.Hour, false)
			defer os.RemoveAll(dataDir)

			newScratch(dataDir, time.Hour).sweep(time.Now())

			Expect(exists(s, childID)).To(BeTrue())
			_, _, ended := s.read(rootID)
			Expect(ended).To(BeTrue())
		})
	})

	It("ignores journals that don't name a call", func() {
		s, dataDir := newTree(0, true)
		defer os.RemoveAll(dataDir)
		Expect(os.WriteFile(filepath.Join(s.dcg, "events"+journalSuffix), []byte(cleanMark+"\n"), 0o600)).To(Succeed())

		s.sweep(time.Now())

		Expect(exists(s, "events")).To(BeTrue())
	})

	It("does nothing without a scratch, as in tests that build a caller directly", func() {
		var s *scratch
		s.begin(rootID, true)
		s.track(rootID, childID)
		s.end(rootID)
		s.sweep(time.Now())
	})
})
