package filter

import (
	"io"
	"github.com/mshafiee/mp4/box"
)

type noopFilter struct{}

// Noop returns a filter that does nothing
func Noop() Filter {
	return &noopFilter{}
}

func (f *noopFilter) FilterMoov(m *box.MoovBox) error {
	return nil
}

func (f *noopFilter) FilterMdat(w io.Writer, m *box.MdatBox) error {
	err := box.EncodeHeader(m, w)
	if err == nil {
		_, err = io.Copy(w, m.Reader())
	}
	return err
}
