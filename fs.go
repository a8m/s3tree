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

// code/dev/config.json
// code/dev/worker/file-a.gz
// code/dev/worker/file-b.gz

// code, dev, config.json
// code
// code/dev
// code/dev/config.json
type Fs struct {
	files map[string]*file
}

func (f *Fs) Stat(path string) (os.FileInfo, error) {
	return f.files[path], nil
}

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

func (fs *Fs) addFile(o *s3.Object) {
	path := strings.Trim(*o.Key, "/")
	dirs := strings.Split(path, "/")
	for i, d := range dirs {
		var f *file
		var filePath = d
		parentPath := strings.Join(dirs[:i], "/")
		if parentPath != "" {
			filePath = parentPath + "/" + filePath
		}
		// it's a file
		if i > 0 && i == len(dirs)-1 {
			f = &file{ /*o*/ nil, d, filePath, nil}
		} else {
			f = &file{nil, d, filePath, make([]*file, 0)}
		}
		// add to parent
		if _, ok := fs.files[filePath]; !ok {
			fs.files[filePath] = f
			if dir, ok := fs.files[parentPath]; ok && i > 0 {
				dir.files = append(dir.files, f)
			}
		}
	}
}
