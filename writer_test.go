package tablewr

import (
	"bytes"
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
