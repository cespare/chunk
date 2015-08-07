package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

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
		Fatalln("One of -end, -len must be given.")
	case *end != 0 && *lenf != 0:
		Fatalln("Only one of -end, -len may be given.")
	case *end < 0:
		Fatalln("-end cannot be negative.")
	case *end < 0:
		Fatalln("-end cannot be negative.")
	case *end > 0 && *end <= off:
		Fatalln("-end must be greater than the offset.")
	}

	n := *lenf
	if n == 0 {
		n = *end - off
	}
	if _, err := io.Copy(os.Stdout, io.NewSectionReader(f, off, n)); err != nil {
		fmt.Printf("Error while reading chunk: %s\n", err)
		os.Exit(1)
	}
}

func Fatalln(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}
