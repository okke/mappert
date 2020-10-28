package mappert

import (
	"strings"
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

func TestStruct2MapAppliesNameConvertor(t *testing.T) {

	mapping, _ := Struct2Map(struct {
		Name string
		Sku  int
	}{
		Name: "chipotle", Sku: 3000,
	}, ConvertNamesUsing(func(name string) string {
		return strings.ToUpper(name)
	}))

	if mapping["NAME"] != "chipotle" {
		t.Error("expected chipotle as name, not", mapping["Name"])
	}

	if mapping["SKU"] != 3000 {
		t.Error("expected sku of 3000 as name, not", mapping["Name"])
	}

}

func TestStruct2MapAppliesFieldIgnoration(t *testing.T) {

	mapping, _ := Struct2Map(struct {
		Name string
		SKU  int
		Like bool
	}{
		Name: "chipotle", SKU: 3000, Like: true,
	}, IgnoreFields("Name", "Like"))

	if _, found := mapping["Name"]; found {
		t.Error("Did not expect a name in", mapping)
	}

	if _, found := mapping["SKU"]; !found {
		t.Error("Did expect an SKU in", mapping)
	}

	if _, found := mapping["Like"]; found {
		t.Error("Did expect an SKU in", mapping)
	}

}

func TestStruct2MapAppliesFieldIgnorationAndNameConversion(t *testing.T) {

	mapping, _ := Struct2Map(struct {
		Name string
		SKU  int
		Like bool
	}{
		Name: "chipotle", SKU: 3000, Like: true,
	}, IgnoreFields("Name", "Like"), ConvertNamesUsing(func(name string) string {
		return strings.ToLower(name)
	}))

	if _, found := mapping["Name"]; found {
		t.Error("Did not expect a name in", mapping)
	}

	if _, found := mapping["sku"]; !found {
		t.Error("Did expect an sku in", mapping)
	}

	if _, found := mapping["Like"]; found {
		t.Error("Did expect an SKU in", mapping)
	}

}
