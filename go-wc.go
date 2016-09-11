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
    bytes int
}

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

func (wc *WordCount) Show() {
    fmt.Printf("%d\t%d\t%s\n", wc.lineCount, wc.bytes, wc.filename)
}

func main() {
    for _, filename := range os.Args[1:] {
        wc := WordCount{filename, 0, 0}

        err := wc.CountLines()
        if err != nil {
            fmt.Println(err)
            continue
        }

        err = wc.CountBytes()
        if err != nil {
            fmt.Println(err)
            continue
        }

        wc.Show()
    }
}
