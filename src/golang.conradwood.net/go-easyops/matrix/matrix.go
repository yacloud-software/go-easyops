package matrix

import "sync"

type amatrix struct {
	sync.Mutex
	rows      []*matrixrow
	col_names []string
}
type matrixrow struct {
	matrix *amatrix
	name   string
	cells  []*matrixcell
}
type matrixcell struct {
	content interface{}
}

func (m *amatrix) GetColumnNames() []string {
	return m.col_names
}
func (m *amatrix) SetCellByName(rowname, colname string, content interface{}) {
	row := m.GetRowByName(rowname)
	row.SetCellByHeader(colname, content)
}

func (m *amatrix) GetRowByName(name string) *matrixrow {
	m.Lock()
	defer m.Unlock()
	for _, r := range m.rows {
		if r.name == name {
			return r
		}
	}
	row := &matrixrow{matrix: m, name: name}
	m.rows = append(m.rows, row)
	return row
}

func (r *matrixrow) SetCellByHeader(header string, status interface{}) {
	r.matrix.Lock()
	defer r.matrix.Unlock()
	col := -1
	for i, head := range r.matrix.col_names {
		if head == header {
			col = i
			break
		}
	}
	if col == -1 {
		r.matrix.col_names = append(r.matrix.col_names, header)
		col = len(r.matrix.col_names) - 1
	}
	r.SetCell(col, status)
}
func (m *amatrix) Rows() []*matrixrow {
	return m.rows
}
func (r *matrixrow) Name() string {
	return r.name
}
func (r *matrixrow) Cells() []*matrixcell {
	return r.cells
}
func (r *matrixrow) SetCell(col int, c interface{}) {
	for len(r.cells) <= col {
		r.cells = append(r.cells, &matrixcell{})
	}
	r.cells[col].content = c
}
func (c *matrixcell) Content() interface{} {
	return c.content
}
