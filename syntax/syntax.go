package syntax

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

func init() {
	Syntaxes["go"] = Syntax{
		Matchers: []Matcher{
			Matcher{"string", regexp.MustCompile(`^(?m)".*?[^\\]"`), termbox.ColorRed, termbox.ColorBlack},
			Matcher{"rune", regexp.MustCompile(`^(?m)'.*?[^\\]'`), termbox.ColorYellow, termbox.ColorBlack},
			Matcher{"comment", regexp.MustCompile(`^(?m)//.*`), termbox.ColorMagenta, termbox.ColorBlack},
			Matcher{"multi line comment", regexp.MustCompile(`^(?s)/[*].*?[*]/`), termbox.ColorMagenta, termbox.ColorBlack},
			Matcher{"trailing spaces", regexp.MustCompile(`^(?m)[ \t]+$`), termbox.ColorBlack, termbox.ColorYellow},
		},
	}
}

type Syntax struct {
	Matchers []Matcher
}

var Syntaxes = make(map[string]Syntax)

func (s Syntax) Parse(text []byte) []Match {
	c := NewCursor(text)
	matches := []Match{}
Loop:
	for {
		for _, matcher := range s.Matchers {
			ms := matcher.Re.FindSubmatch(c.Remain())
			if ms != nil {
				m := ms[0]
				if len(ms) == 2 {
					m = ms[1]
					fmt.Printf("%s\n", m)
				}
				if string(m) == "" {
					continue
				}
				start := c.Pos()
				c.Skip(len(m))
				end := c.Pos()
				matches = append(matches, matcher.NewMatch(start, end))
				continue Loop
			}
		}
		if !c.Advance() {
			break
		}
	}
	return matches
}

type Cursor struct {
	text []byte
	b    int // byte offset
	l    int // line offset
	o    int // byte in line offset
}

func NewCursor(text []byte) *Cursor {
	return &Cursor{text: text}
}

func (c *Cursor) Pos() Pos {
	return Pos{c.l, c.o}
}

func (c *Cursor) Remain() []byte {
	if c.l == len(c.text) {
		return []byte("")
	}
	return c.text[c.b:]
}

func (c *Cursor) Advance() bool {
	if c.b == len(c.text) {
		return false
	}
	c.next()
	return true
}

func (c *Cursor) Skip(b int) {
	i := 0
	for i < b {
		_, size := c.next()
		i += size
	}
}

func (c *Cursor) next() (r rune, size int) {
	r, size = utf8.DecodeRune(c.Remain())
	c.b += size
	c.o += size
	if r == '\n' {
		c.l += 1
		c.o = 0
	}
	return r, size
}

type Matcher struct {
	Name string
	Re   *regexp.Regexp
	Fg   termbox.Attribute
	Bg   termbox.Attribute
}

func (m Matcher) NewMatch(start, end Pos) Match {
	return Match{Name: m.Name, Start: start, End: end, Fg: m.Fg, Bg: m.Bg}
}

type Match struct {
	Name  string
	Start Pos
	End   Pos
	Fg    termbox.Attribute
	Bg    termbox.Attribute
}

type Pos struct {
	L int
	O int
}
