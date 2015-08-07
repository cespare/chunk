package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/dustin/go-humanize"
)

func main() {
	var (
		off  Int64
		end  Int64
		lenf Int64
	)
	flag.Var(&end, "end", "Ending offset")
	flag.Var(&lenf, "len", "Length of chunk")
	flag.Usage = func() {
		fmt.Printf(`Usage: %s [OPTIONS] FILENAME OFFSET
where OPTIONS are:
`, os.Args[0])
		flag.PrintDefaults()
		fmt.Println(`Exactly one of -end, -len must be given.
Numbers may be written as 1000, 1e3, or 1kB.`)
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	if err := off.Set(flag.Arg(1)); err != nil {
		Fatalf("Error with offset: %s\n", err)
	}

	switch {
	case end == 0 && lenf == 0:
		Fatalln("One of -end, -len must be given.")
	case end != 0 && lenf != 0:
		Fatalln("Only one of -end, -len may be given.")
	case end < 0:
		Fatalln("-end cannot be negative.")
	case end < 0:
		Fatalln("-end cannot be negative.")
	case end > 0 && end <= off:
		Fatalln("-end must be greater than the offset.")
	}

	n := lenf
	if n == 0 {
		n = end - off
	}
	sr := io.NewSectionReader(f, int64(off), int64(n))
	if _, err := io.Copy(os.Stdout, sr); err != nil {
		fmt.Printf("Error while reading chunk: %s\n", err)
		os.Exit(1)
	}
}

type Int64 int64

func (n *Int64) String() string { return fmt.Sprintf("%d", *n) }
func (n *Int64) Set(s string) error {
	nn, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		*n = Int64(nn)
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		if f > float64(math.MaxInt64) {
			return fmt.Errorf("float value too large for float64: %g", f)
		}
		*n = Int64(f)
		return nil
	}
	u, err := humanize.ParseBytes(s)
	if err == nil {
		*n = Int64(u)
		return nil
	}
	return fmt.Errorf("cannot parse %q", s)
}

func Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func Fatalln(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}
