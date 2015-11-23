package main

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"strconv"
	"strings"
)

func parseEvent(ev term.Event, t *Text, sel *Selection, mode *string) []*Action {
	if ev.Type != term.EventKey {
		panic(fmt.Sprintln("what the..", ev.Type, "event?"))
	}

	switch ev.Key {
	case term.KeyCtrlW:
		return []*Action{{kind: "selection", value: "off"}, {kind: "exit"}}
	case term.KeyCtrlS:
		return []*Action{{kind: "selection", value: "off"}, {kind: "save"}}
	case term.KeyCtrlK:
		return []*Action{{kind: "selection", value: "off"}}
	// move
	case term.KeyArrowLeft:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "left"}}
	case term.KeyArrowRight:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "right"}}
	case term.KeyArrowUp:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "up"}}
	case term.KeyArrowDown:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "down"}}
	case term.KeyPgup:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "pageup"}}
	case term.KeyPgdn:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "pagedown"}}
	case term.KeyHome:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "bol"}}
	case term.KeyEnd:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "eol"}}
	// insert
	case term.KeyEnter:
		return []*Action{{kind: "deleteSelection"}, {kind: "selection", value: "off"}, {kind: "insert", value: "\n"}}
	case term.KeyCtrlN:
		return []*Action{{kind: "deleteSelection"}, {kind: "selection", value: "off"}, {kind: "insert", value: "\n"}, {kind: "insert", value: "autoIndent"}}
	case term.KeySpace:
		return []*Action{{kind: "deleteSelection"}, {kind: "selection", value: "off"}, {kind: "insert", value: " "}}
	case term.KeyTab:
		tab := "\t"
		if t.tabToSpace {
			tab = strings.Repeat(" ", t.tabWidth)
		}
		return []*Action{{kind: "deleteSelection"}, {kind: "selection", value: "off"}, {kind: "insert", value: tab}}
	case term.KeyCtrlU:
		return []*Action{{kind: "removeTab"}}
	case term.KeyCtrlO:
		return []*Action{{kind: "insertTab"}}
	// delete : value will added after actual deletion.
	case term.KeyDelete:
		if sel.on {
			return []*Action{{kind: "deleteSelection"}, {kind: "selection", value: "off"}}
		} else {
			return []*Action{{kind: "delete"}}
		}
	case term.KeyBackspace, term.KeyBackspace2:
		if sel.on {
			return []*Action{{kind: "deleteSelection"}, {kind: "selection", value: "off"}}
		} else {
			return []*Action{{kind: "backspace"}}
		}
	// undo, redo
	case term.KeyCtrlZ:
		return []*Action{{kind: "undo"}}
	case term.KeyCtrlY:
		return []*Action{{kind: "redo"}}
	// copy, paste, cut
	case term.KeyCtrlC:
		return []*Action{{kind: "copy"}, {kind: "selection", value: "off"}}
	case term.KeyCtrlV:
		if sel.on {
			return []*Action{{kind: "deleteSelection"}, {kind: "selection", value: "off"}, {kind: "paste"}}
		}
		return []*Action{{kind: "paste"}}
	case term.KeyCtrlJ:
		if sel.on {
			return []*Action{{kind: "deleteSelection"}, {kind: "selection", value: "off"}, {kind: "replace"}}
		}
		return []*Action{}
	case term.KeyCtrlX:
		if sel.on {
			return []*Action{{kind: "copy"}, {kind: "deleteSelection"}, {kind: "selection", value: "off"}}
		} else {
			return []*Action{{kind: "copy"}, {kind: "delete"}}
		}
	// find
	case term.KeyCtrlD, term.KeyF3:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "findNextSelect"}}
	case term.KeyCtrlB, term.KeyF2:
		return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "findPrevSelect"}}
	case term.KeyCtrlF:
		return []*Action{{kind: "modeChange", value: "find"}}
	case term.KeyCtrlR:
		return []*Action{{kind: "modeChange", value: "replace"}}
	case term.KeyCtrlG:
		return []*Action{{kind: "modeChange", value: "gotoline"}}
	case term.KeyCtrlL:
		return []*Action{{kind: "selectLine"}}
	default:
		if ev.Ch == 0 {
			return []*Action{{kind: "none"}}
		}
		if ev.Mod&term.ModAlt != 0 {
			switch ev.Ch {
			case 'j':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "left"}}
			case 'J':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "left"}}
			case 'l':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "right"}}
			case 'L':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "right"}}
			case 'i':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "up"}}
			case 'I':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "up"}}
			case 'k':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "down"}}
			case 'K':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "down"}}
			case 'm':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "prevBowEow"}}
			case 'M':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "prevBowEow"}}
			case '.':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "nextBowEow"}}
			case '>':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "nextBowEow"}}
			case 'u':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "bocBolRepeat"}}
			case 'U':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "bocBolRepeat"}}
			case 'y':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "bol"}}
			case 'Y':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "bol"}}
			case 'o':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "eol"}}
			case 'O':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "eol"}}
			case 'w':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "pageup"}}
			case 'W':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "pageup"}}
			case 's':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "pagedown"}}
			case 'S':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "pagedown"}}
			case 'q':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "bof"}}
			case 'Q':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "bof"}}
			case 'a':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "eof"}}
			case 'A':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "eof"}}
			case 'n':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "nextGlobal"}}
			case 'N':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "nextGlobal"}}
			case 'h':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "prevGlobal"}}
			case 'H':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "prevGlobal"}}
			case ']', 'x':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "nextArg"}}
			case '}', 'X':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "nextArg"}}
			case '[', 'z':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "prevArg"}}
			case '{', 'Z':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "prevArg"}}
			case 'd':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "findNext"}}
			case 'b':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "findPrev"}}
			case 'c':
				return []*Action{{kind: "selection", value: "off"}, {kind: "move", value: "matchingBracket"}}
			case 'C':
				return []*Action{{kind: "selection", value: "on"}, {kind: "move", value: "matchingBracket"}}
			default:
				return []*Action{{kind: "none"}}
			}
		}
		if sel.on {
			return []*Action{{kind: "deleteSelection"}, {kind: "insert", value: string(ev.Ch)}}
		} else {
			return []*Action{{kind: "insert", value: string(ev.Ch)}}
		}
	}
}

func do(a *Action, t *Text, c *Cursor, sel *Selection, history *History, status *string, holdStatus *bool, findstr string) {
	defer func() {
		if sel.on {
			sel.SetEnd(c)
		}
	}()
	switch a.kind {
	case "none":
		return
	case "selection":
		if a.value == "on" && !sel.on {
			sel.on = true
			sel.SetStart(c)
		} else if a.value == "off" {
			sel.on = false
		}
	case "move":
		switch a.value {
		case "left":
			c.MoveLeft()
		case "right":
			c.MoveRight()
		case "up":
			c.MoveUp()
		case "down":
			c.MoveDown()
		case "prevBowEow":
			c.MovePrevBowEow()
		case "nextBowEow":
			c.MoveNextBowEow()
		case "bol":
			c.MoveBol()
		case "eol":
			c.MoveEol()
		case "bocBolRepeat":
			c.MoveBocBolRepeat()
		case "pageup":
			c.PageUp()
		case "pagedown":
			c.PageDown()
		case "bof":
			c.MoveBof()
		case "eof":
			c.MoveEof()
		case "nextGlobal":
			c.GotoNextGlobalLineWithout(" \t#/{}()")
		case "prevGlobal":
			c.GotoPrevGlobalLineWithout(" \t#/{}()")
		case "nextArg":
			c.GotoNextAny("{(,)}")
			r, _ := c.RuneAfter()
			if r == '(' || r == '{' {
				c.MoveRight()
			}
		case "prevArg":
			r, _ := c.RuneBefore()
			if r == '(' || r == '{' {
				c.MoveLeft()
			}
			c.GotoPrevAny("{(,)}")
			r, _ = c.RuneAfter()
			if r == '(' || r == '{' {
				c.MoveRight()
			}
		case "matchingBracket":
			c.GotoMatchingBracket()
		case "findPrev":
			ok := c.GotoPrev(findstr)
			if !ok {
				c.GotoLast(findstr)
			}
		case "findNext":
			ok := c.GotoNext(findstr)
			if !ok {
				c.GotoFirst(findstr)
			}
		case "findPrevWord":
			c.GotoPrevWord(findstr)
		case "findNextWord":
			c.GotoNextWord(findstr)
		// TODO: "findPrevSelect" and "findNextSelect" are hack. make separate action.
		case "findPrevSelect":
			ok := c.GotoPrev(findstr)
			if !ok {
				ok = c.GotoLast(findstr)
			}
			if ok {
				sel.on = true
				for range findstr {
					c.MoveRight()
				}
				sel.SetStart(c)
				for range findstr {
					c.MoveLeft()
				}
				sel.SetEnd(c)
			}
		case "findNextSelect":
			ok := c.GotoNext(findstr)
			if !ok {
				ok = c.GotoFirst(findstr)
			}
			if ok {
				sel.on = true
				for range findstr {
					c.MoveRight()
				}
				sel.SetStart(c)
				for range findstr {
					c.MoveLeft()
				}
				sel.SetEnd(c)
			}
		default:
			panic(fmt.Sprintln("what the..", a.value, "move?"))
		}
	case "insert":
		if a.value == "autoIndent" {
			prevline := t.lines[c.l-1].data
			trimed := strings.TrimLeft(prevline, " \t")
			indent := prevline[:len(prevline)-len(trimed)]
			c.Insert(indent)
			a.value = indent
			return
		}
		c.Insert(a.value)
	case "delete":
		a.value = c.Delete()
	case "insertTab":
		tab := "\t"
		if t.tabToSpace {
			tab = strings.Repeat(" ", t.tabWidth)
		}
		lines := make([]int, 0)
		if sel.on {
			min, max := sel.MinMax()
			if max.o == 0 {
				max.l--
			}
			for l := min.l; l <= max.l; l++ {
				lines = append(lines, l)
			}
		} else {
			lines = append(lines, c.l)
		}
		tabedLine := ""
		for _, l := range lines {
			t.Line(l).Insert(tab, 0)
			if tabedLine != "" {
				tabedLine += ","
			}
			tabedLine += strconv.Itoa(l) + ":" + tab
			if l == c.l {
				c.SetB(c.b + len(tab))
			}
		}
		a.value = tabedLine
	case "removeTab":
		// removeTab is slightly differ from insertTab.
		// removeTab should remember what is removed, not tab string it self.
		lines := make([]int, 0)
		if sel.on {
			min, max := sel.MinMax()
			if max.o == 0 {
				max.l--
			}
			for l := min.l; l <= max.l; l++ {
				lines = append(lines, l)
			}
		} else {
			lines = append(lines, c.l)
		}
		untabedLine := ""
		for _, l := range lines {
			removed := ""
			if strings.HasPrefix(t.Line(l).data, "\t") {
				removed += t.Line(l).Remove(0, 1)
			} else {
				for i := 0; i < t.tabWidth; i++ {
					if len(t.Line(l).data) == 0 {
						break
					}
					if !strings.HasPrefix(t.Line(l).data, " ") {
						break
					}
					removed += t.Line(l).Remove(0, 1)
				}
			}
			if untabedLine != "" {
				untabedLine += ","
			}
			untabedLine += strconv.Itoa(l) + ":" + removed
			if l == c.l && !c.AtBol() {
				c.SetB(c.b-len(removed))
			}
		}
		a.value = untabedLine
	case "backspace":
		a.value = c.Backspace()
	case "deleteSelection":
		if sel.on {
			a.value = c.DeleteSelection(sel)
			sel.on = false
		}
	case "selectLine":
		c.MoveBol()
		if !sel.on {
			sel.on = true
			sel.SetStart(c)
		}
		if c.OnLastLine() {
			c.MoveEol()
		} else {
			c.MoveDown()
		}
		sel.SetEnd(c)
	case "selectWord":
		if !c.AtBow() {
			c.MovePrevBowEow()
		}
		if !sel.on {
			sel.on = true
			sel.SetStart(c)
		}
		c.MoveNextBowEow()
		sel.SetEnd(c)
	case "undo":
		if history.head == 0 {
			return
		}
		sel.on = false
		history.head--
		action := history.At(history.head)
		switch action.kind {
		case "insert":
			c.Copy(action.afterCursor)
			for range action.value {
				c.Backspace()
			}
		case "insertTab":
			lineInfos := strings.Split(action.value, ",")
			for _, li := range lineInfos {
				if li == "" {
					continue
				}
				lis := strings.Split(li, ":")
				lstr := lis[0]
				tab := lis[1]
				l, err := strconv.Atoi(lstr)
				if err != nil {
					panic(err)
				}
				for _, r := range tab {
					rr := t.Line(l).Remove(0, 1)
					if rr != string(r) {
						panic("removed and current is not matched")
					}
				}
			}
			c.Copy(action.beforeCursor)
		case "paste", "replace":
			c.Copy(action.beforeCursor)
			for range action.value {
				c.Delete()
			}
		case "backspace":
			c.Copy(action.afterCursor)
			c.Insert(action.value)
		case "delete", "deleteSelection":
			c.Copy(action.afterCursor)
			c.Insert(action.value)
		case "removeTab":
			lineInfos := strings.Split(action.value, ",")
			for _, li := range lineInfos {
				if li == "" {
					continue
				}
				lis := strings.Split(li, ":")
				lstr := lis[0]
				removed := lis[1]
				l, err := strconv.Atoi(lstr)
				if err != nil {
					panic(err)
				}
				t.Line(l).Insert(removed, 0)
			}
			c.Copy(action.beforeCursor)
		default:
			panic(fmt.Sprintln("what the..", action.kind, "history?"))
		}
	case "redo":
		if history.head == history.Len() {
			return
		}
		sel.on = false
		action := history.At(history.head)
		history.head++
		switch action.kind {
		case "insert":
			c.Copy(action.beforeCursor)
			c.Insert(action.value)
		case "insertTab":
			lineInfos := strings.Split(action.value, ",")
			for _, li := range lineInfos {
				if li == "" {
					continue
				}
				lis := strings.Split(li, ":")
				lstr := lis[0]
				tab := lis[1]
				l, err := strconv.Atoi(lstr)
				if err != nil {
					panic(err)
				}
				t.Line(l).Insert(tab, 0)
			}
			c.Copy(action.afterCursor)
		case "paste", "replace":
			c.Copy(action.beforeCursor)
			c.Insert(action.value)
		case "backspace":
			c.Copy(action.beforeCursor)
			for range action.value {
				c.Backspace()
			}
		case "delete", "deleteSelection":
			c.Copy(action.beforeCursor)
			for range action.value {
				c.Delete()
			}
		case "removeTab":
			lineInfos := strings.Split(action.value, ",")
			for _, li := range lineInfos {
				if li == "" {
					continue
				}
				lis := strings.Split(li, ":")
				lstr := lis[0]
				removed := lis[1]
				l, err := strconv.Atoi(lstr)
				if err != nil {
					panic(err)
				}
				for _, r := range removed {
					rr := t.Line(l).Remove(0, 1)
					if rr != string(r) {
						panic("removed and current is not matched")
					}
				}
			}
			c.Copy(action.afterCursor)
		default:
			panic(fmt.Sprintln("what the..", action.kind, "history?"))
		}
	default:
		panic(fmt.Sprintln("what the..", a.kind, "action?"))
	}
}
