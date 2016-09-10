package main

import (
    "fmt"
    "io/ioutil"
    "bytes"
    "os"
)

type WordCount struct {
    filename string
    lineCount int
}

func (wc *WordCount) CountLines() error {
    b, err := ioutil.ReadFile(wc.filename)
    if err != nil {
        return err
    }
    wc.lineCount = bytes.Count(b, []byte{'\n'})
    return nil
}

func (wc *WordCount) Show() {
    fmt.Printf("%d\t%s\n", wc.lineCount, wc.filename)
}

func main() {
    filename := os.Args[1]
    wc := WordCount{filename, 0}
    err := wc.CountLines()
    if err != nil {
        fmt.Println(err)
    } else {
        wc.Show()
    }
}
