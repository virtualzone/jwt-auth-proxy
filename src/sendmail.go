package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/smtp"
)

var (
	smtpClient = func(addr string) (dialer, error) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		host, _, _ := net.SplitHostPort(addr)
		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
)

func SendMail(recv string, body string) (dialer, error) {
	c, err := smtpClient(GetConfig().SMTPServer)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer c.Close()
	err = c.Mail(GetConfig().SMTPSenderAddr)
	if err != nil {
		log.Println(err)
		return c, err
	}
	err = c.Rcpt(recv)
	if err != nil {
		log.Println(err)
		return c, err
	}
	wc, err := c.Data()
	if err != nil {
		log.Println(err)
		return c, err
	}
	defer wc.Close()
	buf := bytes.NewBufferString(body)
	if _, err := buf.WriteTo(wc); err != nil {
		log.Println(err)
		return c, err
	}
	return c, nil
}

type dialer interface {
	Close() error
	Hello(localName string) error
	Mail(from string) error
	Rcpt(to string) error
	Data() (io.WriteCloser, error)
}
