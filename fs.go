package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"strings"
	"syscall"
	"time"
)

// S3-Object example
// { ETag: "\"3...-15\"",
// Key: "code/..../ironBeastFile.tar.gz",
// LastModified: 2014-06-25 14:32:12 +0000 UTC,
// Owner: { DisplayName: "...", ID: "07e122e2b7.... 0b6a31" },
// Size: 227849710,
// StorageClass: "GLACIER" }
type file struct {
	*s3.Object
	name  string
	path  string
	files []*file
}

func (f *file) Name() string { return f.name }
func (f *file) Size() int64 {
	if f.Object != nil {
		return *f.Object.Size
	}
	return 0
}
func (f *file) Mode() (o os.FileMode) { return }
func (f *file) ModTime() time.Time    { return *f.LastModified }
func (f *file) IsDir() bool {
	return f.files != nil && len(f.files) > 0
}
func (f *file) Sys() interface{} {
	var s *syscall.Stat_t
	return s
}

// Fs implement the tree.Fs interface
type Fs struct {
	files map[string]*file
}

// NewFs return new Fs instance.
func NewFs() *Fs { return &Fs{make(map[string]*file)} }

// Stat return "file" by the given path.
// "file" implemented the os.FileInfo
func (f *Fs) Stat(path string) (os.FileInfo, error) {
	if file, ok := f.files[path]; ok {
		return file, nil
	}
	return nil, os.ErrNotExist

}

// ReadDir return the list of files in the given dir-path.
func (f *Fs) ReadDir(path string) ([]string, error) {
	keys := []string{}
	dir, ok := f.files[path]
	if ok {
		for _, val := range dir.files {
			keys = append(keys, val.name)
		}
	}
	return keys, nil
}

// get s3.Object, split its path(Key) to dirs,
// and for each of them create a "file" and add
// it to Fs if not exists.
func (fs *Fs) addFile(path string, o *s3.Object) {
	path = strings.Trim(path, "/")
	dirs := strings.Split(path, "/")
	for i, d := range dirs {
		var f *file
		var filePath = d
		parentPath := strings.Join(dirs[:i], "/")
		if parentPath != "" {
			filePath = parentPath + "/" + filePath
		}
		if i > 0 && i == len(dirs)-1 {
			f = &file{o, d, filePath, nil}
		} else {
			f = &file{nil, d, filePath, make([]*file, 0)}
		}
		if _, ok := fs.files[filePath]; !ok {
			fs.files[filePath] = f
			if dir, ok := fs.files[parentPath]; ok && i > 0 {
				dir.files = append(dir.files, f)
			}
		}
	}
}

func (fs *Fs) isEmpty() bool {
	return len(fs.files) == 0
}
