package handlers_test

import (
	"archive/tar"
	"bytes"
	"github.com/maxnordlund/handlers"
	"io"
	"testing"
)

func TestFile_Close(t *testing.T) {
	file := handlers.SampleFile()
	if err := file.Close(); err != nil {
		t.Fatalf("TestFile_Close: err should always be nil, was %s", err)
	}
}

func TestFile_Readdir_lessThenZero(t *testing.T) {
	file := handlers.SampleFile()
	files, err := file.Readdir(-1)
	if err != nil {
		t.Fatalf("TestFile_Readdir < 0: err should be nil, was %s", err)
	}
	if len(files) != len(file.Files()) {
		t.Fatalf("TestFile_Readdir < 0: len(files) should be 1, was %d", len(files))
	}
	for i, facit := range file.Files() {
		if files[i] != facit {
			t.Fatalf("TestFile_Readdir < 0: files[%d] should equal file.files[%d], was %v", i, i, files[i])
		}
	}
	files, err = file.Readdir(0)
	if err != nil {
		t.Fatalf("TestFile_Readdir < 0: err should be nil, was %s", err)
	}
	if len(files) != 0 {
		t.Fatalf("TestFile_Readdir < 0: len(files) should be 0, was %d", len(files))
	}
}

func TestFile_Readdir_greaterThenZero(t *testing.T) {
	file := handlers.SampleFile()
	files, err := file.Readdir(1)
	if err != nil {
		t.Fatalf("TestFile_Readdir > 0: err should be nil, was %s", err)
	}
	if len(files) != 1 {
		t.Fatalf("TestFile_Readdir > 0: len(files) should be 1, was %d", len(files))
	}
	if files[0] != file.Files()[0] {
		t.Fatalf("TestFile_Readdir > 0: files[0] should equal file.files[0], was %v", files[0])
	}
	files, err = file.Readdir(3)
	if err != io.EOF {
		t.Fatalf("TestFile_Readdir > 0: err should be io.EOF, was %s", err)
	}
	if len(files) != 2 {
		t.Fatalf("TestFile_Readdir > 0: len(files) should be 2, was %d", len(files))
	}
	if files[0] != file.Files()[1] {
		t.Fatalf("TestFile_Readdir > 0: files[0] should equal file.files[1], was %v", files[0])
	}
	if files[1] != file.Files()[2] {
		t.Fatalf("TestFile_Readdir > 0: files[1] should equal file.files[2], was %v", files[1])
	}
	files, err = file.Readdir(1)
	if err != nil {
		t.Fatalf("TestFile_Readdir > 0: err should be nil, was %s", err)
	}
	if len(files) != 0 {
		t.Fatalf("TestFile_Readdir > 0: len(files) should be 0, was %d", len(files))
	}
}

func TestFile_Stat(t *testing.T) {
	file := handlers.SampleFile()
	info, err := file.Stat()
	if err != nil {
		t.Fatalf("TestFile_Stat: err should always be nil, was %s", err)
	}
	if info != file.Info() {
		t.Fatalf("TestFile_Stat: info should equal file.info, was %v", info)
	}
}

func TestFileSystem_goodArchive(t *testing.T) {
	files, fs := handlers.SampleTarFileSystem()
	for _, name := range files {
		file, ok := fs[name]
		if !ok {
			t.Fatalf("TestFileSystem_goodArchive: should contain the %q file, none found", name)
		}
		info := file.Files()[0]
		if info != file.Info() {
			t.Fatalf("TestFileSystem_goodArchive: file.files[0].info should equal file.info, was %v", info)
		}
	}
}

func TestFileSystem_badArchive(t *testing.T) {
	_, err := handlers.NewTarFileSystem(tar.NewReader(bytes.NewReader([]byte{0})))
	if err == nil {
		t.Fatalf("TestFileSystem_badArchive: err should not be nil, was %s", err)
	}
}

func TestFileSystem_Open(t *testing.T) {
	files, fs := handlers.SampleTarFileSystem()
	for _, name := range files {
		file, err := fs.Open(name)
		if err != nil {
			t.Fatalf("TestFileSystem_Open: err should be nil, was %s", err)
		}
		info, err := file.Stat()
		if err != nil {
			t.Fatalf("TestFileSystem_Open: err should always be nil, was %s", err)
		}
		if info != fs[name].Info() {
			t.Fatalf("TestFileSystem_Open: file.Stat() should equal file.info, was %v", info)
		}
	}
	if file, err := fs.Open("fubar"); err == nil {
		t.Fatalf("TestFileSystem_Open: should not contain \"fubar\", was %v", file)
	}
}
