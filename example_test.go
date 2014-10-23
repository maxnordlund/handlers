package handlers_test

import (
	"archive/tar"
	"compress/gzip"
	"github.com/maxnordlund/handlers"
	"log"
	"net/http"
	"os"
)

func ExampleNewTarFileSystem() {
	// Open some tar with the files to be served.
	gz, err := os.Open("website.tar.gz")
	if err != nil {
		log.Fatalln(err)
	}

	// Make sure to decompress it, if needed.
	archive, err := gzip.NewReader(gz)
	if err != nil {
		log.Fatalln(err)
	}

	// Pass the tar to the handler.
	fs, err := handlers.NewTarFileSystem(tar.NewReader(archive))
	if err != nil {
		log.Fatalln(err)
	}

	// Finally pass the handler to http.FileSystem to start serving the files.
	log.Fatalln(http.ListenAndServe(":8080", http.FileServer(fs)))
}
