package main

type selection struct {
	on bool
	start Point
	end Point
}

func NewSelection() *selection {
	return &selection{}
}

func (s *selection) SetStart(c *cursor) {
	s.start = Point{c.line, c.offset()}
}

func (s *selection) SetEnd(c *cursor) {
	s.end = Point{c.line, c.offset()}
}

func (s *selection) Contains(p Point) bool {
	min := s.start
	max := s.end
	if (s.start.l > s.end.l) || (s.start.l == s.end.l && s.start.o > s.end.o) {
		min = s.end
		max = s.start
	}
	if min.l <= p.l && p.l <= max.l {
		if p.l == min.l && p.o < min.o {
			return false
		} else if p.l == max.l && p.o >= max.o {
			return false
		}
		return true
	}
	return false
}


func withShift(ch rune) bool {
	shifts := "QWERTYUIOP{}|ASDFGHJKL:ZXCVBNM<>?!@#$%^&*()_+"
	for _, sch := range shifts {
		if ch == sch {
			return true
		}
	}
	return false
}
