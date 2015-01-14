package handlers

import (
	"archive/tar"
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
)

// A file represents the external view if a file in a tar.
// It implements the http.File interface.
type file struct {
	// File part
	*bytes.Reader
	info os.FileInfo

	// Directory part
	files []os.FileInfo
	index int
}

// See http://godoc.org/io#Closer
func (f *file) Close() error {
	return nil
}

// See http://godoc.org/os#File.Readdir
func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	index := f.index
	if index > len(f.files) {
		return []os.FileInfo{}, nil
	}
	if count > 0 {
		f.index += count
		if f.index > len(f.files) {
			return f.files[index:], io.EOF
		} else {
			return f.files[index:f.index], nil
		}
	} else {
		f.index = len(f.files)
		return f.files[index:], nil
	}
}

// See http://godoc.org/os#File.Stat
func (f *file) Stat() (os.FileInfo, error) {
	return f.info, nil
}

// A tarFile represents the internal view of a file in a tar.
type tarFile struct {
	// File part
	info   os.FileInfo
	buffer []byte

	// Directory part
	files []os.FileInfo
}

// A fileSystem represents all files within a tar.
// It the implements http.FileSystem interface.
type fileSystem map[string]*tarFile

// NewTarFileSystem returns a http.FileSystem backed by the provided tar file.
func NewTarFileSystem(tr *tar.Reader) (http.FileSystem, error) {
	fs := make(fileSystem)
	dirs := make(map[string][]os.FileInfo)
	for {
		buf := new(bytes.Buffer)
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return nil, err
		}
		io.Copy(buf, tr)
		dir := path.Dir(hdr.Name)
		info := hdr.FileInfo()
		dirs[dir] = append(dirs[dir], info)
		fs[hdr.Name] = &tarFile{
			info:   info,
			buffer: buf.Bytes(),
		}
	}
	for name, tf := range fs {
		tf.files = dirs[path.Dir(name)]
	}
	return fs, nil
}

// See http://godoc.org/net/http#FileSystem
func (fs fileSystem) Open(name string) (http.File, error) {
	tf, ok := fs[name]
	if !ok {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  os.ErrNotExist,
		}
	}
	return &file{
		Reader: bytes.NewReader(tf.buffer),
		info:   tf.info,
		files:  tf.files,
	}, nil
}
