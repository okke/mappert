package mappert

import (
	"testing"
)

func TestStruct2MapFollowsPointerToStruct(t *testing.T) {

	if _, err := Struct2Map(&struct{ Name string }{
		Name: "chipotle",
	}); err != nil {
		t.Error(err)
	}
}

func TestStruct2MapExpectsStruct(t *testing.T) {

	if _, err := Struct2Map(&[]int{}); err == nil {
		t.Error("expect an error when providing a pointer to something not a struct")
	}
}

func TestStruct2MapConvertsPrimitiveFieldNames(t *testing.T) {

	mapping, _ := Struct2Map(struct {
		Name string
		SKU  int
	}{
		Name: "chipotle", SKU: 3000,
	})

	if mapping["Name"] != "chipotle" {
		t.Error("expected chipotle as name, not", mapping["Name"])
	}

	if mapping["SKU"] != 3000 {
		t.Error("expected sku of 3000 as name, not", mapping["Name"])
	}

}

func TestStruct2MapConvertsStructFieldToMap(t *testing.T) {

	type nameStruct struct {
		First  string
		Second string
	}

	mapping, _ := Struct2Map(struct {
		Name nameStruct
	}{
		Name: nameStruct{First: "chi", Second: "potle"},
	})

	if found := mapping["Name"].(map[string]interface{})["First"]; found != "chi" {
		t.Error("expected chi as first name, not", found)
	}

	if found := mapping["Name"].(map[string]interface{})["Second"]; found != "potle" {
		t.Error("expected pottle as second name, not", found)
	}

}

func TestStruct2MapConvertsSliceFieldToMap(t *testing.T) {

	type nameStruct struct {
		First  string
		Second string
	}

	mapping, _ := Struct2Map(struct {
		Names []nameStruct
	}{
		Names: []nameStruct{{First: "chi", Second: "potle"}, {First: "jala", Second: "peno"}},
	})

	if found := mapping["Names"].([]interface{})[0].(map[string]interface{})["First"]; found != "chi" {
		t.Error("expected chi as first name of first element, not", found)
	}

	if found := mapping["Names"].([]interface{})[1].(map[string]interface{})["Second"]; found != "peno" {
		t.Error("expected peno as second name of second element, not", found)
	}

}

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
