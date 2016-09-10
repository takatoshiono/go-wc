package main

import (
    "fmt"
    "io/ioutil"
    "bytes"
)

type WordCount struct {
    filename string
    lineCount int
}

func (wc *WordCount) CountLines() {
    b, err := ioutil.ReadFile(wc.filename)
    if err != nil {
        fmt.Printf("Cannot read %s.\n", wc.filename)
    }
    wc.lineCount = bytes.Count(b, []byte{'\n'})
}

func (wc *WordCount) Show() {
    fmt.Printf("%d\t%s\n", wc.lineCount, wc.filename)
}

func main() {
    wc := WordCount{"./go-wc.go", 0}
    wc.CountLines()
    wc.Show()
}
