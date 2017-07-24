package main

import (
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type file struct {
	*s3.Object
	name  string
	path  string
	files []*file
}

// Name return the s3-file name.
func (f *file) Name() string {
	return f.name
}

// Size return the size of the s3-file,
// if it's dir, it will return the sum of its childs
func (f *file) Size() (size int64) {
	if f.Object != nil {
		size = *f.Object.Size
	} else {
		for _, ff := range f.files {
			size += ff.Size()
		}
	}
	return
}

// Mode return os.ModeDir if it's dir
func (f *file) Mode() (o os.FileMode) {
	if f.IsDir() {
		o = os.ModeDir
	}
	return
}

// ModTime return the last modification
func (f *file) ModTime() (t time.Time) {
	if f.Object != nil {
		t = *f.LastModified
	}
	return
}

// Test if the given file is a directory
func (f *file) IsDir() bool {
	return f.files != nil && len(f.files) > 0
}

// Sys return dummy-stat_t
func (f *file) Sys() interface{} {
	return new(syscall.Stat_t)
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
