package main

import (
	"io/ioutil"
	"log"
	"text/template"
)

type ConfirmMailVars struct {
	From      string
	To        string
	ConfirmID string
}

type PasswordMailVars struct {
	From     string
	To       string
	Password string
}

var TemplateSignup *template.Template
var TemplateChangeEmail *template.Template
var TemplateResetPassword *template.Template
var TemplateNewPassword *template.Template

func readMailTemplatesFromFile() {
	content, err := ioutil.ReadFile(GetConfig().TemplateChangeEmail)
	if err != nil {
		log.Fatal(err)
	}
	TemplateChangeEmail, _ = template.New("TemplateChangeEmail").Parse(string(content))

	content, err = ioutil.ReadFile(GetConfig().TemplateSignup)
	if err != nil {
		log.Fatal(err)
	}
	TemplateSignup, _ = template.New("TemplateSignup").Parse(string(content))

	content, err = ioutil.ReadFile(GetConfig().TemplateResetPassword)
	if err != nil {
		log.Fatal(err)
	}
	TemplateResetPassword, _ = template.New("TemplateResetPassword").Parse(string(content))

	content, err = ioutil.ReadFile(GetConfig().TemplateNewPassword)
	if err != nil {
		log.Fatal(err)
	}
	TemplateNewPassword, _ = template.New("TemplateNewPassword").Parse(string(content))

}
