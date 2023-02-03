package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type TextFormatter struct {
	table    *Table
	colSizes []int
}

func (t *TextFormatter) tableWidth() int {
	res := 0
	for _, cs := range t.colSizes {
		res = res + cs
	}
	return res
}
func (t *Table) ToPrettyString() string {
	tf := &TextFormatter{table: t}
	return tf.ToPrettyString()
}
func (tf *TextFormatter) adjustColsToMatchTerminal() {
	di, err := TerminalSize()
	if err != nil {
		fmt.Printf("[go-easyops] unable to determine terminal size (%s)\n", err)
		return
	}
	if di.Columns() < len(tf.colSizes) {
		fmt.Printf("[go-easyops] terminal has %d columns, but table has %d columns\n", di.Columns(), len(tf.colSizes))
		return
	}

	extraSize := len(tf.colSizes)*3 + 2 // borders and seperators
	if di.Columns() < extraSize {
		return
	}
	// make sure table fits into terminal
	shrank := false
	for tf.tableWidth() > di.Columns()-extraSize {

		// find column to shrink
		col_shrink := -1
		for col, size := range tf.colSizes {
			if col_shrink == -1 || tf.colSizes[col_shrink] < size {
				col_shrink = col
			}
		}
		if col_shrink == -1 {
			// no columns to shrink at all??
			return
		}
		tf.colSizes[col_shrink]--
		shrank = true
	}
	for col, size := range tf.colSizes {
		tf.table.SetMaxLen(col, size)
	}
	if shrank {
		/*
			fmt.Printf("Shrank to %d\n", di.Columns())
			for i, size := range tf.colSizes {
				fmt.Printf("Colsizes %d. %d\n", i, size)
			}
		*/
	}

}
func (tf *TextFormatter) ToPrettyString() string {
	r := tf.table.GetPrintingRows()
	if tf.table.headerRow != nil {
		r = append(r, tf.table.headerRow)
	}
	tf.colSizes = getCharWidth(r)
	tf.adjustColsToMatchTerminal()
	var sb strings.Builder
	if tf.table.headerRow != nil {
		sb.WriteString(tf.renderRow(tf.table.headerRow, 2))
		sb.WriteString(tf.seperatorRow())
	}
	for _, row := range tf.table.GetPrintingRows() {
		sb.WriteString(tf.renderRow(row, 0))
	}
	sb.WriteString(tf.seperatorRow())
	return sb.String()
}
func (tf *TextFormatter) seperatorRow() string {
	s := ""
	for _, cs := range tf.colSizes {
		xs := "|"
		for len(xs) < (cs + 3) {
			xs = xs + "="
		}
		s = s + xs
	}
	return s + "|\n"

}

// pos: 0 left, 1 right, 2 center
func (tf *TextFormatter) renderRow(row *Row, pos int) string {
	res := "| "
	l := false
	for i, c := range row.Cells() {
		s := c.String()
		for len(s) < tf.colSizes[i] {
			if pos == 2 {
				if l {
					l = false
					s = " " + s
				} else {
					l = true
					s = s + " "
				}
			} else if pos == 0 {
				s = s + " "
			} else {
				s = " " + s
			}
		}
		res = res + s + " | "
	}
	res = res + "\n"
	return res

}

// calculate max width of columns
func getCharWidth(rows []*Row) []int {
	maxCols := 0
	for _, r := range rows {
		if r.Cols() > maxCols {
			maxCols = r.Cols()
		}
	}
	sz := make(map[int]int) // colidx -> size
	for _, r := range rows {
		for i := 0; i < maxCols; i++ {
			c := r.GetCell(i)
			if c == nil {
				continue
			}
			w := utf8.RuneCountInString(c.String())
			if w > sz[i] {
				sz[i] = w
			}
		}
	}

	// turn map into array
	var res []int
	for i := 0; i < len(sz); i++ {
		res = append(res, sz[i])
	}
	return res
}
