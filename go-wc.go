package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "bufio"
)

type WordCount struct {
    filename string
    lineCount int
    bytes int
}

func (wc *WordCount) CountLines() error {
    var lineCount int
    var fp *os.File
    var err error

    fp, err = os.Open(wc.filename)
    if err != nil {
        return err
    }
    defer fp.Close()

    scanner := bufio.NewScanner(fp)
    for scanner.Scan() {
        scanner.Text()
        lineCount++
    }
    if err = scanner.Err(); err != nil {
        return err
    }

    wc.lineCount = lineCount
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

    return nil
}

func (wc *WordCount) Show() {
    fmt.Printf("%d\t%d\t%s\n", wc.lineCount, wc.bytes, wc.filename)
}

func main() {
    for _, filename := range os.Args[1:] {
        wc := WordCount{filename, 0, 0}
        err := wc.CountAll()
        if err != nil {
            fmt.Println(err)
            continue
        }
        wc.Show()
    }
}
