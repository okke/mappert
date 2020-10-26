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

func TestMap2StructWithSliceOfStrings(t *testing.T) {

	mapping := map[string]interface{}{"Names": []string{"chi", "potle"}}

	into := struct {
		Names []string
	}{}

	Map2Struct(mapping, &into)

	if len(into.Names) != 2 {
		t.Error("exptected 2 names")
	}

	if into.Names[0] != "chi" {
		t.Error("expected first name to be chi, not", into.Names[0])
	}

}

func TestMap2StructWithSliceOfStucts(t *testing.T) {

	type nameStruct struct {
		First  string
		Second string
	}

	mapping := map[string]interface{}{
		"Names": []map[string]interface{}{
			{"First": "chi", "Second": "potle"},
			{"First": "jala", "Second": "peno"},
		},
	}

	into := struct {
		Names []nameStruct
	}{}

	Map2Struct(mapping, &into)

	if len(into.Names) != 2 {
		t.Error("exptected 2 names")
	}

	if name := into.Names[0].First; name != "chi" {
		t.Error("expected first name of first element to be chi, not", name)
	}

	if name := into.Names[0].Second; name != "potle" {
		t.Error("expected second name of first element to be potle, not", name)
	}

	if name := into.Names[1].First; name != "jala" {
		t.Error("expected first name of second element to be jala, not", name)
	}

	if name := into.Names[1].Second; name != "peno" {
		t.Error("expected second name of second element to be peno, not", name)
	}
}
