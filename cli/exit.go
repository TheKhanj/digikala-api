package main

import (
	"fmt"
	"os"
)

const (
	CODE_SUCCESS int = iota
	CODE_GENERAL_ERROR
	CODE_BAD_INVOCATION_ERROR
	CODE_READING_CONFIG_ERROR
	CODE_NOT_IMPLEMENTED_ERROR
	CODE_INVALID_CONFIG_ERROR
)

type Exit struct {
	err      error
	code     int
	afterMsg string
}

var (
	SuccessfulExit    = NewExit(nil, CODE_SUCCESS)
	BadInvocationExit = NewHelpfulExit(nil, CODE_BAD_INVOCATION_ERROR)
)

var HelpMessage = `Usage:
  digi [-c config.json] fetch products
  digi [-c config.json] gen-schema products`

func (this *Exit) Run() {
	if this.err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", this.err)
	}

	if this.afterMsg != "" {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, this.afterMsg)
	}

	os.Exit(this.code)
}

func NewExit(err error, code int) Exit {
	return NewMsgExit(err, code, "")
}

func NewHelpfulExit(err error, code int) Exit {
	return NewMsgExit(err, code, HelpMessage)
}

func NewMsgExit(err error, code int, afterMsg string) Exit {
	return Exit{
		err:      err,
		code:     code,
		afterMsg: afterMsg,
	}
}
