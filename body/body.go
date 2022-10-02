package body

import (
	"go-pdf/object"
	"go-pdf/shared"
)

type Body struct {
	Obj []*object.IndirectVal
}

func NewBody() *Body {
	return &Body{}
}

func (b *Body) Resolve(i shared.ID) shared.Object {
	for _, obj := range b.Obj {
		if obj.ID.Equal(i) {
			return obj.Value()
		}
	}

	return nil
}
