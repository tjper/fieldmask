package fieldmask

import (
	"google.golang.org/genproto/protobuf/field_mask"
)

type (
	GetMasker interface {
		GetMask() *field_mask.FieldMask
	}
	Applicator interface {
		Apply()
	}
	MaskApplicator interface {
		ApplyMask([]MaskUpdate)
	}
)

type MaskUpdate func() error

type Update struct {
	updates    []MaskUpdate
	applicator MaskApplicator
}

func NewUpdate(a MaskApplicator) Applicator {
	return &Update{
		applicator: a,
	}
}

func (u *Update) SetMaskFunc(m GetMasker, path string, update MaskUpdate) {
	for _, mask := range m.GetMask().GetPaths() {
		if mask == path {
			u.updates = append(u.updates, update)
		}
	}
}

func (u *Update) Apply() {
	u.applicator.ApplyMask(u.updates)
}
