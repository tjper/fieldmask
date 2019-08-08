package fieldmask

import (
	"google.golang.org/genproto/protobuf/field_mask"
)

type GetMasker interface {
	GetMask() *field_mask.FieldMask
}

type Applicator interface {
	Apply([]MaskUpdate)
}

type MaskUpdate func() error

type Update struct {
	updates []MaskUpdate
}

func (u *Update) SetMaskFunc(m GetMasker, path string, update MaskUpdate) {
	for _, mask := range m.GetMask().GetPaths() {
		if mask == path {
			u.updates = append(u.updates, update)
		}
	}
}

func (u *Update) Apply(a Applicator) {
	a.Apply(u.updates)
}
