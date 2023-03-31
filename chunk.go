package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

func main() {
	log.SetFlags(0)
	var (
		start int64
		end   int64
		n     int64
	)
	parseFunc := func(p *int64) func(string) error {
		return func(s string) error {
			n, err := parseNumber(s)
			if err != nil {
				return err
			}
			*p = n
			return nil
		}
	}
	flag.Func("start", "Chunk start `offset`", parseFunc(&start))
	flag.Func("end", "Chunk end `offset`", parseFunc(&end))
	flag.Func("len", "Chunk `length`", parseFunc(&n))
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage:

  chunk [flags ...] filename

where the flags are:
`)
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Numbers may be written as 1000, 1e3, or 1kB. The -start and -end flags
may be negative (to indicate offsets relative to the end of the file).

All flags are optional, but -end and -len are mutually exclusive.
If -start is not given, it defaults to the beginning of the file.
If -end and -len are not given, they default to the end of the file.
`)
	}
	flag.Parse()

	if end > 0 && n > 0 {
		log.Fatal("-end and -len are mutually exclusive")
	}
	if start > 0 && end > 0 && end < start {
		log.Fatal("-end cannot be before -start")
	}
	if n < 0 {
		log.Fatalln("-len cannot be negative")
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if start < 0 {
		start += stat.Size()
		if start < 0 {
			log.Println("Start offset would be before beginning of file; setting start to 0")
			start = 0
		}
	}
	if end == 0 && n == 0 {
		end = stat.Size()
	}
	if end < 0 {
		end += stat.Size()
		if end < 0 {
			log.Println("End offset would be before beginning of file; setting end to 0")
			end = 0
		}
	}
	if n == 0 {
		n = end - start
	}

	log.Printf("Selecting chunk of size %d starting at %d", n, start)
	sr := io.NewSectionReader(f, start, n)
	if _, err := io.Copy(os.Stdout, sr); err != nil {
		log.Fatalf("Error while reading chunk: %s", err)
	}
}

func parseNumber(s string) (int64, error) {
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		if f > float64(math.MaxInt64) || f < float64(math.MinInt64) {
			return 0, fmt.Errorf("float value out of bounds for int64: %g", f)
		}
		return int64(f), nil
	}
	s, neg := strings.CutPrefix(s, "-")
	if u, err := humanize.ParseBytes(s); err == nil {
		n := int64(u)
		if neg {
			n = -n
		}
		return n, nil
	}
	return 0, fmt.Errorf("cannot parse %q", s)
}
