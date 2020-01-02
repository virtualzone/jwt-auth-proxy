package main

import (
	"io"
	"testing"
)

var smtpMockContent smtpDialerMockContent

func TestSendmail(t *testing.T) {
	SendMail("foo@bar.com", "Hello World!")
	checkTestString(t, "no-reply@localhost", smtpMockContent.FromValue)
	checkTestString(t, "foo@bar.com", smtpMockContent.RcptValue)
	checkTestString(t, "Hello World!", smtpMockContent.Buffer.DataValue)
}

type smtpDialerMockContent struct {
	HelloValue string
	FromValue  string
	RcptValue  string
	Buffer     writeCloserMock
}

type smtpDialerMock struct {
}

func (r *smtpDialerMock) Close() error {
	return nil
}
func (r *smtpDialerMock) Hello(localName string) error {
	smtpMockContent = smtpDialerMockContent{}
	smtpMockContent.HelloValue = localName
	return nil
}
func (r *smtpDialerMock) Mail(from string) error {
	smtpMockContent.FromValue = from
	return nil
}
func (r *smtpDialerMock) Rcpt(to string) error {
	smtpMockContent.RcptValue = to
	return nil
}

func (r *smtpDialerMock) Data() (io.WriteCloser, error) {
	smtpMockContent.Buffer = writeCloserMock{}
	return &smtpMockContent.Buffer, nil
}

type writeCloserMock struct {
	DataValue string
}

func (r *writeCloserMock) Write(buf []byte) (int, error) {
	r.DataValue += string(buf)
	return len(buf), nil
}

func (r *writeCloserMock) Close() error {
	return nil
}
