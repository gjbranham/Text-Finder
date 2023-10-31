package output

import (
	"bytes"
	"fmt"
	"os"
)

var P *Printer

type Printer struct {
	Printer MyPrinter
}

type MyPrinter interface {
	Print(text string)
}

type Buffer struct {
	Buf *bytes.Buffer
}

func (b *Buffer) Print(text string) {
	fmt.Fprint(b.Buf, text)
}

func (b *Buffer) ReadAndClear() string {
	s := b.Buf.String()
	b.Buf.Reset()
	return s
}

type Stdout struct{}

func (s *Stdout) Print(text string) {
	fmt.Fprint(os.Stdout, text)
}

func SetPrinter(p MyPrinter) {
	P = &Printer{Printer: p}
}

func Print(text string) {
	P.Printer.Print(text + "\n")
}
