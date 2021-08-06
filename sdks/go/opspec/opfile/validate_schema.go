// Code generated by "esc -pkg=opfile -o validate_schema.go -private ../../../../opspec/opfile/jsonschema.json"; DO NOT EDIT.

package opfile

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/opspec/opfile/jsonschema.json": {
		name:    "jsonschema.json",
		local:   "../../../../opspec/opfile/jsonschema.json",
		size:    36027,
		modtime: 1628286479,
		compressed: `
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
Lry7jAAA
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
