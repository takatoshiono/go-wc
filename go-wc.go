package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
)

type FlagOptions struct {
	printLines bool
	printBytes bool
	printWords bool
}

type WordCount struct {
	filename  string
	lineCount int
	bytes     int
	wordCount int
}

type WordCountList []WordCount

func (wc *WordCount) CountLines() error {
	b, err := ioutil.ReadFile(wc.filename)
	if err != nil {
		return err
	}
	wc.lineCount = bytes.Count(b, []byte{'\n'})
	return nil
}

func (wc *WordCount) CountBytes() error {
	b, err := ioutil.ReadFile(wc.filename)
	if err != nil {
		return err
	}
	wc.bytes = len(b)
	return nil
}

func (wc *WordCount) CountWords() error {
	b, err := ioutil.ReadFile(wc.filename)
	if err != nil {
		return err
	}
	wc.wordCount = len(bytes.Fields(b))
	return nil
}

func (wc *WordCount) CountAll() error {
	var err error

	err = wc.CountLines()
	if err != nil {
		return err
	}

	err = wc.CountBytes()
	if err != nil {
		return err
	}

	err = wc.CountWords()
	if err != nil {
		return err
	}

	return nil
}

func (wc *WordCount) Show(opts FlagOptions) {
	if opts.printLines {
		fmt.Printf(" %7d", wc.lineCount)
	}
	if opts.printWords {
		fmt.Printf(" %7d", wc.wordCount)
	}
	if opts.printBytes {
		fmt.Printf(" %7d", wc.bytes)
	}
	fmt.Printf(" %s\n", wc.filename)
}

func (list WordCountList) Show(opts FlagOptions) {
	var lines, bytes, words int
	for _, r := range list {
		lines += r.lineCount
		bytes += r.bytes
		words += r.wordCount
	}

	if opts.printLines {
		fmt.Printf(" %7d", lines)
	}
	if opts.printWords {
		fmt.Printf(" %7d", words)
	}
	if opts.printBytes {
		fmt.Printf(" %7d", bytes)
	}
	fmt.Println(" total")
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

func main() {
	opts := parseFlagOptions()

	results := make(WordCountList, 0, len(flag.Args()))

	for _, filename := range flag.Args() {
		wc := WordCount{filename, 0, 0, 0}
		err := wc.CountAll()
		if err != nil {
			fmt.Println(err)
			continue
		}
		wc.Show(opts)
		results = append(results, wc)
	}

	if len(results) > 1 {
		results.Show(opts)
	}
}
