package main

import (
	"unicode"
	"unicode/utf8"
	"strings"
	"github.com/mattn/go-runewidth"
)

var (
	taboffset = 4
	pageoffset = 16
)

//type place int
//
//var (
//	NONE place = iota
//	BOC
//	EOC
//	EOL
//)

type Cursor struct {
	l int // line offset
	o int // cursor offset - When MoveUp or MoveDown, it will calculated from visual offset.
	v int // visual offset - When MoveLeft of MoveRight, it will matched to cursor offset.
	b int // byte offset
	t *Text
	// stick place - will implement later
}

func NewCursor(t *Text) *Cursor {
	return &Cursor{0, 0, 0, 0, t}
}

func (c *Cursor) Copy(c2 Cursor) {
	c.l = c2.l
	c.o = c2.o
	c.v = c2.v
	c.b = c2.b
}

func (c *Cursor) SetOffsets(b int) {
	c.b = b
	c.v = c.VFromB(b)
	c.o = c.v
}

// Before shifting, visual offset will matched to cursor offset.
func (c *Cursor) ShiftOffsets(b, v int) {
	c.v = c.o
	c.b += b
	c.v += v
	c.o += v
}

// After MoveUp or MoveDown, we need reclaculate cursor offsets (except visual offset).
func (c *Cursor) RecalculateOffsets() {
	c.o = c.OFromV(c.v)
	c.b = c.BFromC(c.o)
}

func (c *Cursor) OFromV(v int) (o int) {
	// Cursor offset cannot go further than line's maximum visual length.
	maxv := c.LineVisualLength()
	if v >  maxv {
		return maxv
	}
	// It's not allowed the cursor is in the middle of multi-length(visual) character.
	// So we need recaculate the cursors offset.
	remain := c.LineData()
	lasto := 0
	for {
		r, rlen := utf8.DecodeRuneInString(remain)
		remain = remain[rlen:]
		lasto = o
		o += RuneVisualLength(r)
		if o==v {
			return o
		} else if o > v {
			return lasto
		}
	}
}


func (c *Cursor) BFromC(o int) (b int) {
	remain := c.LineData()
	for o>0 {
		r, rlen := utf8.DecodeRuneInString(remain)
		remain = remain[rlen:]
		b+= rlen
		o-= RuneVisualLength(r)
	}
	return
}

func BFromC(line string, o int) (b int) {
	remain := line
	for o>0 {
		r, rlen := utf8.DecodeRuneInString(remain)
		remain = remain[rlen:]
		b+= rlen
		o-= RuneVisualLength(r)
	}
	return
}

func (c *Cursor) VFromB(b int) (v int){
	remain := c.LineData()[:b]
	for len(remain) > 0 {
		r, rlen := utf8.DecodeRuneInString(remain)
		remain = remain[rlen:]
		v += RuneVisualLength(r)
	}
	return
}

func (c *Cursor) Position() Point {
	return Point{c.l, c.o}
}

// TODO : relativePosition(p Point) Point ?
func (c *Cursor) PositionInWindow(w *Window) Point {
	return c.Position().Sub(w.min)
}

func (c *Cursor) LineData() string {
	return c.t.lines[c.l].data
}

func (c *Cursor) LineDataUntilCursor() string {
	return c.LineData()[:c.b]
}

func (c *Cursor) LineDataFromCursor() string {
	return c.LineData()[c.b:]
}

func (c *Cursor) ExceededLineLimit() bool {
	return c.b > len(c.LineData())
}

func (c *Cursor) RuneAfter() (rune, int) {
	return utf8.DecodeRuneInString(c.LineData()[c.b:])
}

func (c *Cursor) RuneBefore() (rune, int) {
	return utf8.DecodeLastRuneInString(c.LineData()[:c.b])
}

// should refine after
// may be use dictionary??
func RuneVisualLength(r rune) int {
	if r=='\t' {
		return taboffset
	}
	return runewidth.RuneWidth(r)
}

func (c *Cursor) LineByteLength() int {
	return len(c.LineData())
}

func (c *Cursor) LineVisualLength() int {
	return c.VFromB(c.LineByteLength())
}

func (c *Cursor) AtBol() bool{
	return c.b == 0
}

func (c *Cursor) AtEol() bool{
	return c.b == c.LineByteLength()
}

func (c *Cursor) OnFirstLine() bool{
	return c.l == 0
}

func (c *Cursor) OnLastLine() bool {
	return c.l == len(c.t.lines)-1
}

func (c *Cursor) AtBow() bool {
	r, _ := c.RuneAfter()
	rb, _ := c.RuneBefore() // is it not panic on bof?
	if (unicode.IsLetter(r) || unicode.IsDigit(r))  && !(unicode.IsLetter(rb) || unicode.IsDigit(rb)) {
		return true
	}
	return false
}

func (c *Cursor) AtEow() bool {
	r, _ := c.RuneAfter()
	rb, _ := c.RuneBefore() // is it not panic on bof?
	if !(unicode.IsLetter(r) || unicode.IsDigit(r)) && (unicode.IsLetter(rb) || unicode.IsDigit(rb)) {
		return true
	}
	return false
}

func (c *Cursor) AtBof() bool {
	return c.OnFirstLine() && c.AtBol()
}

func (c *Cursor) AtEof() bool {
	return c.OnLastLine() && c.AtEol()
}

func (c *Cursor) InStrings() bool {
	instr := false
	var starter rune
	var old rune
	var oldold rune
	for _, r := range c.LineDataUntilCursor() {
		if !instr && strings.ContainsAny(string(r), "'\"") && !(old == '\\' && oldold != '\\') {
			instr = true
			starter = r
		} else if instr && (r == starter) && !(old == '\\' && oldold != '\\') {
			instr = false
			starter = ' '
		}
		oldold = old
		old = r
	}
	return instr
}

func (c *Cursor) MoveLeft() {
	if c.AtBof() {
		return
	} else if c.AtBol() {
		c.l--
		c.SetOffsets(c.LineByteLength())
		return
	}
	r, rlen := c.RuneBefore()
	vlen := RuneVisualLength(r)
	c.ShiftOffsets(-rlen, -vlen)
}

func (c *Cursor) MoveRight() {
	if c.AtEof() {
		return
	} else if c.AtEol() || c.ExceededLineLimit(){
		c.l++
		c.SetOffsets(0)
		return
	}
	r, rlen := c.RuneAfter()
	vlen := RuneVisualLength(r)
	c.ShiftOffsets(rlen, vlen)
}

func (c *Cursor) MoveUp() {
	if c.OnFirstLine() {
		return
	}
	c.l--
	c.RecalculateOffsets()
}

func (c *Cursor) MoveDown() {
	if c.OnLastLine() {
		return
	}
	c.l++
	c.RecalculateOffsets()
}

func (c *Cursor) MovePrevBowEow() {
	if c.AtBof() {
		return
	}
	if c.AtBow() {
		for {
			c.MoveLeft()
			if c.AtEow() || c.AtBof() {
				return
			}
		}
	} else if c.AtEow() || c.AtBof() {
		for {
			c.MoveLeft()
			if c.AtBow() {
				return
			}
		}
	} else {
		r, _ := c.RuneAfter()
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			// we are in the middle of non-words. find eow.
			for {
				c.MoveLeft()
				if c.AtEow() || c.AtBof() {
					return
				}
			}
		} else {
			// we are in the middle of a word. find bow.
			for {
				c.MoveLeft()
				if c.AtBow() || c.AtBof() {
					return
				}
			}
		}
	}
}

func (c *Cursor) MoveNextBowEow() {
	if c.AtEof() {
		return
	}
	if c.AtBow() {
		for {
			c.MoveRight()
			if c.AtEow() || c.AtEof() {
				return
			}
		}
	} else if c.AtEow() {
		for {
			c.MoveRight()
			if c.AtBow() || c.AtEof() {
				return
			}
		}
	} else {
		r, _ := c.RuneAfter()
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			// we are in the middle of non-words. find bow.
			for {
				c.MoveRight()
				if c.AtBow() || c.AtEof() {
					return
				}
			}
		} else {
			// we are in the middle of a word. find eow.
			for {
				c.MoveRight()
				if c.AtEow() || c.AtEof() {
					return
				}
			}
		}
	}
}

func (c *Cursor) MoveBol() {
	c.SetOffsets(0)
}

func (c *Cursor) MoveBocBolAdvance() {
	// if already bol, move cursor to prev line
	if c.AtBol() && !c.OnFirstLine() {
		c.MoveUp()
		return
	}

	remain := c.LineData()
	b := 0 // where line contents start
	for len(remain)>0 {
		r, rlen := utf8.DecodeRuneInString(remain)
		remain = remain[rlen:]
		if !unicode.IsSpace(r) {
			break
		}
		b += rlen
	}
	if c.b > b {
		c.SetOffsets(b)
		return
	}
	c.SetOffsets(0)
}

func (c *Cursor) MoveEol() {
	c.SetOffsets(len(c.LineData()))
}

func (c *Cursor) MoveEolAdvance() {
	// if already eol, move to next line
	if c.AtEol() && !c.OnLastLine() {
		c.MoveDown()
	}

	c.SetOffsets(c.LineByteLength())
}

func (c *Cursor) PageUp() {
	for i:=0; i < pageoffset; i++ {
		if c.OnFirstLine() {
			break
		}
		c.MoveUp()
	}
}

func (c *Cursor) PageDown() {
	for i:=0; i < pageoffset; i++ {
		if c.OnLastLine() {
			break
		}
		c.MoveDown()
	}
}

func (c *Cursor) MoveBof() {
	for {
		if c.OnFirstLine() {
			break
		}
		c.MoveUp()
	}
	c.MoveBol()
}

func (c *Cursor) MoveEof() {
	for {
		if c.OnLastLine() {
			break
		}
		c.MoveDown()
	}
	c.MoveEol()
}

func (c *Cursor) SplitLine() {
	c.t.SplitLine(c.l, c.b)
	c.MoveDown()
	c.SetOffsets(0)
}

func (c *Cursor) Insert(str string) {
	for _, r := range str {
		if r == '\n' {
			c.SplitLine()
			continue
		}
		c.t.Insert(string(r), c.l, c.b)
		c.MoveRight()
	}
}

func (c *Cursor) Tab(sel *Selection) []int {
	tabed := make([]int, 0)
	if sel == nil {
		c.t.lines[c.l].InsertTab()
		tabed = append(tabed, c.l)
		c.SetOffsets(c.b+1)
		return tabed
	}
	min, max := sel.MinMax()
	for l := min.l; l < max.l + 1; l++ {
		if l == min.l && min.b == len(c.t.lines[min.l].data) {
			continue
		} else if l == max.l && max.b == 0 {
			continue
		}
		c.t.lines[l].InsertTab()
		tabed = append(tabed, l)
	}
	for _, l := range tabed {
		if l == c.l && !c.AtBol() {
			c.SetOffsets(c.b+1)
		}
	}
	return tabed
}

func (c *Cursor) UnTab(sel *Selection) []int {
	untabed := make([]int, 0)
	if sel == nil {
		if err := c.t.lines[c.l].RemoveTab(); err == nil {
			untabed = append(untabed, c.l)
		}
		c.SetOffsets(c.b-1)
		return untabed
	}
	min, max := sel.MinMax()
	for l := min.l; l < max.l+1; l++ {
		if l == min.l && min.b == len(c.t.lines[min.l].data) {
			continue
		} else if l == max.l && max.b == 0 {
			continue
		}
		if err := c.t.lines[l].RemoveTab(); err == nil {
			untabed = append(untabed, l)
		}
	}
	for _, l := range untabed {
		if l == c.l && !c.AtBol() {
			c.SetOffsets(c.b-1)
		}
	}
	return untabed
}

func (c *Cursor) Delete() string {
	if c.AtEof() {
		return ""
	}
	if c.AtEol() {
		c.t.JoinNextLine(c.l)
		return "\n"
	}
	_, rlen := c.RuneAfter()
	return c.t.Remove(c.l, c.b, c.b+rlen)
}

func (c *Cursor) Backspace() string {
	if c.AtBof() {
		return ""
	}
	c.MoveLeft()
	return c.Delete()
}

func (c *Cursor) DeleteSelection(sel *Selection) string {
	min, max := sel.MinMax()
	bmin := Point{min.l, BFromC(c.t.lines[min.l].data, min.o)}
	bmax := Point{max.l, BFromC(c.t.lines[max.l].data, max.o)}
	// TODO : should get deleted strings from RemoveRange
	deleted := c.t.RemoveRange(bmin, bmax)
	c.l = min.l
	c.SetOffsets(bmin.o)
	return deleted
}

func (c *Cursor) GotoNext(find string) {
	for l := c.l; l < len(c.t.lines); l++ {
		linedata := string(c.t.lines[l].data)
		offset := 0
		if l == c.l {
			if c.b == len(linedata) {
				continue
			}
			linedata = linedata[c.b+1:]
			offset = c.b+1
		}
		b := strings.Index(linedata, find)
		if b != -1 {
			c.l = l
			c.SetOffsets(b+offset)
			break
		}
	}
}

func (c *Cursor) GotoPrev(find string) {
	for l := c.l; l >= 0; l-- {
		linedata := string(c.t.lines[l].data)
		if l == c.l {
			linedata = linedata[:c.b]
		}
		b := strings.LastIndex(linedata, find)
		if b != -1 {
			c.l = l
			c.SetOffsets(b)
			break
		}
	}
}

func (c *Cursor) GotoFirst(find string) {
	for l := 0; l < len(c.t.lines); l++ {
		linedata := string(c.t.lines[l].data)
		b := strings.Index(linedata, find)
		if b != -1 {
			c.l = l
			c.SetOffsets(b)
			break
		}
	}
}

func (c *Cursor) GotoLast(find string) {
	for l := len(c.t.lines)-1; l >= 0; l-- {
		linedata := string(c.t.lines[l].data)
		b := strings.LastIndex(linedata, find)
		if b != -1 {
			c.l = l
			c.SetOffsets(b)
			break
		}
	}
}

func (c *Cursor) GotoNextAny(chars string) {
	for l := c.l; l < len(c.t.lines); l++ {
		linedata := string(c.t.lines[l].data)
		offset := 0
		if l == c.l {
			if c.b == len(linedata) {
				continue
			}
			linedata = linedata[c.b+1:]
			offset = c.b+1
		}		
		b := strings.IndexAny(linedata, chars)
		if b != -1 {
			c.l = l
			c.SetOffsets(b+offset)
			break
		}
	}
}

func (c *Cursor) GotoPrevAny(chars string) {
	for l := c.l; l >= 0; l-- {
		linedata := string(c.t.lines[l].data)
		if l == c.l {
			linedata = linedata[:c.b]
		}
		b := strings.LastIndexAny(linedata, chars)
		if b != -1 {
			c.l = l
			c.SetOffsets(b)
			break
		}
	}
}

func (c *Cursor) GotoNextGlobalLineWithout(str string) {
	findLine := -1
	for l := c.l+1; l < len(c.t.lines); l++ {
		find := true
		for _, r := range str {
			if c.t.lines[l].data == "" {
				find = false
			} else if strings.HasPrefix(c.t.lines[l].data, string(r)) {
				find = false
			}
		}
		if find {
			findLine = l
			break
		}
	}
	if findLine != -1 {
		c.l = findLine
		c.SetOffsets(0)
		return
	}
}

func (c *Cursor) GotoPrevGlobalLineWithout(str string) {
	var startLine int
	if c.b == 0 {
		startLine = c.l - 1
	} else {
		startLine = c.l
	}
	findLine := -1
	for l := startLine; l >= 0; l-- {
		find := true
		for _, r := range str {
			if c.t.lines[l].data == "" {
				find = false
			} else if strings.HasPrefix(c.t.lines[l].data, string(r)) {
				find = false
			}
		}
		if find {
			findLine = l
			break
		}
	}
	if findLine != -1 {
		c.l = findLine
		c.SetOffsets(0)
		return
	}
}

func (c *Cursor) GotoNextDefinition(defn []string) {
	nextLines := c.t.lines[c.l+1:]
	for i, line := range nextLines {
		l := c.l + 1 + i
		find := false
		for _, d := range defn {
			if strings.HasPrefix(string(line.data), d) {
				find = true
				break
			}
		}
		if find {
			c.l = l
			c.SetOffsets(0)
			break
		}
	}
}

func (c *Cursor) GotoPrevDefinition(defn []string) {
	var startLine int
	if c.b == 0 {
		startLine = c.l - 1
	} else {
		startLine = c.l
	}
	find := false
	for l := startLine; l >= 0; l-- {
		for _, d := range defn {
			if strings.HasPrefix(string(c.t.lines[l].data), d) {
				find = true
				break
			}
		}
		if find {
			c.l = l
			c.SetOffsets(0)
			break
		}
	}
}

func (c *Cursor) GotoMatchingBracket() {
	rb, _ := c.RuneBefore()
	ra, _ := c.RuneAfter()
	var r rune
	dir := ""
	if strings.Contains("{[(", string(rb)) {
		r = rb
		dir = "right"
	}
	if strings.Contains("}])", string(ra)) {
		r = ra
		dir = "left"
	}
	if dir == "" {
		return
	}
	// rune for matching.
	var m rune
	switch r {
	case '{':
		m = '}'
	case '}':
		m = '{'
	case '[':
		m = ']'
	case ']':
		m = '['
	case '(':
		m = ')'
	case ')':
		m = '('
	}
	if dir == "left" && rb == m {
		return
	} else if dir == "right" && ra == m {
		return
	}
	set := string(r) + string(m)
	depth := 0
	origc := *c
	for {
		bc := *c
		if dir == "right" {
			c.GotoNextAny(set)
		} else {
			c.GotoPrevAny(set)
		}
		if c.l == bc.l && c.o == bc.o {
			// did not find next set.
			c.Copy(origc)
			return
		}
		if c.InStrings() {
			continue
		}
		cr, _ := c.RuneAfter()
		if cr == r {
			depth++
		} else if cr == m {
			if depth == 0 {
				if dir == "left" {
					c.MoveRight()
				}
				return
			}
			depth--
		}
	}
}

func (c *Cursor) GotoLine(l int) {
	if l >= len(c.t.lines) {
		l = len(c.t.lines)-1
	}
	c.l = l
	c.SetOffsets(0)
}

func (c *Cursor) Word() string {
	// check cursor is on a word
	r, _ := c.RuneAfter()
	if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
		return ""
	}
	// find min byte offset
	bmin := c.b
	remain := c.LineDataUntilCursor()
	for {
		if len(remain) == 0 {
			break
		}
		r, rlen := utf8.DecodeLastRuneInString(remain)
		remain = remain[:len(remain)-rlen]
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			break
		}
		bmin -= rlen
	}
	// find max byte offset
	bmax := c.b
	remain = c.LineDataFromCursor()
	for {
		if len(remain) == 0 {
			break
		}
		r, rlen := utf8.DecodeRuneInString(remain)
		remain = remain[rlen:]
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			break
		}
		bmax += rlen
	}
	return c.LineData()[bmin:bmax]
}
