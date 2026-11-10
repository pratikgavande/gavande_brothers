package bill_documents

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/phpdave11/gofpdf"
)

func RenderPDF(document Document) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(18, 18, 18)
	pdf.SetAutoPageBreak(true, 18)
	pdf.AddPage()

	pdf.SetFillColor(31, 78, 121)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(174, 14, "Gavande Bros", "", 1, "L", true, 0, "")
	pdf.SetTextColor(32, 32, 32)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(174, 11, "BILL", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(87, 7, "Bill number: "+strconv.FormatInt(document.Bill.BillNo, 10), "", 0, "L", false, 0, "")
	pdf.CellFormat(87, 7, "Date: "+document.Bill.Date.Format("02 Jan 2006 15:04"), "", 1, "R", false, 0, "")
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(174, 8, "Customer", "B", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	customer := document.Customer
	pdf.MultiCell(174, 6, fmt.Sprintf("%s\nMobile: %s\nGST number: %s\nEmail: %s\nAddress: %s",
		customer.Name, customer.Mobile, customer.GSTNo, customer.Email, customer.Address), "", "L", false)
	pdf.Ln(3)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(87, 7, "Site: "+document.Bill.SiteName, "", 0, "L", false, 0, "")
	pdf.CellFormat(87, 7, "Location: "+document.Bill.Location, "", 1, "R", false, 0, "")
	pdf.CellFormat(174, 7, "Engineer: "+document.Bill.EngineersName, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	pdf.SetFillColor(31, 78, 121)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 9)
	widths := []float64{64, 34, 34, 42}
	headers := []string{"Item", "Rate", "Unit", "Amount"}
	for i, header := range headers {
		pdf.CellFormat(widths[i], 8, header, "1", 0, "L", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(32, 32, 32)
	pdf.SetFont("Arial", "", 9)
	for _, item := range document.Items {
		if pdf.GetY() > 265 {
			pdf.AddPage()
		}
		pdf.CellFormat(widths[0], 8, item.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 8, item.Rate, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, item.Unit, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, strconv.FormatInt(int64(item.Amount), 10), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(4)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(174, 7, "GST status (stored): "+strconv.FormatInt(int64(document.Bill.GSTStatus), 10), "", 1, "R", false, 0, "")
	pdf.CellFormat(174, 7, "GST amount (stored): "+strconv.FormatInt(int64(document.Bill.GSTAmount), 10), "", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(174, 8, "Total amount (stored): "+strconv.FormatInt(int64(document.Bill.Amount), 10), "T", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(174, 7, "Amount received (stored): "+strconv.FormatInt(int64(document.Bill.AmountReceived), 10), "", 1, "R", false, 0, "")
	pdf.CellFormat(174, 7, "Amount pending (stored): "+strconv.FormatInt(int64(document.Bill.AmountPending), 10), "", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 100, 100)
	pdf.Ln(7)
	pdf.CellFormat(174, 6, "Amounts and GST values are shown as stored in the billing system.", "", 1, "L", false, 0, "")

	var output bytes.Buffer
	if err := pdf.Output(&output); err != nil {
		return nil, fmt.Errorf("render bill PDF: %w", err)
	}
	return output.Bytes(), nil
}
