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

type Counter struct {
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

func (c *Counter) Count(r io.Reader) (bool, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return false, err
	}
	c.lines = bytes.Count(b, []byte{'\n'})
	c.bytes = len(b)
	c.words = len(bytes.Fields(b))
	return true, nil
}

func (c Counter) Show(opts FlagOptions, filename string) {
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

func (c *Counter) Add(src Counter) {
	c.lines += src.lines
	c.bytes += src.bytes
	c.words += src.words
}

func main() {
	opts := parseFlagOptions()

	var totalCount Counter

	filenames := flag.Args()
	if len(filenames) == 0 {
		var c Counter
		_, err := c.Count(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.Show(opts, "")
		return
	}

	for _, filename := range filenames {
		var c Counter
		fp, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = c.Count(fp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		totalCount.Add(c)
		c.Show(opts, filename)
		fp.Close()
	}

	if len(filenames) > 1 {
		totalCount.Show(opts, "total")
	}
}
