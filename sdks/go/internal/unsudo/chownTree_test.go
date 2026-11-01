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

var _ = Context("RestoreCacheTree", func() {
	It("should add a+rwX, preserve directory special bits, and strip file setuid/setgid", func() {
		GinkgoT().Setenv("SUDO_UID", "")
		GinkgoT().Setenv("SUDO_GID", "")
		rootPath := GinkgoT().TempDir()
		tests := []struct {
			name string
			mode os.FileMode
			want os.FileMode
		}{
			{"plain", 0400, 0666},
			{"unreadable", 0000, 0666},
			{"executable", 0100, 0777},
			{"group-executable", 0010, 0777},
			{"other-executable", 0001, 0777},
			{"privileged", 0755 | os.ModeSetuid | os.ModeSetgid, 0777},
			{"privileged-writable", 0777 | os.ModeSetuid | os.ModeSetgid, 0777},
			{"sticky-file", 0600 | os.ModeSticky, 0666 | os.ModeSticky},
		}
		for _, test := range tests {
			entryPath := filepath.Join(rootPath, test.name)
			Expect(os.WriteFile(entryPath, nil, 0600)).To(Succeed())
			Expect(os.Chmod(entryPath, test.mode)).To(Succeed())
		}
		dirPath := filepath.Join(rootPath, "nested")
		Expect(os.Mkdir(dirPath, 0700)).To(Succeed())
		Expect(os.Chmod(dirPath, 0700|os.ModeSetgid|os.ModeSticky)).To(Succeed())

		Expect(RestoreCacheTree(rootPath)).To(Succeed())

		for _, test := range tests {
			info, err := os.Stat(filepath.Join(rootPath, test.name))
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Mode()).To(Equal(test.want), test.name)
		}
		info, err := os.Stat(dirPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.Mode()).To(Equal(os.ModeDir | 0777 | os.ModeSetgid | os.ModeSticky))
	})

	It("should skip symlinks and leave their targets untouched", func() {
		GinkgoT().Setenv("SUDO_UID", "")
		rootPath := GinkgoT().TempDir()
		outsidePath := GinkgoT().TempDir()
		outsideFile := filepath.Join(outsidePath, "file")
		Expect(os.WriteFile(outsideFile, nil, 0600)).To(Succeed())
		for name, target := range map[string]string{
			"dir-link":    outsidePath,
			"file-link":   outsideFile,
			"broken-link": filepath.Join(outsidePath, "missing"),
		} {
			Expect(os.Symlink(target, filepath.Join(rootPath, name))).To(Succeed())
		}
		before := map[string]os.FileInfo{}
		for _, entryPath := range []string{outsidePath, outsideFile, filepath.Join(rootPath, "dir-link"), filepath.Join(rootPath, "file-link"), filepath.Join(rootPath, "broken-link")} {
			info, err := os.Lstat(entryPath)
			Expect(err).NotTo(HaveOccurred())
			before[entryPath] = info
		}
		time.Sleep(20 * time.Millisecond)
		Expect(RestoreCacheTree(rootPath)).To(Succeed())
		for entryPath, info := range before {
			after, err := os.Lstat(entryPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(after.Sys()).To(Equal(info.Sys()), entryPath)
		}
	})

	It("should reject a symlink as the tree root", func() {
		linkPath := filepath.Join(GinkgoT().TempDir(), "link")
		Expect(os.Symlink(GinkgoT().TempDir(), linkPath)).To(Succeed())
		Expect(RestoreCacheTree(linkPath)).To(MatchError(ContainSubstring("must be a real directory")))
	})

	It("should not chmod an outside symlink target when an entry is replaced after stat", func() {
		rootPath := GinkgoT().TempDir()
		entryPath := filepath.Join(rootPath, "entry")
		targetPath := filepath.Join(GinkgoT().TempDir(), "target")
		Expect(os.WriteFile(entryPath, nil, 0600)).To(Succeed())
		Expect(os.WriteFile(targetPath, nil, 0600)).To(Succeed())
		info, err := os.Lstat(entryPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Remove(entryPath)).To(Succeed())
		Expect(os.Symlink(targetPath, entryPath)).To(Succeed())
		root, err := os.OpenRoot(rootPath)
		Expect(err).NotTo(HaveOccurred())
		defer root.Close()

		Expect(makeCacheEntryWritable(root, "entry", info)).To(Succeed())

		after, err := os.Stat(targetPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(after.Mode().Perm()).To(Equal(os.FileMode(0600)))
	})

	It("should leave already-owned entries with the target mode untouched", func() {
		GinkgoT().Setenv("SUDO_UID", strconv.Itoa(os.Geteuid()))
		GinkgoT().Setenv("SUDO_GID", strconv.Itoa(os.Getegid()))
		rootPath := GinkgoT().TempDir()
		filePath := filepath.Join(rootPath, "file")
		linkPath := filepath.Join(rootPath, "link")
		brokenLinkPath := filepath.Join(rootPath, "broken-link")
		Expect(os.WriteFile(filePath, []byte("cached"), 0666)).To(Succeed())
		Expect(os.Symlink(filePath, linkPath)).To(Succeed())
		Expect(os.Symlink(filepath.Join(rootPath, "missing"), brokenLinkPath)).To(Succeed())

		Expect(os.Chmod(rootPath, 0777)).To(Succeed())
		Expect(os.Chmod(filePath, 0666)).To(Succeed())

		before := map[string]os.FileInfo{}
		for _, entryPath := range []string{filePath, linkPath, brokenLinkPath} {
			info, err := os.Lstat(entryPath)
			Expect(err).To(BeNil())
			before[entryPath] = info
		}
		// Even a same-owner chown or same-mode chmod updates ctime, so an unchanged ctime proves the entry was skipped.
		time.Sleep(20 * time.Millisecond)
		Expect(RestoreCacheTree(rootPath)).To(Succeed())

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
		if uid != 0 {
			Skip("ownership restoration requires sudo")
		}
		uid, gid = 1234, 2345
		GinkgoT().Setenv("SUDO_UID", strconv.Itoa(uid))
		GinkgoT().Setenv("SUDO_GID", strconv.Itoa(gid))

		Expect(RestoreCacheTree(rootPath)).To(Succeed())

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

	It("should leave ownership unchanged without sudo", func() {
		GinkgoT().Setenv("SUDO_UID", "")
		GinkgoT().Setenv("SUDO_GID", "")
		rootPath := GinkgoT().TempDir()
		uid, gid := os.Geteuid(), os.Getegid()
		if uid == 0 {
			uid, gid = 1234, 2345
			Expect(os.Chown(rootPath, uid, gid)).To(Succeed())
		}
		Expect(RestoreCacheTree(rootPath)).To(Succeed())
		info, err := os.Stat(rootPath)
		Expect(err).To(BeNil())
		stat := info.Sys().(*syscall.Stat_t)
		Expect(stat.Uid).To(Equal(uint32(uid)))
		Expect(stat.Gid).To(Equal(uint32(gid)))
	})

	It("should return filesystem errors", func() {
		Expect(os.IsNotExist(RestoreCacheTree(filepath.Join(GinkgoT().TempDir(), "missing")))).To(BeTrue())
	})
})
