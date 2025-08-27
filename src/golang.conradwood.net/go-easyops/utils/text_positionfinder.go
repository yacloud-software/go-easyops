package utils

import (
	"slices"
	"strings"
	"sync"
)

type TextPositionFinder interface {
	// advance position to line after pattern found. return true if found
	FindLineContaining(s string) bool
	// insert line at position
	AddLine(s string)
	// return current content
	Content() []byte
}
type textPositionFinder struct {
	sync.Mutex
	lines    []string
	line_pos int // only advances forward, points to the next line
}

/*
a position finder helps finding certain positions in pieces of text. for example: find the '}' after the first occurence of line "foo" and "bar". Each call advances the position further. AddLine inserts a new line at current position.
*/
func NewTextPositionFinder(ct []byte) TextPositionFinder {
	lines := strings.Split(string(ct), "\n")
	res := &textPositionFinder{lines: lines}
	return res
}
func (pf *textPositionFinder) remaining_lines() []string {
	if pf.line_pos >= len(pf.lines) {
		return nil
	}
	return pf.lines[pf.line_pos+1:]
}

// advance to position after a line containing...
func (pf *textPositionFinder) FindLineContaining(s string) bool {
	pf.Lock()
	defer pf.Unlock()
	for i, l := range pf.remaining_lines() {
		//	fmt.Printf("Line %d: \"%s\"\n", i, l)
		if strings.Contains(l, s) {
			pf.line_pos = pf.line_pos + i + 1
			return true
		}
	}
	return false
}
func (pf *textPositionFinder) AddLine(line string) {
	pf.Lock()
	defer pf.Unlock()
	pos := pf.line_pos + 1
	pf.lines = slices.Insert(pf.lines, pos, line)

}
func (pf *textPositionFinder) Content() []byte {
	s := strings.Join(pf.lines, "\n")
	return []byte(s)
}
