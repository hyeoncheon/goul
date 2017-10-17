package goul

import (
	"errors"
	"fmt"
	"os"
)

//*** pipeline router -----------------------------------------------

// Pipeline is a structure of pipeline router
type Pipeline struct {
	Router
	err   error
	pipes []Pipe
}

// Run implements Router
func (r *Pipeline) Run() (ctrl, done chan Item, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Pipeline#Run recovered from panic!\n")
			fmt.Fprintf(os.Stderr, "Probably an inheritance problem of pipeline instance.\n")
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			err = errors.New("panic")
		}
	}()

	if r.getReader() == nil || r.getWriter() == nil {
		return nil, nil, errors.New(ErrRouterNoReaderOrWriter)
	}

	ctrl = make(chan Item)
	var ch chan Item
	ch, r.err = r.getReader().Read(ctrl, nil)
	if r.err != nil {
		return nil, nil, r.err
	}
	for _, p := range r.GetPipes() {
		if p.IsConverter() {
			ch, r.err = p.Convert(ch, nil)
		} else {
			ch, r.err = p.Revert(ch, nil)
		}
		if r.err != nil {
			return nil, nil, r.err
		}
	}
	done, r.err = r.getWriter().Write(ch, nil)
	if r.err != nil {
		return nil, nil, r.err
	}

	Log(r.getLogger(), "router", "started ---------------------------------")
	return ctrl, done, nil
}

// AddPipe implements Router
func (r *Pipeline) AddPipe(pipe Pipe) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Pipeline#AddPipe recovered from panic!\n")
			fmt.Fprintf(os.Stderr, "Probably an inheritance problem of pipeline instance.\n")
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			err = errors.New("panic")
		}
	}()

	pipe.SetLogger(r.getLogger())
	r.pipes = append(r.pipes, pipe)
	return nil
}

// GetPipes implements Router
func (r *Pipeline) GetPipes() []Pipe {
	if r.pipes == nil {
		return []Pipe{}
	}
	return r.pipes
}
