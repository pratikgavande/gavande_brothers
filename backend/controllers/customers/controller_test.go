package customers

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCreateCustomerRequestToModel(t *testing.T) {
	name, mobile, gstNo, email, address := "Customer", "123", "", "", ""
	status := int8(0)
	request := createCustomerRequest{
		Name: &name, Mobile: &mobile, GSTNo: &gstNo,
		Email: &email, Address: &address, Status: &status,
	}

	customer, ok := request.toModel()
	if !ok {
		t.Fatal("valid schema-complete customer was rejected")
	}
	if customer.Name != name || customer.Mobile != mobile || customer.Status != status {
		t.Fatalf("unexpected customer mapping: %+v", customer)
	}

	request.Email = nil
	if _, ok := request.toModel(); ok {
		t.Fatal("customer with a missing non-null email field was accepted")
	}
}

func TestCreateCustomerRequestValidatesRequiredTextAndLengths(t *testing.T) {
	name, mobile, gstNo, email, address := " ", "123", "", "", ""
	status := int8(1)
	request := createCustomerRequest{
		Name: &name, Mobile: &mobile, GSTNo: &gstNo,
		Email: &email, Address: &address, Status: &status,
	}
	if _, ok := request.toModel(); ok {
		t.Fatal("blank customer name was accepted")
	}

	name = "Customer"
	email = strings.Repeat("x", 251)
	if _, ok := request.toModel(); ok {
		t.Fatal("email longer than the database varchar column was accepted")
	}
}

func TestCustomerPatchOnlyIncludesSuppliedValues(t *testing.T) {
	updates, err := customerPatch(map[string]json.RawMessage{
		"email": json.RawMessage(`""`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(updates) != 1 || updates["email"] != "" {
		t.Fatalf("unexpected sparse customer updates: %#v", updates)
	}
}

func TestCustomerPatchRejectsNullAndImmutableFields(t *testing.T) {
	for _, fields := range []map[string]json.RawMessage{
		{"name": json.RawMessage(`null`)},
		{"customer_id": json.RawMessage(`1`)},
		{"name": json.RawMessage(`"  "`)},
	} {
		if _, err := customerPatch(fields); err == nil {
			t.Fatalf("expected patch validation error for %#v", fields)
		}
	}
}
