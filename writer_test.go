package tablewr

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTableWriter_Write(t *testing.T) {

	tests := []struct {
		Name     string
		Expected string
		Options  []func(o *tableWriterOpts)
	}{{Name: "default settings",
		Expected: "\n name    age  city   \n---------------------\n foo     23   blr    \n bar     25   france \n foobar  100  ohio   \n\n",
		Options:  nil},
		{
			Name:     "with options",
			Expected: "\n\n  name    |  age  |  city    |\n----------|-------|----------|\n  foo     |  23   |  blr     |\n  bar     |  25   |  france  |\n  foobar  |  100  |  ohio    |\n\n\n",
			Options:  []func(o *tableWriterOpts){WithSep(), WithTablePadding(2, 2), WithColPadding(2, 2)},
		}}

	for _, test := range tests {
		var out bytes.Buffer
		tr := New(&out, 0, test.Options...)
		rows := [][]string{{"name", "age", "city"},
			{"foo", "23", "blr"},
			{"bar", "25", "france"},
			{"foobar", "100", "ohio"}}

		if err := tr.Write(rows); err != nil {
			t.Error(err)
			return
		}

		if test.Expected != out.String() {
			t.Errorf("want:\n%q \n got:\n%q", test.Expected, out.String())
		}
	}
}

func BenchmarkTableWriter_Write(b *testing.B) {
	sizes := []int{
		10,
		100,
		1000,
		10000,
		100000,
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size:%d", size), func(b *testing.B) {
			b.ReportAllocs()
			rows := makeRowsWithHeader(size)
			var buff bytes.Buffer
			tr := New(&buff, 0, WithSep())
			if err := tr.Write(rows); err != nil {
				return
			}
		})
	}
}

func makeRowsWithHeader(size int) [][]string {
	rows := make([][]string, size+1)
	for i := range rows {
		rows[i] = []string{"data", "dataa", "dataaa", "dataaaa", "dataaaaaa"}
	}
	rows[0] = []string{"header", "headere", "headerrr", "headerrrrr", "headerrrrrrr"}
	return rows
}
