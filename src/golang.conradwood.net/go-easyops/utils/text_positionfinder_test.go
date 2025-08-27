package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestPosFinderLine2(t *testing.T) {
	testval := `
line1
line2
line3
line4
`
	new_testval := `
line1
line2
line2b
line3
line4
`

	ct := []byte(testval)
	nct := []byte(new_testval)
	pf := NewTextPositionFinder(ct).(*textPositionFinder)
	x := pf.Content()
	if !bytes.Equal(x, ct) {
		t.Fatalf("content mismatch before any modification")
	}
	if !pf.FindLineContaining("line2") {
		t.Fatalf("failed to find pattern")
	}
	pf.AddLine("line2b")
	expect_content(t, nct, pf)
}

func TestPosFinderLine0(t *testing.T) {
	testval := `
line1
line2
line3
line4
`
	new_testval := `
line0
line1
line2
line3
line4
`

	ct := []byte(testval)
	nct := []byte(new_testval)
	pf := NewTextPositionFinder(ct).(*textPositionFinder)
	x := pf.Content()
	if !bytes.Equal(x, ct) {
		t.Fatalf("content mismatch before any modification")
	}
	pf.AddLine("line0")
	expect_content(t, nct, pf)
}

func TestPosFinderLine5(t *testing.T) {
	testval := `
line1
line2
line3
line4
`
	new_testval := `
line1
line2
line3
line4
line4b
`

	ct := []byte(testval)
	nct := []byte(new_testval)
	pf := NewTextPositionFinder(ct).(*textPositionFinder)
	x := pf.Content()
	if !bytes.Equal(x, ct) {
		t.Fatalf("content mismatch before any modification")
	}
	if !pf.FindLineContaining("line4") {
		t.Fatalf("failed to find pattern")
	}
	pf.AddLine("line4b")
	expect_content(t, nct, pf)
}

func TestPosFinderMultiLine(t *testing.T) {
	testval := `
line1
line2
line3
line4
`
	new_testval := `
line1
line2
line2b
line3
line4
`

	ct := []byte(testval)
	nct := []byte(new_testval)
	pf := NewTextPositionFinder(ct).(*textPositionFinder)
	x := pf.Content()
	if !bytes.Equal(x, ct) {
		t.Fatalf("content mismatch before any modification")
	}
	pf.FindLineContaining("line")
	pf.FindLineContaining("line")
	pf.AddLine("line2b")
	expect_content(t, nct, pf)
}

func expect_content(t *testing.T, ct []byte, pf *textPositionFinder) {
	if bytes.Equal(ct, pf.Content()) {
		return
	}
	lines_expected := strings.Split(string(ct), "\n")
	lines_got := strings.Split(string(pf.Content()), "\n")
	max_lines := len(lines_expected)
	if len(lines_got) > max_lines {
		max_lines = len(lines_got)
	}
	ta := &Table{}
	ta.AddHeaders("got", "expected")
	for i := 0; i < max_lines; i++ {
		l1 := ""
		l2 := ""
		if i < len(lines_got) {
			l1 = lines_got[i]
		}
		if i < len(lines_expected) {
			l2 = lines_expected[i]
		}
		ta.AddString(l1)
		ta.AddString(l2)
		ta.NewRow()
	}
	t.Errorf("modification produced incorrect result")
	t.Logf("\n" + ta.ToPrettyString())
	//	t.Logf("New Content:\n%s\n", string(pf.Content()))
}
