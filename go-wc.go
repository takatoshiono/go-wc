package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

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

func (list WordCountList) Show() {
	var lines, bytes, words int
	for _, r := range list {
		lines += r.lineCount
		bytes += r.bytes
		words += r.wordCount
	}
	fmt.Printf("%d\t%d\t%d\t%s\n", lines, words, bytes, "total")
}

func (wc *WordCount) Show() {
	fmt.Printf("%d\t%d\t%d\t%s\n", wc.lineCount, wc.wordCount, wc.bytes, wc.filename)
}

func main() {
	results := make(WordCountList, 0, len(os.Args[1:]))

	for _, filename := range os.Args[1:] {
		wc := WordCount{filename, 0, 0, 0}
		err := wc.CountAll()
		if err != nil {
			fmt.Println(err)
			continue
		}
		wc.Show()
		results = append(results, wc)
	}

	results.Show()
}
