package utils

import (
	"strings"
)

type TextFormatter struct {
	table    *Table
	colSizes []int
}

func (t *Table) ToPrettyString() string {
	tf := &TextFormatter{table: t}
	return tf.ToPrettyString()
}
func (tf *TextFormatter) ToPrettyString() string {
	r := tf.table.rows
	if tf.table.headerRow != nil {
		r = append(r, tf.table.headerRow)
	}
	tf.colSizes = getCharWidth(r)

	var sb strings.Builder
	if tf.table.headerRow != nil {
		sb.WriteString(tf.renderRow(tf.table.headerRow, true))
		sb.WriteString(tf.seperatorRow())
	}
	for _, row := range tf.table.rows {
		sb.WriteString(tf.renderRow(row, false))
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
func (tf *TextFormatter) renderRow(row *Row, center bool) string {
	res := "| "
	l := false
	for i, c := range row.Cells() {
		s := c.String()
		for len(s) < tf.colSizes[i] {
			if center {
				if l {
					l = false
					s = " " + s
				} else {
					l = true
					s = s + " "
				}
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
			w := len(c.String())
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
