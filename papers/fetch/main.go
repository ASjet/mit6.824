package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

const PageURL = "https://pdos.csail.mit.edu/6.824/schedule.html"
const BaseURL = "https://pdos.csail.mit.edu/6.824/"

type File struct {
	Name string
	Url  string
}

func main() {
	resp, err := http.Get(PageURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch %v\n", err)
		os.Exit(-1)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing %s as HTML: %v", PageURL, err)
		os.Exit(-1)
	}
	allLinks := walkLinks(nil, doc)
	pdfs := []File{}
	for _, l := range allLinks {
		if strings.HasSuffix(l, ".pdf") {
			name := strings.TrimPrefix(l, "papers/")
			name = strings.TrimPrefix(name, "notes/")
			pdfs = append(pdfs, File{name, BaseURL + l})
		}
	}
	fmt.Println(pdfs)
	wg := sync.WaitGroup{}
	for _, f := range pdfs {
		wg.Add(1)
		go func(f File) {
			defer wg.Done()
			fmt.Printf("fetching %s\n", f.Url)
			resp, err := http.Get(f.Url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch %v\n", err)
			}
			file, err := os.Create(f.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			io.Copy(file, resp.Body)
			resp.Body.Close()
			file.Close()
		}(f)
	}
	wg.Wait()
}

func walkLinks(links []string, node *html.Node) []string {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = walkLinks(links, c)
	}
	return links
}
