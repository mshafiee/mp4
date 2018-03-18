package filter

import (
	"io"
	"github.com/mshafiee/mp4"
	"github.com/mshafiee/mp4/box"
)

type Filter interface {
	// Updates the moov box
	FilterMoov(m *box.MoovBox) error
	// Filters the Mdat data and writes it to w
	FilterMdat(w io.Writer, m *box.MdatBox) error
}

// EncodeFiltered encodes a media to a writer, filtering the media using the specified filter
func EncodeFiltered(w io.Writer, m *mp4.MP4, f Filter) error {
	err := m.Ftyp.Encode(w)
	if err != nil {
		return err
	}
	err = f.FilterMoov(m.Moov)
	if err != nil {
		return err
	}
	err = m.Moov.Encode(w)
	if err != nil {
		return err
	}
	for _, b := range m.Boxes() {
		if b.Type() != "ftyp" && b.Type() != "moov" && b.Type() != "mdat" {
			err = b.Encode(w)
			if err != nil {
				return err
			}
		}
	}
	return f.FilterMdat(w, m.Mdat)
}
