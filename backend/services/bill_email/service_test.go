package bill_email

import (
	"bytes"
	"net/mail"
	"strings"
	"testing"
)

func TestBuildMessageIncludesPDFAttachment(t *testing.T) {
	from, _ := mail.ParseAddress("billing@example.com")
	to, _ := mail.ParseAddress("customer@example.com")
	message, err := buildMessage(from, to, 42, []byte("%PDF-1.4 test"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(message, []byte("application/pdf")) ||
		!bytes.Contains(message, []byte("gavande-bros-bill-42.pdf")) ||
		!strings.Contains(string(message), "customer@example.com") {
		t.Fatalf("MIME email is missing expected attachment metadata or recipient")
	}
}
