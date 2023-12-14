package utils

// a "table", like a spreadsheet that may be rendered on screen or so

import (
	"fmt"
	"strings"
)

type Table struct {
	addingRow int // row we're currently "writing" to (0...n)
	rows      []*Row
	headerRow *Row
	hidden    map[int]bool
	maxlenCol map[int]int // for each column(0..n) there _might_ be a maximum length set
}

type Row struct {
	t     *Table
	cells []*Cell
}

type Cell struct {
	typ  int // 0=empty,1=string,2=uint64,3=timestamp,4=float64,5=bool,6=int64
	txt  string
	num  uint64
	ts   uint32
	f    float64
	b    bool
	snum int64
}

func (c *Cell) String() string {
	if c.typ == 0 {
		return ""
	} else if c.typ == 1 {
		return c.txt
	} else if c.typ == 2 {
		return fmt.Sprintf("%d", c.num)
	} else if c.typ == 3 {
		return TimestampString(c.ts)
	} else if c.typ == 7 {
		return TimestampAgeString(c.ts) + " (" + TimestampString(c.ts) + ")"
	} else if c.typ == 4 {
		return fmt.Sprintf("%0.2f", c.f)
	} else if c.typ == 5 {
		return fmt.Sprintf("%v", c.b)
	} else if c.typ == 6 {
		return fmt.Sprintf("%d", c.snum)
	}
	return fmt.Sprintf("type %d", c.typ)
}

// create a new row (writing will commence at a new row
func (t *Table) NewRow() {
	t.addingRow++
}
func (t *Table) getHeaderRow() *Row {
	if t.headerRow == nil {
		t.headerRow = &Row{t: t}
	}
	return t.headerRow
}
func (t *Table) AddHeader(s string) {
	t.getHeaderRow().AddCell(&Cell{typ: 1, txt: s})
}
func (t *Table) AddHeaders(s ...string) {
	for _, a := range s {
		t.getHeaderRow().AddCell(&Cell{typ: 1, txt: a})
	}
}

func (t *Table) GetRowOrCreate(num int) *Row {
	for len(t.rows) <= num {
		t.rows = append(t.rows, &Row{t: t})
	}
	return t.rows[num]
}
func (t *Table) AddBool(b bool) *Table {
	r := t.GetRowOrCreate(t.addingRow)
	r.AddCell(&Cell{typ: 5, b: b})
	return t
}
func (t *Table) AddString(s string) *Table {
	r := t.GetRowOrCreate(t.addingRow)
	r.AddCell(&Cell{typ: 1, txt: s})
	return t
}
func (t *Table) AddStrings(sts ...string) *Table {
	for _, s := range sts {
		r := t.GetRowOrCreate(t.addingRow)
		r.AddCell(&Cell{typ: 1, txt: s})
	}
	return t
}
func (t *Table) AddTimestamp(ts uint32) *Table {
	r := t.GetRowOrCreate(t.addingRow)
	r.AddCell(&Cell{typ: 3, ts: ts})
	return t
}
func (t *Table) AddTimestampWithAge(ts uint32) *Table {
	r := t.GetRowOrCreate(t.addingRow)
	r.AddCell(&Cell{typ: 7, ts: ts})
	return t
}
func (t *Table) AddFloat64(f float64) *Table {
	r := t.GetRowOrCreate(t.addingRow)
	r.AddCell(&Cell{typ: 4, f: f})
	return t
}
func (t *Table) AddUint32(i uint32) *Table {
	t.AddUint64(uint64(i))
	return t
}
func (t *Table) AddInt(i int) *Table {
	t.AddUint64(uint64(i))
	return t
}
func (t *Table) AddInt64(i int64) *Table {
	r := t.GetRowOrCreate(t.addingRow)
	r.AddCell(&Cell{typ: 6, snum: i})
	return t
}
func (t *Table) AddUint64(i uint64) *Table {
	r := t.GetRowOrCreate(t.addingRow)
	r.AddCell(&Cell{typ: 2, num: i})
	return t

}
func (r *Row) AddCell(cell *Cell) {
	r.cells = append(r.cells, cell)
}

// return # of cells (considering the col<->idx mapping)
func (r *Row) Cols() int {
	return len(r.Cells())
}

// return all cells (considering the col<->idx mapping)
func (r *Row) Cells() []*Cell {
	if r.t.hidden == nil {
		r.t.hidden = make(map[int]bool)
	}
	var res []*Cell
	for i := 0; i < len(r.cells); i++ {
		if r.t.hidden[i] {
			continue
		}
		res = append(res, r.cells[i])
	}
	return res
}

// return a cell (considering the col<->idx mapping)
func (r *Row) GetCell(idx int) *Cell {
	col := r.t.idx2col(idx)
	//	fmt.Printf("Want cell %d, using %d\n", idx, col)
	if len(r.cells) <= col {
		return nil
	}
	return r.cells[col]
}

func (t *Table) ToCSV() string {
	rows := len(t.rows)
	sb := strings.Builder{}
	for i := 0; i < rows; i++ {
		row := t.GetRowOrCreate(i)
		if row.Cols() == 0 {
			continue
		}
		line := ""
		deli := ""
		for cn := 0; cn < row.Cols(); cn++ {
			cel := row.GetCell(cn)
			s := escapeCell(cel.String())
			line = line + deli + s
			deli = ","
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}
func escapeCell(s string) string {
	s = strings.ReplaceAll(s, ",", "\\,")
	return s
}

// column 0..n
func (t *Table) DisableColumn(col int) {
	if t.hidden == nil {
		t.hidden = make(map[int]bool)
	}
	t.hidden[col] = true
}

// column 0..n
func (t *Table) EnableColumn(col int) {
	if t.hidden == nil {
		return
	}
	t.hidden[col] = false
}

// column 0..n
func (t *Table) EnableAllColumns() {
	t.hidden = nil
}

// calculates the column offset (considering hidden columns)
func (t *Table) idx2col(idx int) int {
	if t.hidden == nil {
		return idx
	}
	off := 0
	for i := 0; i < idx; i++ {
		if t.hidden[idx] {
			off++
		}
	}
	return idx + off
}
func (t *Table) GetMaxLen(col int) int {
	if t.maxlenCol == nil {
		t.maxlenCol = make(map[int]int)
	}
	f, found := t.maxlenCol[col]
	if !found {
		return 0xFFFFFFFF
	}
	return f

}
func (t *Table) SetMaxLen(col, width int) {
	if t.maxlenCol == nil {
		t.maxlenCol = make(map[int]int)
	}
	t.maxlenCol[col] = width
}

// gets "printing" rows. Multi-line text or text that is wrapped will create an extra row
func (t *Table) GetPrintingRows() []*Row {
	var res []*Row
	for _, r := range t.rows {
		res = append(res, r.toPrintingRows()...)
	}
	return res
}

// wrap text and return multiple rows
func (r *Row) toPrintingRows() []*Row {
	var res []*Row
	for cn, c := range r.Cells() {
		if c.typ != 1 {
			// it's not a string, insert cell as is to row
			if len(res) == 0 {
				res = append(res, &Row{t: r.t})
			}
			r := res[0]
			for len(r.cells) <= cn {
				r.cells = append(r.cells, &Cell{})
			}
			r.cells[cn] = c
			continue
		}

		// it's a string.. wrap text
		ml := r.t.GetMaxLen(cn)
		s := c.String()
		lines := strings.Split(s, "\n")
		var tlines []string
		for _, l := range lines {
			for {
				if len(l) <= ml {
					tlines = append(tlines, l)
					break
				}
				splitAt := findBetterSplitAt(l, ml)
				nl := l[:splitAt]
				tlines = append(tlines, nl)
				l = l[splitAt:]
			}
		}
		for len(res) < len(tlines) {
			res = append(res, &Row{t: r.t})
		}
		for i, l := range tlines {
			r := res[i]
			for len(r.cells) <= cn {
				r.cells = append(r.cells, &Cell{})
			}
			r.cells[cn] = &Cell{typ: c.typ, txt: l}

		}
	}
	return res
}
func findBetterSplitAt(line string, proposed int) int {
	splits := []byte{' ', ':'}
	i := proposed
	j := 0
	for {
		j++
		if j > 20 {
			return proposed
		}
		i--
		if i < 0 {
			return proposed
		}
		for _, sp := range splits {
			if line[i] == sp {
				if i < proposed {
					return i + 1
				}
				return i
			}
		}

	}

}
