package jq

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/thekhanj/digikala-api/cli/internal"
	"golang.org/x/sys/unix"
)

type keyValOption struct {
	key string
	val string
}

type option interface {
	string | keyValOption
}

type Jq struct {
	filter    io.Reader
	filePaths []string
	files     []io.Reader
	options   []interface{}
	ranOnce   bool
}

func (this *Jq) Start() ([]byte, error) {
	if this.ranOnce {
		return nil, errors.New("cannot call function more than once")
	}
	this.ranOnce = true

	filterFifo, err := this.newFifo(this.filter)
	if err != nil {
		return nil, err
	}

	files := this.filePaths
	for _, file := range this.files {
		f, err := this.newFifo(file)
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	options := this.getOptions()

	args := make([]string, 0)
	args = append(args, options...)
	args = append(args, "-f", filterFifo)
	args = append(args, files...)

	cmd := exec.Command("jq", args...)

	return cmd.Output()
}

func (this *Jq) getOptions() []string {
	ret := make([]string, 0)

	for _, o := range this.options {
		switch v := o.(type) {
		case string:
			ret = append(ret, v)
		case keyValOption:
			ret = append(ret, v.key, v.val)
		default:
			panic(fmt.Errorf("invalid option: %v", v))
		}
	}

	return ret
}

func (this *Jq) newFifo(r io.Reader) (string, error) {
	err := internal.AssertDir(internal.GetProcDir())
	if err != nil {
		return "", err
	}

	tmp, err := os.CreateTemp(internal.GetProcDir(), "jq-fifo-*")
	if err != nil {
		return "", err
	}

	path := tmp.Name()
	err = tmp.Close()
	if err != nil {
		return "", err
	}

	err = os.Remove(path)
	if err != nil {
		return "", err
	}

	err = unix.Mkfifo(path, 0600)
	if err != nil {
		return "", err
	}

	go func() {
		file, err := os.OpenFile(path, os.O_WRONLY, 0600)
		if err != nil {
			log.Printf("jq: opening fifo failed: %v", err)
			return
		}

		defer func() {
			err = file.Close()
			if err != nil {
				log.Printf("jq: closing fifo failed: %v", err)
			}

			err = os.Remove(path)
			if err != nil {
				log.Printf("jq: closing fifo failed: %v", err)
			}
		}()

		b, err := io.Copy(file, r)
		if err != nil {
			log.Printf(
				"jq: write failed: wrote %d bytes before error: %v",
				b, err,
			)
		}
	}()

	return path, nil
}

type JqBuilder struct {
	jq Jq
}

func (this *JqBuilder) WithFlag(flag string) *JqBuilder {
	this.jq.options = append(this.jq.options, flag)

	return this
}

func (this *JqBuilder) WithOption(key, val string) *JqBuilder {
	this.jq.options = append(this.jq.options, keyValOption{key, val})

	return this
}

func (this *JqBuilder) WithFilter(filter io.Reader) *JqBuilder {
	this.jq.filter = filter

	return this
}

func (this *JqBuilder) WithFilterString(filter string) *JqBuilder {
	r := strings.NewReader(filter)

	return this.WithFilter(r)
}

func (this *JqBuilder) WithFile(file io.Reader) *JqBuilder {
	this.jq.files = append(this.jq.files, file)

	return this
}

func (this *JqBuilder) WithFilePath(path string) *JqBuilder {
	this.jq.filePaths = append(this.jq.filePaths, path)

	return this
}

func (this *JqBuilder) WithFileData(data []byte) *JqBuilder {
	r := bytes.NewReader(data)

	return this.WithFile(r)
}

func (this *JqBuilder) Build() (*Jq, error) {
	if this.jq.filter == nil {
		return nil, errors.New("no filter provided")
	}

	if len(this.jq.files)+len(this.jq.filePaths) == 0 {
		return nil, errors.New("no input files provided")
	}

	return &this.jq, nil
}

func NewJqBuilder() *JqBuilder {
	b := JqBuilder{
		jq: Jq{
			filter:    nil,
			filePaths: make([]string, 0),
			files:     make([]io.Reader, 0),
			options:   make([]interface{}, 0),
			ranOnce:   false,
		},
	}

	return &b
}

func NewJq(
	data []byte, filter string, flags ...string,
) (*Jq, error) {
	b := NewJqBuilder().
		WithFileData(data).
		WithFilterString(filter)

	for _, f := range flags {
		b.WithFlag(f)
	}

	return b.Build()
}
