package bill_documents

import (
	"bytes"
	"testing"
	"time"

	"billing-app/models"
)

func TestRenderPDFUsesBillDocument(t *testing.T) {
	pdf, err := RenderPDF(Document{
		Bill: models.Bill{
			BillNo: 8,
			Date:   time.Date(2026, 11, 10, 9, 0, 0, 0, time.UTC),
			Amount: 1234,
		},
		Customer: models.Customer{Name: "Test Customer"},
		Items:    []models.BillItem{{Name: "Service", Rate: "1234", Unit: "one", Amount: 1234}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatalf("generated output is not a PDF: %q", pdf[:min(len(pdf), 10)])
	}
}
