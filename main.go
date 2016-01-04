package main

import (
	"flag"
	"github.com/a8m/tree"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	fmt "github.com/k0kubun/pp"
	"os"
	"strings"
)

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
	fs := &Fs{make(map[string]*file)}
	if err != nil {
		fmt.Println(err)
	} else {
		if l := len(resp.Contents); l != 0 {
			for i := 0; i < l; i++ {
				if strings.HasPrefix(*resp.Contents[i].Key, "code/kinesis-scaling-utils") {
					continue
				}
				fs.addFile(resp.Contents[i])
			}
		}
	}
	var nd, nf int
	inf := tree.New("code")
	opts := &tree.Options{Fs: fs, UnitSize: *u, LastMod: *D, OutFile: os.Stdout, Colorize: true}
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
