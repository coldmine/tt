package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	term "github.com/nsf/termbox-go"
)

func main() {
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var newFlag bool
	flagset.BoolVar(&newFlag, "new", false, "let tor to edit a new file.")

	// sort args, so let flags always placed ahead of file arg.
	args := os.Args[1:]
	sort.Strings(args)
	flagset.Parse(args)

	fileArgs := flagset.Args()
	if len(fileArgs) == 0 {
		fmt.Println("please, set text file")
		os.Exit(1)
	}
	farg := fileArgs[len(fileArgs)-1]

	f, initL, initO, err := parseFileArg(farg)
	if err != nil {
		fmt.Println("file arg is invalid: ", err)
		os.Exit(1)
	}

	exist := true
	if _, err := os.Stat(f); os.IsNotExist(err) {
		exist = false
	}
	if !exist && !newFlag {
		fmt.Println("file not exist. please retry with -new flag.")
		os.Exit(1)
	} else if exist && newFlag {
		fmt.Println("file already exist.")
		os.Exit(1)
	}

	var text *Text
	if exist {
		text, err = open(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		lines := make([]Line, 0)
		lines = append(lines, Line{""})
		text = &Text{lines: lines, tabToSpace: false, tabWidth: 4, edited: false}
	}

	err = term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	term.SetInputMode(term.InputAlt)
	// term.SetOutputMode(term.Output256)
	term.Clear(term.ColorDefault, term.ColorDefault)
	term.Flush()

	termw, termh := term.Size()
	mainarea := NewArea(Point{0, 0}, Point{termh - 1, termw})
	win := NewWindow(mainarea.Size())
	cursor := &Cursor{}
	selection := &Selection{}
	history := newHistory()

	mode := &ModeSelector{}
	mode.normal = &NormalMode{
		text:      text,
		cursor:    cursor,
		selection: selection,
		history:   history,
		f:         f,
		mode:      mode,
		copied:    loadConfig("copy"),
	}
	mode.find = &FindMode{
		mode: mode,
		str:  loadConfig("find"),
	}
	mode.replace = &ReplaceMode{
		mode: mode,
		str:  loadConfig("replace"),
	}
	mode.gotoline = &GotoLineMode{
		cursor: cursor,
		mode:   mode,
	}
	mode.exit = &ExitMode{
		f:      f,
		cursor: cursor,
		mode:   mode,
	}
	mode.current = mode.normal // will start tor as normal mode.

	selection.text = mode.normal.text

	// Set cursor.
	cursor.text = mode.normal.text
	if initL != -1 {
		l := initL
		// to internal line number
		if l != 0 {
			l--
		}
		cursor.GotoLine(l)
		if initO != -1 {
			cursor.SetO(initO)
		}
	} else {
		l, b := loadLastPosition(f)
		cursor.GotoLine(l)
		cursor.SetCloseToB(b)
	}

	events := make(chan term.Event, 20)
	go func() {
		for {
			events <- term.PollEvent()
		}
	}()
	for {
		win.Follow(cursor, 3)
		clearScreen(mainarea)
		drawScreen(mainarea, win, mode.normal.text, selection, cursor)
		if mode.current.Error() != "" {
			printErrorStatus(mode.current.Error())
		} else {
			printStatus(mode.current.Status())
		}
		if mode.current == mode.normal {
			winP := cursor.Position().Sub(win.min)
			term.SetCursor(mainarea.min.o+winP.o, mainarea.min.l+winP.l)
		} else {
			term.SetCursor(vlen(mode.current.Status(), mode.normal.text.tabWidth), termh)
		}
		term.Flush()

		// wait for keyboard input
		select {
		case ev := <-events:
			switch ev.Type {
			case term.EventKey:
				mode.current.Handle(ev)
			case term.EventResize:
				term.Clear(term.ColorDefault, term.ColorDefault)
				termw, termh = term.Size()
				resizeScreen(mainarea, win, termw, termh)
			}
		}
	}
}