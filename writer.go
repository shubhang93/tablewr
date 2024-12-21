package tablewr

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type TableWriter struct {
	wr             *tabwriter.Writer
	sb             *strings.Builder
	headerRowIndex int
	tableWriterOpts
}

type tableWriterOpts struct {
	TableTopPadding    int
	TableBottomPadding int
	ColLeftPadding     int
	ColRightPadding    int
	OptsFlags          uint
}

var defaultOptions = tableWriterOpts{
	TableTopPadding:    1,
	TableBottomPadding: 1,
	ColLeftPadding:     1,
	ColRightPadding:    1,
}

func WithSep() func(opts *tableWriterOpts) {
	return func(opts *tableWriterOpts) {
		opts.OptsFlags |= tabwriter.Debug
	}
}

func WithColPadding(left, right int) func(opts *tableWriterOpts) {
	return func(opts *tableWriterOpts) {
		opts.ColLeftPadding = left
		opts.ColRightPadding = right
	}
}

func WithTablePadding(top, bottom int) func(opts *tableWriterOpts) {
	return func(opts *tableWriterOpts) {
		opts.TableTopPadding = top
		opts.TableBottomPadding = bottom
	}
}

func New(writer io.Writer, headerRowIndex int, opts ...func(opts *tableWriterOpts)) *TableWriter {
	tropts := defaultOptions
	for _, opt := range opts {
		opt(&tropts)
	}
	tr := tabwriter.NewWriter(writer, 0, 0, 0, '\t', tropts.OptsFlags)
	return &TableWriter{
		wr:              tr,
		sb:              &strings.Builder{},
		headerRowIndex:  headerRowIndex,
		tableWriterOpts: tropts,
	}
}

func (twr *TableWriter) Write(rows [][]string) error {
	if len(rows) < 1 {
		return nil
	}
	cols := rows[0]
	if len(cols) == 0 {
		return nil
	}

	maxWidths := make([]int, len(cols))
	for _, row := range rows {
		for i, col := range row {
			maxWidths[i] = max(maxWidths[i], len(col))
		}
	}

	defer twr.wr.Flush()

	if _, err := io.WriteString(twr.wr, strings.Repeat("\n", twr.TableTopPadding)); err != nil {
		return err
	}

	for i, row := range rows {
		for j, col := range row {
			mw := maxWidths[j]
			numLeftOver := mw - len(col)
			rightPadding := strings.Repeat(" ", numLeftOver+twr.ColRightPadding)
			leftPadding := strings.Repeat(" ", twr.ColLeftPadding)

			twr.sb.WriteString(leftPadding)
			twr.sb.WriteString(col)
			twr.sb.WriteString(rightPadding)
			twr.sb.WriteByte('\t')

			_, err := fmt.Fprint(twr.wr, twr.sb.String())
			if err != nil {
				return err
			}
			twr.sb.Reset()
		}

		_, err := fmt.Fprint(twr.wr, "\n")
		if err != nil {
			return err
		}

		if i == twr.headerRowIndex {
			if err := twr.writeRowDelim(maxWidths, twr.ColRightPadding+twr.ColRightPadding); err != nil {
				return err
			}
		}
	}
	if _, err := io.WriteString(twr.wr, strings.Repeat("\n", twr.TableBottomPadding)); err != nil {
		return err
	}
	return nil
}

func (twr *TableWriter) writeRowDelim(colWidths []int, padding int) error {
	defer twr.sb.Reset()
	for _, width := range colWidths {
		twr.sb.WriteString(strings.Repeat("-", width))
		twr.sb.WriteString(strings.Repeat("-", padding))
		twr.sb.WriteByte('\t')
	}
	twr.sb.WriteByte('\n')
	if _, err := fmt.Fprint(twr.wr, twr.sb.String()); err != nil {
		return err
	}
	return nil
}
