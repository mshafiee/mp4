package mp4

import (
	"io"
	"github.com/mshafiee/mp4/box"
)

/*
A MPEG-4 media contains three main boxes :
	ftyp : the file type box
	moov : the movie box (meta-data)
	mdat : the media data (chunks and samples)
Other boxes can also be present (pdin, moof, mfra, free, ...), but are not decoded.
*/

type MP4 struct {
	Ftyp  *box.FtypBox
	Moov  *box.MoovBox
	Mdat  *box.MdatBox
	boxes []box.Box
}

// Decode decodes a media from a Reader
func Decode(r io.Reader) (*MP4, error) {
	var v = new(MP4)

LoopBoxes:
	for {
		h, err := box.DecodeHeader(r)
		if err != nil {
			return nil, err
		}
		b, err := box.DecodeBox(h, r)
		if err != nil {
			return nil, err
		}
		v.boxes = append(v.boxes, b)
		switch h.Type {
		case "ftyp":
			v.Ftyp = b.(*box.FtypBox)
		case "moov":
			v.Moov = b.(*box.MoovBox)
		case "mdat":
			v.Mdat = b.(*box.MdatBox)
			v.Mdat.ContentSize = h.Size - box.BoxHeaderSize
			break LoopBoxes
		}
	}
	return v, nil
}

// Dump displays some information about a media
func (m *MP4) Dump() {
	m.Ftyp.Dump()
	m.Moov.Dump()
}

// Boxes lists the top-level boxes from a media
func (m *MP4) Boxes() []box.Box {
	return m.boxes
}

// Encode encodes a media to a Writer
func (m *MP4) Encode(w io.Writer) error {
	err := m.Ftyp.Encode(w)
	if err != nil {
		return err
	}
	err = m.Moov.Encode(w)
	if err != nil {
		return err
	}
	for _, b := range m.boxes {
		if b.Type() != "ftyp" && b.Type() != "moov" {
			err = b.Encode(w)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
