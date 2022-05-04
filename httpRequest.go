package main

import (
	"fmt"
  "io"
	"io/ioutil"
	"net/http"
	"os"
  "time"
)

func main() {
  exampleConcurrentRequest()
  // exampleGetRequest()
}

/**
 * go run . http://www.google.it http://www.microsoft.it
 */
func exampleConcurrentRequest() {
  start := time.Now()
  ch := make(chan string)
  for _, url := range os.Args[1:] {
    go fetch(url, ch) // start a goroutine
  }
  for range os.Args[1:] {
    fmt.Println(<-ch) // receive from channel ch
  }
  fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
  start := time.Now()
  resp, err := http.Get(url)
  if err != nil {
    ch <- fmt.Sprint(err) // send to channel ch
    return
  }
  nbytes, err := io.Copy(ioutil.Discard, resp.Body)
  resp.Body.Close() // don't leak resources
  if err != nil {
    ch <- fmt.Sprintf("while reading %s: %v", url, err)
    return
  }
  secs := time.Since(start).Seconds()
  ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}

/**
 * go run . http://www.google.it
 */
func exampleGetRequest() {
	for _, url := range os.Args[1:] {
    resp, err := http.Get(url)
    if err != nil {
      fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
      os.Exit(1)
    }
    b, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
      fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
      os.Exit(1)
    }
    fmt.Printf("%s", b)
  }
}
