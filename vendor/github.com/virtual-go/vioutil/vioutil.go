package vioutil

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ VIOUtil

import (
	"bytes"
	"github.com/virtual-go/fs"
	"io"
	"os"
	"sort"
)

// virtual filesystem interface
type VIOUtil interface {
	// ReadDir reads the directory named by dirname and returns
	// a list of directory entries sorted by filename.
	ReadDir(dirname string) ([]os.FileInfo, error)

	// ReadFile reads the file named by filename and returns the contents.
	// A successful call returns err == nil, not err == EOF. Because ReadFile
	// reads the whole file, it does not treat an EOF from Read as an error
	// to be reported.
	ReadFile(filename string) ([]byte, error)

	// WriteFile writes data to a file named by filename.
	// If the file does not exist, WriteFile creates it with permissions perm;
	// otherwise WriteFile truncates it before writing.
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

func New(fs fs.FS) VIOUtil {
	return _VIOUtil{
		fs: fs,
	}
}

type _VIOUtil struct {
	fs fs.FS
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (
	b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}

// byName implements sort.Interface.
type byName []os.FileInfo

func (f byName) Len() int           { return len(f) }
func (f byName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func (this _VIOUtil) ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := this.fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Sort(byName(list))
	return list, nil
}

func (this _VIOUtil) ReadFile(filename string) ([]byte, error) {
	f, err := this.fs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// It's a good but not certain bet that FileInfo will tell us exactly how much to
	// read, so let's try it but be prepared for the answer to be wrong.
	var n int64

	if fi, err := f.Stat(); err == nil {
		// Don't preallocate a huge buffer, just in case.
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	}
	// As initial capacity for readAll, use n + a little extra in case Size is zero,
	// and to avoid another allocation after Read has filled the buffer. The readAll
	// call will read into its allocated internal buffer cheaply. If the size was
	// wrong, we'll either waste some space off the end or reallocate as needed, but
	// in the overwhelmingly common case we'll get it just right.
	return readAll(f, n+bytes.MinRead)
}

func (this _VIOUtil) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := this.fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
