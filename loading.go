package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tj/go-spin"
)

type Option struct {
	Interval       int
	LoadingMessage string
	FinishMessage  string
}

type Reader struct {
	io.Reader
	spinner *spin.Spinner
	i       int
	opt     *Option
}

func NewReader(r io.Reader, opt Option) *Reader {
	return &Reader{
		Reader:  r,
		spinner: spin.New(),
		i:       0,
		opt:     &opt,
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if err != nil {
		if err == io.EOF {
			fmt.Fprintf(os.Stderr, "\033[1;32m%s\033[m\n", r.opt.FinishMessage)
		}
		return
	}

	if r.i == r.opt.Interval {
		fmt.Fprintf(os.Stderr, "\r  \033[36m%s\033[m %s ",
			r.opt.LoadingMessage, r.spinner.Next())
		r.i = 0
	}

	r.i++
	return
}
