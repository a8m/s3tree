package main

import (
	"flag"
	"fmt"
	"github.com/a8m/tree"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"strings"
	"syscall"
	"time"
)

type file struct {
	*s3.Object
	path  string
	size  int64
	isDir bool
	files []*file
}

func (f *file) Name() string          { return f.path }
func (f *file) Size() int64           { return f.size }
func (f *file) Mode() (o os.FileMode) { return }
func (f *file) ModTime() time.Time    { return *f.LastModified }
func (f *file) IsDir() bool           { return f.isDir }
func (f *file) Sys() interface{} {
	var s *syscall.Stat_t
	return s
}

var storage = map[string]*file{}

type fs struct{}

func (f *fs) Stat(path string) (os.FileInfo, error) {
	return storage[path], nil
}

func (f *fs) ReadDir(path string) ([]string, error) {
	keys := []string{}
	for key, val := range storage {
		if key == (path + "/" + val.path) {
			keys = append(keys, val.path)
		}
	}
	return keys, nil
}

var (
	a      = flag.Bool("a", false, "")
	d      = flag.Bool("d", false, "")
	f      = flag.Bool("f", false, "")
	s      = flag.Bool("s", false, "")
	h      = flag.Bool("h", false, "")
	p      = flag.Bool("p", false, "")
	u      = flag.Bool("u", false, "")
	g      = flag.Bool("g", false, "")
	Q      = flag.Bool("Q", false, "")
	D      = flag.Bool("D", false, "")
	inodes = flag.Bool("inodes", false, "")
	device = flag.Bool("device", false, "")
)

func main() {
	flag.Parse()
	r := "us-east-1"
	b := "l2r"
	pre := "code"
	svc := s3.New(&aws.Config{Region: &r})
	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &b,
		Prefix: &pre,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		if l := len(resp.Contents); l != 0 {
			for i := 0; i < l; i++ {
				if strings.HasPrefix(*resp.Contents[i].Key, "code/kinesis-scaling-utils") {
					continue
				}
				paths := strings.Split(*resp.Contents[i].Key, "/")
				for j, path := range paths {
					if path == "" {
						continue
					}
					if j == len(paths)-1 {
						fmt.Println(resp.Contents[i])
						storage[*resp.Contents[i].Key] = &file{
							Object: resp.Contents[i],
							size:   *resp.Contents[i].Size,
							path:   path}
					} else {
						storage[strings.Join(paths[:j+1], "/")] = &file{
							path:  path,
							isDir: true,
						}
					}
				}
			}
		}
		//		fmt.Println(storage)
	}
	var nd, nf int
	inf := tree.New("code")
	opts := &tree.Options{Fs: new(fs), UnitSize: *u, LastMod: *D, OutFile: os.Stdout}
	if d, f := inf.Visit(opts); f != 0 {
		nd, nf = nd+d-1, nf+f
	}
	inf.Print("", opts)
	// print footer
	footer := fmt.Sprintf("\n%d directories", nd)
	if !opts.DirsOnly {
		footer += fmt.Sprintf(", %d files", nf)
	}
	fmt.Println(footer)
}
