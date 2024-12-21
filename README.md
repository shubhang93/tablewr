# tablewr

## A very minimal table writer in Golang

- Uses the stdlib's `text/tabwriter` to render a table
- Auto-scales column widths based on the column's maximum length
- Tiny API

### Usage

```go
package main

import (
	"github.com/shubhang93/tablewr"
	"os"
)

func main() {
	wr := tablewr.New(os.Stderr, 0, tablewr.WithSep())
	data := [][]string{
		{"title", "price", "sold out", "rating"},
		{"The Shining", "30$", "yes", "*****"},
		{"The Mask", "10$", "no", "***"},
		{"Godfather", "40$", "no", "*****"},
		{"Godfather-2", "40$", "no", "*****"},
		{"Shawshank Redemption", "30$", "yes", "*****"},
	}
	if err := wr.Write(data); err != nil {
		panic("write err:" + err.Error())
	}
}

```

### Output

```text

 title                | price | sold out | rating |
----------------------|-------|----------|--------|
 The Shining          | 30$   | yes      | *****  |
 The Mask             | 10$   | no       | ***    |
 Godfather            | 40$   | no       | *****  |
 Godfather-2          | 40$   | no       | *****  |
 Shawshank Redemption | 30$   | yes      | *****  |


```

