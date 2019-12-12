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
		size:    41475,
		modtime: 1576191262,
		compressed: `
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
jPKY8fr/AQAA//8RwGmCA6IAAA==
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
