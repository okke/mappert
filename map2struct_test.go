package mappert

import "testing"

func TestMap2StructWithPrimitives(t *testing.T) {

	mapping := map[string]interface{}{"Name": "chipotle", "SKU": 3000}

	into := struct {
		Name string
		SKU  int
	}{}

	Map2Struct(mapping, &into)

	if into.Name != "chipotle" {
		t.Error("expected chipotle, not", into.Name)
	}

	if into.SKU != 3000 {
		t.Error("expected an sku of 3000, not", into.SKU)
	}

}

func TestMap2StructWithInnerStruct(t *testing.T) {

	type nameStruct struct {
		First  string
		Second string
	}

	mapping := map[string]interface{}{"Name": map[string]interface{}{"First": "chi", "Second": "potle"}}

	into := struct {
		Name nameStruct
	}{}

	Map2Struct(mapping, &into)

	if into.Name.First != "chi" {
		t.Error("expected first name to be chi, not", into.Name.First)
	}

	if into.Name.Second != "potle" {
		t.Error("expected fsecond name to be potle, not", into.Name.Second)
	}

}
