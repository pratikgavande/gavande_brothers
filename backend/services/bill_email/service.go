package bill_email

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"billing-app/config"
)

var ErrNotConfigured = errors.New("SMTP is not configured")
var ErrInvalidEmail = errors.New("invalid email address")

type Service struct {
	smtp config.SMTPConfig
}

func NewService(smtpConfig config.SMTPConfig) *Service {
	return &Service{smtp: smtpConfig}
}

func (s *Service) SendBill(ctx context.Context, recipient string, billNo int64, pdf []byte) error {
	if !s.smtp.Configured() {
		return ErrNotConfigured
	}

	from, err := mail.ParseAddress(s.smtp.From)
	if err != nil {
		return fmt.Errorf("invalid SMTP_FROM address")
	}
	to, err := mail.ParseAddress(recipient)
	if err != nil || strings.TrimSpace(to.Address) != strings.TrimSpace(recipient) {
		return ErrInvalidEmail
	}
	port, err := strconv.Atoi(s.smtp.Port)
	if err != nil || port < 1 || port > 65535 {
		return errors.New("invalid SMTP port")
	}

	message, err := buildMessage(from, to, billNo, pdf)
	if err != nil {
		return err
	}
	return s.send(ctx, port, from.Address, to.Address, message)
}

func buildMessage(from, to *mail.Address, billNo int64, pdf []byte) ([]byte, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	headers := make(textproto.MIMEHeader)
	headers.Set("From", from.String())
	headers.Set("To", to.String())
	headers.Set("Subject", fmt.Sprintf("Your Gavande Bros bill #%d", billNo))
	headers.Set("MIME-Version", "1.0")
	headers.Set("Content-Type", "multipart/mixed; boundary="+writer.Boundary())
	var message bytes.Buffer
	for key, values := range headers {
		for _, value := range values {
			fmt.Fprintf(&message, "%s: %s\r\n", key, value)
		}
	}
	message.WriteString("\r\n")

	textHeader := make(textproto.MIMEHeader)
	textHeader.Set("Content-Type", "text/plain; charset=utf-8")
	textHeader.Set("Content-Transfer-Encoding", "quoted-printable")
	textPart, err := writer.CreatePart(textHeader)
	if err != nil {
		return nil, err
	}
	if _, err := textPart.Write([]byte("Please find your Gavande Bros bill attached.\r\n")); err != nil {
		return nil, err
	}

	pdfHeader := make(textproto.MIMEHeader)
	pdfHeader.Set("Content-Type", "application/pdf")
	pdfHeader.Set("Content-Transfer-Encoding", "base64")
	filename := mime.FormatMediaType("attachment", map[string]string{
		"filename": fmt.Sprintf("gavande-bros-bill-%d.pdf", billNo),
	})
	pdfHeader.Set("Content-Disposition", filename)
	pdfPart, err := writer.CreatePart(pdfHeader)
	if err != nil {
		return nil, err
	}
	encoder := base64.NewEncoder(base64.StdEncoding, pdfPart)
	if _, err := encoder.Write(pdf); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	message.Write(body.Bytes())
	return message.Bytes(), nil
}

func (s *Service) send(ctx context.Context, port int, from, to string, message []byte) error {
	address := net.JoinHostPort(s.smtp.Host, strconv.Itoa(port))
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect to SMTP server: %w", err)
	}
	defer conn.Close()
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
	}

	if port == 465 {
		tlsConn := tls.Client(conn, &tls.Config{
			ServerName: s.smtp.Host,
			MinVersion: tls.VersionTLS12,
		})
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			return fmt.Errorf("establish SMTP TLS: %w", err)
		}
		conn = tlsConn
	}

	client, err := smtp.NewClient(conn, s.smtp.Host)
	if err != nil {
		return fmt.Errorf("start SMTP session: %w", err)
	}
	defer client.Close()

	if port != 465 {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return errors.New("SMTP server does not support required STARTTLS")
		}
		if err := client.StartTLS(&tls.Config{
			ServerName: s.smtp.Host,
			MinVersion: tls.VersionTLS12,
		}); err != nil {
			return fmt.Errorf("establish SMTP STARTTLS: %w", err)
		}
	}
	if err := client.Auth(smtp.PlainAuth("", s.smtp.Username, s.smtp.Password, s.smtp.Host)); err != nil {
		return fmt.Errorf("authenticate with SMTP server: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("set SMTP sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("set SMTP recipient: %w", err)
	}
	data, err := client.Data()
	if err != nil {
		return fmt.Errorf("open SMTP message: %w", err)
	}
	if _, err := data.Write(message); err != nil {
		_ = data.Close()
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := data.Close(); err != nil {
		return fmt.Errorf("SMTP server rejected message: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("complete SMTP delivery: %w", err)
	}
	return nil
}
