package object_test

import (
	"testing"

	"go-pdf/object"
	"go-pdf/testutil"
)

func TestArray(t *testing.T) {
	in := "[549 3.14 false (Ralph) /SomeName]"
	rdr := testutil.NewPdf(in)

	is, err := object.ReadArray(rdr, nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(is)

	if len(is.Elems) != 5 {
		t.Errorf("wrong array length, expected 5, got %d", len(is.Elems))
		t.FailNow()
	}

	e0, ok := is.Elems[0].(object.Integer)
	if !ok {
		t.Error("element 0 is not an integer!")
	}

	if int(e0) != 549 {
		t.Errorf("element 0: expected int(549), got int(%d)", e0)
	}

	e1, ok := is.Elems[1].(object.Floating)
	if !ok {
		t.Error("element 1 is not a real!")
	}

	if e1 != 3.14 {
		t.Errorf("element 1: expected float(3.14), got float(%f)", e1)
	}

	e2, ok := is.Elems[2].(object.Bool)
	if !ok {
		t.Error("element 2 is not a bool!")
	}

	if bool(e2) != false {
		t.Error("element 2: expected bool(false), got bool(true)")
	}

	e3, ok := is.Elems[3].(*object.String)
	if !ok {
		t.Error("element 3 is not a string!")
	}

	if e3.S != "Ralph" {
		t.Errorf("element 3: expected string(Ralph), got string(%s)", e3.S)
	}

	e4, ok := is.Elems[4].(*object.Name)
	if !ok {
		t.Error("element 4 is not a name!")
	}

	if e4.Name != "SomeName" {
		t.Errorf("element 4: expected name(SomeName), got name(%s)", e4.Name)
	}
}

func TestArrayWithIndirectObj(t *testing.T) {
	in := `[12 1 22 obj
	(cheap)
	endobj 34]`
	rdr := testutil.NewPdf(in)

	is, err := object.ReadArray(rdr, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(is)

	if len(is.Elems) != 3 {
		t.Errorf("wrong array length, expected 3, got %d", len(is.Elems))
	}

	e0, ok := is.Elems[0].(object.Integer)
	if !ok {
		t.Error("element 0 is not an integer")
	}

	if e0 != 12 {
		t.Errorf("element 0: expected int(12), got int(%d)", e0)
	}

	e1, ok := is.Elems[1].(*object.IndirectVal)
	if !ok {
		t.Error("element 1 is not an integer")
	}

	if e1.ID.ID != 1 || e1.ID.Gen != 22 {
		t.Errorf("element 1: wrong indirect id, expected 1/22, got %s", e1.ID.String())
	}

	e1v, ok := e1.Value().(*object.String)
	if !ok {
		t.Errorf("element 1: indirect value is not a string")
	}

	if e1v.S != "cheap" {
		t.Errorf("element 1: value expected string(cheap), got string(%s)", e1v.S)
	}

	e2, ok := is.Elems[2].(object.Integer)
	if !ok {
		t.Error("element 2 is not an integer")
	}

	if e2 != 34 {
		t.Errorf("element 2: expected int(34), got int(%d)", e2)
	}
}
