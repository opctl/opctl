//go:build unix

package caches

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("cache dir permissions", func() {
	for _, noCache := range []bool{false, true} {
		for _, entryType := range []string{"symlink", "file"} {
			It(fmt.Sprintf("should reject a %s mount source without changing it (no-cache=%t)", entryType, noCache), func() {
				dataDirPath := GinkgoT().TempDir()
				cacheRootPath := filepath.Join(dataDirPath, "caches")
				noCacheRootPath := ""
				if noCache {
					noCacheRootPath = filepath.Join(dataDirPath, "nocache", "run")
					cacheRootPath = noCacheRootPath
				}
				Expect(os.MkdirAll(cacheRootPath, 0700)).To(Succeed())
				hostPath := filepath.Join(cacheRootPath, "shared")
				targetPath := hostPath
				if entryType == "symlink" {
					targetPath = GinkgoT().TempDir()
					Expect(os.Symlink(targetPath, hostPath)).To(Succeed())
				} else {
					Expect(os.WriteFile(hostPath, nil, 0600)).To(Succeed())
				}
				before, err := os.Stat(targetPath)
				Expect(err).NotTo(HaveOccurred())

				dirs, err := Interpret(nil, map[string]interface{}{"/cache": "shared"}, nil, dataDirPath, noCacheRootPath)

				Expect(err).To(MatchError(ContainSubstring("must be a real directory")))
				Expect(dirs).To(BeNil())
				after, err := os.Stat(targetPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(after.Sys()).To(Equal(before.Sys()))
			})
		}

		It(fmt.Sprintf("should create mode 0777 dirs despite a restrictive umask (no-cache=%t)", noCache), func() {
			dataDirPath := GinkgoT().TempDir()
			noCacheRootPath := ""
			if noCache {
				noCacheRootPath = filepath.Join(dataDirPath, "nocache", "run")
			}
			previousUmask := syscall.Umask(0077)
			defer syscall.Umask(previousUmask)

			dirs, err := Interpret(nil, map[string]interface{}{"/cache": "shared"}, nil, dataDirPath, noCacheRootPath)

			Expect(err).NotTo(HaveOccurred())
			info, err := os.Stat(dirs["/cache"])
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Mode().Perm()).To(Equal(os.FileMode(0777)))
			parentInfo, err := os.Stat(filepath.Dir(dirs["/cache"]))
			Expect(err).NotTo(HaveOccurred())
			Expect(parentInfo.Mode().Perm()).To(Equal(os.FileMode(0700)))
		})

		for _, mode := range []os.FileMode{0700, 0777} {
			It(fmt.Sprintf("should preserve directory special bits with mode %04o (no-cache=%t)", mode, noCache), func() {
				dataDirPath := GinkgoT().TempDir()
				noCacheRootPath := ""
				cacheRootPath := filepath.Join(dataDirPath, "caches")
				if noCache {
					noCacheRootPath = filepath.Join(dataDirPath, "nocache", "run")
					cacheRootPath = noCacheRootPath
				}
				hostPath := filepath.Join(cacheRootPath, "shared")
				Expect(os.MkdirAll(hostPath, 0700)).To(Succeed())
				Expect(os.Chmod(hostPath, mode|os.ModeSetgid|os.ModeSticky)).To(Succeed())
				contentsPath := filepath.Join(hostPath, "contents")
				Expect(os.WriteFile(contentsPath, []byte("cached"), 0600)).To(Succeed())
				before, err := os.Stat(hostPath)
				Expect(err).NotTo(HaveOccurred())

				dirs, err := Interpret(nil, map[string]interface{}{"/cache": "shared"}, nil, dataDirPath, noCacheRootPath)

				Expect(err).NotTo(HaveOccurred())
				Expect(dirs["/cache"]).To(Equal(hostPath))
				info, err := os.Stat(hostPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.Mode().Perm()).To(Equal(os.FileMode(0777)))
				Expect(info.Mode() & (os.ModeSetgid | os.ModeSticky)).To(Equal(os.ModeSetgid | os.ModeSticky))
				if mode == 0777 {
					Expect(info.Sys()).To(Equal(before.Sys()))
				}
				contents, err := os.ReadFile(contentsPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(Equal("cached"))
			})
		}
	}
})

func TestInterpretCacheDirUnprivileged(t *testing.T) {
	if dataDirPath := os.Getenv("OPCTL_TEST_UNPRIVILEGED_CACHE_DIRS"); dataDirPath != "" {
		for _, noCacheRootPath := range []string{"", filepath.Join(dataDirPath, "nocache", "run")} {
			cacheRootPath := noCacheRootPath
			if cacheRootPath == "" {
				cacheRootPath = filepath.Join(dataDirPath, "caches")
			}
			for name, mode := range map[string]os.FileMode{"restricted": 0755, "writable": 0777} {
				hostPath := filepath.Join(cacheRootPath, name)
				before, err := os.Stat(hostPath)
				if err != nil {
					if errors.Is(err, fs.ErrPermission) {
						t.Skipf("unprivileged helper cannot access fixture in this environment: %v", err)
					}
					t.Fatal(err)
				}
				if err := os.Chmod(hostPath, 0777); !errors.Is(err, syscall.EPERM) {
					t.Fatalf("expected EPERM for another user's directory, got %v", err)
				}
				dirs, err := Interpret(nil, map[string]interface{}{"/cache": name}, nil, dataDirPath, noCacheRootPath)
				if err != nil {
					t.Fatalf("another user's cache directory should still bind: %v", err)
				}
				if dirs["/cache"] != hostPath {
					t.Fatalf("got cache path %q, want %q", dirs["/cache"], hostPath)
				}
				after, err := os.Stat(hostPath)
				if err != nil {
					t.Fatal(err)
				}
				if after.Mode().Perm() != mode || !reflect.DeepEqual(after.Sys(), before.Sys()) {
					t.Errorf("foreign directory %s was changed", hostPath)
				}
			}
		}
		return
	}
	if os.Geteuid() != 0 {
		t.Skip("requires root to create foreign-owned fixtures and drop privileges")
	}
	fixturePath, err := os.MkdirTemp("/tmp", "opctl-cache-dirs-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fixturePath)
	for _, relativePath := range []string{".", "caches", "nocache", "nocache/run"} {
		path := filepath.Join(fixturePath, relativePath)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(path, 0755); err != nil {
			t.Fatal(err)
		}
	}
	for _, rootPath := range []string{"caches", "nocache/run"} {
		for name, mode := range map[string]os.FileMode{"restricted": 0755, "writable": 0777} {
			path := filepath.Join(fixturePath, rootPath, name)
			if err := os.Mkdir(path, mode); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, mode); err != nil {
				t.Fatal(err)
			}
		}
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	helperPath := filepath.Join(fixturePath, "test-helper")
	if err := os.WriteFile(helperPath, contents, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(helperPath, 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(helperPath, "-test.run=^TestInterpretCacheDirUnprivileged$", "-test.v")
	cmd.Env = append(os.Environ(), "OPCTL_TEST_UNPRIVILEGED_CACHE_DIRS="+fixturePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Credential: &syscall.Credential{Uid: 65534, Gid: 65534}}
	if output, err := cmd.CombinedOutput(); err != nil {
		if errors.Is(err, fs.ErrPermission) {
			t.Skipf("cannot execute unprivileged helper in this environment (check /tmp noexec and directory access): %v", err)
		}
		t.Fatalf("unprivileged helper: %v\n%s", err, output)
	} else {
		t.Logf("unprivileged helper:\n%s", output)
	}
}
