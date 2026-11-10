package bill_items

import (
	"encoding/json"
	"testing"
)

func TestCreateBillItemRequestToModel(t *testing.T) {
	name, rate, unit := "Installation", "500", "service"
	amount, status := int32(1000), int8(1)
	date := "legacy-text-date"
	request := createBillItemRequest{
		Name: &name, Rate: &rate, Unit: &unit,
		Amount: &amount, Status: &status, Date: &date,
	}

	item, ok := request.toModel()
	if !ok {
		t.Fatal("complete bill item request was rejected")
	}
	if item.Name != name || item.Rate != rate || item.Date == nil || *item.Date != date {
		t.Fatalf("unexpected bill item mapping: %+v", item)
	}

	request.Date = nil
	if _, ok := request.toModel(); !ok {
		t.Fatal("nullable text date should be optional")
	}
	request.Unit = nil
	if _, ok := request.toModel(); ok {
		t.Fatal("bill item with missing non-null unit field was accepted")
	}
}

func TestBillItemPatchAllowsClearingNullableDate(t *testing.T) {
	updates, err := billItemPatch(map[string]json.RawMessage{
		"date":   json.RawMessage(`null`),
		"amount": json.RawMessage(`0`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if updates["date"] != nil || updates["amount"] != int32(0) {
		t.Fatalf("unexpected patch values: %#v", updates)
	}
	if _, err := billItemPatch(map[string]json.RawMessage{"bill_id": json.RawMessage(`9`)}); err == nil {
		t.Fatal("bill_id must not be changed by item PATCH")
	}
}
