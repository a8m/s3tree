package main

// filesystem
// ---------
// fs = {"path": "node"}
// break the problem into pieces.

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
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
		// it's a file
		if i > 0 && i == len(dirs)-1 {
			f = &file{ /*o*/ nil, d, filePath, nil}
		} else {
			f = &file{nil, d, filePath, make([]*file, 0)}
		}
		// add to parent
		fs.files[filePath] = f
		if dir, ok := fs.files[parentPath]; ok && i > 0 {
			dir.files = append(dir.files, f)
		}
	}
}
