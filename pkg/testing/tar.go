package testing

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"log"
)

// TestFile is a simple struct used in TarArchive function
type TestFile struct {
	Name string
	Body string
}

// TarArchive creates a tar buffer with given files
func TarArchive(files []TestFile) io.ReadCloser {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}
	return ioutil.NopCloser(&buf)
}
