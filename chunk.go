package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

func copyChunk(w io.Writer, r io.ReadSeeker, off, n int64) error {
	_, err := r.Seek(off, os.SEEK_SET)
	if err != nil {
		return err
	}
	_, err = io.CopyN(w, r, n)
	return err
}

func main() {
	end := flag.Int64("end", 0, "Ending offset")
	lenf := flag.Int64("len", 0, "Length of chunk")
	flag.Usage = func() {
		fmt.Printf(`Usage: %s [OPTIONS] FILENAME OFFSET
where OPTIONS are:
`, os.Args[0])
		flag.PrintDefaults()
		fmt.Println("(Exactly one of -end, -len must be given.)")
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
	off, err := strconv.ParseInt(flag.Arg(1), 10, 64)
	if err != nil {
		fmt.Printf("Cannot parse offset: %s\n", err)
		os.Exit(1)
	}

	switch {
	case *end == 0 && *lenf == 0:
		fmt.Println("One of -end, -len must be given.")
	case *end != 0 && *lenf != 0:
		fmt.Println("Only one of -end, -len may be given.")
	case *end < 0:
		fmt.Println("-end cannot be negative.")
	case *end < 0:
		fmt.Println("-end cannot be negative.")
	case *end > 0 && *end <= off:
		fmt.Println("-end must be greater than the offset.")
	default:
		goto after
	}
	os.Exit(1)

after:
	n := *lenf
	if n == 0 {
		n = *end - off
	}
	if err := copyChunk(os.Stdout, f, off, n); err != nil {
		fmt.Printf("Error while reading chunk: %s\n", err)
		os.Exit(1)
	}
}
