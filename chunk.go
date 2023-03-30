package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/dustin/go-humanize"
)

func main() {
	log.SetFlags(0)
	var (
		off  int64
		end  int64
		lenf int64
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
	flag.Func("end", "Chunk end `offset`", parseFunc(&end))
	flag.Func("len", "Chunk `length`", parseFunc(&lenf))
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
	if err := parseFunc(&off)(flag.Arg(1)); err != nil {
		log.Fatalf("Error with offset: %s\n", err)
	}

	switch {
	case end == 0 && lenf == 0:
		log.Fatal("One of -end, -len must be given")
	case end != 0 && lenf != 0:
		log.Fatal("Only one of -end, -len may be given")
	case end < 0:
		log.Fatal("-end cannot be negative")
	case end < 0:
		log.Fatal("-end cannot be negative")
	case end > 0 && end <= off:
		log.Fatal("-end must be greater than the offset")
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

func parseNumber(s string) (int64, error) {
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		if f > float64(math.MaxInt64) {
			return 0, fmt.Errorf("float value too large for int64: %g", f)
		}
		return int64(f), nil
	}
	if u, err := humanize.ParseBytes(s); err == nil {
		return int64(u), nil
	}
	return 0, fmt.Errorf("cannot parse %q", s)
}
