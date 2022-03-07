package main

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

type BlockWriter struct {
	Width      int
	Indent     int
	TextBorder int

	w     prompt.ConsoleWriter
	stack []prompt.Color
}

func NewBlockWriter(width, indent, textBorder int) *BlockWriter {
	bw := &BlockWriter{
		w:          prompt.NewStdoutWriter(),
		Width:      width,
		Indent:     indent,
		TextBorder: textBorder,
	}
	bw.w.WriteStr("\n")
	return bw
}

func (bw *BlockWriter) Open(bgColor, textColor prompt.Color, title string) {
	bw.renderLeftBorder()
	size := bw.Width - bw.Indent*(len(bw.stack)) - bw.TextBorder
	bw.EmptyTitleLn(bgColor, textColor, size, title)
	bw.stack = append(bw.stack, bgColor)
}

func (bw *BlockWriter) Close() {
	if len(bw.stack) > 0 {
		bw.renderLeftBorder()

		current := bw.stack[len(bw.stack)-1]
		size := bw.Width - bw.Indent*(len(bw.stack)+1)
		bw.EmptyLn(current, size)
		bw.stack = bw.stack[:len(bw.stack)-1]
	}
}

func (bw *BlockWriter) Empty(c prompt.Color, size int) {
	ln := ""
	for i := 0; i < size; i++ {
		ln += " "
	}
	bw.w.SetColor(c, c, false)
	bw.w.WriteStr(ln)
	bw.w.Flush()
}

func (bw *BlockWriter) EmptyTitleLn(bgColor, textColor prompt.Color, size int, title string) {
	ind := ""
	for i := 0; i < bw.TextBorder; i++ {
		ind += " "
	}
	ind += title

	ln := ""
	for i := 0; i < size-len(ind); i++ {
		ln += " "
	}
	ln += "\n"

	bw.w.SetColor(textColor, bgColor, true)
	bw.w.WriteStr(ind)

	bw.w.Flush()
	bw.w.SetColor(bgColor, bgColor, true)
	bw.w.WriteStr(ln)
	bw.w.Flush()
}

func (bw *BlockWriter) EmptyLn(c prompt.Color, size int) {
	ln := ""
	for i := 0; i < size; i++ {
		ln += " "
	}
	ln += "\n"
	bw.w.SetColor(c, c, false)
	bw.w.WriteStr(ln)
	bw.w.Flush()
}

func (bw *BlockWriter) renderLeftBorder() {
	for _, indent := range bw.stack {
		bw.Empty(indent, bw.Indent)
	}
}

func (bw *BlockWriter) WriteText(fg, bg prompt.Color, bold bool, text string) {
	for _, s := range strings.Split(text, "\n") {
		bw.WriteLine(fg, bg, bold, s+"\n")
	}
}

func (bw *BlockWriter) WriteLine(fg, bg prompt.Color, bold bool, line string) {
	bw.renderLeftBorder()
	bw.Empty(bg, bw.TextBorder)
	bw.w.SetColor(fg, bg, bold)
	bw.w.WriteStr(line)
	bw.w.Flush()
}
