package main

import (
	"errors"
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
	L = flag.Int("L", 0, "")

	// S3 args
	bucket = flag.String("b", "l2r", "")
	prefix = flag.String("p", "code", "")
	region = flag.String("region", "us-east-1", "")
)

func main() {
	flag.Parse()
	var noPrefix = len(*prefix) == 0
	if len(*bucket) == 0 {
		err := errors.New("-b(s3 bucket) is required.")
		errAndExit(err)
	}
	svc := s3.New(&aws.Config{Region: region})
	spin := NewSpin()
	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: bucket,
		Prefix: prefix,
	})
	spin.Done()
	var fs = NewFs()
	if err != nil {
		errAndExit(err)
	} else {
		// Loop over s3 object
		for _, obj := range resp.Contents {
			key := *obj.Key
			if noPrefix {
				key = fmt.Sprintf("%s/%s", *bucket, key)
			}
			fs.addFile(key, obj)
		}
	}
	var nd, nf int
	rootDir := *prefix
	if noPrefix {
		rootDir = *bucket
	}
	/*	if fs.isEmpty() {
		err := errors.New("no objects found in path: " + rootDir)
		errAndExit(err)
	}*/
	opts := &tree.Options{
		Fs:        fs,
		UnitSize:  *u,
		LastMod:   *D,
		OutFile:   os.Stdout,
		DeepLevel: *L,
	}
	inf := tree.New(rootDir)
	if d, f := inf.Visit(opts); f != 0 {
		nd, nf = nd+d-1, nf+f
	}
	inf.Print(opts)
	// print footer
	footer := fmt.Sprintf("\n%d directories", nd)
	if !opts.DirsOnly {
		footer += fmt.Sprintf(", %d files", nf)
	}
	fmt.Println(footer)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func errAndExit(err error) {
	fmt.Fprintf(os.Stderr, "s3tree: \"%s\"\n", err)
	os.Exit(1)
}
