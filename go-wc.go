package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"unicode"
)

const version = "0.0.1"

type FlagOptions struct {
	printLines bool
	printBytes bool
	printWords bool
}

type Counter struct {
	lines int
	words int
	bytes int
	mux   sync.Mutex
}

func (c *Counter) Count(r io.Reader) (bool, error) {
	reader := bufio.NewReader(r)
	var wg sync.WaitGroup
	for {
		p := make([]byte, 4*1024)
		n, err := reader.Read(p)
		if n == 0 {
			break
		}
		c.AddBytes(n)

		wg.Add(1)
		go func() {
			var localCounter = &Counter{}
			bytesRead := p[:n]
			inField := false
			for i := 0; i < len(bytesRead); {
				b, size := bytesRead[i : i+1][0], 1
				wasInField := inField
				inField = !unicode.IsSpace(rune(b))
				if inField && !wasInField {
					localCounter.words += 1
				}
				if rune(b) == '\n' {
					localCounter.lines += 1
				}
				i += size
			}
			c.Add(localCounter)
			wg.Done()
		}()

		if err == io.EOF {
			break
		}

		// fix word count between the read buffer
		next, err := reader.Peek(1)
		if err != nil && err != io.EOF {
			return false, err
		}
		if !unicode.IsSpace(rune(p[n-1 : n][0])) && !unicode.IsSpace(rune(next[0])) {
			c.AddWords(-1)
		}
	}
	wg.Wait()
	return true, nil
}

func (c *Counter) Show(opts *FlagOptions, filename string) {
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

func (c *Counter) Add(src *Counter) {
	c.mux.Lock()
	c.lines += src.lines
	c.bytes += src.bytes
	c.words += src.words
	c.mux.Unlock()
}

func (c *Counter) AddLines(n int) {
	c.mux.Lock()
	c.lines += n
	c.mux.Unlock()
}

func (c *Counter) AddBytes(n int) {
	c.mux.Lock()
	c.bytes += n
	c.mux.Unlock()
}

func (c *Counter) AddWords(n int) {
	c.mux.Lock()
	c.words += n
	c.mux.Unlock()
}

func parseFlagOptions() *FlagOptions {
	var opts = &FlagOptions{false, false, false}

	flag.BoolVar(&opts.printLines, "l", false, "print lines")
	flag.BoolVar(&opts.printBytes, "c", false, "print bytes")
	flag.BoolVar(&opts.printWords, "w", false, "print words")
	flag.Parse()

	if !opts.printLines && !opts.printBytes && !opts.printWords {
		opts.printLines = true
		opts.printBytes = true
		opts.printWords = true
	}

	return opts
}

func main() {
	opts := parseFlagOptions()

	var totalCount = &Counter{}

	filenames := flag.Args()
	if len(filenames) == 0 {
		var c = &Counter{}
		_, err := c.Count(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "stdin: count: ", err)
			os.Exit(1)
		}
		c.Show(opts, "")
		os.Exit(0)
	}

	for _, filename := range filenames {
		var c = &Counter{}
		fp, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: open: %s\n", filename, err)
			continue
		}
		_, err = c.Count(fp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: count: %s\n", filename, err)
			continue
		}
		totalCount.Add(c)
		c.Show(opts, filename)
		fp.Close()
	}

	if len(filenames) > 1 {
		totalCount.Show(opts, "total")
	}

	os.Exit(0)
}
