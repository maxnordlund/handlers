package handlers

import (
	"archive/tar"
	"bytes"
	"log"
	"os"
	"time"
)

var (
	sampleTar      []byte
	sampleTarFiles []string
)

func init() {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	files := []struct {
		Path, Body string
	}{
		{"index.html", `<ul>
											<li><a href="./about/">About</a></li>
											<li><a href="./contact/">Contact</a></li>
										</ul>`},
		{"about/index.html", "We are four adventurers."},
		{"contact/index.html", "Somewhere over the rainbow."},
	}

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Path,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}

	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}
	sampleTar = buf.Bytes()
	sampleTarFiles = []string{
		"index.html",
		"about/index.html",
		"contact/index.html",
	}
}

func SampleTarFileSystem() ([]string, fileSystem) {
	fs, err := NewTarFileSystem(tar.NewReader(bytes.NewReader(sampleTar)))
	if err != nil {
		log.Fatalln(err)
	}
	return sampleTarFiles, fs.(fileSystem)
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi fileInfo) Name() string       { return fi.name }
func (fi fileInfo) Size() int64        { return fi.size }
func (fi fileInfo) Mode() os.FileMode  { return fi.mode }
func (fi fileInfo) ModTime() time.Time { return fi.modTime }
func (fi fileInfo) IsDir() bool        { return fi.mode.IsDir() }
func (fi fileInfo) Sys() interface{}   { return nil }

func SampleFile() *file {
	fooTime, err := time.Parse(time.Kitchen, "1:41PM")
	if err != nil {
		log.Fatalln(err)
	}
	barTime, err := time.Parse(time.Kitchen, "2:17PM")
	if err != nil {
		log.Fatalln(err)
	}
	bazTime, err := time.Parse(time.Kitchen, "5:39PM")
	if err != nil {
		log.Fatalln(err)
	}

	return &file{
		files: []os.FileInfo{
			fileInfo{"foo", 1337, 0666, fooTime},
			fileInfo{"bar", 1453, 0666, barTime},
			fileInfo{"baz", 4711, 0666, bazTime},
		},
	}
}

func (f *file) Files() []os.FileInfo {
	return f.files
}

func (f *file) Info() os.FileInfo {
	return f.info
}

func (tf *tarFile) Files() []os.FileInfo {
	return tf.files
}

func (tf *tarFile) Info() os.FileInfo {
	return tf.info
}
