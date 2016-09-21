package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
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
	reader := bufio.NewReader(r)
	p := make([]byte, 4*1024)
	for {
		n, err := reader.Read(p)
		if n == 0 {
			break
		}
		bytesRead := p[0:n]
		c.lines += bytes.Count(bytesRead, []byte{'\n'})
		c.bytes += n

		// len(bytes.Fields(bytesRead)) だと遅いので `bytes.Fields` の前半部分を真似してカウントする
		inField := false
		for i := 0; i < len(bytesRead); {
			r, size := utf8.DecodeRune(bytesRead[i:])
			wasInField := inField
			inField = !unicode.IsSpace(r)
			if inField && !wasInField {
				c.words += 1
			}
			i += size
		}

		if err == io.EOF {
			break
		}
		// fix word count between the read buffer
		next, err := reader.Peek(1)
		if err != nil && err != io.EOF {
			return false, err
		}
		if !unicode.IsSpace(rune(p[n-1 : n][0])) && !unicode.IsSpace(rune(next[0])) {
			c.words -= 1
		}
	}
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
