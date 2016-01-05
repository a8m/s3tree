package main

import (
	"flag"
	"fmt"
	"github.com/a8m/tree"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

var (
	a = flag.Bool("a", false, "")
	d = flag.Bool("d", false, "")
	f = flag.Bool("f", false, "")
	s = flag.Bool("s", false, "")
	h = flag.Bool("h", false, "")
	u = flag.Bool("u", false, "")
	g = flag.Bool("g", false, "")
	Q = flag.Bool("Q", false, "")
	D = flag.Bool("D", false, "")

	// S3 args
	bucket = flag.String("b", "l2r", "")
	prefix = flag.String("p", "code", "")
	region = flag.String("region", "us-east-1", "")
)

func main() {
	flag.Parse()
	svc := s3.New(&aws.Config{Region: region})
	spin := NewSpin()
	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: bucket,
		Prefix: prefix,
	})
	spin.Done()
	var fs = NewFs()
	if err != nil {
		fmt.Println(err)
	} else {
		if l := len(resp.Contents); l != 0 {
			for i := 0; i < l; i++ {
				fs.addFile(resp.Contents[i])
			}
		}
	}
	var nd, nf int
	inf := tree.New(*prefix)
	opts := &tree.Options{Fs: fs, UnitSize: *u, LastMod: *D, OutFile: os.Stdout, DeepLevel: 3}
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
