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
	o          = flag.String("o", "", "")            // Output to file
	a          = flag.Bool("a", false, "")           // All files
	d          = flag.Bool("d", false, "")           // Dirs only
	f          = flag.Bool("f", false, "")           // Full path
	s          = flag.Bool("s", false, "")           // Show byte size
	h          = flag.Bool("h", false, "")           // Show SI size
	Q          = flag.Bool("Q", false, "")           // Quote filename
	D          = flag.Bool("D", false, "")           // Show last mod
	C          = flag.Bool("C", false, "")           // Colorize
	L          = flag.Int("L", 0, "")                // Deep level
	U          = flag.Bool("U", false, "")           // No sort
	v          = flag.Bool("v", false, "")           // Version sort
	t          = flag.Bool("t", false, "")           // Last modification sort
	r          = flag.Bool("r", false, "")           // Reverse sort
	P          = flag.String("P", "", "")            // Matching pattern
	I          = flag.String("I", "", "")            // Ignoring pattern
	ignorecase = flag.Bool("ignore-case", false, "") // Ignore case-senstive
	dirsfirst  = flag.Bool("dirsfirst", false, "")   // Dirs first sort
	sort       = flag.String("sort", "", "")         // Sort by name or size
	// S3 args
	bucket = flag.String("b", "", "")
	prefix = flag.String("p", "", "")
	region = flag.String("region", "us-east-1", "")
)

var usage = `Usage: s3tree [options...] [paths...]

Options:
    --------- S3 options ----------
    -b		    s3 bucket(required).
    -p		    s3 prefix.
    --region name   aws region(default to us-east-1).
    ------- Listing options -------
    -a		    All files are listed.
    -d		    List directories only.
    -f		    Print the full path prefix for each file.
    -L		    Descend only level directories deep.
    -P		    List only those files that match the pattern given.
    -I		    Do not list files that match the given pattern.
    --ignore-case   Ignore case when pattern matching.
    -o filename	    Output to file instead of stdout.
    -------- File options ---------
    -Q		    Quote filenames with double quotes.
    -s		    Print the size in bytes of each file.
    -h		    Print the size in a more human readable way.
    -D		    Print the date of last modification or (-c) status change.
    ------- Sorting options -------
    -v		    Sort files alphanumerically by version.
    -t		    Sort files by last modification time.
    -U		    Leave files unsorted.
    -r		    Reverse the order of the sort.
    --dirsfirst	    List directories before files (-U disables).
    --sort X	    Select sort: name,size,version.
    ------- Graphics options ------
    -i		    Don't print indentation lines.
    -C		    Turn colorization on always.
`

func main() {
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
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
	// Output file
	var outFile = os.Stdout
	if *o != "" {
		outFile, err = os.Create(*o)
		if err != nil {
			errAndExit(err)
		}
	}
	defer outFile.Close()
	opts := &tree.Options{
		Fs:        fs,
		OutFile:   outFile,
		All:       *a,
		DirsOnly:  *d,
		FullPath:  *f,
		ByteSize:  *s,
		UnitSize:  *h,
		Quotes:    *Q,
		LastMod:   *D,
		Colorize:  *C,
		DeepLevel: *L,
		NoSort:    *U,
		ReverSort: *r,
		DirSort:   *dirsfirst,
		VerSort:   *v || *sort == "version",
		NameSort:  *sort == "name",
		SizeSort:  *sort == "size",
		ModSort:   *t,
		Pattern:   *P,
		IPattern:  *I,
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
	fmt.Fprintf(outFile, footer)
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
