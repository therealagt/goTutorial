package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
    start := time.Now()
    ch := make(chan string)

    file, err := os.Open("argumentList.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error opening argumentList.txt: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    var urls []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        url := scanner.Text()
        if url != "" {
            urls = append(urls, url)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "Error reading argumentList.txt: %v\n", err)
        os.Exit(1)
    }

	outFile, err := os.Create("outputList.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
        os.Exit(1)
    }
    defer outFile.Close()

    for _, url := range urls {
        go fetch(url, ch)
    }
    for range urls {
        fmt.Fprintln(outFile, <-ch) 
    }
    fmt.Fprintf(outFile, "%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return 
	}
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return 
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}