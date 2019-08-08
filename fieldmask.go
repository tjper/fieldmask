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
		SetPathFunc(string, MaskUpdate)
	}
	MaskApplicator interface {
		ApplyMask([]MaskUpdate)
	}
)

type MaskUpdate func() error

type Update struct {
	updates    []MaskUpdate
	applicator MaskApplicator
	masker     GetMasker
}

func NewUpdate(m GetMasker, a MaskApplicator) Applicator {
	return &Update{
		applicator: a,
		masker:     m,
	}
}

func (u *Update) SetPathFunc(path string, update MaskUpdate) {
	for _, mask := range u.masker.GetMask().GetPaths() {
		if mask == path {
			u.updates = append(u.updates, update)
		}
	}
}

func (u *Update) Apply() {
	u.applicator.ApplyMask(u.updates)
}
