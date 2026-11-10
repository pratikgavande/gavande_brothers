package bills

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCreateBillRequestToModel(t *testing.T) {
	billNo, customerID := int64(123), int32(7)
	date := time.Date(2026, 11, 10, 9, 0, 0, 0, time.FixedZone("IST", 5*60*60+30*60))
	gstStatus, gstAmount := int8(0), int32(0)
	amount, received, pending, itemCount := int32(1000), int32(0), int32(1000), int32(1)
	engineer, site, location := "Engineer", "Site", "Pune"
	status := int8(1)
	request := createBillRequest{
		BillNo: &billNo, CustomerID: &customerID, Date: &date,
		GSTStatus: &gstStatus, GSTAmount: &gstAmount, Amount: &amount,
		AmountReceived: &received, AmountPending: &pending,
		EngineersName: &engineer, SiteName: &site, Location: &location,
		ItemCount: &itemCount, Status: &status,
	}

	bill, ok := request.toModel()
	if !ok {
		t.Fatal("complete bill request was rejected")
	}
	if bill.BillNo != billNo || bill.CustomerID != customerID || bill.AmountPending != pending {
		t.Fatalf("unexpected bill mapping: %+v", bill)
	}

	request.CustomerID = nil
	if _, ok := request.toModel(); ok {
		t.Fatal("bill with missing customer_id was accepted")
	}
	request.CustomerID = &customerID
	customerID = 0
	if _, ok := request.toModel(); ok {
		t.Fatal("bill with non-positive customer_id was accepted")
	}
}

func TestBillPatchPreservesExplicitZeroAndRejectsOwnershipChanges(t *testing.T) {
	updates, err := billPatch(map[string]json.RawMessage{
		"gst_amount": json.RawMessage(`0`),
		"status":     json.RawMessage(`0`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if updates["gstAmount"] != int32(0) || updates["status"] != int8(0) {
		t.Fatalf("explicit zero values were not preserved: %#v", updates)
	}

	for _, field := range []string{"bill_id", "bill_no", "customer_id", "amount", "amount_received", "amount_pending"} {
		if _, err := billPatch(map[string]json.RawMessage{field: json.RawMessage(`1`)}); err == nil {
			t.Errorf("expected %s to be rejected by bill PATCH", field)
		}
	}
}
