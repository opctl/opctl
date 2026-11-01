package unsudo

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ChownTree", func() {
	It("should leave already-owned entries untouched", func() {
		GinkgoT().Setenv("SUDO_UID", strconv.Itoa(os.Geteuid()))
		GinkgoT().Setenv("SUDO_GID", strconv.Itoa(os.Getegid()))
		rootPath := GinkgoT().TempDir()
		filePath := filepath.Join(rootPath, "file")
		linkPath := filepath.Join(rootPath, "link")
		brokenLinkPath := filepath.Join(rootPath, "broken-link")
		Expect(os.WriteFile(filePath, []byte("cached"), 0600)).To(Succeed())
		Expect(os.Symlink(filePath, linkPath)).To(Succeed())
		Expect(os.Symlink(filepath.Join(rootPath, "missing"), brokenLinkPath)).To(Succeed())

		before := map[string]os.FileInfo{}
		for _, entryPath := range []string{filePath, linkPath, brokenLinkPath} {
			info, err := os.Lstat(entryPath)
			Expect(err).To(BeNil())
			before[entryPath] = info
		}
		// Even a same-owner chown updates ctime, so an unchanged ctime proves the entry was skipped.
		time.Sleep(20 * time.Millisecond)
		Expect(ChownTree(rootPath)).To(Succeed())

		for entryPath, info := range before {
			after, err := os.Lstat(entryPath)
			Expect(err).To(BeNil())
			Expect(after.Sys()).To(Equal(info.Sys()), entryPath)
		}
	})

	It("should chown a tree without following symlinks", func() {
		rootPath := GinkgoT().TempDir()
		dirPath := filepath.Join(rootPath, "nested")
		filePath := filepath.Join(dirPath, "file")
		Expect(os.Mkdir(dirPath, 0700)).To(Succeed())
		Expect(os.WriteFile(filePath, []byte("cached"), 0600)).To(Succeed())

		outsidePath := GinkgoT().TempDir()
		outsideFilePath := filepath.Join(outsidePath, "file")
		Expect(os.WriteFile(outsideFilePath, nil, 0600)).To(Succeed())
		linkPath := filepath.Join(rootPath, "link")
		brokenLinkPath := filepath.Join(rootPath, "broken-link")
		Expect(os.Symlink(outsidePath, linkPath)).To(Succeed())
		Expect(os.Symlink(filepath.Join(outsidePath, "missing"), brokenLinkPath)).To(Succeed())

		uid, gid := os.Geteuid(), os.Getegid()
		if uid == 0 {
			uid, gid = 1234, 2345
		}
		GinkgoT().Setenv("SUDO_UID", strconv.Itoa(uid))
		GinkgoT().Setenv("SUDO_GID", strconv.Itoa(gid))

		Expect(ChownTree(rootPath)).To(Succeed())

		for _, entryPath := range []string{rootPath, dirPath, filePath, linkPath, brokenLinkPath} {
			info, err := os.Lstat(entryPath)
			Expect(err).To(BeNil())
			stat := info.Sys().(*syscall.Stat_t)
			Expect(stat.Uid).To(Equal(uint32(uid)))
			Expect(stat.Gid).To(Equal(uint32(gid)))
		}
		for _, entryPath := range []string{outsidePath, outsideFilePath} {
			info, err := os.Stat(entryPath)
			Expect(err).To(BeNil())
			stat := info.Sys().(*syscall.Stat_t)
			Expect(stat.Uid).To(Equal(uint32(os.Geteuid())))
			Expect(stat.Gid).To(Equal(uint32(os.Getegid())))
		}
	})

	It("should fall back to the effective user and group", func() {
		GinkgoT().Setenv("SUDO_UID", "")
		GinkgoT().Setenv("SUDO_GID", "")
		rootPath := GinkgoT().TempDir()
		Expect(ChownTree(rootPath)).To(Succeed())
		info, err := os.Stat(rootPath)
		Expect(err).To(BeNil())
		stat := info.Sys().(*syscall.Stat_t)
		Expect(stat.Uid).To(Equal(uint32(os.Geteuid())))
		Expect(stat.Gid).To(Equal(uint32(os.Getegid())))
	})

	It("should return filesystem errors", func() {
		Expect(os.IsNotExist(ChownTree(filepath.Join(GinkgoT().TempDir(), "missing")))).To(BeTrue())
	})
})
