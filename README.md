# fastgogi
A Golang gitignore.io fetching library using fasthttp 🏃

[![Build Status](https://travis-ci.org/gofunky/fastgogi.svg?branch=master)](https://travis-ci.org/gofunky/fastgogi)
[![Go Report Card](https://goreportcard.com/badge/github.com/gofunky/fastgogi)](https://goreportcard.com/report/github.com/gofunky/fastgogi)
[![GoDoc](https://godoc.org/github.com/gofunky/fastgogi?status.svg)](https://godoc.org/github.com/gofunky/fastgogi)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/587d4f2b02a54750a73987f58d16ff24)](https://www.codacy.com/app/gofunky/fastgogi?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=gofunky/fastgogi&amp;utm_campaign=Badge_Grade)

### Example

```go
package main

import "github.com/gofunky/fastgogi"
import "fmt"

func main() {
	myGogi := fastgogi.NewClient()
	
	// Get the list of available types.
	availableTypes, err := myGogi.List()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Available Types:\n%s", availableTypes)
	
	// Get the gitignore template for Go and IntelliJ.
	gitignoreContent, err := myGogi.Get("go", "IntelliJ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Available Types:\n%s", gitignoreContent)
}
```

Try it online in the <a href="https://play.golang.org/p/tOPjCHi3eVs">playground</a>.

### Why is there no cmd version?

Check out this generator: <a href="https://github.com/gofunky/gogigen">gogigen</a>
