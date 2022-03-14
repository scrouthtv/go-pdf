package body

import "github.com/scrouthtv/go-pdf/object"

type Body struct {
	Obj []object.Indirect
}

func (b *Body) Resolve(i *object.IndirectRef) object.IndirectVal {
	// TODO
	panic("not implemented")
}
