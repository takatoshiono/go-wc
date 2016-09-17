package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type FlagOptions struct {
	printLines bool
	printBytes bool
	printWords bool
}

type Count struct {
	lines int
	words int
	bytes int
}

func parseFlagOptions() FlagOptions {
	var opts FlagOptions
	flag.BoolVar(&opts.printLines, "l", false, "print lines")
	flag.BoolVar(&opts.printBytes, "c", false, "print bytes")
	flag.BoolVar(&opts.printLines, "w", false, "print words")
	flag.Parse()

	if !opts.printLines && !opts.printBytes && !opts.printWords {
		opts.printLines = true
		opts.printBytes = true
		opts.printWords = true
	}

	return opts
}

func count(r io.Reader) (Count, error) {
	var c Count
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return c, err
	}
	c.lines = bytes.Count(b, []byte{'\n'})
	c.bytes = len(b)
	c.words = len(bytes.Fields(b))
	return c, nil
}

func (c Count) Show(opts FlagOptions, filename string) {
	if opts.printLines {
		fmt.Printf(" %7d", c.lines)
	}
	if opts.printWords {
		fmt.Printf(" %7d", c.words)
	}
	if opts.printBytes {
		fmt.Printf(" %7d", c.bytes)
	}
	fmt.Printf(" %s\n", filename)
}

func (c *Count) Add(src Count) {
	c.lines += src.lines
	c.bytes += src.bytes
	c.words += src.words
}

func main() {
	opts := parseFlagOptions()

	var totalCount Count

	filenames := flag.Args()
	if len(filenames) == 0 {
		r, err := count(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		r.Show(opts, "")
		return
	}

	for _, filename := range filenames {
		fp, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}
		r, err := count(fp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		totalCount.Add(r)
		r.Show(opts, filename)
		fp.Close()
	}

	if len(filenames) > 1 {
		totalCount.Show(opts, "total")
	}
}
