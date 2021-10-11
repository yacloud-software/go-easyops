package utils

// a "table", like a spreadsheet that may be rendered on screen or so

import (
	"fmt"
)

type Table struct {
	addingRow int // row we're currently "writing" to (0...n)
	rows      []*Row
	headerRow *Row
}

type Row struct {
	cells []*Cell
}

type Cell struct {
	typ int // 0=empty,1=string,2=uint64,3=timestamp,4=float64,5=bool
	txt string
	num uint64
	ts  uint32
	f   float64
	b   bool
}

func (c *Cell) String() string {
	if c.typ == 0 {
		return ""
	} else if c.typ == 1 {
		return c.txt
	} else if c.typ == 2 {
		return fmt.Sprintf("%d", c.num)
	} else if c.typ == 3 {
		return TimestampAgeString(c.ts) + " (" + TimestampString(c.ts) + ")"
	} else if c.typ == 4 {
		return fmt.Sprintf("%0.2f", c.f)
	} else if c.typ == 5 {
		return fmt.Sprintf("%v", c.b)
	}
	return fmt.Sprintf("type %d", c.typ)
}

// create a new row (writing will commence at a new row
func (t *Table) NewRow() {
	t.addingRow++
}
func (t *Table) getHeaderRow() *Row {
	if t.headerRow == nil {
		t.headerRow = &Row{}
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
		t.rows = append(t.rows, &Row{})
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
func (t *Table) AddUint64(i uint64) *Table {
	r := t.GetRowOrCreate(t.addingRow)
	r.AddCell(&Cell{typ: 2, num: i})
	return t

}
func (r *Row) AddCell(cell *Cell) {
	r.cells = append(r.cells, cell)
}
func (r *Row) Cols() int {
	return len(r.cells)
}
func (r *Row) GetCell(idx int) *Cell {
	if len(r.cells) <= idx {
		return nil
	}
	return r.cells[idx]
}
