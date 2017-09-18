package goul

import "errors"

//*** pipeline router -----------------------------------------------

// Pipeline is a structure of pipeline router
type Pipeline struct {
	Router
	err   error
	pipes []Pipe
}

// Run implements Router
func (r *Pipeline) Run() (chan Item, chan Item, error) {
	if r.getReader() == nil || r.getWriter() == nil {
		return nil, nil, errors.New(ErrRouterNoReaderOrWriter)
	}

	cntl := make(chan Item)
	var ch, tx chan Item

	ch, r.err = r.getReader().Read(cntl, nil)
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
	tx, r.err = r.getWriter().Write(ch, nil)
	if r.err != nil {
		return nil, nil, r.err
	}

	Log(r.getLogger(), "PipeRT", "started ---------------------------------")
	return cntl, tx, nil
}

// AddPipe implements Router
func (r *Pipeline) AddPipe(pipe Pipe) error {
	pipe.SetLogger(r.getLogger())
	r.pipes = append(r.pipes, pipe)
	return nil
}

// GetPipes implements Router
func (r *Pipeline) GetPipes() []Pipe {
	if r.pipes != nil {
		return r.pipes
	}
	return []Pipe{}
}
