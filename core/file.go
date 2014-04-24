package space

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v1"
)

type File struct {
	Path     string
	Buffer   *bytes.Buffer
	Info     FileInfo
	realpath string
	status   int
	Metadata *Metadata
}

func (f *File) Status(i int) {
	f.status = i
}

func (f *File) Written() bool {
	return f.status != 0
}

func (f *File) Write() (err error) {
	path, _ := filepath.Split(f.Path)
	ospath := filepath.FromSlash(path)

	if ospath != "" {
		err = os.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Create(f.Path)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, f.Buffer)
	if err != nil {
		f.status = 200
	}
	return err
}

func (f *File) Read() (err error) {
	fh, err := os.Open(f.realpath)
	if err != nil {
		return err
	}
	r := bufio.NewReader(fh)
	defer fh.Close()

	// parse front-matter
	contents, metadata, err := FrontMatterParser(r)

	if err != nil {
		return err
	}

	f.Metadata = &Metadata{}
	if metadata != nil {
		// parse yaml
		err = yaml.Unmarshal(metadata.Bytes(), f.Metadata)
		if err != nil {
			return err
		}
	}

	_, err = f.Buffer.ReadFrom(contents)
	return err
}
