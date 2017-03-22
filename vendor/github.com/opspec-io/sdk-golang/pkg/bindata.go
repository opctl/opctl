// Code generated by go-bindata.
// sources:
// pkg/data/packageManifest.json
// DO NOT EDIT!

package pkg

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _pkgDataPackagemanifestJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5c\x5f\x6f\xdb\x38\x12\x7f\xf7\xa7\x20\xdc\x3e\x24\xd7\x58\x4e\xdb\x34\x77\x97\xe2\x10\xe4\xba\x7b\x40\x0f\x77\x6d\xd0\x00\x3d\x60\xed\xb4\xa0\xa5\xb1\xcd\x8d\x44\xea\x48\xca\xa9\xbb\xe8\x77\x5f\x90\xd4\x5f\x8a\x92\x2d\x45\x69\xb3\xed\xee\xcb\xd6\x22\x67\x34\xf3\xe3\x8f\xc3\xe1\x90\xca\x6f\x23\x84\xc6\x8f\x85\xbf\x86\x08\x8f\xcf\xd0\x78\x2d\x65\x7c\x36\x9d\xfe\x2a\x18\x9d\x98\xa7\x1e\xe3\xab\x69\xc0\xf1\x52\x4e\x8e\x4f\xa6\xe6\xd9\xa3\xf1\x91\x92\x23\x41\x26\x22\xce\xa6\x53\x16\x8b\x18\x7c\x8f\xb0\xa9\xfa\xff\xf4\xd8\x7b\xea\x3d\x4f\xfb\x4f\x63\xec\xdf\xe0\x15\xfc\x17\x53\xb2\x04\x21\x3d\xa5\xdf\xe8\x90\x44\x86\xa0\xd4\x58\x5d\x4c\x6b\x00\xc2\xe7\x24\x96\x84\x51\xd5\x27\x6b\x44\x4b\xc6\x11\x46\xa9\x88\xe9\x1a\x73\x16\x03\x97\x04\xc4\xf8\x0c\x29\xb7\x10\x1a\x53\x1c\x41\xfe\xab\xae\xee\x0d\x8e\x00\xb1\x25\x92\x6b\x40\x2c\xd6\x6a\x74\x37\xb9\x8d\xb5\x49\x42\x72\x42\x57\x63\xfd\xf8\x8b\x69\xb5\x54\x34\x69\xfe\xa9\xf8\xd9\xf5\x05\x84\xc6\x89\x14\x65\xdd\x8f\x39\x2c\x55\xef\x47\xd3\x00\x96\x84\x12\xa5\x55\x4c\x63\xcc\x71\x54\x15\x65\x89\xec\x2d\xcb\x13\xba\x5b\x4e\xf8\x96\xb1\x1b\xe0\xa2\x1d\x89\xf7\xa6\x87\x0b\x85\x86\x77\x40\xf4\x1e\x78\xfa\x9a\x51\xfa\xaa\x31\x87\xff\x27\x84\x83\x22\xdc\xac\x34\xb6\x23\x84\xae\x75\x3b\x0e\x02\x2d\x8f\xc3\xcb\x32\x0f\x96\x38\x14\x90\x32\x29\x7f\x45\xc1\x8f\x80\xf0\x4b\x0d\x45\xc9\xfe\x9c\x90\x79\xe3\x51\xd3\x20\x13\x0e\xbe\x64\x7c\x8b\x34\x9e\x20\x81\x2b\x37\x31\x75\x8e\x35\x5b\xfc\x0a\xbe\x2c\x9e\x3b\xf8\x9a\xd9\x54\x79\xd0\xdc\xb5\x85\x8f\x79\xb3\x8b\x69\xd9\x7f\x5f\x8e\x6c\x55\x4b\x9c\x84\xd2\xa5\xa6\x46\x6f\xdd\x13\x6d\x70\x98\xc0\x4b\x84\x11\x87\x10\x4b\xb2\x01\x14\x63\xb9\x46\x84\x4a\xe0\x31\x07\x09\x01\x5a\x72\x16\xe9\x91\x0f\x08\x47\x84\xa2\xdb\x35\xf1\xd7\x29\x15\x10\x11\x48\xf1\xee\xe8\x2e\x56\x13\x71\x05\x3e\x87\x7d\xcc\x7e\xbd\xcc\x2c\x49\xc7\x8d\x08\x24\x8c\x70\xa3\x09\x0b\xc6\x42\xc0\xd4\xb2\x61\xd4\x60\x4f\x2b\x0f\x47\xb6\x78\x2e\xba\x93\xbe\xba\x53\x6d\x0e\xa0\x8c\x2f\xe9\xaf\xeb\xca\xdc\x5c\x92\x10\x9a\xd9\x5d\xb4\x36\xd1\xfb\x5f\x24\x84\x41\x99\xad\x5e\xf9\x27\xb5\xef\x95\xda\x0a\xe2\xef\x82\xd5\x9a\x2b\x4e\x5a\xd3\x24\x5a\x00\x7f\xc5\xa8\x90\x1c\x13\x5a\x5d\xee\x72\x7a\xd7\x7b\xf5\xa6\x2d\x0e\xc3\xb7\x4b\x9b\xb7\x16\xfc\xff\xbe\x7a\xfb\x06\x5d\xe9\x5c\x07\xcd\xb4\x00\xba\x81\xed\x2d\xe3\xc1\xf5\x41\x96\x1f\x49\xc6\x42\xe1\x11\x90\x4b\x9d\x52\xad\x65\x14\xa6\x79\xd5\xe7\x6d\x3c\x29\x65\x5c\x93\xe3\x93\x47\x02\x7c\xa5\x7b\xf2\xc2\x7b\xf6\xec\xb0\x32\x8e\xb9\xfd\x98\x73\xbc\xad\x36\x11\x09\x91\x63\x22\xb9\x17\xda\x3a\x44\x23\x17\x17\x4a\x4c\x18\x63\xba\xed\x08\x85\x12\x18\x0e\x8a\xe7\x0f\x07\x0a\xa0\x49\xd4\x05\x09\xd5\x7f\x38\x20\x8e\xef\x0c\x44\x26\x61\x5c\xdf\xed\xef\x92\xf1\x08\xdb\x91\x68\xcc\x28\x68\x42\xcc\x2a\xca\xeb\xa1\x3a\x9b\x96\x2a\x7c\xae\x80\xd7\x03\x93\x05\xdd\x3b\x13\x14\x84\x0e\x6a\xc6\x44\xb4\x00\xb5\x00\x35\x6a\xb0\x02\x6b\xad\x3d\x1d\xb0\x99\xf5\x1c\x15\x46\x59\x2d\xd7\x8d\xc1\xf1\xda\x89\x50\x84\x3f\x91\xa8\x1b\x29\x52\x91\xe1\x78\xd1\x40\x0b\x7b\x90\x2b\x66\x13\xda\xd9\x6c\x23\x32\x98\xd9\x27\x7d\xcc\x4e\x42\x49\xe2\x10\xba\xc5\xa3\x42\x6a\x30\xe3\x9f\xf6\x30\x9e\xb2\xda\x4c\x6a\xb3\x9a\x32\x39\x1c\x45\x5e\x58\xf6\x76\x8e\x88\x65\x47\xb2\xf9\xbf\xb7\x2b\x5a\x60\x38\x67\x9a\x98\xf3\xb5\x16\x84\x4e\x29\x8f\x23\x8d\x69\xce\xcf\xcb\xed\x4d\x19\xfa\x1b\x13\x18\x87\xcc\xd1\x53\xd6\x3e\xb8\x2c\xfd\x7e\xd3\xe6\x74\x85\xe9\x9d\x38\x5b\x06\xf8\xce\xe4\x34\x6f\xee\x41\xb9\x87\x95\x9c\x57\x43\x5b\x35\x3d\x17\xcc\xbf\x01\xd9\xcc\xeb\x72\xfb\x4e\x96\x5a\xa3\x75\xa5\x65\x5b\xf9\xde\xc4\x6b\xf3\xda\x6f\xc4\xeb\xee\x84\x34\xe6\x7e\x17\x3b\xb9\x14\x79\x37\x59\x34\x78\xbb\xf6\x72\xf5\x5e\x7f\xee\xe5\xec\xa2\x69\x0d\xa2\x1f\x76\x2f\xd7\x03\x8a\xef\x65\x2f\x57\x8f\x45\x5d\xf6\x72\x76\xb0\x4d\x62\xe0\x02\xa4\x0a\xb2\x15\xef\x8d\xf4\x40\xfe\xff\xd5\x72\xbe\xeb\x86\x32\xc0\x12\x26\x92\x44\xb0\x73\x4b\x59\xf1\x21\x17\x43\xc6\x9b\xbb\x7a\xe1\x3d\xb7\xf7\x01\xae\x81\xe9\xb0\x27\x2d\xfc\x6a\xdf\x95\x1e\xed\x8b\x93\x8a\xc2\x7c\x42\x22\xbc\x82\x89\x9a\x37\xbb\xe0\xba\x40\x46\x04\x69\x11\xc4\x61\x09\x1c\xa8\x0f\x08\x0b\xa4\xa7\x1b\x04\x68\xb1\x45\xb3\x15\x91\xeb\x64\xe1\xf9\x2c\x9a\x1a\x81\x69\x40\x94\xbb\x8b\x44\x69\x9a\xe6\x72\x05\xc2\x3b\x24\x24\x07\xc8\x1a\x9e\x7a\x4f\x9f\x17\x2a\x86\x05\xd8\x06\x64\x18\x9c\x21\xc2\x24\xec\xc6\x45\x2d\x32\x1c\x0f\x9f\x0d\x0a\x93\xf1\x67\x18\x6c\xd6\x4c\x48\x7d\x6a\xd7\x09\x9e\x4c\x6a\x38\x84\x9e\x0f\x8a\x50\xee\xd5\x30\x20\x91\x78\x73\xd2\x0d\x20\x25\x31\x1c\x38\x27\x83\x82\xa3\xbd\x19\x0c\x98\xd3\xce\xc0\x9c\x0e\x07\xcc\x8b\xa1\x81\x39\x1d\x08\x98\x84\x93\x6e\xb8\x24\x9c\x0c\x07\xcb\xe9\xa0\xb0\x28\x5f\x86\x41\x45\x40\xb4\xd9\xa3\xd0\x7c\x81\x04\x44\x98\x4a\xe2\xa3\xf4\x2e\x83\xbd\xc0\x19\x45\x0a\x15\x83\xd6\xd9\x74\x5a\x3c\x9a\x0e\xea\x7d\x6a\x73\x3b\x00\x23\x57\x8b\x55\x87\xfe\x0f\xd0\x95\x5c\x77\x2c\xe9\x1a\xa1\xc1\xf2\xda\xd3\x86\xb4\xd6\x55\xc0\x2f\x15\xa1\x9f\xba\x7d\xca\xcc\xfb\xb6\x3e\xd9\xd9\xea\x9e\x3e\x1d\x1f\x55\x4d\xce\xaa\x5c\xc7\x4e\x5f\x1f\x7e\x71\xb8\x65\x8b\xf5\x23\x14\x87\x7b\xec\x30\x63\x2c\x25\x70\xbb\x90\xd4\x0a\x46\x2a\x32\x18\x1c\x7f\x6b\x40\xc3\x11\xa6\x8a\x0d\xe2\x98\xc3\x0a\x3e\x8d\x7b\x56\x86\x1c\xb5\x9e\x96\xc2\x60\xa9\xbd\x73\x61\x50\xcb\xf6\x2b\x0c\x1a\xf7\x7f\xb0\x82\xb7\xd1\xfa\x6d\x0b\xde\xad\xd3\xe8\x81\xd5\x30\x2b\x83\x50\xad\x61\xc6\x36\xa3\x2d\xc4\x2f\xfb\x9c\xce\x98\xb9\x7f\xd9\xc0\xd9\xd9\xe4\xa3\x87\x27\x9f\x2f\x26\xbf\x1c\x4f\xfe\x7e\xfd\xa4\xe7\xa1\xbc\x7b\x54\xf2\x5b\x8f\x9d\xd2\x2e\xb7\xae\xe2\x8e\xd9\x00\xca\xca\x07\x62\x03\xa8\x2b\x9f\x43\x0c\xa1\xae\x14\xbd\xf6\xc9\xd7\xee\x10\x47\xfd\x95\x3b\x7e\xfa\xab\x96\xf8\x88\x55\x86\xeb\xe3\x30\x44\x2b\x8e\xe3\x75\xce\xc5\x97\x48\x00\xa0\x6c\x69\x01\xea\xdd\x92\x1b\x12\x43\x40\xcc\xbd\x6f\xf5\x6b\xfa\x0a\x87\xe1\x47\x2d\x56\xbc\xa0\xce\xb1\xdf\x76\xe7\x0d\x3e\xa3\x12\x13\x0a\x5c\x69\x74\xe6\x0c\x7b\x28\x61\xf1\x5d\xa4\xd5\x6c\x0d\x43\x08\xef\xa2\x43\x00\x27\xd8\xd6\xe0\x8c\x0d\x55\x87\x5d\xa3\x56\xed\xd1\xfb\x5c\x23\x57\xd3\x65\x21\xf3\xa3\x60\x8f\x15\xe3\x15\x8b\x22\x4c\x03\xc4\x13\xaa\xf6\x43\x18\xe5\xef\x7a\x89\xd8\x06\x38\x27\x01\x08\x84\xe9\x16\x09\x90\x08\x4b\xbd\xbe\x98\xa2\x61\x08\x1b\x70\x14\xc3\x9a\xf3\x30\xd4\x9c\x8b\xb5\x8c\x87\xb9\x93\xc9\x42\x2c\x21\xb8\x72\x2c\x97\xd5\x49\xe8\x58\x85\x09\x77\x2e\x5b\x0d\x83\x50\x60\xdb\x1a\xa5\xf3\x6e\x1f\x0e\x66\x26\x58\x5f\x9f\x1d\x9e\xab\xd0\x3d\x9f\x4f\x4b\xd1\xfb\xb1\x53\xca\xb1\x53\x5d\x08\x16\x26\x12\x0a\xfc\xf5\xbd\x53\x7d\x2b\x35\xc2\x71\x0c\x01\x92\x0c\x61\xfd\x70\x83\x39\xc1\x8b\x10\x10\xe3\x08\x67\x72\xe9\xd7\x12\xb9\x94\x87\x5e\x2f\x11\x65\xb5\xe7\x28\xe6\x6c\x43\x02\x08\x8e\x5c\xaf\xba\x25\x61\x88\x16\x80\x12\x01\x81\x57\xc3\x04\xed\x4a\x5d\xea\xa3\x51\x1b\x8f\x7d\x57\x78\xd7\x50\x02\xdd\xbc\xc7\xf7\x3b\x9a\xb3\x0f\xff\xe8\x30\x68\x14\x01\xdd\x10\xce\x68\x04\x54\x16\xe3\x42\x68\x79\x22\xb9\x61\xec\x4b\xf6\x7b\x05\x58\x2d\xe9\x7f\xcc\xc9\xb2\x34\xd7\xd9\xad\xd9\xa2\x9f\xb6\x4e\x97\x5c\xce\x9e\x2f\x85\x42\xd7\x84\x29\x5a\x1f\xf6\x8c\xd1\x91\xba\xd7\x80\xee\x18\x49\xc3\xde\x7d\x06\xeb\x5d\x7e\xb6\xa4\x86\x84\x9a\xc5\x63\xf0\x39\x51\xd7\x37\x8e\x93\x30\x7c\x1d\x00\x95\x44\x6e\xf7\xb4\x35\xeb\xae\x47\x53\xd9\xab\x74\x14\x4b\xde\xd7\xb2\xba\x71\xbf\xe7\xb2\xd9\x74\xfe\x5a\x16\xef\xa4\xaa\x73\x8b\x55\x6a\xed\x76\x14\x97\xde\x2f\xb9\xcf\x98\x34\x3b\x6b\xde\x6c\xe5\x9d\x6a\x59\x53\x16\x08\xd2\xab\x3c\x38\x08\x38\x08\x51\x09\x3e\x69\x53\x16\x7e\x1e\x5c\x6c\x10\x32\xf8\x99\xdb\x29\xe5\x90\xc8\x7e\xf0\xf6\x8e\xe8\x25\x40\x65\x00\x9c\xa3\x98\xc3\x92\x7c\xaa\xe2\x69\xaa\x1a\x0f\x19\xcf\xb7\x49\x6b\xb5\xe7\x1b\xe1\xc9\x12\xf9\x47\xc4\xf3\x96\xf1\x9b\x9f\x6a\x9f\x4b\xba\xbc\xfd\x1f\xe3\x37\xca\x95\xa0\xf4\xc9\xa6\x5c\xa3\x83\xea\x06\xa6\x74\xec\xa3\x63\xe3\xee\xc3\x9d\xc6\xad\x7e\xb5\x62\xd5\x18\xf2\xd2\xe5\xb7\x5c\x23\x18\xa2\xd4\xe5\xae\x62\x15\xf9\xe6\xc8\x7a\xd7\xfe\xd5\x07\x5d\xed\xfa\x27\xa1\x41\xb5\x68\x5a\xfe\x8c\xbb\xd4\xa1\xa9\x1e\x71\x61\xbe\xe3\xe0\x13\xe1\xb3\x18\xd0\xc2\x74\x7f\x89\x6e\xd7\x40\x11\x4d\xc2\xf0\xc8\x54\x29\x71\x04\x48\x9f\xfd\x67\x39\x14\x16\x22\x89\x20\x50\xf9\xf3\x82\xc9\x35\xd2\xf2\xfb\xdc\x07\xfc\x1a\x05\xb5\xe2\xe6\x7f\x18\x76\xab\x29\xed\xcb\xaa\x41\x0a\x48\x69\x01\xc5\x35\x78\x69\x53\xef\x32\x04\x8b\x6d\xe8\x5a\x02\x5b\x6b\x6d\xc2\x9d\x3d\xee\xca\x1c\x4b\xb5\xd5\xfd\xb0\xb5\x93\x61\xfb\xdb\xfb\xbc\xa5\xe5\x3b\xfa\x8c\xec\xad\x9a\xeb\x9f\xe6\xf7\x57\xdd\x23\xcc\x58\x59\xd5\x3d\x06\x19\x16\xdf\x31\xba\xe4\xd5\xb9\xa6\xe8\x92\x77\xe8\x4d\xd3\x4c\x4b\x13\x59\xef\x7c\x40\xe8\x37\x5d\xc2\xdc\x13\xc4\xdc\xc0\xfe\x50\x96\x8a\x94\xce\x6a\x71\xd1\xdc\x1b\x46\xa3\xe3\xc1\x82\x98\x9a\xd7\x1f\x42\xc7\x76\xc7\x05\xa5\xa3\xdb\x8e\x73\x4a\x91\xc4\x31\xe3\x52\xfd\xf3\xf1\xc1\xeb\x37\x1f\xaf\x5e\xbd\xbd\xfc\xf9\xe3\xfb\x8b\x77\x87\x48\x24\x0b\x21\x89\xd4\x57\x21\x91\xd8\x52\x89\x3f\xa9\x55\x91\x43\xad\x67\xb6\x24\x72\x88\x43\xec\x9b\xac\x45\x2d\x9a\xfa\xc3\x75\xc4\x96\xa8\xdc\x1d\x61\x89\x78\x42\x25\x89\xc0\xdb\xef\xcf\x90\xa4\x7f\x86\xc3\xcd\x1d\xdd\xd4\xe4\xe4\xe6\x99\x77\xec\x1d\xd7\xef\xd5\x1c\x64\x27\x0d\xd5\x1b\x34\xfa\x8f\xc5\x18\x19\x6f\x2d\xa3\xf0\xb0\xd1\x3e\x7b\x35\x57\x4d\x1f\x0e\x66\x7a\xc9\x3e\x9c\xcf\x3d\xc7\x3f\x0f\xce\xcf\x0e\xe6\xf3\x89\xfa\x75\x31\xf9\x05\x4f\x3e\x4f\xae\x9f\x1c\x9c\x9f\xcd\xe7\x5e\xe5\xd1\xe1\x5f\x0e\x0f\xcf\xf5\xf3\x27\xa5\xe7\xf3\xf9\x64\x3e\xf7\xae\x9f\x1c\x9e\x3f\x2e\xfd\x39\x92\xd1\x97\xd1\xe8\xf7\x00\x00\x00\xff\xff\x9c\xa4\x13\x0b\x36\x47\x00\x00")

func pkgDataPackagemanifestJsonBytes() ([]byte, error) {
	return bindataRead(
		_pkgDataPackagemanifestJson,
		"pkg/data/packageManifest.json",
	)
}

func pkgDataPackagemanifestJson() (*asset, error) {
	bytes, err := pkgDataPackagemanifestJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pkg/data/packageManifest.json", size: 18230, mode: os.FileMode(420), modTime: time.Unix(1489886942, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"pkg/data/packageManifest.json": pkgDataPackagemanifestJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"pkg": &bintree{nil, map[string]*bintree{
		"data": &bintree{nil, map[string]*bintree{
			"packageManifest.json": &bintree{pkgDataPackagemanifestJson, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}